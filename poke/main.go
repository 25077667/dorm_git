package main

import (
	"crypto/rand"
	"encoding/hex"
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

func gen_random_tok() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

var poke = NewRegularIntMap()

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		random_tok := "/poke/user/" + gen_random_tok()
		c.Redirect(http.StatusFound, c.Request.URL.Path+random_tok)
	})

	r.GET("/poke", func(c *gin.Context) {
		random_tok := "/user/" + gen_random_tok()
		c.Redirect(http.StatusFound, c.Request.URL.RequestURI()+"/poke/"+random_tok)
	})

	r.GET("/poke/user/:tok", func(c *gin.Context) {
		tok := c.Param("tok")
		if val, ok := poke.Load(tok); ok {
			poke.Store(tok, val+1)
		} else {
			poke.Store(tok, 0)
		}
		c.String(http.StatusOK, "%d", func() int16 {
			ret, _ := poke.Load(tok)
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
