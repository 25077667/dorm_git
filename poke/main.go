package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type RegularIntMap struct {
	sync.RWMutex
	internal map[string]int16
}

func NewRegularIntMap() *RegularIntMap {
	return &RegularIntMap{
		internal: make(map[string]int16),
	}
}

func (rm *RegularIntMap) Load(key string) (value int16, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *RegularIntMap) Store(key string, value int16) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}

var poke = NewRegularIntMap()

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		ip := c.ClientIP()

		if val, ok := poke.Load(ip); ok {
			poke.Store(ip, val+1)
		} else {
			poke.Store(ip, 0)
		}
		c.String(http.StatusOK, "%d", func() int16 {
			ret, _ := poke.Load(ip)
			return ret
		}())
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
