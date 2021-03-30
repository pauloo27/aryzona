package logger

type Level struct {
	Name, Color string
	// Error = print to stderr and print stacktrace
	// Fatal = exit the program
	Error, Fatal bool
}
