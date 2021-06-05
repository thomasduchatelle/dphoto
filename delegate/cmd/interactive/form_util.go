package interactive

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"os"
	"strings"
	"time"
)

func scanString(label string, defaultValue string) (string, bool) {
	reader := bufio.NewReader(os.Stdin)
	printedDefaultValue := ""
	if defaultValue != "" {
		printedDefaultValue = fmt.Sprintf(" [%s]", defaultValue)
	}
	fmt.Printf("%s%s: ", aurora.Bold(aurora.White(label)), aurora.Gray(12, printedDefaultValue))

	value, err := reader.ReadString('\n')
	value = strings.Trim(strings.TrimSuffix(value, "\n"), " ")
	if (err != nil || value == "") && defaultValue != "" {
		return defaultValue, true
	}
	return value, err == nil
}

func scanDate(label string, defaultValue time.Time) (time.Time, bool) {
	reader := bufio.NewReader(os.Stdin)
	printedDefaultValue := ""
	if !defaultValue.IsZero() {
		printedDefaultValue = fmt.Sprintf(" [%s]", defaultValue.Format("2006-01-02"))
	}
	fmt.Printf("%s%s: ", aurora.Bold(aurora.White(label)), aurora.Gray(12, printedDefaultValue))

	value, err := reader.ReadString('\n')
	value = strings.Trim(strings.TrimSuffix(value, "\n"), " ")
	if err != nil || value == "" {
		return defaultValue, !defaultValue.IsZero()
	}

	date, err := parseDate(value)
	return date, err == nil
}

func parseDate(value string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02"} {
		parse, err := time.Parse(layout, value)
		if err == nil {
			return parse, nil
		}
	}

	return time.Time{}, errors.Errorf("'%s' is not a valid date, or datetime, format.", value)
}

// scanBool reads a boolean, notation can be [Y/n] if (_, false) is interpreted as positive
func scanBool(label string, notation string) (bool, bool) {
	reader := bufio.NewReader(os.Stdin)
	printedDefaultValue := ""
	if notation != "" {
		printedDefaultValue = fmt.Sprintf(" [%s]", notation)
	}
	fmt.Printf("%s%s: ", aurora.Bold(aurora.White(label)), aurora.Gray(12, printedDefaultValue))

	value, err := reader.ReadString('\n')
	if err != nil {
		return false, false
	}

	answer := strings.ToLower(strings.Trim(strings.TrimSuffix(value, "\n"), " "))
	fmt.Println("answer", answer)
	switch answer {
	case "yes", "oui", "true", "y", "o", "1":
		return true, true

	case "no", "non", "false", "n", "0":
		return false, true
	}

	fmt.Printf("switch failed '%s'", answer)
	return false, false
}
