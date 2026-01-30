package server

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/VinneB/gym-webapp/internal/sql"
	"github.com/VinneB/gym-webapp/internal/structapi"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("pages/template_test.html")),
	}
}

var validPages []string = []string{"/addexercise"}

func StartServer() {
	fmt.Println("Entering start server")
	mux := http.NewServeMux()
	mux.HandleFunc("/", htmlTemplateHandler)
	mux.HandleFunc("GET /htmx/exercises", ExercisesGetHandler)
	mux.HandleFunc("POST /htmx/exercises", ExercisesPostHandler)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test/favicon.ico")
	})
	err := sql.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer sql.CloseDatabase()
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func htmlTemplateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("html " + r.URL.Path)
	renderer := newTemplate()
	data, err := getData(r.URL.Path)
	if err != nil {
		log.Println(err)
		SendError(w, r, http.StatusInternalServerError, "Sorry")
	}
	if slices.Contains(validPages, r.URL.Path) {
		data.Page = r.URL.Path
		renderer.Render(w, "index", data)
	} else {
		data.Page = "error"
		fmt.Println("Error Page")
		renderer.Render(w, "index", data)
	}
}

func htmxHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("htmx")
	renderer := newTemplate()
	data := structapi.Data{}
	renderer.Render(w, "error_section", data)
	endpoint := strings.TrimPrefix(r.URL.Path, "/htmx/")
	log.Printf("htmx endpoint - %s\n", endpoint)
	switch endpoint {
	case "exercises":

	}
}

func SendError(w http.ResponseWriter, r *http.Request, statusCode int, err string) {
	w.WriteHeader(statusCode)
	data := structapi.Data{Errors: []string{err}, Page: "error"}
	newTemplate().Render(w, "index", data)
	log.Println("Respondede")
}

func getData(path string) (structapi.Data, error) {
	data := structapi.Data{}
	data.Page = path
	data.Muscles = structapi.Muscles
	if path == "/addexercise" {
		exercises, err := sql.GetExercises()
		if err != nil {
			return structapi.Data{}, err
		}
		data.Exercises = exercises
	} else if path == "/htmx/exercises" {
		exercises, err := sql.GetExercises()
		if err != nil {
			return structapi.Data{}, err
		}
		data.Exercises = exercises
	}
	return data, nil
}
