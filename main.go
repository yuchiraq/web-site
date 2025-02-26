// main.go
package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Project struct {
	Title string
	Image string
}

func main() {
	r := mux.NewRouter()

	// Данные для примера
	projects := []Project{
		{"Коттедж 'Премиум'", "/static/images/project1.jpg"},
		{"Офисный комплекс", "/static/images/project2.jpg"},
		{"Реконструкция здания", "/static/images/project3.jpg"},
	}

	// Обработчики
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(
			"templates/base.html",
			"templates/index.html",
		))
		tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
			"Projects": projects[:3],
		})
	})

	r.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(
			"templates/base.html",
			"templates/projects.html",
		))
		tmpl.ExecuteTemplate(w, "base", projects)
	})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	http.Handle("/", r)
	http.ListenAndServe(":8088", nil)
}
