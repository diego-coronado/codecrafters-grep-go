package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		fmt.Println("error")
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var ok bool

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	if pattern == "\\d" {
		ok = bytes.ContainsAny(line, "0123456789")
	} else if pattern == "\\w" {
		strLine := string(line)
		checkStrLn := 0
		for _, char := range strLine {
			if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' {
				checkStrLn++
			}
		}
		ok = checkStrLn > 0
	} else if len(pattern) > 2 && pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
		if pattern[1] == '^' { // negative character group
			strLine := string(line)
			ptrn := pattern[2 : len(pattern)-1]
			for _, char := range ptrn {
				strLine = strings.ReplaceAll(strLine, string(char), "")
			}
			ok = len(strLine) > 0
		} else { // positive character group
			ok = bytes.ContainsAny(line, pattern[1:len(pattern)-1])
		}
	} else {
		ok = bytes.ContainsAny(line, pattern)

	}

	return ok, nil
}
