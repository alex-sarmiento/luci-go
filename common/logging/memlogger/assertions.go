// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package memlogger

import (
	"fmt"

	"github.com/luci/luci-go/common/logging"
)

// ShouldHaveLog is a goconvey custom assertion which asserts that the logger has
// received a log. It takes up to 3 arguments.
// argument 1 (expected[0]) is the log level.
// argument 2 (expected[1]) is the message. If omitted or the empty string, this value is not checked.
// argument 3 (expected[2]) is the fields data. If omitted, this value is not checked.
func ShouldHaveLog(actual interface{}, expected ...interface{}) string {
	m, ok := actual.(*MemLogger)
	if !ok {
		return fmt.Sprintf("actual value must be a *MemLogger, not %T", actual)
	}

	level := logging.Level(0)
	msg := ""
	data := map[string]interface{}(nil)

	switch len(expected) {
	case 3:

		data, ok = expected[2].(map[string]interface{})
		if !ok {
			return fmt.Sprintf(
				"Third argument to this assertion must be an map[string]interface{} (was %T)", expected[2])
		}
		fallthrough
	case 2:

		msg, ok = expected[1].(string)
		if !ok {
			return fmt.Sprintf(
				"Second argument to this assertion must be a string (was %T)", expected[1])
		}
		fallthrough
	case 1:
		level, ok = expected[0].(logging.Level)
		if !ok {
			return "First argument to this assertion must be a logging.Level"
		}

	default:
		return fmt.Sprintf(
			"This assertion requires at least 1 comparison value (you provided %d)", len(expected))
	}

	if m.Has(level, msg, data) {
		return ""
	}

	logString := fmt.Sprintf("Level %s", level)
	if msg != "" {
		logString += fmt.Sprintf(" message '%s'", msg)
	}

	if data != nil {
		logString += fmt.Sprintf(" data '%v'", data)
	}

	return fmt.Sprintf("No log matching %s", logString)
}

// ShouldNotHaveLog is the inverse of ShouldHaveLog. It asserts that the logger
// has not seen such a log.
func ShouldNotHaveLog(actual interface{}, expected ...interface{}) string {
	res := ShouldHaveLog(actual, expected...)

	if res != "" {
		return ""
	}

	return "Found a log, but wasn't supposed to."
}
