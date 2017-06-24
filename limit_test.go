package limit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestParam(t *testing.T) {
	assert.Panics(t, func() {
		Limit(0)
	})
}

func TestLimit(t *testing.T) {
	router := gin.New()
	router.Use(Limit(1))
	router.GET("/", func(*gin.Context) {
		time.Sleep(500 * time.Microsecond)
	})

	w := performRequest("GET", "/", router)
	assert.Equal(t, 200, w.Code)

}

func TestHandlerPanic(t *testing.T) {
	assert.Panics(t, func() {
		router := gin.New()
		router.Use(Limit(1))
		router.GET("/err", func(*gin.Context) {
			panic("err")
		})

		performRequest("GET", "/err", router)
	})
}

func TestFulled(t *testing.T) {
	const max = 5

	attempts := 1000
	var failed int
	var wg sync.WaitGroup

	router := gin.New()
	router.Use(Limit(max))
	router.GET("/", func(*gin.Context) {
		time.Sleep(5 * time.Microsecond)
	})

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := performRequest("GET", "/", router)
			if w.Code == 502 {
				failed++
			}
		}()
	}
	wg.Wait()

	// We expect some Gets to fail
	assert.True(t, failed > attempts/2)
}

func qqqq(max int) *gin.Engine {
	router := gin.New()

	router.Use(Limit(max))
	router.GET("/", func(*gin.Context) {
		time.Sleep(500 * time.Microsecond)
		fmt.Println("get")
	})
	router.GET("/err", func(*gin.Context) {
		//time.Sleep(500 * time.Microsecond)
		panic("foo err")
	})
	return router
}

func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
