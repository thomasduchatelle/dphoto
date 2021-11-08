package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		panic("Usage: go run . <reject files> <dir>")
	}

	rejectFile, _ := filepath.Abs(os.Args[1])
	dest, _ := filepath.Abs(os.Args[2])

	err := os.MkdirAll(dest, 0744)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(rejectFile)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewScanner(file)

	for reader.Scan() {
		s3Path := reader.Text()
		filename := path.Base(s3Path)
		dir := path.Dir(strings.TrimPrefix(s3Path, "s3://dush-photos/"))

		//fmt.Printf("aws s3 cp '%s' '%s'\n", line, path.Join(dest, dir, filename))
		_ = os.MkdirAll(path.Join(dest, dir), 0744)
		output, err := exec.Command("aws", "s3", "cp", s3Path, path.Join(dest, dir, filename)).Output()
		if err != nil {
			fmt.Printf("[%s] Error: %s\n", s3Path, err)
			continue
		} else {
			fmt.Print(string(output))
		}

		output, err = exec.Command("aws", "s3", "rm", s3Path).Output()
		if err != nil {
			fmt.Printf("[%s] Error: %s\n", s3Path, err)
			continue
		} else {
			fmt.Print(string(output))
		}
	}
	if reader.Err() != nil {
		panic(reader.Err())
	}

	fmt.Println("All done.")
}
