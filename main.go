package main

import (
	// "encoding/json"
	"html/template"
	// "io"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
			port = "8080"
			log.Printf("defaulting to port %s", port)
	}

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.Execute(w, nil)
	})

	log.Println("Cerebral running on localhost:" + port)
	log.Fatal(http.ListenAndServe(":"  +port, nil))
}