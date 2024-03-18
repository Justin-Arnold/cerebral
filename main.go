package main

import (
	"cerebral/internal/config"
	"html/template"
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

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/fragments/services.html"))
		data, err := config.LoadConfig("config.toml")
		if err != nil {
			log.Fatal(err)
		}

		tmpl.Execute(w, data)
	})

	http.HandleFunc("/add-service", func(rw http.ResponseWriter, req *http.Request) {
		err := req.ParseForm() // Parses the form data
		if err != nil {
			log.Fatal(err)
		}

		name := req.FormValue("name") // Access the "name" field
		url := req.FormValue("url")   // Access the "url" field

		tmpl := template.Must(template.ParseFiles("./templates/fragments/services.html"))
		data, err := config.UpdateConfig("config.toml", name, url)

		log.Print(data)
		if err != nil {
			log.Fatal(err)
		}

		tmpl.Execute(rw, data)
	})

	log.Println("Cerebral running on localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
