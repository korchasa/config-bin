package responder

import (
    "encoding/json"
    "fmt"
    log "github.com/sirupsen/logrus"
    "html/template"
    "net/http"
)

type TemplatesProvider interface {
    MustGet(name string) *template.Template
}

type Responder struct {
    tplProvider TemplatesProvider
}

func New(tplProvider TemplatesProvider) *Responder {
    return &Responder{tplProvider: tplProvider}
}

func (res *Responder) JSONError(req *http.Request, w http.ResponseWriter, httpStatus int, code string, err error) {
    if err == nil {
        err = fmt.Errorf("unknown error")
    }
    log.
        WithFields(log.Fields{"code": code, "url": req.URL.String()}).
        Warnf("ERROR %s: %v", req.RequestURI, err)
    if w != nil {
        res.sendResponseJSON(w, httpStatus, &ResultResponseStruct{
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

// HTMLError sends error response
func (res *Responder) HTMLError(req *http.Request, w http.ResponseWriter, httpStatus int, code string, err error) {
    if err == nil {
        err = fmt.Errorf("unknown error")
    }
    log.
        WithFields(log.Fields{"code": code, "url": req.URL.String()}).
        Warnf("ERROR(html page) %s: %v", req.RequestURI, err)
    if w == nil {
        return
    }

    w.WriteHeader(httpStatus)
    err = res.tplProvider.MustGet("error.gohtml").Execute(w, struct {
        Code    string
        Message string
    }{
        Code:    code,
        Message: err.Error(),
    })
    if err != nil {
        log.Errorf("errors encoding error body to html, body: %v", err)
    }
}

// JSONSuccess sends a response with a success status
func (res *Responder) JSONSuccess(w http.ResponseWriter, httpStatus int, result interface{}) {
    res.sendResponseJSON(w, httpStatus, &ResultResponseStruct{
        Success: true,
        Result:  result,
    })
}

// sendResponseJSON method send body like json
func (res *Responder) sendResponseJSON(w http.ResponseWriter, httpStatus int, body *ResultResponseStruct) {
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
