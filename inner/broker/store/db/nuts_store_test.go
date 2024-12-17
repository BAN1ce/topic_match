package db

import (
	"context"
	"github.com/nutsdb/nutsdb"
	"testing"
)

func TestLocal_ReadPrefixKey(t *testing.T) {

	store := NewNutsDBStore(nutsdb.Options{
		EntryIdxMode: nutsdb.HintKeyValAndRAMIdxMode,
		SegmentSize:  nutsdb.MB * 256,
		NodeNum:      1,
		RWMode:       nutsdb.MMap,
		SyncEnable:   true,
	}, "testing", nutsdb.WithDir("./../data/testing"))

	store.PutKey(context.TODO(), "prefix/1", "value1")
	store.PutKey(context.TODO(), "prefix/2", "value2")
	store.PutKey(context.TODO(), "prefix/3", "value3")

	s, err := store.ReadPrefixKey(context.TODO(), "prefix")
	if err != nil {
		t.Error(err)
	}
	if len(s) != 3 {
		t.Error("ReadPrefixKey failed")
	}

}
