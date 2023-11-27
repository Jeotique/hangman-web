package hangman

import (
	"log"
	"math/rand"
	"os"
	"slices"
	"strings"
)

func GenerateWord(difficulty string) string {
	println("generate word")
	lines, err := os.ReadFile(difficulty + ".txt")
	if err != nil {
		log.Fatal(err)
		return ""
	}
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
	println(finalWord)
	return strings.Trim(finalWord, "\n")
}

func TestLetterOrWord(data string, user GameData) string {
	if len(data) > 1 {
		if data == user.WordToGuess {
			return "win"
		} else {
			user.Errors += 2
			if user.Errors >= 10 {
				return "lose"
			} else {
				return "refresh"
			}
		}
	} else if len(data) == 1 {
		println("test letter len 1")
		if user.GivenLetters[data] {
			return "refresh"
		} else {
			user.GivenLetters[data] = true
			if strings.Contains(user.WordToGuess, data) {
				user.Guessed = strings.Split(user.WordToGuess, "")
				for i, s := range user.WordToGuess {
					if user.GivenLetters[string(s)] {
						user.Guessed[i] = string(s)
					} else {
						user.Guessed[i] = "_"
					}
				}
				if !slices.Contains(user.Guessed, "_") {
					return "win"
				} else {
					return "refresh"
				}
			} else {
				user.Errors += 1
				println("erreur")
				println(user.Errors)
				return "refresh"
			}
		}
	} else {
		return "refresh"
	}
	return "refresh"
}

func GetGameData(difficulty string, username string) GameData {
	WordToGuess := GenerateWord(difficulty)
	guess := strings.Repeat("_", len(WordToGuess))
	Guessed := strings.Split(guess, "")
	letter1 := rand.Intn(len(WordToGuess))
	letter2 := rand.Intn(len(WordToGuess))
	letter3 := rand.Intn(len(WordToGuess))
	Users[username] = GameData{
		Difficulty:   difficulty,
		Errors:       0,
		WordToGuess:  WordToGuess,
		Guessed:      Guessed,
		GivenLetters: make(map[string]bool),
		GuessWord:    strings.Repeat("_", len(WordToGuess)),
		Username:     username,
	}
	Users[username].GivenLetters[string(rune(WordToGuess[letter1]))] = true
	Users[username].Guessed[letter1] = string(rune(WordToGuess[letter1]))
	if difficulty == "easy" || difficulty == "medium" {
		Users[username].GivenLetters[string(rune(WordToGuess[letter2]))] = true
		Users[username].Guessed[letter2] = string(rune(WordToGuess[letter2]))
	}
	if difficulty == "easy" {
		Users[username].GivenLetters[string(rune(WordToGuess[letter3]))] = true
		Users[username].Guessed[letter3] = string(rune(WordToGuess[letter3]))
	}
	return Users[username]
}

func GetGameExistingData(username string) GameData {
	return Users[username]
}
