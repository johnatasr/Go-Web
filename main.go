package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var client *redis.Client
var templates *template.Template

var bootStrapURL string = `<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" 
								integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">`

type indexStruct struct {
	Languages map[string]string
	bootstrap string
}

var store = sessions.NewCookieStore([]byte("G0lAng"))

func handleGetLogin(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func handlePostLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Values["password"] = password
	session.Save(r, w)
}

func testGetHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	untyped, ok := session.Values["username"]

	if !ok {
		return
	}

	username, ok := untyped.(string)
	if !ok {
		return
	}

	w.Write([]byte(username))
}

func handleGetMain(w http.ResponseWriter, r *http.Request) {
	comments, err := client.LRange("comments", 0, 10).Result()

	if err != nil {
		return
	}

	date := time.Now()
	templates.ExecuteTemplate(w, "index.html", date.Format("01-02-2006 15:04:05"))

	if comments != nil {
		templates.ExecuteTemplate(w, "index.html", comments)
	}
	//fmt.Fprint(w, "Aqui é a pagina principal")
}

func handlePostMain(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	client.LPush("comments", comment)
	http.Redirect(w, r, "/", 302)
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

	staticDir := http.Dir("./static/")
	fileServer := http.FileServer(staticDir)

	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/login", handleGetLogin).Methods("GET")
	r.HandleFunc("/login", handlePostLogin).Methods("POST")
	r.HandleFunc("/", handleGetMain).Methods("GET")
	r.HandleFunc("/", handlePostMain).Methods("POST")
	r.HandleFunc("/detalhes", handleDetais).Methods("GET")
	r.HandleFunc("/contato", handleContact).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
