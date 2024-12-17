package app

import (
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/inner/broker/store/db"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/store"
	"github.com/nutsdb/nutsdb"
	"io"
)

type Store interface {
	broker.TopicMessageStore
	store.KeyStore
	io.Closer
}

func newStore(bucketName string) Store {
	switch config.GetConfig().Store.Default {
	case "redis":
		return LoadRedisStore()

	default:
		return LoadLocalStore(bucketName)
	}
}

func LoadRedisStore() *db.Redis {
	return db.NewRedis()
}

func LoadLocalStore(bucketName string) *db.NutsDB {
	return db.NewNutsDBStore(nutsdb.Options{
		EntryIdxMode: nutsdb.HintKeyValAndRAMIdxMode,
		SegmentSize:  nutsdb.MB * 256,
		NodeNum:      1,
		RWMode:       nutsdb.MMap,
		SyncEnable:   true,
	}, bucketName, nutsdb.WithDir("./../data/nutsdb/"+bucketName))
}
