// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:generate protoc --go_out=. acquisition_network_device.proto acquisition_task.proto metrics.proto

// Package ts_mon_proto contains ts_mon protobuf source and generated protobuf
// data.
//
// The package name here must match the protobuf package name, as the generated
// files will reside in the same directory.
package ts_mon_proto

import (
	"github.com/golang/protobuf/proto"
)

var _ = proto.Marshal