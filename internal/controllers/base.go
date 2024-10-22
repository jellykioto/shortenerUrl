package controllers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi"
)

type Logger interface {
	Info(args ...interface{})
}

type BaseController struct {
	logger Logger
	Urls   map[string]string
}

func NewBaseController(logger Logger) *BaseController {
	return &BaseController{
		logger: logger,
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", c.handleShortenURL)
	r.Post("/", c.handleShortenURL)
	// r.Get("/{name}", c.handleName)
	return r
}

// func (c *BaseController) handleMain(writer http.ResponseWriter, request *http.Request) {
// 	c.logger.Info("main")
// 	writer.Write([]byte("Hello"))
// }

// func (c *BaseController) handleName(writer http.ResponseWriter, request *http.Request) {
// 	c.logger.Info("name")
// 	writer.Write([]byte("hello "))
// }

func (c *BaseController) handleShortenURL(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("shortUrl")
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannon read request body: %s", err), http.StatusBadRequest)
			return
		}
		if string(body) == "" {
			http.Error(w, "Empty POST request body!", http.StatusBadRequest)
			return
		}
		url := string(body)
		id := generateID()
		c.Urls[id] = url
		response := fmt.Sprintf("http://localhost:8080/%s", id)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(response))
		if err != nil {
			return
		}
	} else if r.Method == http.MethodGet {
		id := r.URL.Path[1:]
		fmt.Println(r.URL.Path[1:])
		url, ok := c.Urls[id]
		if !ok {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func generateID() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	uniqueId := make([]rune, 8)
	for i := range uniqueId {
		uniqueId[i] = letters[rand.Intn(len(letters))]
	}
	return string(uniqueId)
}
