package persistence

import (
	"github.com/project-flogo/services/flow-store/persistence/api"
	"github.com/project-flogo/services/flow-store/persistence/cache"
)

var storage api.Storage

func init() {
	storage = cache.NewCacheStorage()
}

func GetStorage() api.Storage {
	return storage
}

//func isRedis() bool {
//	v := os.Getenv("FLOGO_PERSISTENCE_DB")
//	if v != "" && v == "redis" {
//		return true
//	}
//	return false
//}
