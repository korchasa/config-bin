package responder

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var ErrUnknown = errors.New("unknown error")

type TemplatesProvider interface {
	MustGet(name string) *template.Template
}

type Responder struct {
	tplProvider TemplatesProvider
}

func New(tplProvider TemplatesProvider) *Responder {
	return &Responder{tplProvider: tplProvider}
}

func (res *Responder) JSONError(req *http.Request, resp http.ResponseWriter, httpStatus int, code string, err error) {
	if err == nil {
		err = ErrUnknown
	}
	log.
		WithFields(log.Fields{"code": code, "url": req.URL.String()}).
		Warnf("ERROR %s: %v", req.RequestURI, err)
	if resp != nil {
		res.sendResponseJSON(resp, httpStatus, &ResultResponseStruct{
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

// HTMLError sends error response.
func (res *Responder) HTMLError(req *http.Request, resp http.ResponseWriter, httpStatus int, code string, err error) {
	if err == nil {
		err = ErrUnknown
	}
	log.
		WithFields(log.Fields{"code": code, "url": req.URL.String()}).
		Warnf("ERROR(html page) %s: %v", req.RequestURI, err)
	if resp == nil {
		return
	}

	resp.WriteHeader(httpStatus)
	err = res.tplProvider.MustGet("error.gohtml").Execute(resp, struct {
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

// JSONSuccess sends a response with a success status.
func (res *Responder) JSONSuccess(resp http.ResponseWriter, httpStatus int, result interface{}) {
	res.sendResponseJSON(resp, httpStatus, &ResultResponseStruct{
		Success: true,
		Result:  result,
	})
}

// sendResponseJSON method send body like json.
func (res *Responder) sendResponseJSON(resp http.ResponseWriter, httpStatus int, body *ResultResponseStruct) {
	data, err := json.Marshal(body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		log.Errorf("errors encoding body to json, body: %v", body)
		return
	}
	log.Infof("response: %s", data)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(httpStatus)
	_, err = resp.Write(data)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
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
