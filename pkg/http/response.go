package http

import (
    "encoding/json"
    log "github.com/sirupsen/logrus"
    "net/http"
)

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
