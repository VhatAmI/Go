package main

import (
	"log"
	"net/http"
)
	
func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<html>
			<head>
				<title> Go Chat</title>
			</head>
			<body>
				Lets chat!
			</body>
		</html>
	`))
})
//start web server here
if err := http.ListenAndServe(":8080", nil); err != nil {
	log.Fatal("ListenAndServe: ", err)
}
}	