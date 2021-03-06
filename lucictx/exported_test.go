// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package lucictx

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/luci/luci-go/common/system/environ"
)

func TestLiveExported(t *testing.T) {
	// t.Parallel() because of os.Environ manipulation

	dir, err := ioutil.TempDir(os.TempDir(), "exported_test")
	if err != nil {
		t.Fatalf("could not create tempdir! %s", err)
	}
	defer os.RemoveAll(dir)

	Convey("LiveExports", t, func() {
		os.Unsetenv(EnvKey)

		tf, err := ioutil.TempFile(dir, "exported_test.liveExport")
		tfn := tf.Name()
		tf.Close()
		So(err, ShouldBeNil)
		defer os.Remove(tfn)

		le := &liveExport{path: tfn}

		Convey("Can only be closed once", func() {
			le.Close()
			So(func() { le.Close() }, ShouldPanic)
		})

		Convey("Removes the file when it is closed", func() {
			_, err := os.Stat(tfn)
			So(err, ShouldBeNil)

			le.Close()
			_, err = os.Stat(tfn)
			So(os.IsNotExist(err), ShouldBeTrue)
		})

		Convey("Can add to command", func() {
			cmd := exec.Command("test", "arg")
			cmd.Env = os.Environ()
			le.SetInCmd(cmd)
			So(len(cmd.Env), ShouldEqual, len(os.Environ())+1)
			So(cmd.Env[len(cmd.Env)-1], ShouldStartWith, EnvKey)
			So(cmd.Env[len(cmd.Env)-1], ShouldEndWith, le.path)
		})

		Convey("Can modify in command", func() {
			cmd := exec.Command("test", "arg")
			cmd.Env = os.Environ()
			cmd.Env[0] = EnvKey + "=helloworld"
			le.SetInCmd(cmd)
			So(len(cmd.Env), ShouldEqual, len(os.Environ()))
			So(cmd.Env[0], ShouldStartWith, EnvKey)
			So(cmd.Env[0], ShouldEndWith, le.path)
		})

		Convey("Can add to environ", func() {
			env := environ.System()
			_, ok := env.Get(EnvKey)
			So(ok, ShouldBeFalse)
			le.SetInEnviron(env)
			val, ok := env.Get(EnvKey)
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, le.path)
		})

		Convey("Can set in global env", func() {
			le.UnsafeSetInGlobalEnviron()
			So(os.Getenv(EnvKey), ShouldEqual, le.path)

			Convey("cannot set it more than once", func() {
				So(le.UnsafeSetInGlobalEnviron, ShouldPanic)
			})

			Convey("closing resets the variable", func() {
				le.Close()
				_, ok := os.LookupEnv(EnvKey)
				So(ok, ShouldBeFalse)
			})
		})

		Convey("Changing global env resores previous", func() {
			os.Setenv(EnvKey, "sup")
			le.UnsafeSetInGlobalEnviron()
			So(os.Getenv(EnvKey), ShouldEqual, le.path)
			le.Close()
			So(os.Getenv(EnvKey), ShouldEqual, "sup")
		})
	})
}
