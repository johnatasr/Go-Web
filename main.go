package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
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

	hash, err := client.Get("user: " + username).Bytes()

	if err == redis.Nil {
		templates.ExecuteTemplate(w, "login.html", "Usuário desconhecido")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))

	if err != nil {
		templates.ExecuteTemplate(w, "login.html", "Login inválido")
		return
	}

	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Values["password"] = password
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
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

func handleGetRegistre(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "registre.html", nil)
}

func handlePostRegistre(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: Erro em gerar senha"))
		return
	}

	err = client.Set("user: "+username, hash, 0).Err()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: Erro ao salver cliente"))
		return
	}

	http.Redirect(w, r, "/login", 302)
}

func handleGetMain(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	_, ok := session.Values["username"]

	if !ok {
		http.Redirect(w, r, "/login", 302)
		return
	}

	comments, err := client.LRange("comments", 0, 10).Result()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao carregar comentários"))
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
	err := client.LPush("comments", comment).Err()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao salvar comentário"))
		return
	}
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

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
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
	r.HandleFunc("/registrar", handleGetRegistre).Methods("GET")
	r.HandleFunc("/registrar", handlePostRegistre).Methods("POST")
	r.HandleFunc("/", AuthRequired(handleGetMain)).Methods("GET")
	r.HandleFunc("/", AuthRequired(handlePostMain)).Methods("POST")
	r.HandleFunc("/detalhes", AuthRequired(handleDetais)).Methods("GET")
	r.HandleFunc("/contato", AuthRequired(handleContact)).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
