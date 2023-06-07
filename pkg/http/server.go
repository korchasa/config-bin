package http

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
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
    //s.router.
    //    HandleFunc("/{bid}/{vid}", s.logRequest(s.handleShowBinVersion())).
    //    Methods("GET")
    //s.router.
    //    HandleFunc("/api/v1/{bid}", s.logRequest(s.handleCreateBin())).
    //    Methods("POST")
    //s.router.
    //    HandleFunc("/api/v1/{bid}/get", s.logRequest(s.handleGetBin())).
    //    Methods("GET")
    //s.router.
    //    HandleFunc("/api/v1/{bid}/update", s.logRequest(s.handleUpdateBin())).
    //    Methods("POST")
    //s.router.
    //    HandleFunc("/v1/{bid}/{vid}/rollback", s.logRequest(s.handleRollbackBin())).
    //    Methods("POST")
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

// handleError sends error response
func (s *Server) handleError(r *http.Request, w http.ResponseWriter, httpStatus int, code string, err error) {
    log.WithFields(log.Fields{"code": code}).Warnf("ERROR %s: %v", r.RequestURI, err)
    if w != nil {
        s.sendResponse(w, httpStatus, &ResultResponseStruct{
            Success: false,
            Result:  nil,
            Error: struct {
                Code    string `json:"code"`
                Message string `json:"message"`
            }{
                Code:    code,
                Message: err.Error(),
            },
        })
    }
}

// sendOKResponse sends a response with a success status
func (s *Server) sendOKResponse(w http.ResponseWriter, httpStatus int, result interface{}) {
    s.sendResponse(w, httpStatus, &ResultResponseStruct{
        Success: true,
        Result:  result,
    })
}

// sendResponse method send body like json
func (s *Server) sendResponse(w http.ResponseWriter, httpStatus int, body *ResultResponseStruct) {
    data, err := json.Marshal(body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Errorf("errors encoding body to json, body: %v", body)
        return
    }
    log.Infof("response: %s", data)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(httpStatus)
    _, err = w.Write(data)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Errorf("failed to write data to http.ResponseWriter, body: %v", data)
    }
}

type ResultResponseStruct struct {
    Success bool        `json:"success"`
    Result  interface{} `json:"result"`
    Error   struct {
        Code    string `json:"code"`
        Message string `json:"message"`
    } `json:"error"`
}
