// Package printer display in standard output simple info/success/error messages
package printer

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	log "github.com/sirupsen/logrus"
)

// Info prints a message in the std output prefixed with 'info'
func Info(format string, args ...interface{}) {
	fmt.Println(aurora.Sprintf("%s: %s", aurora.Blue("info"), fmt.Sprintf(format, args...)))
}

// Success prints a message in sto output in green
func Success(format string, args ...interface{}) {
	fmt.Println(aurora.Sprintf(aurora.Green(format), args...))
}

// Error prints the error with description
func Error(err error, format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf("%s: %s", aurora.Red("error"), err.Error()))
	fmt.Println(aurora.Sprintf(aurora.Red(format), args...))
}

// ErrorText prints an error message
func ErrorText(format string, args ...interface{}) {
	fmt.Println(aurora.Red(fmt.Sprintf("error: "+format, args...)))
}

// FatalIfError prints the error and exits the program if err is not nil
func FatalIfError(err error, code int) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%s: %s", aurora.Red("error"), err.Error()))
		log.WithError(err).Fatal("fatal error")
		log.Exit(code)
	}
}

// FatalWithMessageIfError does the same than Error, then exit, if err is not nil
func FatalWithMessageIfError(err error, code int, format string, args ...interface{}) {
	if err != nil {
		Error(err, format, args...)
		log.WithError(err).Fatalf(format, args...)
		log.Exit(code)
	}
}
