package httpway

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func New() *Router {
	return &Router{
		Router: httprouter.New(),
	}
}

type Router struct {
	*httprouter.Router
	SessionManager SessionManager
	Logger         Logger

	prev   *Router
	handle httprouter.Handle
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, handle httprouter.Handle) {
	r.Handle("GET", path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, handle httprouter.Handle) {
	r.Handle("HEAD", path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, handle httprouter.Handle) {
	r.Handle("OPTIONS", path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, handle httprouter.Handle) {
	r.Handle("POST", path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, handle httprouter.Handle) {
	r.Handle("PUT", path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, handle httprouter.Handle) {
	r.Handle("PATCH", path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handle httprouter.Handle) {
	r.Handle("DELETE", path, handle)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handle httprouter.Handle) {
	newHandle := r.GenerateChainHandler(handle)

	r.Router.Handle(method, path, newHandle)
}

//Add a middleware before (and after) the handler run
//   router := httpway.New()
//   public := router.Middleware(AccessLogger)
//   private := public.Middleware(AuthCheck)
//  
//   public.GET("/public", somePublicHandler)
//   private.GET("/private", somePrivateHandler)
//
//  func AccessLogger(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//  	startTime:=time.Now()
//
//	httpway.GetContext(r).Next(w, r, ps)
//
//  	fmt.Printf("Request: %s duration: %s\n", r.URL.EscapedPath(), time.Since(startTime))
//  }
//  
//  func AuthCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	ctx := httpway.GetContext(r)
//
//  	if !ctx.Session().IsAuth() {
//		http.Error(w, "Auth required", 401)
//		return
//  	}
//	ctx.Next(w, r, ps)
//  }
//
func (r *Router) Middleware(handle httprouter.Handle) *Router {
	rt := &Router{
		prev:   r,
		handle: handle,
		Router: r.Router,
		Logger: r.Logger,
		SessionManager: r.SessionManager,
	}

	return rt
}

//get handler with all the middlewares chained
func (router *Router) GenerateChainHandler(handle httprouter.Handle) httprouter.Handle {
	if router.prev == nil {
		return handle
	}

	var (
		lastMiddleware httprouter.Handle
		middlewareList = make([]httprouter.Handle, 0)
	)

	mid := router
	middlewareList = append(middlewareList, handle)

	for mid.prev != nil {
		if mid.prev.handle == nil {
			lastMiddleware = mid.handle
			break
		}
		middlewareList = append(middlewareList, mid.handle)
		mid = mid.prev
	}
	middlewareListLen := len(middlewareList)

	httprouterHandler := func(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
		w = CreateContext(router, w, r, &middlewareList, &middlewareListLen)

		lastMiddleware(w, r, pr)
	}

	return httprouterHandler
}
