package str

const (
	// LowerLetters A String for lower letters "a-z"
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// Letters A String for upper letters "A-Z"
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Letters A String for letters "a-zA-Z"
	Letters = LowerLetters + UpperLetters

	// Numbers A String for numbers "0123456789"
	Numbers = "0123456789"

	// Symbols A String for symbols "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Symbols = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

	// LetterNumbers A String for letters and numbers
	LetterNumbers = Letters + Numbers

	// SymbolNumbers A String for symbols and numbers "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~0123456789"
	SymbolNumbers = Symbols + Numbers

	// LetterDigitSymbols A String for letters, numbers and symbols
	LetterDigitSymbols = Symbols + Letters + Numbers

	// Base64 a-zA-Z0-9+/
	Base64 = LetterNumbers + "+/"

	// Base64URL a-zA-Z0-9-_
	Base64URL = LetterNumbers + "-_"
)
