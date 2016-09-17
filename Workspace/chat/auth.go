package main

import(
	"net/http"
	"log"
	"fmt"
	"strings"
	"github.com/stretchr/objx"
	"github.com/stretchr/gomniauth"
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
		provider, err := gomniauth.Provider(provider)
		if err != nil {
		log.Fatalln("Error authenticating with", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil,nil)
		if err != nil {
		  log.Fatalln("Error when trying to GetBeginAuthURL for ", provider,"-", err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
	
	provider, err :=gomniauth.Provider(provider)
	if err !=nil {
		log.Fatalln("Error when trying to callback with provider : ", provider, "-", err )
	}
	
	creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
		log.Fatalln("Error when trying to get user from ", provider, "-", err)
		}
	
	user, err := provider.GetUser(creds)
	if err != nil {
		log.Fatalln("Couldn't get user from ", provider, "-", err)
	}
		
	authCookieValue := objx.New(map[string]interface{}{
		"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
		Name: "auth",
		Value: authCookieValue,
		Path: "/"})
	
	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)

		
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
} 