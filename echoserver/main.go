package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	mux := http.NewServeMux()
	// Define the routes
	mux.HandleFunc("/public", loggingMiddleware(publicHandler))
	mux.HandleFunc("/", loggingMiddleware(indexHandler))
	mux.HandleFunc("/private", loggingMiddleware(privateHandler))
	mux.HandleFunc("/echo", loggingMiddleware(echoHandler))

	// Start the HTTP server on port 8080
	port := 8080
	logrus.Infof("Server is running on :%d", port)
	s := &http.Server{
		Handler: mux,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			logrus.Infof("conn context: %v <> %v", c.RemoteAddr(), c.LocalAddr())
			return ctx
		},
		Addr: fmt.Sprintf(":%d", port),
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func loggingMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"path":        r.URL.Path,
			"headers":     r.Header,
		}).Info("handling request")
		fn(w, r)
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Index!\n")
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is a public route\n")
}

func privateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is a private route\n")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"remote_ip": "%s", "time": "%s"}`, r.RemoteAddr, time.Now().Format(time.RFC3339))
}
