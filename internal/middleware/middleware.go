package middleware

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"go_kit_project/internal/static"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Request-Method", "DELETE, POST, GET, OPTIONS")
		if r.Method == static.OPTIONS {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}

func MiddlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		fmt.Println("Executing middlewareOne again")
	})
}

// Logger
func InitializeLogger(logg log.Logger) (l log.Logger, f *os.File) {
	f, _ = os.OpenFile(RootDir()+"/"+static.LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	wrt := io.MultiWriter(os.Stdout, f)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(wrt)
		logger = log.With(logger, static.KeyTs, log.DefaultTimestampUTC)
		logger = log.With(logger, static.KeyCaller, log.DefaultCaller)
	}
	logg = logger
	return logg, f
}

//ROOT DIR
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

func LoggingOperation(logg log.Logger, values ...interface{}) {
	logg, f := InitializeLogger(logg)
	_ = logg.Log(values...)
	defer f.Close()
}
