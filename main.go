package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var templates *template.Template

var bootStrapURL string = `<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" 
								integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">`

type indexStruct struct {
	Languages map[string]string
	bootstrap string
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	date := time.Now()
	templates.ExecuteTemplate(w, "index.html", date.Format("01-02-2006 15:04:05"))
	//fmt.Fprint(w, "Aqui é a pagina principal")
}

func handleDetais(w http.ResponseWriter, r *http.Request) {

	lang := indexStruct{
		Languages: map[string]string{
			"Python":     "É que mais domino",
			"Javascript": "Gosto bastante principalmente do React",
			"Typescript": "Deveria utilizar mais e estudar",
			"PHP":        "Trabalho com mais não gosto muito",
		},
		bootstrap: bootStrapURL,
	}

	templates.ExecuteTemplate(w, "details", lang)
	//fmt.Fprint(w, "Aqui é a página de detalhes")
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	// const gitHubUrl string = `<a href="https://github.com/johnatasr?tab=repositories">johnatasr?tab=repositories</a>`
	templates.ExecuteTemplate(w, "contact.html", nil)
	//fmt.Fprint(w, "Olhai meu Github é https://github.com/johnatasr?tab=repositories")
}

func main() {

	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", handleMain).Methods("GET")
	r.HandleFunc("/detalhes", handleDetais).Methods("GET")
	r.HandleFunc("/contato", handleContact).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
