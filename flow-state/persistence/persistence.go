package persistence

import (
	"github.com/project-flogo/services/flow-state/persistence/api"
	"github.com/project-flogo/services/flow-state/persistence/cache"
)

var storage api.Storage

func init() {
	storage = cache.NewCacheStorage()
}

func GetStorage() api.Storage {
	return storage
}
