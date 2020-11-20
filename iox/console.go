package iox

// ConsoleColor console color
var ConsoleColor = defineColor()

type color struct {
	Red     string
	Magenta string
	Yellow  string
	Blue    string
	White   string
	Gray    string
	Reset   string
}

func defineColor() *color {
	return &color{
		Red:     "\x1b[91m",
		Magenta: "\x1b[95m",
		Yellow:  "\x1b[93m",
		Blue:    "\x1b[94m",
		White:   "\x1b[97m",
		Gray:    "\x1b[90m",
		Reset:   "\x1b[0m",
	}
}
