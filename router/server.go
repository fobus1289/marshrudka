package router

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

type shutdown struct {
	*sync.Mutex
	request chan *http.Request
	counter bool
}

type server struct {
	clients           int32
	shutdown          *shutdown
	services          reflectMap
	HandlersInterface handlersInterface
	routers           routers
}

func (s *shutdown) start() {
	for request := range s.request {
		s.Lock()
		s.counter = true
		<-request.Context().Done()
		s.counter = false
		s.Unlock()
	}
}

func (s *server) HasClient() bool {
	log.Println(s.clients)
	return s.shutdown.counter
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//s.shutdown.request <- r
	//atomic.AddInt32(&s.clients, 1)
	defer atomic.AddInt32(&s.clients, -1)
	if !s.routers.Find(w, r) {
		http.NotFound(w, r)
	}
}

func (s *server) Use(handlers ...interface{}) {
	s.HandlersInterface = interfaceJoin(s.HandlersInterface, handlers)
}

func (s *server) Group(path string, handlers ...interface{}) IGroup {
	return &group{
		Path:              path,
		HandlersInterface: interfaceJoin(s.HandlersInterface, handlers),
		Server:            s,
	}
}

func (s *server) GET(path string, handlers ...interface{}) IMatch {
	return initRouter(&options{
		server:   s,
		path:     path,
		handlers: interfaceJoin(s.HandlersInterface, handlers),
		methods: map[string]bool{
			http.MethodGet: true,
		},
	})
}

func (s *server) POST(path string, handlers ...interface{}) IMatch {
	return initRouter(&options{
		server:   s,
		path:     path,
		handlers: interfaceJoin(s.HandlersInterface, handlers),
		methods: map[string]bool{
			http.MethodPost: true,
		},
	})
}

func (s *server) PUT(path string, handlers ...interface{}) IMatch {
	return initRouter(&options{
		server:   s,
		path:     path,
		handlers: interfaceJoin(s.HandlersInterface, handlers),
		methods: map[string]bool{
			http.MethodPut: true,
		},
	})
}

func (s *server) PATCH(path string, handlers ...interface{}) IMatch {
	return initRouter(&options{
		server:   s,
		path:     path,
		handlers: interfaceJoin(s.HandlersInterface, handlers),
		methods: map[string]bool{
			http.MethodPatch: true,
		},
	})
}

func (s *server) DELETE(path string, handlers ...interface{}) IMatch {
	return initRouter(&options{
		server:   s,
		path:     path,
		handlers: interfaceJoin(s.HandlersInterface, handlers),
		methods: map[string]bool{
			http.MethodPatch: true,
		},
	})
}

func (s *server) ANY(path string, handlers ...interface{}) IMatch {
	return initRouter(&options{
		server:   s,
		path:     path,
		handlers: interfaceJoin(s.HandlersInterface, handlers),
		methods: map[string]bool{
			http.MethodGet:    true,
			http.MethodPost:   true,
			http.MethodPut:    true,
			http.MethodPatch:  true,
			http.MethodDelete: true,
		},
	})
}

func (s *server) MATCH(path string, methods []string, handlers ...interface{}) IMatch {

	if len(methods) < 1 {
		panic("MATCH error")
	}

	var mapMethods = map[string]bool{}

	for _, method := range methods {
		mapMethods[strings.ToUpper(method)] = true
	}

	return initRouter(&options{
		server:   s,
		path:     path,
		handlers: interfaceJoin(s.HandlersInterface, handlers),
		methods:  mapMethods,
	})
}

func (s *server) FileServer(method, path, dir string, BeforeHandlers ...interface{}) {
	initRouter(&options{
		server:   s,
		path:     path,
		handlers: interfaceJoin(interfaceJoin(s.HandlersInterface, BeforeHandlers), []interface{}{http.FileServer(http.Dir(dir))}),
		methods: map[string]bool{
			strings.ToUpper(method): true,
		},
	})
}

func (s *server) Run(addr string) error {
	rebuildRouters(s)
	go s.shutdown.start()
	showUrl(s, addr)
	return http.ListenAndServe(addr, s)
}

func (s *server) RunTLS(addr, certFile, keyFile string) error {
	panic("implement me")
}

func (s *server) RunAsync(addr string) error {
	panic("implement me")
}

func (s *server) RunAsyncTLS(addr, certFile, keyFile string) error {
	panic("implement me")
}
