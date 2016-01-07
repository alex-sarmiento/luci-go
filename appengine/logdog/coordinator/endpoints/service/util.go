// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package service

import (
	"fmt"
)

func impossible(err error) {
	if err != nil {
		panic(fmt.Errorf("impossible condition: %v", err))
	}
}