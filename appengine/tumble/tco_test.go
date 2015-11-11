// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package tumble

import (
	"testing"

	"github.com/luci/gae/service/datastore"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

type BigObjectGroup struct {
	_id int64 `gae:"$id,1"`

	Count int64
}

type BigObject struct {
	ID    int64          `gae:"$id"`
	Group *datastore.Key `gae:"$parent"`

	Data []byte `gae:",noindex"`
}

type SlowMutation struct {
	Count int64
}

func (s *SlowMutation) Root(c context.Context) *datastore.Key {
	return datastore.Get(c).MakeKey("BigObjectGroup", 1)
}

var JunkBytes = make([]byte, 512*1024)

func init() {
	for i := 0; i < 32768; i += 16 {
		copy(JunkBytes[i:i+16], []byte("BuhNerfCatBlarty"))
	}
}

func (s *SlowMutation) RollForward(c context.Context) (muts []Mutation, err error) {
	ds := datastore.Get(c)

	grp := &BigObjectGroup{}
	err = ds.Get(grp)
	if err == datastore.ErrNoSuchEntity {
		err = ds.Put(grp)
	}
	if err != nil {
		return
	}

	grp.Count++
	bo := &BigObject{ID: grp.Count, Group: ds.KeyForObj(grp), Data: JunkBytes}
	err = ds.PutMulti([]interface{}{grp, bo})

	retMut := *s
	retMut.Count--
	if retMut.Count > 0 {
		muts = append(muts, &retMut)
	}
	return
}

func init() {
	Register((*SlowMutation)(nil))
}

func TestTailCallOptimization(t *testing.T) {
	if testing.Short() {
		t.Skip("TCO test skipped in -short mode")
	}

	t.Parallel()

	Convey("TCO test", t, func() {
		testing := Testing{}

		c := testing.Context()

		// This will start a chain. We should only be able to handle about 8
		// of these in a single datastore transaction.
		So(RunMutation(c, &SlowMutation{40}), ShouldBeNil)

		// 6 processed shards shows that we did 5 transactions (of 8 each). The last
		// shard is empty.
		So(testing.Drain(c), ShouldEqual, 6)

		ds := datastore.Get(c)

		bog := &BigObjectGroup{}
		So(ds.Get(bog), ShouldBeNil)

		// Hey look at that, bog shows all 40 :)
		So(bog.Count, ShouldEqual, 40)
	})
}
