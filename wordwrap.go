// Package wordwrap provides methods for breaking a string on a number of bytes
// in a UTF-8 safe manner such that a rune will never be cut.
package wordwrap

import (
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
	return join(SplitString(s, byteLimit), "\n")
}

// Copied from https://github.com/golang/go/blob/91911e39/src/strings/strings.go#L351-L370 to avoid a large import tree
//
// Copyright 2009 The Go Authors. All rights reserved
// Use of this source code is governed by a BSD-style
//
// Join concatenates the elements of a to create a single string.   The separator string
// sep is placed between elements in the resulting string.
func join(a []string, sep string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return a[0]
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n)
	bp := copy(b, a[0])
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}
