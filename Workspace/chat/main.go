package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"trace"
	"os"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once 	 sync.Once
	filename string
	templ 	 *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
		})
		t.templ.Execute(w, r)
}	
	
func main() {
	//-addr="hostvalue" as commandline arguement, can also be ip:port if local host is not desired
	var addr = flag.String("addr", ":8080", "Address of host, defaults to local port 8080")
	flag.Parse()
	//set up gomniauth for OAuth2
	gomniauth.SetSecurityKey("blahblahblah")
	gomniauth.WithProviders(
		facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret", "http://localhost:8080/auth/callback/github"),
		google.New("AIzaSyDzsizedt9cTXnD2_ahSK6XUmjYNRwyCjM", "VpBY9ZTNrVYV6nKywIFRUNiD", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()
	//start web server here
	log.Println("Starting server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}	