package example

import (
	"net/http"
	"sync"
	"time"

	"github.com/rur/good/baseline/routes_test/page"
	"github.com/rur/good/baseline/routes_test/service"
	"github.com/rur/treetop"
)

// resources that are request-scoped data bound to page handlers
type resources struct {
	user service.User
	// EDITME: add request specific resources for this page
}

// loadResources constructs a resources struct for use by the request handlers of this page
//
// If an error or failure occurs, this function is responsible for responding to the client using the ResponseWriter.
func loadResources(env *service.Env, w http.ResponseWriter, req *http.Request) (rsc *resources, ok bool) {
	// EDITME: setup your handler resources here
	ok = true
	rsc = &resources{
		user: service.User{
			Name:  "!unauthenticated!",
			Email: "unauthenticated@example.com",
		},
	}
	return
}

// teardownResources happens after the response has been written for the request,
// this hook exists in case any special teardown is needed for your resources instance
func teardownResources(rsc *resources, env *service.Env) {
	// EDITME: Add teardown logic for your resources
}

//
// resource caching mechanics
//

var (
	rscCache     map[uint32]*resources
	rscCacheLock sync.RWMutex
)

type handlerWithResources func(*resources, *service.Env, treetop.Response, *http.Request) interface{}

// bindResources is middleware with memoization which loads example page resources for local request handlers
func bindResources(f handlerWithResources) page.ViewHandlerWithEnv {
	return func(env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
		key := rsp.ResponseID()
		rscCacheLock.RLock()
		rsc, ok := rscCache[key]
		rscCacheLock.RUnlock()
		if ok {
			return f(rsc, env, rsp, req)
		}
		// expensive path
		rsc, ok = loadResources(env, rsp, req)
		if !ok {
			// an error ocurred, abort
			return nil
		}

		// Cache resources for this response ID and setup a
		// callback to tear it down when the response is finished
		rscCacheLock.Lock()
		rscCache[key] = rsc
		rscCacheLock.Unlock()
		go func() {
			defer func() {
				rscCacheLock.Lock()
				rsc, ok := rscCache[key]
				delete(rscCache, key)
				rscCacheLock.Unlock() // release lock *before* attempting teardown
				if ok {
					teardownResources(rsc, env)
				}
			}()

			// block until the response is concluded, defer will clear the rscCache entry
			select {
			case <-rsp.Context().Done():
			case <-time.After(10 * time.Minute): // failsafe timeout to protect shared resources
				if env.ErrorLog != nil {
					env.ErrorLog.Println("failsafe request timeout expired after 10min for example page")
				}
				// since rsp Context is bound to the HTTP request, I assume that this is safe to do
				http.Error(rsp, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
			}
		}()
		return f(rsc, env, rsp, req)
	}
}

func init() {
	rscCache = make(map[uint32]*resources)
}
