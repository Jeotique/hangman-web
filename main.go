package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"project/bdd"
	"project/hangman"
)

var temp *template.Template
var Bdd *bdd.QuickDB

func main() {
	tmp, err := template.ParseGlob("./templates/*.html")
	Bdd = bdd.NewQuickDB("database.json")
	errrr := Bdd.LoadUsers()
	if errrr != nil {
		log.Fatal(errrr)
		return
	}
	temp = tmp
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/play", playPage)
	http.HandleFunc("/play/post", execPost)
	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/asset"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))
	errr := http.ListenAndServe(":8080", nil)
	if errr != nil {
		log.Fatal(errr)
		return
	}
}

func execPost(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("word")
	username := r.FormValue("username")
	result := hangman.TestLetterOrWord(data, username)
	SaveUsers()
	switch result {
	case "win":
		WinPage(w, username)
		break
	case "lose":
		LosePage(w, username)
		break
	case "refresh":
		RefreshPlayPage(w, username)
		break
	}
}

func WinPage(w http.ResponseWriter, username string) {
	err := temp.ExecuteTemplate(w, "win", struct {
		Word       string
		Difficulty string
	}{
		Word:       hangman.Users[username].WordToGuess,
		Difficulty: hangman.Users[username].Difficulty,
	})
	SaveUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LosePage(w http.ResponseWriter, username string) {
	err := temp.ExecuteTemplate(w, "lose", struct {
		Word       string
		Difficulty string
	}{
		Word:       hangman.Users[username].WordToGuess,
		Difficulty: hangman.Users[username].Difficulty,
	})
	SaveUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RefreshPlayPage(w http.ResponseWriter, username string) {
	userData := hangman.GetGameExistingData(username)
	finalWord := ""
	for _, i := range userData.Guessed {
		finalWord += i
	}
	userData.GuessWord = finalWord
	SaveUsers()
	err := temp.ExecuteTemplate(w, "play", userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	err := temp.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func playPage(w http.ResponseWriter, r *http.Request) {
	difficulty := r.FormValue("difficulty")
	username := r.FormValue("username")

	userData := hangman.GetGameData(difficulty, username)
	finalWord := ""
	for _, i := range userData.Guessed {
		finalWord += i
	}
	userData.GuessWord = finalWord
	SaveUsers()
	err := temp.ExecuteTemplate(w, "play", userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SaveUsers() {
	var users []hangman.GameData
	for _, i := range hangman.Users {
		users = append(users, i)
	}
	err := Bdd.SaveUsers(users)
	if err != nil {
		log.Fatal(err)
	}
}
