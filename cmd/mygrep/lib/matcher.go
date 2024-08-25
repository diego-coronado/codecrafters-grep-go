package lib

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Matcher struct {
	Chs []*Ch
}

func InitMatcher() *Matcher {
	return &Matcher{}
}

func (m *Matcher) ScanPattern(pattern string) *Matcher {

	chs := []*Ch{}

	var i int

	if pattern[0] == '^' {
		chs = append(chs, &Ch{Type: CharStartAnchor})
		i++
	}

	for i < len(pattern) {
		ch := &Ch{}

		currentChar := pattern[i]
		var nextChar byte

		if i+1 < len(pattern) {
			nextChar = pattern[i+1]
		}

		if currentChar == '\\' && nextChar != '\\' {
			ch.Type = CharEscape
			ch.Value = string(currentChar) + string(nextChar)
			i += 2
		} else if nextChar == '+' {
			ch.Type = CharOneOrMoreTimes
			ch.Value = string(currentChar)
			i += 2
		} else if nextChar == '?' {
			fmt.Printf("? character!!\n")
			ch.Type = CharZeroOrOneTimes
			ch.Value = string(currentChar)
			i += 2
		} else if currentChar == '.' {
			ch.Type = CharWildCard
			ch.Value = ""
			i++
		} else if currentChar == '$' && i == len(pattern)-1 {
			ch.Type = CharEndAnchor
			ch.Value = ""
			chs = append(chs, ch)
			break
		} else if currentChar == '(' {
			endPos := strings.Index(pattern[i:], ")")
			if endPos == -1 {
				ch.Type = CharLiteral
				ch.Value = string(currentChar)
				i++
			} else {
				ch.Type = CharAlternation
				ch.Value = pattern[i+1 : endPos]
				alternateStrArr := strings.Split(pattern[i+1:i+endPos], "|")
				ch.AlternateOptions = alternateStrArr
				i += endPos + 1
			}
		} else if currentChar == '[' {
			closeSqrBracketPos := strings.IndexAny(pattern[i:], "]")
			matchPattern := pattern[i+1 : closeSqrBracketPos]
			if i+1 < len(pattern) && pattern[i+1] == '^' {
				ch.Type = CharNegGroup
				ch.Value = matchPattern[1:]
				i += closeSqrBracketPos + 1
			} else {
				ch.Type = CharPosGroup
				ch.Value = matchPattern
				i += closeSqrBracketPos + 1
			}
		} else {
			ch.Type = CharLiteral
			ch.Value = string(currentChar)
			i++
		}

		chs = append(chs, ch)
	}

	m.Chs = chs
	fmt.Printf("chs len %v\n", len(chs))
	return m
}

func (m *Matcher) Match(line []byte) bool {
	if m.Chs[0].Type == CharStartAnchor {
		// m.Chs = m.Chs[1:]
		return m.MatchLine(line, m.Chs[1:])
	}

	for i := 0; i < len(line); i++ {
		if m.MatchLine(line[i:], m.Chs) {
			return true
		}
	}
	return false
}

func (m *Matcher) MatchLine(line []byte, chs []*Ch) bool {
	fmt.Printf("start matchline\n")
	var i int

	for ind, ch := range chs {
		// check for $ if it is the end of the line
		// fmt.Printf("current ch val %v\n", ch.Value)
		fmt.Printf("current ch type %v\n", ch.Type)

		if i >= len(line) {
			fmt.Printf("end of line %v\n", ch.Value)

			if ch.Type == CharZeroOrOneTimes {
				return true
			}

			return ch.Type == CharEndAnchor
		}

		currLineChar := line[i]
		fmt.Printf("current line char %v\n", string(currLineChar))
		// fmt.Printf("current line char %v\n", string(currLineChar) != ch.Value)

		switch ch.Type {
		case CharLiteral:
			if string(currLineChar) != ch.Value {
				return false
			}
		case CharEscape:
			if ch.Value == "\\w" {
				if !(unicode.IsDigit(rune(currLineChar)) || unicode.IsLetter(rune(currLineChar)) || currLineChar == '_') {
					return false
				}
			}

			if ch.Value == "\\d" {
				if !unicode.IsDigit(rune(currLineChar)) {
					return false
				}
			}
		case CharPosGroup:
			if !bytes.ContainsAny([]byte{currLineChar}, ch.Value) {
				return false
			}
		case CharNegGroup:
			if bytes.ContainsAny([]byte{currLineChar}, ch.Value) {
				return false
			}
		case CharEndAnchor:
			return false
		case CharOneOrMoreTimes:
			if currLineChar != ch.Value[0] {
				return false
			}

			for j := i; j < len(line) && string(currLineChar) == ch.Value; j++ {
				if m.MatchLine(line[j+1:], m.Chs[ind+1:]) {
					return true
				}
			}
			return false
		case CharZeroOrOneTimes:
			if m.MatchLine(line[i:], m.Chs[ind+1:]) {
				return true
			}
			if string(currLineChar) == ch.Value && m.MatchLine(line[i+1:], m.Chs[ind+1:]) {
				return true
			}
			return false
		case CharAlternation:
			for _, alternation := range ch.AlternateOptions {

				if i+len(alternation) > len(line) {
					continue
				}

				if string(line[i:i+len(alternation)]) == alternation {
					if m.MatchLine(line[i+len(alternation):], chs[ind+1:]) {
						return true
					}
				}
			}
			return false
		}
		i++
	}
	return true
}
