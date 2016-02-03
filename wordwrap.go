// Package wordwrap provides methods for breaking a string on a number of bytes
// in a UTF-8 safe manner such that a rune will never be cut.
package wordwrap

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type charPos struct {
	pos, size int
}

// SplitString splits a string at a certain number of bytes without breaking
// UTF-8 runes and on Unicode space characters when possible.
//
// SplitString will panic if it is forced to split a multibyte rune.
//
// For example if the rune `し` (3 bytes) is given, yet we ask it to break on a
// byte limit of 2, it will panic.
func SplitString(s string, byteLimit uint) []string {
	workingLine := ""
	finishedLines := []string{}

	spacePos := charPos{}
	lastPos := charPos{}

	for _, r := range s {
		rl := utf8.RuneLen(r)

		workingLine += string(r)

		if unicode.IsSpace(r) {
			spacePos = charPos{len(workingLine), rl}
		}

		if len(workingLine) >= int(byteLimit) {
			if spacePos.size > 0 {
				finishedLines = append(finishedLines, workingLine[0:spacePos.pos])

				workingLine = workingLine[spacePos.pos:]
			} else {
				if len(workingLine) > int(byteLimit) {
					finishedLines = append(finishedLines, workingLine[0:lastPos.pos])
					workingLine = workingLine[lastPos.pos:]
				} else {
					finishedLines = append(finishedLines, workingLine)
					workingLine = ""
				}
			}

			if len(finishedLines[len(finishedLines)-1]) > int(byteLimit) {
				panic("attempted to cut character")
			}

			spacePos = charPos{}
		}

		lastPos = charPos{len(workingLine), rl}
	}

	if workingLine != "" {
		finishedLines = append(finishedLines, workingLine)
	}

	return finishedLines
}

// WrapString splits a string as with SplitString and joins together with a \n
func WrapString(s string, byteLimit uint) string {
	return strings.Join(SplitString(s, byteLimit), "\n")
}
