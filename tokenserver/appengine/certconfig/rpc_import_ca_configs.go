// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package certconfig

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	ds "github.com/luci/gae/service/datastore"
	"github.com/luci/gae/service/info"
	"github.com/luci/luci-go/common/config"
	"github.com/luci/luci-go/common/data/stringset"
	"github.com/luci/luci-go/common/errors"
	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/common/proto/google"

	"github.com/luci/luci-go/tokenserver/api/admin/v1"
	"github.com/luci/luci-go/tokenserver/appengine/utils"
)

// ImportCAConfigsRPC implements Admin.ImportCAConfigs RPC method.
type ImportCAConfigsRPC struct {
}

// ImportCAConfigs fetches CA configs from from luci-config right now.
func (r *ImportCAConfigsRPC) ImportCAConfigs(c context.Context, _ *google.Empty) (*admin.ImportedConfigs, error) {
	cfg, err := fetchConfigFile(c, "tokenserver.cfg")
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "can't read config file - %s", err)
	}
	logging.Infof(c, "Importing tokenserver.cfg at rev %s", cfg.Revision)

	// Read list of CAs.
	msg := admin.TokenServerConfig{}
	if err = proto.UnmarshalText(cfg.Content, &msg); err != nil {
		return nil, grpc.Errorf(codes.Internal, "can't parse config file - %s", err)
	}

	seenIDs, err := LoadCAUniqueIDToCNMap(c)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "can't load unique_id map - %s", err)
	}
	if seenIDs == nil {
		seenIDs = map[int64]string{}
	}
	seenIDsDirty := false

	// There should be no duplicates.
	seenCAs := stringset.New(len(msg.GetCertificateAuthority()))
	for _, ca := range msg.GetCertificateAuthority() {
		if seenCAs.Has(ca.Cn) {
			return nil, grpc.Errorf(codes.Internal, "duplicate entries in the config")
		}
		seenCAs.Add(ca.Cn)
		// Check unique ID is not being reused.
		if existing, seen := seenIDs[ca.UniqueId]; seen {
			if existing != ca.Cn {
				return nil, grpc.Errorf(
					codes.Internal, "duplicate unique_id %d in the config: %q and %q",
					ca.UniqueId, ca.Cn, existing)
			}
		} else {
			seenIDs[ca.UniqueId] = ca.Cn
			seenIDsDirty = true
		}
	}

	// Update the mapping CA unique_id -> CA CN. Unique integer ids are used in
	// various tokens in place of a full CN name to save space. This mapping is
	// additive (all new CAs should have different IDs).
	if seenIDsDirty {
		if err := StoreCAUniqueIDToCNMap(c, seenIDs); err != nil {
			return nil, grpc.Errorf(codes.Internal, "can't store unique_id map - %s", err)
		}
	}

	// Add new CA datastore entries or update existing ones.
	wg := sync.WaitGroup{}
	me := errors.NewLazyMultiError(len(msg.GetCertificateAuthority()))
	for i, ca := range msg.GetCertificateAuthority() {
		wg.Add(1)
		go func(i int, ca *admin.CertificateAuthorityConfig) {
			defer wg.Done()
			certFileCfg, err := fetchConfigFile(c, ca.CertPath)
			if err != nil {
				logging.Errorf(c, "Failed to fetch %q: %s", ca.CertPath, err)
				me.Assign(i, err)
			} else if err := importCA(c, ca, certFileCfg.Content, cfg.Revision); err != nil {
				logging.Errorf(c, "Failed to import %q: %s", ca.Cn, err)
				me.Assign(i, err)
			}
		}(i, ca)
	}
	wg.Wait()
	if err = me.Get(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "can't import CA - %s", err)
	}

	// Find CAs that were removed from the config.
	toRemove := []string{}
	q := ds.NewQuery("CA").Eq("Removed", false).KeysOnly(true)
	err = ds.Run(c, q, func(k *ds.Key) {
		if !seenCAs.Has(k.StringID()) {
			toRemove = append(toRemove, k.StringID())
		}
	})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "datastore error - %s", err)
	}

	// Mark them as inactive in the datastore.
	wg = sync.WaitGroup{}
	me = errors.NewLazyMultiError(len(toRemove))
	for i, name := range toRemove {
		wg.Add(1)
		go func(i int, name string) {
			defer wg.Done()
			if err := removeCA(c, name, cfg.Revision); err != nil {
				logging.Errorf(c, "Failed to remove %q: %s", name, err)
				me.Assign(i, err)
			}
		}(i, name)
	}
	wg.Wait()
	if err = me.Get(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "datastore error - %s", err)
	}

	return &admin.ImportedConfigs{
		ImportedConfigs: []*admin.ImportedConfigs_ConfigFile{
			{
				Name:     "tokenserver.cfg",
				Revision: cfg.Revision,
			},
		},
	}, nil
}

// fetchConfigFile fetches a file from this services' config set.
func fetchConfigFile(c context.Context, path string) (*config.Config, error) {
	configSet := "services/" + info.AppID(c)
	logging.Infof(c, "Reading %q from config set %q", path, configSet)
	c, _ = context.WithTimeout(c, 30*time.Second) // URL fetch deadline
	return config.GetConfig(c, configSet, path, false)
}

// importCA imports CA definition from the config (or updates an existing one).
func importCA(c context.Context, ca *admin.CertificateAuthorityConfig, certPem string, rev string) error {
	// Read CA certificate file, convert it to der.
	certDer, err := utils.ParsePEM(certPem, "CERTIFICATE")
	if err != nil {
		return fmt.Errorf("bad PEM - %s", err)
	}

	// Check the certificate makes sense.
	cert, err := x509.ParseCertificate(certDer)
	if err != nil {
		return fmt.Errorf("bad cert - %s", err)
	}
	if !cert.IsCA {
		return fmt.Errorf("not a CA cert")
	}
	if cert.Subject.CommonName != ca.Cn {
		return fmt.Errorf("bad CN in the certificate, expecting %q, got %q", ca.Cn, cert.Subject.CommonName)
	}

	// Serialize the config back to proto to store it in the entity.
	cfgBlob, err := proto.Marshal(ca)
	if err != nil {
		return err
	}

	// Create or update the entity.
	return ds.RunInTransaction(c, func(c context.Context) error {
		existing := CA{CN: ca.Cn}
		err := ds.Get(c, &existing)
		if err != nil && err != ds.ErrNoSuchEntity {
			return err
		}
		// New one?
		if err == ds.ErrNoSuchEntity {
			logging.Infof(c, "Adding new CA %q", ca.Cn)
			return ds.Put(c, &CA{
				CN:         ca.Cn,
				Config:     cfgBlob,
				Cert:       certDer,
				AddedRev:   rev,
				UpdatedRev: rev,
			})
		}
		// Exists already? Check whether we should update it.
		if !existing.Removed &&
			bytes.Equal(existing.Config, cfgBlob) &&
			bytes.Equal(existing.Cert, certDer) {
			return nil
		}
		logging.Infof(c, "Updating CA %q", ca.Cn)
		existing.Config = cfgBlob
		existing.Cert = certDer
		existing.Removed = false
		existing.UpdatedRev = rev
		existing.RemovedRev = ""
		return ds.Put(c, &existing)
	}, nil)
}

// removeCA marks the CA in the datastore as removed.
func removeCA(c context.Context, name string, rev string) error {
	return ds.RunInTransaction(c, func(c context.Context) error {
		existing := CA{CN: name}
		if err := ds.Get(c, &existing); err != nil {
			return err
		}
		if existing.Removed {
			return nil
		}
		logging.Infof(c, "Removing CA %q", name)
		existing.Removed = true
		existing.RemovedRev = rev
		return ds.Put(c, &existing)
	}, nil)
}
