package ui

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type FormUseCase struct {
	TerminalPort PrintReadTerminalPort
}

// NewSimpleForm creates a simple form using standard output and input
func NewSimpleForm() *FormUseCase {
	return &FormUseCase{TerminalPort: &FormFmtAdapter{}}
}

// ReadString read a string from the standard input
func (f *FormUseCase) ReadString(label string, defaultValue string) (string, bool) {
	printedDefaultValue := ""
	if defaultValue != "" {
		printedDefaultValue = fmt.Sprintf(" [%s]", defaultValue)
	}
	f.TerminalPort.Print(fmt.Sprintf("%s%s: ", aurora.Bold(aurora.White(label)), aurora.Gray(12, printedDefaultValue)))

	value, err := f.TerminalPort.ReadAnswer()
	value = strings.Trim(strings.TrimSuffix(value, "\n"), " ")
	if (err != nil || value == "") && defaultValue != "" {
		return defaultValue, true
	}
	return value, err == nil
}

// ReadDate reads a date from standard input, return true if reading was a success.
func (f *FormUseCase) ReadDate(label string, defaultValue time.Time) (time.Time, bool) {
	printedDefaultValue := ""
	if !defaultValue.IsZero() {
		printedDefaultValue = fmt.Sprintf(" [%s]", defaultValue.Format("2006-01-02"))
	}
	f.TerminalPort.Print(fmt.Sprintf("%s%s: ", aurora.Bold(aurora.White(label)), aurora.Gray(12, printedDefaultValue)))

	value, err := f.TerminalPort.ReadAnswer()
	value = strings.Trim(strings.TrimSuffix(value, "\n"), " ")
	if err != nil || value == "" {
		return defaultValue, !defaultValue.IsZero()
	}

	date, err := f.parseDate(value)
	return date, err == nil
}

func (f *FormUseCase) parseDate(value string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02"} {
		parse, err := time.Parse(layout, value)
		if err == nil {
			return parse, nil
		}
	}

	return time.Time{}, errors.Errorf("'%s' is not a valid date, or datetime, format.", value)
}

// ReadBool reads a boolean, notation can be [Y/n] if (_, false) is interpreted as positive
func (f *FormUseCase) ReadBool(label string, notation string) (bool, bool) {
	printedDefaultValue := ""
	if notation != "" {
		printedDefaultValue = fmt.Sprintf(" [%s]", notation)
	}
	f.TerminalPort.Print(fmt.Sprintf("%s%s: ", aurora.Bold(aurora.White(label)), aurora.Gray(12, printedDefaultValue)))

	value, err := f.TerminalPort.ReadAnswer()
	if err != nil {
		return false, false
	}

	answer := strings.ToLower(strings.Trim(strings.TrimSuffix(value, "\n"), " "))
	switch answer {
	case "yes", "oui", "true", "y", "o", "1":
		return true, true

	case "no", "non", "false", "n", "0":
		return false, true
	}

	return false, false
}
