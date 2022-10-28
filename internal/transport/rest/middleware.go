package rest

import (
	"errors"
	"fmt"
	"github.com/MikhailSolovev/URLAsker/pkg/logger"
	"net/http"
)

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

const (
	allowedOrigins     = "*"
	allowedMethods     = "GET, POST, DELETE, PUT, OPTIONS"
	allowedHeaders     = "*"
	exposedHeaders     = "*"
	requestHeaders     = "*"
	allowedCredentials = "true"
)

func SetCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		w.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
		w.Header().Set("Access-Control-Request-Headers", requestHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", allowedCredentials)

		next.ServeHTTP(w, r)
	})
}

func HandleError(h HandlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log.LogDebug(fmt.Sprintf("method: %v endpoint: %v query: %v", r.Method, r.URL.Path, r.URL.RawQuery))

		var Err *Error

		err := h(w, r)
		if err == nil {
			return
		}

		if errors.As(err, &Err) {
			logger.Log.LogWarn(fmt.Sprintf("client: %v debug: %v code: %v", Err.Message, Err.DebugMessage, Err.httpCode))
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(Err.httpCode)
			w.Write([]byte(Err.Message))
			return
		}
	}
}
