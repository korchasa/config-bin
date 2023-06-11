package server

import (
    "configBin/pkg"
    "configBin/pkg/server/responder"
    "fmt"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "html/template"
    "net/http"
    "time"
)

// Server struct
type Server struct {
    router      *mux.Router
    store       Storage
    resp        *responder.Responder
    tplProvider TemplatesProvider
    metrics     Metrics
}

type Metrics interface {
    IncRequestsCount()
    IncEventSendsCount(eventType string, appID string, appInstance string, appEnv string, project string, ack bool)
    IncSuccessfulRequestsCount(httpCode int)
    IncFailedRequestsCount(httpCode int)
    ObserveSuccessfulRequestDuration(dur time.Duration)
    IncKafkaRequestsCount()
    IncKafkaErrorsCount()
    ObserveKafkaRequestDuration(dur time.Duration)
}

type Storage interface {
    CreateBin(id uuid.UUID, pass string, unencryptedData string) error
    // GetBin returns the bin with the given ID.
    GetBin(id uuid.UUID, pass string) (*pkg.Bin, error)
    UpdateBin(id uuid.UUID, pass string, unencryptedData string) error
    IsReady() bool
    Close()
}

type TemplatesProvider interface {
    MustGet(name string) *template.Template
}

type Responder interface {
    JSONError(req *http.Request, w http.ResponseWriter, httpStatus int, code string, err error)
    HTMLError(req *http.Request, w http.ResponseWriter, httpStatus int, code string, err error)
    JSONSuccess(w http.ResponseWriter, httpStatus int, result interface{})
}

// New returns a new server instance by provided parameters
func New(store Storage, resp *responder.Responder, tplProvider TemplatesProvider, metrics Metrics) *Server {
    s := &Server{
        router:      mux.NewRouter(),
        store:       store,
        resp:        resp,
        tplProvider: tplProvider,
        metrics:     metrics,
    }
    s.routes()
    return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    s.router.ServeHTTP(w, req)
}

// Run runs server instance on provided port
func (s *Server) Run(address string) error {
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

func (s *Server) Close() {
    s.store.Close()
}
