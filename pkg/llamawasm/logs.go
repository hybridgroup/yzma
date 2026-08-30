//go:build js && wasm

package llamawasm

// The llama.cpp module cannot take a Go function as its log callback, so the
// shim holds the callback and this package only sets how much it prints. The
// values below stand in for the function pointer that llama.LogSet takes, so
// that the same program builds for both.
const (
	// LogNormal lets llama.cpp print its messages to the console.
	LogNormal uintptr = 0

	// logSilentValue stops every message. LogSilent returns it.
	logSilentValue uintptr = 1

	// The levels are the same as the ggml_log_level values.
	logLevelWarn   = 3
	logLevelSilent = 99
)

// LogSet sets how much llama.cpp prints. Pass LogSilent to stop the messages,
// or LogNormal to let them go to the console of the browser.
func LogSet(cb uintptr) {
	if !Loaded() {
		return
	}

	level := logLevelWarn
	if cb == logSilentValue {
		level = logLevelSilent
	}
	callVoid("_yzma_log_set_verbosity", level)
}

// LogSilent gives the value that stops the messages of llama.cpp. Pass it to
// LogSet.
func LogSilent() uintptr {
	return logSilentValue
}
