package acl

import (
	"log"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/go-openapi/runtime/middleware"
)

func NewAclMiddleware(next http.Handler) http.Handler {
	e, _ := casbin.NewEnforcerSafe("./model.conf", "./policy.csv")

	skipPaths := map[string]bool{
		"/dummyLogin": true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skipPaths[r.URL.Path] {
			log.Printf("%v %v %v", next, w, r)
			next.ServeHTTP(w, r)
			return
		}

		role := "rol"
		res, _ := e.EnforceSafe(role, r.URL.Path, r.Method)
		log.Printf("path=%s role=%s access=%v", r.URL.Path, role, res)
		if res {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	})
}

func GetEntMw() middleware.Builder {
	e, _ := casbin.NewEnforcerSafe("./model.conf", "./policy.csv")

	skipPaths := map[string]bool{
		"/dummyLogin": true,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if skipPaths[r.URL.Path] {
				log.Printf("%v %v %v", next, w, r)
				next.ServeHTTP(w, r)
				return
			}
	
			role := ""
			res, _ := e.EnforceSafe(role, r.URL.Path, r.Method)
			log.Printf("path=%s role=%s access=%v", r.URL.Path, role, res)
			if res {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		})
	}
}
