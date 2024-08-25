package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/grep-starter-go/cmd/mygrep/lib"
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

	// ok, err := MatchLine(string(line), pattern)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	// 	os.Exit(2)
	// }

	matcher := lib.InitMatcher().ScanPattern(pattern)
	ok := matcher.Match(line)

	if !ok {
		fmt.Println("error")
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

// func MatchLine(line string, pattern string) (bool, error) {
// 	if utf8.RuneCountInString(pattern) == 0 {
// 		return false, fmt.Errorf("unsupported pattern: %q", pattern)
// 	}

// 	if pattern[0] == '^' {
// 		return matchPatternFunc(line, pattern[1:], 0), nil
// 	}

// 	for i := 0; i < len(line); i++ {
// 		if matchPatternFunc(line, pattern, i) {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }

// func matchPatternFunc(line string, pattern string, pos int) bool {
// 	patternLength := len(pattern)
// 	lineLength := len(line)
// 	lineIndex := pos
// 	var captures []string

// 	for i := 0; i < patternLength; i++ {
// 		if lineIndex >= lineLength {
// 			return pattern[i] == '$'
// 		}

// 		if i+1 < patternLength && pattern[i+1] == '+' { // handling the + pattern
// 			currentChar := pattern[i]
// 			if !matchChar(line, lineIndex, currentChar) {
// 				return false
// 			}
// 			for lineIndex < lineLength && matchChar(line, lineIndex, currentChar) {
// 				lineIndex++
// 			}
// 			i++ // Skip the '+'
// 		} else if i+1 < patternLength && pattern[i+1] == '?' {
// 			currentChar := pattern[i]
// 			if matchChar(line, lineIndex, currentChar) {
// 				lineIndex++
// 			}
// 			i++
// 		} else if pattern[i] == '\\' && i+1 < patternLength {
// 			ptrChr := pattern[i+1]
// 			if ptrChr == 'w' && !(unicode.IsLetter(rune(line[lineIndex])) || unicode.IsDigit(rune(line[lineIndex])) || line[lineIndex] == '_') {
// 				return false
// 			} else if ptrChr == 'd' && !unicode.IsDigit(rune(line[lineIndex])) {
// 				return false
// 			} else if ptrChr >= '1' && ptrChr <= '9' { // handle backreference \1, \2, etc.
// 				refIndex := int(ptrChr - '1')
// 				if refIndex >= len(captures) {
// 					return false
// 				}
// 				backreference := captures[refIndex]
// 				if !strings.HasPrefix(line[lineIndex:], backreference) {
// 					return false
// 				}
// 				lineIndex += len(backreference)
// 			} else {
// 				i++ // Skip the letter
// 				lineIndex++
// 			}
// 		} else if pattern[i] == '[' {
// 			closeSqrBracketPos := strings.IndexAny(pattern[i:], "]")
// 			matchPattern := pattern[i+1 : closeSqrBracketPos]
// 			if i+1 < patternLength && pattern[i+1] == '^' {
// 				if strings.Contains(matchPattern, string(line[lineIndex])) {
// 					return false
// 				}
// 			} else {
// 				if !strings.Contains(matchPattern, string(line[lineIndex])) {
// 					return false
// 				}
// 			}
// 			i = closeSqrBracketPos
// 		} else if pattern[i] == '.' {
// 			lineIndex++
// 		} else if pattern[i] == '(' {
// 			endPos := strings.Index(pattern[i:], ")")
// 			if endPos == -1 {
// 				return matchChar(line, lineIndex, pattern[i])
// 			}
// 			substr := pattern[i+1 : endPos]
// 			parts := strings.Split(substr, "|")
// 			var found bool
// 			for _, subPattern := range parts {
// 				if matchPatternFunc(line, subPattern, lineIndex) {
// 					captures = append(captures, line[pos:lineIndex])
// 					// return true
// 					found = true
// 					break
// 				}
// 			}
// 			if !found {
// 				return false
// 			}
// 		} else if pattern[i] == '\\' && i+1 < patternLength && (pattern[i+1] == 'w' || pattern[i+1] == 'd') {
// 			start := lineIndex
// 			for lineIndex < lineLength && matchChar(line, lineIndex, pattern[i+1]) {
// 				lineIndex++
// 			}
// 			captures = append(captures, line[start:lineIndex])
// 			i++ // Skip the 'w' or 'd'
// 		} else { // base case
// 			if !matchChar(line, lineIndex, pattern[i]) {
// 				return false
// 			}
// 			lineIndex++
// 		}
// 	}

// 	return true
// }

// func matchChar(line string, lineIndex int, char byte) bool {
// 	if char == '\\' {
// 		return false // Should handle escapes like \w, \d, etc.
// 	}
// 	return line[lineIndex] == char
// }
