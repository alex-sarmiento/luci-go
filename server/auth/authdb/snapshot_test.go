// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package authdb

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"

	"golang.org/x/net/context"

	"github.com/luci/luci-go/server/auth/identity"
	"github.com/luci/luci-go/server/auth/internal"
	"github.com/luci/luci-go/server/auth/service/protocol"
	"github.com/luci/luci-go/server/auth/signing"
	"github.com/luci/luci-go/server/auth/signing/signingtest"
	"github.com/luci/luci-go/server/secrets"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSnapshotDB(t *testing.T) {
	Convey("IsAllowedOAuthClientID works", t, func() {
		c := context.Background()
		db, err := NewSnapshotDB(&protocol.AuthDB{
			OauthClientId: strPtr("primary-client-id"),
			OauthAdditionalClientIds: []string{
				"additional-client-id-1",
				"additional-client-id-2",
			},
		}, "http://auth-service", 1234)
		So(err, ShouldBeNil)

		call := func(email, clientID string) bool {
			res, err := db.IsAllowedOAuthClientID(c, email, clientID)
			So(err, ShouldBeNil)
			return res
		}

		So(call("abc@appspot.gserviceaccount.com", "anonymous"), ShouldBeTrue)
		So(call("dude@example.com", ""), ShouldBeFalse)
		So(call("dude@example.com", googleAPIExplorerClientID), ShouldBeTrue)
		So(call("dude@example.com", "primary-client-id"), ShouldBeTrue)
		So(call("dude@example.com", "additional-client-id-2"), ShouldBeTrue)
		So(call("dude@example.com", "unknown-client-id"), ShouldBeFalse)
	})

	Convey("IsMember works", t, func() {
		c := context.Background()
		db, err := NewSnapshotDB(&protocol.AuthDB{
			Groups: []*protocol.AuthGroup{
				{
					Name:    strPtr("direct"),
					Members: []string{"user:abc@example.com"},
				},
				{
					Name:  strPtr("via glob"),
					Globs: []string{"user:*@example.com"},
				},
				{
					Name:   strPtr("via nested"),
					Nested: []string{"direct"},
				},
				{
					Name:   strPtr("cycle"),
					Nested: []string{"cycle"},
				},
				{
					Name:   strPtr("unknown nested"),
					Nested: []string{"unknown"},
				},
			},
		}, "http://auth-service", 1234)
		So(err, ShouldBeNil)

		call := func(ident string, groups ...string) bool {
			res, err := db.IsMember(c, identity.Identity(ident), groups...)
			So(err, ShouldBeNil)
			return res
		}

		So(call("user:abc@example.com", "direct"), ShouldBeTrue)
		So(call("user:another@example.com", "direct"), ShouldBeFalse)

		So(call("user:abc@example.com", "via glob"), ShouldBeTrue)
		So(call("user:abc@another.com", "via glob"), ShouldBeFalse)

		So(call("user:abc@example.com", "via nested"), ShouldBeTrue)
		So(call("user:another@example.com", "via nested"), ShouldBeFalse)

		So(call("user:abc@example.com", "cycle"), ShouldBeFalse)
		So(call("user:abc@example.com", "unknown"), ShouldBeFalse)
		So(call("user:abc@example.com", "unknown nested"), ShouldBeFalse)

		So(call("user:abc@example.com"), ShouldBeFalse)
		So(call("user:abc@example.com", "unknown", "direct"), ShouldBeTrue)
		So(call("user:abc@example.com", "via glob", "direct"), ShouldBeTrue)
	})

	Convey("SharedSecrets works", t, func() {
		c := context.Background()
		db, err := NewSnapshotDB(&protocol.AuthDB{
			Secrets: []*protocol.AuthSecret{
				{
					Name: strPtr("secret-1"),
					Values: [][]byte{
						[]byte("current"),
					},
				},
				{
					Name: strPtr("secret-2"),
					Values: [][]byte{
						[]byte("current"),
						[]byte("prev1"),
						[]byte("prev2"),
					},
				},
				{
					Name: strPtr("empty"),
				},
			},
		}, "http://auth-service", 1234)
		So(err, ShouldBeNil)

		s, err := db.SharedSecrets(c)
		So(err, ShouldBeNil)
		So(s, ShouldResemble, secrets.StaticStore{
			"secret-1": {
				Current: secrets.NamedBlob{Blob: []byte("current")},
			},
			"secret-2": {
				Current: secrets.NamedBlob{Blob: []byte("current")},
				Previous: []secrets.NamedBlob{
					{Blob: []byte("prev1")},
					{Blob: []byte("prev2")},
				},
			},
		})
	})

	Convey("GetCertificates works", t, func(c C) {
		db, err := NewSnapshotDB(&protocol.AuthDB{
			OauthClientId: strPtr("primary-client-id"),
			OauthAdditionalClientIds: []string{
				"additional-client-id-1",
				"additional-client-id-2",
			},
			TokenServerUrl: strPtr("http://token-server"),
		}, "http://auth-service", 1234)
		So(err, ShouldBeNil)

		authService := signingtest.NewSigner(0, &signing.ServiceInfo{
			AppID:              "auth-service",
			ServiceAccountName: "auth-service-account@example.com",
		})

		tokenService := signingtest.NewSigner(1, &signing.ServiceInfo{
			AppID:              "token-server",
			ServiceAccountName: "token-server-account@example.com",
		})

		calls := 0

		ctx := context.Background()
		ctx = internal.WithTestTransport(ctx, func(r *http.Request, body string) (int, string) {
			calls++
			var certs *signing.PublicCertificates
			var err error
			switch r.URL.String() {
			case "http://auth-service/auth/api/v1/server/certificates":
				certs, err = authService.Certificates(ctx)
			case "http://token-server/auth/api/v1/server/certificates":
				certs, err = tokenService.Certificates(ctx)
			default:
				return 404, "Wrong URL"
			}
			if err != nil {
				panic(err)
			}
			blob, err := json.Marshal(certs)
			if err != nil {
				panic(err)
			}
			return 200, string(blob)
		})

		good := []identity.Identity{
			"service:auth-service",
			"service:token-server",
			"user:auth-service-account@example.com",
			"user:token-server-account@example.com",
		}
		for _, ident := range good {
			certs, err := db.GetCertificates(ctx, ident)
			So(err, ShouldBeNil)
			So(certs, ShouldNotBeNil)
		}

		// Fetched two bundles, once.
		So(calls, ShouldEqual, 2)

		// For unknown signer returns (nil, nil).
		certs, err := db.GetCertificates(ctx, "service:unknown")
		So(err, ShouldBeNil)
		So(certs, ShouldBeNil)
	})

	Convey("IsInWhitelist works", t, func() {
		c := context.Background()
		db, err := NewSnapshotDB(&protocol.AuthDB{
			IpWhitelistAssignments: []*protocol.AuthIPWhitelistAssignment{
				{
					Identity:    strPtr("user:abc@example.com"),
					IpWhitelist: strPtr("whitelist"),
				},
			},
			IpWhitelists: []*protocol.AuthIPWhitelist{
				{
					Name: strPtr("whitelist"),
					Subnets: []string{
						"1.2.3.4/32",
						"10.0.0.0/8",
					},
				},
				{
					Name: strPtr("empty"),
				},
			},
		}, "http://auth-service", 1234)
		So(err, ShouldBeNil)

		wl, err := db.GetWhitelistForIdentity(c, "user:abc@example.com")
		So(err, ShouldBeNil)
		So(wl, ShouldEqual, "whitelist")

		wl, err = db.GetWhitelistForIdentity(c, "user:unknown@example.com")
		So(err, ShouldBeNil)
		So(wl, ShouldEqual, "")

		call := func(ip, whitelist string) bool {
			ipaddr := net.ParseIP(ip)
			So(ipaddr, ShouldNotBeNil)
			res, err := db.IsInWhitelist(c, ipaddr, whitelist)
			So(err, ShouldBeNil)
			return res
		}

		So(call("1.2.3.4", "whitelist"), ShouldBeTrue)
		So(call("10.255.255.255", "whitelist"), ShouldBeTrue)
		So(call("9.255.255.255", "whitelist"), ShouldBeFalse)
		So(call("1.2.3.4", "empty"), ShouldBeFalse)
	})
}

func strPtr(s string) *string { return &s }

func BenchmarkIsMember(b *testing.B) {
	c := context.Background()
	db, _ := NewSnapshotDB(&protocol.AuthDB{
		Groups: []*protocol.AuthGroup{
			{
				Name:   strPtr("outer"),
				Nested: []string{"A", "B"},
			},
			{
				Name:   strPtr("A"),
				Nested: []string{"A_A", "A_B"},
			},
			{
				Name:   strPtr("B"),
				Nested: []string{"B_A", "B_B"},
			},
			{
				Name:   strPtr("A_A"),
				Nested: []string{"A_A_A"},
			},
			{
				Name:   strPtr("A_A_A"),
				Nested: []string{"A_A_A_A"},
			},
			{
				Name: strPtr("A_A_A_A"),
			},
			{
				Name: strPtr("A_B"),
			},
			{
				Name: strPtr("B_A"),
			},
			{
				Name: strPtr("B_B"),
			},
		},
	}, "http://auth-service", 1234)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.IsMember(c, "user:somedude@example.com", "outer")
	}
}
