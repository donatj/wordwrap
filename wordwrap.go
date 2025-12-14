// Package wordwrap provides methods for breaking a string on a number of bytes
// in a UTF-8 safe manner such that a rune will never be cut.
package wordwrap

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

type charPos struct {
	pos, size int
}

// SplitString splits a string at a certain number of bytes without breaking
// UTF-8 runes, grapheme clusters, and on Unicode space characters when possible.
//
// SplitString will panic if it is forced to split a grapheme cluster that is
// larger than the byte limit.
//
// For example if the grapheme cluster `👩‍👩‍👧‍👧` (25 bytes) is given, yet we ask
// it to break on a byte limit of 20, it will panic.
func SplitString(s string, byteLimit uint) []string {
	var workingLine strings.Builder
	finishedLines := []string{}

	spacePos := charPos{}
	lastPos := charPos{}

	// Use grapheme cluster iterator
	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		cluster := gr.Str()
		clusterSize := len(cluster)

		workingLine.WriteString(cluster)

		// Check if the cluster contains a space (check first rune)
		firstRune, _ := utf8.DecodeRuneInString(cluster)
		if unicode.IsSpace(firstRune) {
			spacePos = charPos{workingLine.Len(), clusterSize}
		}

		if workingLine.Len() >= int(byteLimit) {
			if spacePos.size > 0 {
				line := workingLine.String()
				finishedLines = append(finishedLines, line[0:spacePos.pos])

				workingLine.Reset()
				workingLine.WriteString(line[spacePos.pos:])
			} else {
				if workingLine.Len() > int(byteLimit) {
					// If there's no valid break point (lastPos.pos is 0),
					// it means we have a single grapheme cluster larger than byteLimit
					if lastPos.pos == 0 {
						panic("attempted to cut grapheme cluster")
					}
					line := workingLine.String()
					finishedLines = append(finishedLines, line[0:lastPos.pos])
					
					workingLine.Reset()
					workingLine.WriteString(line[lastPos.pos:])
				} else {
					finishedLines = append(finishedLines, workingLine.String())
					workingLine.Reset()
				}
			}

			if len(finishedLines[len(finishedLines)-1]) > int(byteLimit) {
				panic("attempted to cut grapheme cluster")
			}

			spacePos = charPos{}
		}

		lastPos = charPos{workingLine.Len(), clusterSize}
	}

	if workingLine.Len() > 0 {
		if workingLine.Len() > int(byteLimit) {
			panic("attempted to cut grapheme cluster")
		}
		finishedLines = append(finishedLines, workingLine.String())
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
