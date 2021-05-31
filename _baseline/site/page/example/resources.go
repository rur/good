package example

import (
	"net/http"
	"sync"

	"github.com/rur/good/_baseline/site/page"
	"github.com/rur/good/_baseline/site/service"
	"github.com/rur/treetop"
)

// resources that are request-scoped data bound to page handlers
type resources struct {
	user service.User
	// EDITME: add request specific resources for this page
}

// loadResources constructs a resources struct for use by the request handlers of this page
//
// If an error or failure occurs this function is responsible for responding  directly to the client.
//
// Invoking the Write method of the ResponseWriter will prevent subsequent functions from
// handling this request.
func loadResources(env *service.Env, w http.ResponseWriter, req *http.Request) (rsc *resources, ok bool) {
	// EDITME: setup your handler resources here
	ok = true
	rsc = &resources{
		user: service.User{
			Name: "test",
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

// bindResources is middleware with memoization which loads home page resources for local request handlers
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
			// wait until the response is concluded and clear the rscCache entry
			<-rsp.Context().Done()
			rscCacheLock.Lock()
			rsc, ok := rscCache[key]
			delete(rscCache, key)
			rscCacheLock.Unlock() // release lock before attempting teardown
			if ok {
				teardownResources(rsc, env)
			}
		}()
		return f(rsc, env, rsp, req)
	}
}

func init() {
	rscCache = make(map[uint32]*resources)
}
