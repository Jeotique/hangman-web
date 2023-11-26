package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
)

var temp *template.Template
var WordToGuess string
var Errors int
var Guessed []string
var GivenLetters = make(map[string]bool)
var difficulty string

func main() {

	temp, _ = template.ParseGlob("./templates/*.html")

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/play", playPage)
	http.HandleFunc("/play/post", execPost)
	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/asset"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func execPost(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("word")
	if len(data) > 1 {
		if data == WordToGuess {
			WinPage(w)
		} else {
			Errors += 2
			if Errors >= 10 {
				LosePage(w)
			} else {
				RefreshPlayPage(w)
			}
		}
	} else if len(data) == 1 {
		if GivenLetters[data] {
			RefreshPlayPage(w)
		} else {
			GivenLetters[data] = true
			if strings.Contains(WordToGuess, data) {
				Guessed = strings.Split(WordToGuess, "")
				for i, s := range WordToGuess {
					if GivenLetters[string(s)] {
						Guessed[i] = string(s)
					} else {
						Guessed[i] = "_"
					}
				}
				if !slices.Contains(Guessed, "_") {
					WinPage(w)
				} else {
					RefreshPlayPage(w)
				}
			} else {
				Errors += 1
				RefreshPlayPage(w)
			}
		}
	} else {
		RefreshPlayPage(w)
	}
}

func WinPage(w http.ResponseWriter) {
	err := temp.ExecuteTemplate(w, "win", struct {
		Word       string
		Difficulty string
	}{
		Word:       WordToGuess,
		Difficulty: difficulty,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LosePage(w http.ResponseWriter) {
	err := temp.ExecuteTemplate(w, "lose", struct {
		Word       string
		Difficulty string
	}{
		Word:       WordToGuess,
		Difficulty: difficulty,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RefreshPlayPage(w http.ResponseWriter) {
	err := temp.ExecuteTemplate(w, "play", struct {
		Difficulty string
		Guesses    []string
		Word       string
		Errors     int
	}{
		Difficulty: difficulty,
		Guesses:    Guessed,
		Word:       WordToGuess,
		Errors:     Errors,
	})
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
	difficulty = r.FormValue("difficulty")
	WordToGuess = generateWord(difficulty)
	guess := strings.Repeat("_", len(WordToGuess))
	Guessed = strings.Split(guess, "")
	letter1 := rand.Intn(len(WordToGuess))
	letter2 := rand.Intn(len(WordToGuess))
	letter3 := rand.Intn(len(WordToGuess))
	GivenLetters[string(rune(WordToGuess[letter1]))] = true
	Guessed[letter1] = string(rune(WordToGuess[letter1]))
	if difficulty == "easy" || difficulty == "medium" {
		GivenLetters[string(rune(WordToGuess[letter2]))] = true
		Guessed[letter2] = string(rune(WordToGuess[letter2]))
	}
	if difficulty == "easy" {
		GivenLetters[string(rune(WordToGuess[letter3]))] = true
		Guessed[letter3] = string(rune(WordToGuess[letter3]))
	}
	Errors = 0
	err := temp.ExecuteTemplate(w, "play", struct {
		Difficulty string
		Guesses    []string
		Word       string
		Errors     int
	}{
		Difficulty: difficulty,
		Guesses:    Guessed,
		Word:       WordToGuess,
		Errors:     Errors,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateWord(difficulty string) string {
	lines, _ := os.ReadFile(difficulty + ".txt")
	all := strings.Split(string(lines), "\n")
	word := all[rand.Intn(len(all))]
	var finalWord string
	for _, i := range word {
		if i != 13 {
			finalWord += string(rune(i))
		}
	}
	finalWord = strings.ReplaceAll(finalWord, " ", "")
	finalWord = strings.ToLower(finalWord)
	return strings.Trim(finalWord, "\n")
}
