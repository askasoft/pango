package str

const (
	// LettersLower A String for lower letters "a-z"
	LettersLower = "abcdefghijklmnopqrstuvwxyz"

	// LettersUpper A String for upper letters "A-Z"
	LettersUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Letters A String for letters "a-zA-Z"
	Letters = LettersLower + LettersUpper

	// Numbers A String for numbers "0123456789"
	Numbers = "0123456789"

	// Symbols A String for symbols "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Symbols = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

	// LetterNumbers A String for letters and numbers
	LetterNumbers = Letters + Numbers

	// SymbolNumbers A String for symbols and numbers "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~0123456789"
	SymbolNumbers = Symbols + Numbers

	// LetterNumberSymbols A String for letters, numbers and symbols
	LetterNumberSymbols = Symbols + Letters + Numbers

	// Base64 a-zA-Z0-9+/
	Base64 = LetterNumbers + "+/"

	// Base64URL a-zA-Z0-9-_
	Base64URL = LetterNumbers + "-_"
)
