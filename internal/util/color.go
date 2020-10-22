package util

var (
	// Green ANSI Escape Code
	Green = "\033[32m"
	// Red ANSI Escape Code
	Red = "\033[31m"
	// Reset ANSI Escape Code
	Reset = "\033[0m"
	// Yellow ANSI Escape Code
	Yellow = "\033[33m"
)

// Colorize wraps a given message in a given color.
func Colorize(color, message string) string {
	return color + message + Reset
}
