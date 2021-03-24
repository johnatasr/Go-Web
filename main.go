package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var client *redis.Client
var templates *template.Template

var bootStrapURL string = `<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" 
								integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">`

type indexStruct struct {
	Languages map[string]string
	bootstrap string
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	comments, err := client.LRange("comments", 0, 10).Result()

	if err != nil {
		return
	}

	date := time.Now()
	templates.ExecuteTemplate(w, "index.html", date.Format("01-02-2006 15:04:05"))
	templates.ExecuteTemplate(w, "index.html", comments)
	//fmt.Fprint(w, "Aqui é a pagina principal")
}

func handleDetais(w http.ResponseWriter, r *http.Request) {

	// languages := map[string]string{
	// 	"Python":     "É que mais domino",
	// 	"Javascript": "Gosto bastante principalmente do React",
	// 	"Typescript": "Deveria utilizar mais e estudar",
	// 	"PHP":        "Trabalho com mais não gosto muito",
	// }

	// bootstrap := bootStrapURL

	// arr := []interface{}{languages, bootstrap}

	nome := "Johnatas"

	templates.ExecuteTemplate(w, "details.html", nome)
	//fmt.Fprint(w, "Aqui é a página de detalhes")
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	const gitHubUrl string = `<a href="https://github.com/johnatasr?tab=repositories">johnatasr?tab=repositories</a>`
	templates.ExecuteTemplate(w, "contact.html", gitHubUrl)
	//fmt.Fprint(w, "Olhai meu Github é https://github.com/johnatasr?tab=repositories")
}

func main() {

	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", handleMain).Methods("GET")
	r.HandleFunc("/detalhes", handleDetais).Methods("GET")
	r.HandleFunc("/contato", handleContact).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
