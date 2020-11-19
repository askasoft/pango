package str

const (
	// Digits A String for digits "0123456789"
	Digits = "0123456789"

	// LowerLetters A String for lower letters "abcdefghijklmnopqrstuvwxyz"
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// UpperLetters A String for upper letters "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Symbols A String for symbols "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Symbols = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

	// DigitLetters A String for digits and letters
	DigitLetters = Digits + LowerLetters + UpperLetters

	// SymbolDigits A String for symbols and digits "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~0123456789"
	SymbolDigits = Symbols + Digits

	// SymbolDigitLetters A String for symbols, digits and letters
	SymbolDigitLetters = Symbols + DigitLetters
)
