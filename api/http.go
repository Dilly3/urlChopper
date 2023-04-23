package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"net/http"

	"github.com/dilly3/urlshortner/internal"
	"github.com/dilly3/urlshortner/serializer/json"
	"github.com/dilly3/urlshortner/serializer/msgpack"

	//"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type RedirectHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	RedirectService *internal.RedirectService
	Logger          zap.Logger
}

func NewHandler(redirectService *internal.RedirectService) RedirectHandler {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("error setting up Logger =>%v\n", err)
		log.Fatal(err)
	}
	return &handler{RedirectService: redirectService, Logger: *logger}

}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	code := r.RequestURI[1:]

	redirect, err := h.RedirectService.Find(code)
	if err != nil {

		if errors.Cause(err) == internal.ErrRedirectNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return

		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
	http.Redirect(w, r, redirect.Url, http.StatusMovedPermanently)

}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return

	}
	redirect, err := h.serializer(string(contentType)).Decode(reqbody)
	if err != nil {

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = h.RedirectService.Store(redirect)
	if err != nil {

		if errors.Cause(err) == internal.ErrRedirectInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	responseBody, err := h.serializer(string(contentType)).Encode(redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("======> 7", err)
		return
	}

	h.setupResponseHeader(w, contentType, http.StatusCreated, responseBody)

}

func (h *handler) setupResponseHeader(w http.ResponseWriter, contentType string, statusCode int, body []byte) {
	w.Header().Add("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}

}

func (h *handler) serializer(contentType string) internal.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &msgpack.Redirect{}
	}
	return &json.Redirect{}
}

func SelectPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}
