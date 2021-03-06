// Copyright 2015 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package logs

import (
	"time"

	ds "github.com/luci/gae/service/datastore"
	"github.com/luci/luci-go/common/config"
	log "github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/common/retry"
	"github.com/luci/luci-go/grpc/grpcutil"
	"github.com/luci/luci-go/logdog/api/endpoints/coordinator/logs/v1"
	"github.com/luci/luci-go/logdog/api/logpb"
	"github.com/luci/luci-go/logdog/appengine/coordinator"
	"github.com/luci/luci-go/logdog/common/storage"
	"github.com/luci/luci-go/logdog/common/storage/archive"
	"github.com/luci/luci-go/logdog/common/types"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

const (
	// getInitialArraySize is the initial amount of log slots to allocate for a
	// Get request.
	getInitialArraySize = 256

	// getBytesLimit is the maximum amount of data that we are willing to query.
	//
	// We will limit byte responses to 16MB, based on the following constraints:
	//	- AppEngine cannot respond with more than 32MB of data. This includes JSON
	//	  overhead, including notation and base64 data expansion.
	//	- `urlfetch`, which is used for Google Cloud Storage (archival) responses,
	//	  cannot handle responses larger than 32MB.
	getBytesLimit = 16 * 1024 * 1024
)

// Get returns state and log data for a single log stream.
func (s *server) Get(c context.Context, req *logdog.GetRequest) (*logdog.GetResponse, error) {
	return s.getImpl(c, req, false)
}

// Tail returns the last log entry for a given log stream.
func (s *server) Tail(c context.Context, req *logdog.TailRequest) (*logdog.GetResponse, error) {
	r := logdog.GetRequest{
		Project: req.Project,
		Path:    req.Path,
		State:   req.State,
	}
	return s.getImpl(c, &r, true)
}

// getImpl is common code shared between Get and Tail endpoints.
func (s *server) getImpl(c context.Context, req *logdog.GetRequest, tail bool) (*logdog.GetResponse, error) {
	log.Fields{
		"project": req.Project,
		"path":    req.Path,
		"index":   req.Index,
		"tail":    tail,
	}.Debugf(c, "Received get request.")

	path := types.StreamPath(req.Path)
	if err := path.Validate(); err != nil {
		log.WithError(err).Errorf(c, "Invalid path supplied.")
		return nil, grpcutil.Errf(codes.InvalidArgument, "invalid path value")
	}

	ls := &coordinator.LogStream{ID: coordinator.LogStreamID(path)}
	lst := ls.State(c)
	log.Fields{
		"id": ls.ID,
	}.Debugf(c, "Loading stream.")

	if err := ds.Get(c, ls, lst); err != nil {
		if isNoSuchEntity(err) {
			log.Errorf(c, "Log stream does not exist.")
			return nil, grpcutil.Errf(codes.NotFound, "path not found")
		}

		log.WithError(err).Errorf(c, "Failed to look up log stream.")
		return nil, grpcutil.Internal
	}

	// If this log entry is Purged and we're not admin, pretend it doesn't exist.
	if ls.Purged {
		if authErr := coordinator.IsAdminUser(c); authErr != nil {
			log.Fields{
				log.ErrorKey: authErr,
			}.Warningf(c, "Non-superuser requested purged log.")
			return nil, grpcutil.Errf(codes.NotFound, "path not found")
		}
	}

	// If nothing was requested, return nothing.
	resp := logdog.GetResponse{}
	if !(req.State || tail) && req.LogCount < 0 {
		return &resp, nil
	}

	if req.State {
		resp.State = buildLogStreamState(ls, lst)

		var err error
		resp.Desc, err = ls.DescriptorValue()
		if err != nil {
			log.WithError(err).Errorf(c, "Failed to deserialize descriptor protobuf.")
			return nil, grpcutil.Internal
		}
	}

	// Retrieve requested logs from storage, if requested.
	if tail || req.LogCount >= 0 {
		var err error
		resp.Logs, err = s.getLogs(c, req, tail, ls, lst)
		if err != nil {
			log.WithError(err).Errorf(c, "Failed to get logs.")
			return nil, grpcutil.Internal
		}
	}

	log.Fields{
		"logCount": len(resp.Logs),
	}.Debugf(c, "Get request completed successfully.")
	return &resp, nil
}

func (s *server) getLogs(c context.Context, req *logdog.GetRequest, tail bool, ls *coordinator.LogStream,
	lst *coordinator.LogStreamState) ([]*logpb.LogEntry, error) {
	byteLimit := int(req.ByteCount)
	if byteLimit <= 0 || byteLimit > getBytesLimit {
		byteLimit = getBytesLimit
	}

	svc := coordinator.GetServices(c)
	var st storage.Storage
	if !lst.ArchivalState().Archived() {
		log.Debugf(c, "Log is not archived. Fetching from intermediate storage.")

		// Logs are not archived. Fetch from intermediate storage.
		var err error
		st, err = svc.IntermediateStorage(c)
		if err != nil {
			return nil, err
		}
	} else {
		log.Fields{
			"indexURL":    lst.ArchiveIndexURL,
			"streamURL":   lst.ArchiveStreamURL,
			"archiveTime": lst.ArchivedTime,
		}.Debugf(c, "Log is archived. Fetching from archive storage.")

		var err error
		gs, err := svc.GSClient(c)
		if err != nil {
			log.WithError(err).Errorf(c, "Failed to create Google Storage client.")
			return nil, err
		}
		defer func() {
			if err := gs.Close(); err != nil {
				log.WithError(err).Warningf(c, "Failed to close Google Storage client.")
			}
		}()

		st, err = archive.New(c, archive.Options{
			IndexURL:  lst.ArchiveIndexURL,
			StreamURL: lst.ArchiveStreamURL,
			Client:    gs,
			Cache:     svc.StorageCache(),
		})
		if err != nil {
			log.WithError(err).Errorf(c, "Failed to create Google Storage storage instance.")
			return nil, err
		}
	}
	defer st.Close()

	project, path := coordinator.Project(c), ls.Path()

	var fetchedLogs []*logpb.LogEntry
	var err error
	if tail {
		fetchedLogs, err = getTail(c, st, project, path)
	} else {
		fetchedLogs, err = getHead(c, req, st, project, path, byteLimit)
	}
	if err != nil {
		log.WithError(err).Errorf(c, "Failed to fetch log records.")
		return nil, err
	}

	return fetchedLogs, nil
}

func getHead(c context.Context, req *logdog.GetRequest, st storage.Storage, project config.ProjectName,
	path types.StreamPath, byteLimit int) ([]*logpb.LogEntry, error) {
	log.Fields{
		"project":       project,
		"path":          path,
		"index":         req.Index,
		"count":         req.LogCount,
		"bytes":         req.ByteCount,
		"noncontiguous": req.NonContiguous,
	}.Debugf(c, "Issuing Get request.")

	// Allocate result logs array.
	logCount := int(req.LogCount)
	asz := getInitialArraySize
	if logCount > 0 && logCount < asz {
		asz = logCount
	}
	logs := make([]*logpb.LogEntry, 0, asz)

	sreq := storage.GetRequest{
		Project: project,
		Path:    path,
		Index:   types.MessageIndex(req.Index),
		Limit:   logCount,
	}

	var ierr error
	count := 0
	err := retry.Retry(c, retry.TransientOnly(retry.Default), func() error {
		// Issue the Get request. This may return a transient error, in which case
		// we will retry.
		return st.Get(sreq, func(e *storage.Entry) bool {
			var le *logpb.LogEntry
			if le, ierr = e.GetLogEntry(); ierr != nil {
				return false
			}

			if count > 0 && byteLimit-len(e.D) < 0 {
				// Not the first log, and we've exceeded our byte limit.
				return false
			}

			sidx, _ := e.GetStreamIndex() // GetLogEntry succeeded, so this must.
			if !(req.NonContiguous || sidx == sreq.Index) {
				return false
			}

			logs = append(logs, le)
			sreq.Index = sidx + 1
			byteLimit -= len(e.D)
			count++
			return !(logCount > 0 && count >= logCount)
		})
	}, func(err error, delay time.Duration) {
		log.Fields{
			log.ErrorKey:   err,
			"delay":        delay,
			"initialIndex": req.Index,
			"nextIndex":    sreq.Index,
			"count":        len(logs),
		}.Warningf(c, "Transient error while loading logs; retrying.")
	})
	switch err {
	case nil:
		if ierr != nil {
			log.WithError(ierr).Errorf(c, "Bad log entry data.")
			return nil, ierr
		}
		return logs, nil

	case storage.ErrDoesNotExist:
		return nil, nil

	default:
		log.Fields{
			log.ErrorKey:   err,
			"initialIndex": req.Index,
			"nextIndex":    sreq.Index,
			"count":        len(logs),
		}.Errorf(c, "Failed to execute range request.")
		return nil, err
	}
}

func getTail(c context.Context, st storage.Storage, project config.ProjectName, path types.StreamPath) (
	[]*logpb.LogEntry, error) {
	log.Fields{
		"project": project,
		"path":    path,
	}.Debugf(c, "Issuing Tail request.")

	var e *storage.Entry
	err := retry.Retry(c, retry.TransientOnly(retry.Default), func() (err error) {
		e, err = st.Tail(project, path)
		return
	}, func(err error, delay time.Duration) {
		log.Fields{
			log.ErrorKey: err,
			"delay":      delay,
		}.Warningf(c, "Transient error while fetching tail log; retrying.")
	})
	switch err {
	case nil:
		le, err := e.GetLogEntry()
		if err != nil {
			log.WithError(err).Errorf(c, "Failed to load tail entry data.")
			return nil, err
		}
		return []*logpb.LogEntry{le}, nil

	case storage.ErrDoesNotExist:
		return nil, nil

	default:
		log.WithError(err).Errorf(c, "Failed to fetch tail log.")
		return nil, err
	}
}
