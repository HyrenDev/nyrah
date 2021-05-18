package local

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var CACHE = cache.New(cache.NoExpiration, 10*time.Second)