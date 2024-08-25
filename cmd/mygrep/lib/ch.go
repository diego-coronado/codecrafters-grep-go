package lib

type CharType string

const (
	CharLiteral        CharType = "literal"
	CharEscape         CharType = "escape"
	CharPosGroup       CharType = "positiveGroup"
	CharNegGroup       CharType = "negativeGroup"
	CharStartAnchor    CharType = "startAnchor"
	CharEndAnchor      CharType = "endAnchor"
	CharOneOrMoreTimes CharType = "oneOrMore"
	CharZeroOrOneTimes CharType = "zeroOrOne"
	CharWildCard       CharType = "wildcard"
	CharAlternation    CharType = "alternation"
)

type Ch struct {
	Type  CharType
	Value string

	PrevCh           *Ch // Pointer to the previous Ch, used for the + operator
	AlternateOptions []string
}
