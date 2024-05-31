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

	_, err := config.LoadConfig("config.toml")
	if err != nil && os.IsNotExist(err) {
		_, createErr := config.CreateNewConfig("config.toml")
		if createErr != nil {
			log.Fatal(createErr)
		} else {
			log.Print("No config detected. New config  created.")
		}
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
		data, err := config.AddServiceToConfig("config.toml", name, url)

		if err != nil {
			log.Fatal(err)
		}

		tmpl.Execute(rw, data)
	})

	http.HandleFunc("/edit-service", func(rw http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		oldName := req.FormValue("oldName")
		name := req.FormValue("name")
		url := req.FormValue("url")

		template := template.Must(template.ParseFiles("./templates/fragments/services.html"))
		editData := config.EditServiceData{
			OldName: oldName,
			Name:    name,
			URL:     url,
		}

		data, editError := config.EditServiceInConfig("config.toml", editData)
		if editError != nil {
			log.Fatal(editError)
		}

		template.Execute(rw, data)

	})

	http.HandleFunc("/delete-service", func(rw http.ResponseWriter, req *http.Request) {
		parseError := req.ParseForm()
		if parseError != nil {
			log.Fatal(parseError)
		}

		name := req.FormValue("name")

		template := template.Must(template.ParseFiles("./templates/fragments/services.html"))
		data, deleteError := config.DeleteServiceFromConfig("config.toml", name)
		if deleteError != nil {
			log.Fatal(deleteError)
		}

		template.Execute(rw, data)
	})

	log.Println("Cerebral running on localhost:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
