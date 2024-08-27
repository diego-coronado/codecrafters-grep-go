package lib

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func popChFromChs(chs []*Ch) ([]*Ch, *Ch) {
	length := len(chs)
	if length == 0 {
		return chs, nil
	}

	return chs[:length-1], chs[length-1]

}

type Result struct {
	Matched      bool
	EndPositions []int
}

type Matcher struct {
	Chs           []*Ch
	CaptureGroups []string
}

func InitMatcher() *Matcher {
	return &Matcher{}
}

func (m *Matcher) ScanPatternInternal(pattern string) []*Ch {
	chs := make([]*Ch, 0)

	i := 0
	groupIndex := 0

	// detect start of string line anchor
	if strings.HasPrefix(pattern, "^") {
		chs = append(chs, &Ch{
			Type:  CharStartAnchor,
			Value: "",
		})
		i++
	}
	for i < len(pattern) {
		currentChar := pattern[i]
		var nextChar byte

		if i+1 < len(pattern) {
			nextChar = pattern[i+1]
		}

		if currentChar == '$' && i == len(pattern)-1 {
			chs = append(chs, &Ch{
				Type:  CharEndAnchor,
				Value: "",
			})
			break
		}

		if currentChar == '+' {
			poppedChs, lastElement := popChFromChs(chs)
			chs = append(poppedChs, &Ch{
				Type:   CharOneOrMoreTimes,
				Value:  "",
				PrevCh: lastElement,
			})
			i++
			continue
		}

		if currentChar == '\\' && nextChar != '\\' {
			if unicode.IsDigit(rune(nextChar)) {
				chs = append(chs, &Ch{
					Type:  CharBackReference,
					Value: fmt.Sprintf("%c", nextChar),
				})
			} else {
				chs = append(chs, &Ch{
					Type:  CharEscape,
					Value: fmt.Sprintf("%c%c", currentChar, nextChar),
				})
			}

			i += 2
			continue
		}

		if nextChar == '?' {
			chs = append(chs, &Ch{
				Type:  CharZeroOrOneTimes,
				Value: string(currentChar),
			})
			i += 2
			continue
		}

		if currentChar == '.' {
			chs = append(chs, &Ch{
				Type:  CharWildCard,
				Value: "",
			})
			i++
			continue
		}

		if currentChar == '(' {
			endPos := strings.Index(pattern[i:], ")")
			if endPos != -1 {
				groupIndex = groupIndex + 1
				alterStrList := strings.Split(pattern[i+1:i+endPos], "|")
				if len(alterStrList) > 1 {
					ch := &Ch{
						Type:             CharAlternation,
						Value:            "",
						AlternateOptions: make([][]*Ch, 0),
						GroupIndex:       groupIndex,
					}

					for _, alterStr := range alterStrList {
						ch.AlternateOptions = append(ch.AlternateOptions, m.ScanPatternInternal(alterStr))
					}

					chs = append(chs, ch)
				} else {
					chs = append(chs, &Ch{
						Type:             CharCaptureGroup,
						Value:            pattern[i+1 : i+endPos],
						AlternateOptions: nil,
						GroupElements:    m.ScanPatternInternal(pattern[i+1 : i+endPos]),
						GroupIndex:       groupIndex,
					})
				}

				m.CaptureGroups = append(m.CaptureGroups, "")

				i = i + endPos + 1

				continue
			}
		}

		if currentChar == '[' {
			endPos := strings.Index(pattern[i:], "]")
			if endPos != -1 {
				charGroup := pattern[i+1 : i+endPos]
				charType := CharPosGroup
				if charGroup[0] == '^' {
					charType = CharNegGroup
					charGroup = charGroup[1:]
				}
				chs = append(chs, &Ch{
					Type:  charType,
					Value: charGroup,
				})
				i = i + endPos + 1
				continue
			}
		}
		chs = append(chs, &Ch{
			Type:  CharLiteral,
			Value: string(currentChar),
		})
		i++
	}
	return chs
}

func (m *Matcher) ScanPattern(pattern string) *Matcher {
	m.CaptureGroups = make([]string, 1)
	m.Chs = m.ScanPatternInternal(pattern)
	return m
}

func (m *Matcher) Match(line []byte) bool {
	if m.Chs[0].Type == CharStartAnchor {
		r := m.MatchLine(line, m.Chs[1:])
		return r.Matched
	}

	for i := 0; i < len(line); i++ {
		if r := m.MatchLine(line[i:], m.Chs); r.Matched {
			return true
		}
	}
	return false
}

func (m *Matcher) MatchLine(line []byte, chs []*Ch) *Result {
	i := 0
	res := &Result{
		Matched:      true,
		EndPositions: make([]int, 0),
	}

	for ind, ch := range chs {
		if i >= len(line) {

			if ch.Type == CharEndAnchor {
				res.EndPositions = append(res.EndPositions, i)
				break
			}

			res.Matched = false
			break
		}

		currLineChar := line[i]

		switch ch.Type {
		case CharLiteral, CharEscape, CharPosGroup, CharNegGroup:
			if !m.MatchBasic(currLineChar, ch) {
				res.Matched = false
			}
		case CharEndAnchor:
			res.Matched = false
		case CharOneOrMoreTimes:
			hasMatched := false
			for j := i; j < len(line) && m.MatchBasic(line[j], ch.PrevCh); j++ {
				if recRes := m.MatchLine(line[j+1:], chs[ind+1:]); recRes.Matched {
					for _, endPos := range recRes.EndPositions {
						res.EndPositions = append(res.EndPositions, j+1+endPos)
					}
					hasMatched = true
				}
			}

			if !hasMatched {
				res.Matched = false
			}
			return res
		case CharZeroOrOneTimes:
			if matchRes := m.MatchLine(line[i:], chs[ind+1:]); matchRes.Matched {
				res = matchRes
				return res
			}
			if string(currLineChar) == ch.Value {
				if mr := m.MatchLine(line[i+1:], chs[ind+1:]); mr.Matched {
					res = mr
					break
				}
				res.Matched = false
				return res
			}
			res.Matched = false
			return res
		case CharWildCard:
			i++
			continue
		case CharAlternation:
			recRes := false
			for _, alterValue := range ch.AlternateOptions {
				if matchedAlter := m.MatchLine(line[i:], alterValue); matchedAlter.Matched {
					for _, endPosAlter := range matchedAlter.EndPositions {
						nextI := i + endPosAlter

						m.CaptureGroups[ch.GroupIndex] = string(line[i:nextI])
						if mr := m.MatchLine(line[nextI:], chs[ind+1:]); mr.Matched {
							for _, endPos := range mr.EndPositions {
								res.EndPositions = append(res.EndPositions, nextI+endPos)
							}
							recRes = true
						}
					}

				}
			}
			if !recRes {
				res.Matched = false
			}
			return res
		case CharCaptureGroup:

			sr := false
			if mg := m.MatchLine(line[i:], ch.GroupElements); mg.Matched {
				for _, mgEnd := range mg.EndPositions {
					m.CaptureGroups[ch.GroupIndex] = string(line[i : i+mgEnd])

					nextI := i + mgEnd
					if mr := m.MatchLine(line[nextI:], chs[ind+1:]); mr.Matched {
						sr = true
						for _, endPos := range mr.EndPositions {
							res.EndPositions = append(res.EndPositions, nextI+endPos)
						}
					}
				}
			}

			res.Matched = sr
			return res
		case CharBackReference:
			groupIndex, _ := strconv.Atoi(ch.Value)
			groupValue := m.CaptureGroups[groupIndex]
			nextInd := i + len(groupValue)

			if string(line[i:nextInd]) == groupValue {
				i = nextInd
				continue
			}
			res.Matched = false
			return res
		}
		i++
	}
	if res.Matched {
		res.EndPositions = append(res.EndPositions, i)
	}
	return res
}

func (m *Matcher) MatchBasic(currChar byte, ch *Ch) bool {
	switch ch.Type {
	case CharLiteral:
		if string(currChar) == ch.Value {
			return true
		}
	case CharEscape:
		switch ch.Value {
		case "\\w":
			if unicode.IsDigit(rune(currChar)) || unicode.IsLetter(rune(currChar)) || currChar == '_' {
				return true
			}
		case "\\d":
			if unicode.IsDigit(rune(currChar)) {
				return true
			}
		}
	case CharPosGroup:
		if bytes.ContainsAny([]byte{currChar}, ch.Value) {
			return true
		}
	case CharNegGroup:
		if !bytes.ContainsAny([]byte{currChar}, ch.Value) {
			return true
		}

	}
	return false
}
