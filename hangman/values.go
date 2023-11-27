package hangman

type GameData struct {
	WordToGuess  string
	Errors       int
	Guessed      []string
	GuessWord    string
	GivenLetters map[string]bool
	Difficulty   string
	Username     string
	IsInGame     bool
}

var Users = make(map[string]GameData)
