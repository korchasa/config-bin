package http

import (
    "fmt"
    "github.com/gorilla/mux"
    "net/http"
    "time"
)

// Server struct
type Server struct {
    router      *mux.Router
    store       Storage
    metrics     Metrics
    tplProvider TemplatesProvider
}

// NewServer returns a new server instance by provided parameters
func NewServer(store Storage, metrics Metrics, tplProvider TemplatesProvider) *Server {
    return &Server{
        router:      mux.NewRouter(),
        store:       store,
        tplProvider: tplProvider,
        metrics:     metrics,
    }
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    s.router.ServeHTTP(w, req)
}

// Run runs server instance on provided port
func (s *Server) Run(address string) error {
    s.router.
        HandleFunc("/", s.logRequest(s.handleRootHTML())).
        Methods("GET")
    s.router.
        HandleFunc("/{bid}", s.logRequest(s.handleShowBin())).
        Methods("GET")
    s.router.
        HandleFunc("/api/v1/{bid}", s.logRequest(s.handleAPICreateBin())).
        Methods("POST")
    s.router.
        HandleFunc("/api/v1/{bid}/get", s.logRequest(s.handleAPIGetBin())).
        Methods("GET")
    s.router.
        HandleFunc("/api/v1/{bid}/update", s.logRequest(s.handleAPIUpdateBin())).
        Methods("POST")
    s.router.
        HandleFunc("/v1/{bid}/{vid}/rollback", s.logRequest(s.handleAPIRollbackBin())).
        Methods("POST")
    s.router.
        HandleFunc("/liveness", s.logRequest(s.handleLiveness())).
        Methods("GET")
    s.router.
        HandleFunc("/readiness", s.logRequest(s.handleReadiness())).
        Methods("GET")
    httpServer := &http.Server{
        Addr:         address,
        ReadTimeout:  5 * time.Minute,
        WriteTimeout: 10 * time.Second,
        Handler:      s.router,
    }
    err := httpServer.ListenAndServe()
    if err != nil {
        return fmt.Errorf("listen and serve error: %v", err)
    }
    return nil
}
