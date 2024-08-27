package main

import (
	"testing"

	"github.com/codecrafters-io/grep-starter-go/cmd/mygrep/lib"
)

func TestMatchLine(t *testing.T) {
	type MatchLineTest struct {
		line     []byte
		pattern  string
		expected bool
	}
	var matchLineTests = []MatchLineTest{
		// {[]byte("dog"), "d", true},
		// {[]byte("dog"), "f", false},
		// {[]byte("123"), "\\d", true},
		// {[]byte("apple"), "\\d", false},
		// {[]byte("word3"), "\\w", true},
		// {[]byte("$!?"), "\\w", false},
		// {[]byte("a"), "[abcd]", true},
		// {[]byte("efgh"), "[abcd]", false},
		// {[]byte("apple"), "[^xyz]", true},
		// {[]byte("banana"), "[^anb]", false},
		// {[]byte("sally has 3 apples"), "\\d apple", true},
		// {[]byte("sally has 1 orange"), "\\d apple", false},
		// {[]byte("sally has 124 apples"), "\\d\\d\\d apple", true},
		// {[]byte("sally has 12 apples"), "\\d\\d\\d apples", false},
		// {[]byte("sally has 3 dogs"), "\\d \\w\\w\\ws", true},
		// {[]byte("sally has 4 dogs"), "\\d \\w\\w\\ws", true},
		// {[]byte("sally has 1 dog"), "\\d \\w\\w\\ws", false},
		// {[]byte("log"), "^log", true},
		// {[]byte("slog"), "^log", false},
		// {[]byte("cat"), "cat$", true},
		// {[]byte("cats"), "cat$", false},
		// {[]byte("caats"), "ca+ts", true},
		// {[]byte("caaaaats"), "ca+ts", true},
		// {[]byte("apple"), "a+", true},
		// {[]byte("Sass"), "a+", true},
		// // {[]byte("d33ts"), "d\\d+ts", true},
		// {[]byte("hello"), "ca+ts", false},
		// {[]byte("dog"), "a+", false},
		// {[]byte("cat"), "ca?t", true},
		// {[]byte("ct"), "ca?t", true},
		// // {[]byte("dog"), "dogs?", true},
		// {[]byte("dog"), "d.g", true},
		// {[]byte("cog"), "d.g", false},
		// {[]byte("a cat"), "a (cat|dog)", true},
		// {[]byte("a dog"), "a (cat|dog)", true},
		// {[]byte("a cow"), "a (cat|dog)", false},
		// {[]byte("a cat b"), "a (cat|dog) a", false},
		{[]byte("cat and cat"), "(cat) and \\1", true},
		{[]byte("cat and cat"), "(\\w+) and \\1", true},
		{[]byte("cat and dog"), "(\\w+) and \\1", false},
	}

	for _, tt := range matchLineTests {
		actual := lib.InitMatcher().ScanPattern(tt.pattern).Match(tt.line)
		if actual != tt.expected {
			t.Errorf("Matchline(%s, %s) - actual: %t, expected: %t", string(tt.line), tt.pattern, actual, tt.expected)
		}
	}
}
