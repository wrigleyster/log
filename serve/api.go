package serve

import (
	"log"
	"net/http"
	"os"
)

type Server struct {
	addr string
}

func newServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Write([]byte("User ID: " + id))
	})
	router.HandleFunc("/admin/welcome", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dat, err := os.ReadFile("./html/index.htm")
		if err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		w.Header().Add("content-type", "text/html; charset=utf-8")
		if _, err := w.Write(dat); err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
		}
	})
	server := http.Server{
		Addr: s.addr,
		Handler: authenticate(router),
	}
	log.Printf("Server has started %s", s.addr)
	return server.ListenAndServe()
}

func authenticatedHandler(next http.Handler) http.Handler {
	authenticatedRouter := http.NewServeMux()
	authenticatedRouter.HandleFunc("admin/", authenticate(next))
	authenticatedRouter.HandleFunc("secure/", authenticate(next))
	return authenticatedRouter
}
func authenticate(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Add("WWW-Authenticate", "Basic realm=zz, charset=\"UTF-8\"")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func Serve() {
	server := newServer(":8080")
	server.Run()

}
