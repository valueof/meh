package server

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type key int

const (
	INBOUND_DIR    string = "/var/tmp/mehserver"
	REQUEST_ID_KEY key    = 0
)

//go:embed html
var templates embed.FS

var tasks TaskPool

func render(w http.ResponseWriter, r *http.Request, name string, data any) {
	ctx := r.Context()
	logger := getLoggerFromContext(ctx)
	rid := getRequestIDFromContext(ctx)

	t, err := template.New("base.html").ParseFS(templates, "html/base.html", "html/"+name)
	if err != nil {
		logger.Printf("ParseFS on html/%s failed with: %v", name, err)
		fmt.Fprintf(w, "Internal Server Error (%s)", rid)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		logger.Printf("Execute on html/%s failed with %v", name, err)
		fmt.Fprintf(w, "Internal Server Error (%s)", rid)
	}
}

func getRequestIDFromContext(ctx context.Context) (rid string) {
	rid, ok := ctx.Value(REQUEST_ID_KEY).(string)
	if !ok {
		rid = "unknown"
	}
	return
}

func getLoggerFromContext(ctx context.Context) (logger *log.Logger) {
	rid := getRequestIDFromContext(ctx)
	return log.New(os.Stdout, fmt.Sprintf("[%s]", rid), log.LstdFlags)
}

func tracing(uuid func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get("X-Request-Id")
			if rid == "" {
				rid = uuid()
			}

			ctx := context.WithValue(r.Context(), REQUEST_ID_KEY, rid)
			w.Header().Set("X-Request-Id", rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rid := getRequestIDFromContext(r.Context())
				logger.Println(rid, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func RunHTTPServer(addr string) {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags)
	logger.Println("Server is starting...")

	logger.Println("Preparing directory to hold uploaded files")
	err := os.MkdirAll(INBOUND_DIR, 0700)
	if err != nil {
		logger.Printf("Failed to create %s: %v", INBOUND_DIR, err)
		logger.Fatalf("Can't proceed")
	}

	logger.Println("Creating task pool")
	tasks = TaskPool{pool: make(map[string]taskStatus)}

	router := http.NewServeMux()
	router.HandleFunc("/", homepage)
	router.HandleFunc("/upload/", upload)
	router.HandleFunc("/result/", result)
	router.HandleFunc("/favicon.ico", favicon)

	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
		Addr:         addr,
		Handler:      tracing(uuid.NewString)(logging(logger)(router)),
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.SetKeepAlivesEnabled(false)
		if err := s.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready at", addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", addr, err)
	}

	<-done
	logger.Println("Goodbye, friend!")
}
