package routes

import (
	"net/http"

	"../middleware"
	"../models"
	"../sessions"
	"../utils"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	staticDir := http.Dir("./static/")
	fileServer := http.FileServer(staticDir)

	r := mux.NewRouter()
	r.HandleFunc("/", middleware.AuthRequired(handleGetIndex)).Methods("GET")
	r.HandleFunc("/", middleware.AuthRequired(handlePostIndex)).Methods("POST")
	r.HandleFunc("/login", handleGetLogin).Methods("GET")
	r.HandleFunc("/login", handlePostLogin).Methods("POST")
	r.HandleFunc("/registrar", handleGetRegistre).Methods("GET")
	r.HandleFunc("/registrar", handlePostRegistre).Methods("POST")
	r.HandleFunc("/detalhes", middleware.AuthRequired(handleDetails)).Methods("GET")
	r.HandleFunc("/contato", middleware.AuthRequired(handleContact)).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	return r
}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {

	comments, err := models.GetComments()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao carregar comentários"))
		return
	}

	utils.ExecuteTemplate(w, "index.html", comments)
}

func handlePostIndex(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	err := models.PostComment(comment)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao salvar comentário"))
		return
	}

	http.Redirect(w, r, "/", 302)
}

func handleGetLogin(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html", nil)
}

func handlePostLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.html", "Usuário Desconhecido")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.html", "Login inválido")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
		}
		return
	}
	session, _ := sessions.Store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func handleGetRegistre(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "registre.html", nil)
}

func handlePostRegistre(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/login", 302)
}

func handleDetails(w http.ResponseWriter, r *http.Request) {

	nome := "Johnatas"
	utils.ExecuteTemplate(w, "details.html", nome)

}

func handleContact(w http.ResponseWriter, r *http.Request) {
	const gitHubUrl string = `<a href="https://github.com/johnatasr?tab=repositories">johnatasr?tab=repositories</a>`
	utils.ExecuteTemplate(w, "contact.html", gitHubUrl)
}
