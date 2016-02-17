// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package types

import (
	"fmt"
)

// ValueType is an enum for the type of a metric.
type ValueType int

// Types of metric values.
const (
	NonCumulativeIntType ValueType = iota
	CumulativeIntType
	NonCumulativeFloatType
	CumulativeFloatType
	StringType
	BoolType
	NonCumulativeDistributionType
	CumulativeDistributionType
)

func (v ValueType) String() string {
	switch v {
	case NonCumulativeIntType:
		return "NonCumulativeIntType"
	case CumulativeIntType:
		return "CumulativeIntType"
	case NonCumulativeFloatType:
		return "NonCumulativeFloatType"
	case CumulativeFloatType:
		return "CumulativeFloatType"
	case StringType:
		return "StringType"
	case BoolType:
		return "BoolType"
	case NonCumulativeDistributionType:
		return "NonCumulativeDistributionType"
	case CumulativeDistributionType:
		return "CumulativeDistributionType"
	}
	panic(fmt.Sprintf("unknown ValueType %d", v))
}
