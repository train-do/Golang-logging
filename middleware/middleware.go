package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/train-do/Golang-Restfull-API/handler"
	"go.uber.org/zap"
)

type Middleware struct {
	Log *zap.Logger
}

func NewMiddleware(log *zap.Logger) Middleware {
	return Middleware{
		Log: log,
	}
}

func (middleware *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || err == http.ErrNoCookie {
			// fmt.Println("NO COOKIE")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if cookie.Value != handler.Token {
			// fmt.Println("GAGAL VALIDASI")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (middleware *Middleware) MinddlewareLogger(handler http.Handler) http.Handler {
	fmt.Println("MASUK MIDDLEWARE LOGGER")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("**********")
		start := time.Now()

		handler.ServeHTTP(w, r)

		duration := time.Since(start)

		middleware.Log.Info("http request", zap.String("url", r.URL.String()), zap.Duration("duration", duration))
	})

}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || err == http.ErrNoCookie {
			// fmt.Println("NO COOKIE")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if cookie.Value != handler.Token {
			// fmt.Println("GAGAL VALIDASI")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func MinddlewareLogger(log *zap.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		fmt.Println("MASUK MIDDLEWARE LOGGER")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("++++++++")
			start := time.Now()

			handler.ServeHTTP(w, r)

			duration := time.Since(start)

			log.Info("http request middleware", zap.String("url", r.URL.String()), zap.Duration("duration", duration))
		})
	}
}
