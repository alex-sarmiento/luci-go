// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package eventlog provides a client which can be used to to collect and send batches of ChromeInfraEvent logs to the eventlog service.
package eventlog

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/golang/protobuf/proto"
	"github.com/luci/luci-go/common/eventlog/internal/logservice"
	logpb "github.com/luci/luci-go/common/eventlog/proto"
)

// TestEndpoint is the server address for test logs.
const TestEndpoint = "https://jmt17.google.com/log"

// ProdEndpoint is the server address for production logs.
const ProdEndpoint = "https://play.googleapis.com/log"

// Client may be used to send ChromeInfraEvent logs to the eventlog service.
type Client struct {
	c *logservice.Client

	// EventSource identifies the log producer. It may be configured before sending logs.
	EventSource *logpb.InfraEventSource
}

// NewClient constructs a Client which can be used to send ChromeInfraEvent logs to the eventlog service.
// Users must call Close when the Client is no longer needed.
func NewClient(serverAddr string, opts ...ClientOption) *Client {
	// TODO(mcgreevy): help users to set EventSource?

	settings := &clientSettings{
		HTTPClient: http.DefaultClient,
	}
	for _, o := range opts {
		o.apply(settings)
	}

	serviceClient := logservice.NewClient(serverAddr, "CHROME_INFRA")
	serviceClient.HTTPClient = settings.HTTPClient
	return &Client{c: serviceClient}
}

// LogSync synchronously logs events to the eventlog service.
// Use NewLogEvent to assist with constructing a well-formed log event.
// LogSync takes ownership of events.
func (c *Client) LogSync(ctx context.Context, events ...*ChromeInfraLogEvent) error {
	var logEvents []*logpb.LogRequestLite_LogEventLite

	for _, event := range events {
		sourceExt, err := proto.Marshal(event.InfraEvent)
		if err != nil {
			return err
		}
		event.LogEvent.SourceExtension = sourceExt
		logEvents = append(logEvents, event.LogEvent)
	}

	return c.c.LogSync(ctx, logEvents)
}

// NewLogEvent constructs a well-formed log event.
// eventTime is the time that the event occurred, typically obtained by calling
// one of Point, Begin or End.
func (c *Client) NewLogEvent(ctx context.Context, eventTime TypedTime) *ChromeInfraLogEvent {
	// TODO(mcgreevy): look into whether we can populate TraceID, SpanId, ParentId from the context.
	// TODO(mcgreevy): support event_code, event_flow_id.
	return &ChromeInfraLogEvent{
		LogEvent: &logpb.LogRequestLite_LogEventLite{
			EventTimeMs: proto.Int64(eventTime.Time.UnixNano() / 1e6),
		},
		InfraEvent: &logpb.ChromeInfraEvent{
			TimestampKind: &eventTime.Kind,
			EventSource:   c.EventSource,
		},
	}
}

// Close flushes any pending logs and releases any resources held by the client.
// Close should be called when the client is no longer needed.
func (c *Client) Close() error {
	return c.c.Close()
}

// ChromeInfraLogEvent stores a pending LogEvent, and the proto used to populate its SourceExtension field.
type ChromeInfraLogEvent struct {
	// LogEvent is the event to be sent to the server. The LogEvent's SourceExtension
	// field will be overwritten with the serialized contents of InfraEvent
	// when the log event is sent to the server.
	LogEvent   *logpb.LogRequestLite_LogEventLite
	InfraEvent *logpb.ChromeInfraEvent
}

// TypedTime is a time, accompanied by an indication of what the time represents.
type TypedTime struct {
	// The time at which an event occurred. Must have been obtained from Time.Now.
	Time time.Time
	Kind logpb.ChromeInfraEvent_TimestampKind
}

// Point returns the current time for use when logging an event with no temporal extent.
func Point() TypedTime {
	return TypedTime{time.Now(), logpb.ChromeInfraEvent_POINT}
}

// Begin returns the current time for use when logging the beginning of an event.
func Begin() TypedTime {
	return TypedTime{time.Now(), logpb.ChromeInfraEvent_BEGIN}
}

// End returns the current time for use when logging the end of an event.
func End() TypedTime {
	return TypedTime{time.Now(), logpb.ChromeInfraEvent_END}
}
