package ui

import (
	"bufio"
	"fmt"
	"os"
)

// FormFmtAdapter is a simple adapter to use fmt package to print and read to standard outputs
type FormFmtAdapter struct{}

func (f FormFmtAdapter) Print(question string) {
	fmt.Print(question)
}

func (f FormFmtAdapter) ReadAnswer() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}
