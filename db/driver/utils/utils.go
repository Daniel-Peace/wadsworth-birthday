package custom_utils

// This type is here to simply enforce types when using the colorizeString function
type AnsiColor string

// These are all of the ansi colors and the reset character
var (
	Reset   AnsiColor = "\033[0m"
	Red     AnsiColor = "\033[31m"
	Green   AnsiColor = "\033[32m"
	Yellow  AnsiColor = "\033[33m"
	Blue    AnsiColor = "\033[34m"
	Magenta AnsiColor = "\033[35m"
	Cyan    AnsiColor = "\033[36m"
	Gray    AnsiColor = "\033[37m"
	White   AnsiColor = "\033[97m"
)

// This function takes a string to colorize and the desired color.
// It returns a string with the color code as a prefix and the reset code as a postfix.
func ColorizeString(s string, c AnsiColor) string {
	return string(c) + s + string(Reset)
}
