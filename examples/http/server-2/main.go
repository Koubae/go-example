package main

import (
	"fmt"
	"log"
	"net/http"
)

type LogLevel int

func main() {
	//Configure logging
	//log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.LUTC | log.Lmicroseconds | log.Lshortfile)
	logMessage(INFO, "Server starting")

	server := http.NewServeMux()
	server.HandleFunc("GET /", index)
	server.HandleFunc("GET /ping", ping)

	var handler http.Handler
	handler = logAccessMiddleware(server)

	// Start the server on port 8080
	logMessage(INFO, "Server starting on port %d...", 8080)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		logMessage(ERROR, "Error starting server: %v", err)
	}

}

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var appLogLevel = INFO

func logMessage(level LogLevel, format string, v ...interface{}) {
	if level < appLogLevel {
		return
	}

	levelName := map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
	}[level]
	log.Printf("[%s] %s", levelName, fmt.Sprintf(format, v...))
}

// ==================================
// Middlewares
// ==================================
// responseWriterWrapper wraps http.ResponseWriter to capture the status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func logAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		logMessage(INFO, "%s %s %s %d", r.Method, r.URL.Path, r.RemoteAddr, rw.statusCode)
	})
}

// ==================================
// Handlers
// ==================================
func index(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		errorHandler(w, req, http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Welcome to Simple-Go-Server\n")
}

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

// handler for 404 not found
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 - Page Not Found")
	logMessage(WARN, "404 - Not Found: %s %s", r.Method, r.URL.Path)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}
