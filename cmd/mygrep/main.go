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

	ok, err := MatchLine(string(line), pattern)
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

func MatchLine(line string, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	if pattern[0] == '^' {
		return matchPattern(line, pattern[1:], 0), nil
	}

	for i := 0; i < len(line); i++ {
		if matchPattern(line, pattern, i) {
			return true, nil
		}
	}

	return false, nil
}

func matchPattern(line string, pattern string, pos int) bool {
	patternLength := len(pattern)
	lineLength := len(line)
	lineIndex := pos

	for i := 0; i < patternLength; i++ {
		fmt.Println(i)
		if lineIndex >= lineLength {
			return pattern[i] == '$'
		}

		if i+1 < patternLength && pattern[i+1] == '+' { // handling the + pattern
			currentChar := pattern[i]
			if !matchChar(line, lineIndex, currentChar) {
				return false
			}
			for lineIndex < lineLength && matchChar(line, lineIndex, currentChar) {
				lineIndex++
			}
			i++ // Skip the '+'
		} else if i+1 < patternLength && pattern[i+1] == '?' {
			currentChar := pattern[i]
			if matchChar(line, lineIndex, currentChar) {
				lineIndex++
			}
			i++
		} else if pattern[i] == '\\' && i+1 < patternLength {
			ptrChr := pattern[i+1]
			if ptrChr == 'w' && !(unicode.IsLetter(rune(line[lineIndex])) || unicode.IsDigit(rune(line[lineIndex])) || line[lineIndex] == '_') {
				return false
			} else if ptrChr == 'd' && !unicode.IsDigit(rune(line[lineIndex])) {
				return false
			} else {
				i++ // Skip the letter
				lineIndex++
			}
		} else if pattern[i] == '[' {
			closeSqrBracketPos := strings.IndexAny(pattern[i:], "]")
			matchPattern := pattern[i+1 : closeSqrBracketPos]
			if i+1 < patternLength && pattern[i+1] == '^' {
				if strings.Contains(matchPattern, string(line[lineIndex])) {
					return false
				}
			} else {
				if !strings.Contains(matchPattern, string(line[lineIndex])) {
					return false
				}
			}
			i = closeSqrBracketPos
		} else { // base case
			if !matchChar(line, lineIndex, pattern[i]) {
				return false
			}
			lineIndex++
		}
	}

	return true
}

func matchChar(line string, lineIndex int, char byte) bool {
	if char == '\\' {
		return false // Should handle escapes like \w, \d, etc.
	}
	return line[lineIndex] == char
}
