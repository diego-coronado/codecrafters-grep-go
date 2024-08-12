package main

import (
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

	ok, err := matchLine(string(line), pattern)
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

func matchLine(line string, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	for i := 0; i < len(line); i++ {
		if matchPattern(line, pattern, i) {
			return true, nil
		}
	}

	return false, nil
}

func matchPattern(line string, pattern string, pos int) bool {
	length := len(pattern)
	linePos := pos

	for i := 0; i < length; i++ {
		if linePos >= len(line) {
			return false
		}

		if pattern[i] == '\\' && i+1 < length {
			ptrChr := pattern[i+1]
			if ptrChr == 'w' && !(unicode.IsLetter(rune(line[linePos])) || unicode.IsDigit(rune(line[linePos])) || line[linePos] == '_') {
				return false
			} else if ptrChr == 'd' && !unicode.IsDigit(rune(line[linePos])) {
				return false
			} else {
				i++
			}
		} else if pattern[i] == '[' {
			closeSqrBracketPos := strings.IndexAny(pattern[i:], "]")
			matchPattern := pattern[i+1 : closeSqrBracketPos]
			if i+1 < length && pattern[i+1] == '^' {
				if strings.Contains(matchPattern, string(line[linePos])) {
					return false
				}
			} else {
				if !strings.Contains(matchPattern, string(line[linePos])) {
					return false
				}
			}
			i = closeSqrBracketPos
		} else { // base case
			if line[linePos] != pattern[i] {
				return false
			}
		}
		linePos++
	}

	return true
}
