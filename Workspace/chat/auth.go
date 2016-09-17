package main

import(
	"net/http"
	"log"
	"fmt"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
	//no authentication
	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
	//some other error ???
	panic(err.Error())
	} else {
	//User authentication successful 
	h.next.ServeHTTP(w,r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

//loginHandler for 3rd party
//format is /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		log.Println("Handle login for ", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
} 