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

func TestLetterOrWord(data string, username string) string {
	User := Users[username]
	if len(data) > 1 {
		if data == User.WordToGuess {
			return "win"
		} else {
			User.Errors += 2
			Users[username] = User
			if User.Errors >= 10 {
				User.IsInGame = false
				Users[username] = User
				return "lose"
			} else {
				return "refresh"
			}
		}
	} else if len(data) == 1 {
		if User.GivenLetters[data] {
			return "refresh"
		} else {
			User.GivenLetters[data] = true
			if strings.Contains(User.WordToGuess, data) {
				User.Guessed = strings.Split(User.WordToGuess, "")
				for i, s := range User.WordToGuess {
					if User.GivenLetters[string(s)] {
						User.Guessed[i] = string(s)
					} else {
						User.Guessed[i] = "_"
					}
				}
				Users[username] = User
				if !slices.Contains(User.Guessed, "_") {
					User.IsInGame = false
					Users[username] = User
					return "win"
				} else {
					return "refresh"
				}
			} else {
				User.Errors += 1
				Users[username] = User
				if User.Errors >= 10 {
					User.IsInGame = false
					Users[username] = User
					return "lose"
				} else {
					return "refresh"
				}
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
	if value, exists := Users[username]; exists {
		if value.IsInGame {
			return value
		}
	}
	Users[username] = GameData{
		Difficulty:   difficulty,
		Errors:       0,
		WordToGuess:  WordToGuess,
		Guessed:      Guessed,
		GivenLetters: make(map[string]bool),
		GuessWord:    strings.Repeat("_", len(WordToGuess)),
		Username:     username,
		IsInGame:     true,
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
