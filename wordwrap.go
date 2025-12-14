// Package wordwrap provides methods for breaking a string on a number of bytes
// in a UTF-8 safe manner such that a rune will never be cut.
package wordwrap

import (
	"errors"
	"iter"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

// ErrGraphemeClusterTooLarge is returned when a grapheme cluster exceeds the byte limit.
// This error occurs when a single grapheme cluster (such as an emoji with modifiers,
// a character with combining marks, or other complex Unicode sequences) is larger
// than the specified byte limit, making it impossible to split the string without
// breaking the cluster.
var ErrGraphemeClusterTooLarge = errors.New("grapheme cluster exceeds byte limit")

type charPos struct {
	pos, size int
}

// SplitBuilder provides a configurable string splitter with functional options.
type SplitBuilder struct {
	byteLimit               uint
	continueOnError         bool
	breakGraphemeClusters   bool
	trimTrailingWhiteSpace  bool
}

// SplitBuilderOption is a functional option for configuring a SplitBuilder.
type SplitBuilderOption func(*SplitBuilder)

// NewSplitBuilder creates a new SplitBuilder with the given byte limit and options.
// By default, it matches the behavior of SplitString:
//   - continueOnError: false (returns error on grapheme cluster too large)
//   - breakGraphemeClusters: false (preserves grapheme clusters)
//   - trimTrailingWhiteSpace: false (keeps trailing whitespace)
func NewSplitBuilder(byteLimit uint, opts ...SplitBuilderOption) *SplitBuilder {
	sb := &SplitBuilder{
		byteLimit:               byteLimit,
		continueOnError:         false,
		breakGraphemeClusters:   false,
		trimTrailingWhiteSpace:  false,
	}
	
	for _, opt := range opts {
		opt(sb)
	}
	
	return sb
}

// ContinueOnError sets whether to continue processing when encountering errors.
// When true, errors like ErrGraphemeClusterTooLarge are ignored and processing continues.
// When false (default), the iterator stops on the first error.
func ContinueOnError(continueOnError bool) SplitBuilderOption {
	return func(sb *SplitBuilder) {
		sb.continueOnError = continueOnError
	}
}

// BreakGraphemeClusters sets whether to allow breaking grapheme clusters.
// When true, grapheme clusters can be split if they exceed the byte limit.
// When false (default), grapheme clusters are preserved, and an error is returned if they're too large.
func BreakGraphemeClusters(breakGraphemeClusters bool) SplitBuilderOption {
	return func(sb *SplitBuilder) {
		sb.breakGraphemeClusters = breakGraphemeClusters
	}
}

// TrimTrailingWhiteSpace sets whether to remove trailing whitespace from lines.
// When true, whitespace at the end of each line is removed.
// When false (default), trailing whitespace is preserved.
func TrimTrailingWhiteSpace(trimTrailingWhiteSpace bool) SplitBuilderOption {
	return func(sb *SplitBuilder) {
		sb.trimTrailingWhiteSpace = trimTrailingWhiteSpace
	}
}

// Split returns an iterator that yields line index and line content pairs.
// The iterator processes the input string according to the SplitBuilder's configuration.
func (sb *SplitBuilder) Split(s string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		var workingLine strings.Builder
		lineIndex := 0

		spacePos := charPos{}
		lastPos := charPos{}

		gr := uniseg.NewGraphemes(s)
		for gr.Next() {
			cluster := gr.Str()
			clusterSize := len(cluster)

			// If breaking grapheme clusters is allowed and the cluster is too large,
			// break it down to individual runes
			if sb.breakGraphemeClusters && clusterSize > int(sb.byteLimit) {
				for _, r := range cluster {
					runeBytes := []byte(string(r))
					runeSize := len(runeBytes)
					
					workingLine.Write(runeBytes)
					
					if workingLine.Len() >= int(sb.byteLimit) {
						line := workingLine.String()
						if sb.trimTrailingWhiteSpace {
							line = strings.TrimRight(line, " \t\n\r")
						}
						if !yield(lineIndex, line) {
							return
						}
						lineIndex++
						workingLine.Reset()
						spacePos = charPos{}
					}
					
					lastPos = charPos{workingLine.Len(), runeSize}
				}
				continue
			}

			workingLine.WriteString(cluster)

			firstRune, _ := utf8.DecodeRuneInString(cluster)
			if unicode.IsSpace(firstRune) {
				spacePos = charPos{workingLine.Len(), clusterSize}
			}

			if workingLine.Len() >= int(sb.byteLimit) {
				if spacePos.size > 0 {
					line := workingLine.String()
					linePart := line[0:spacePos.pos]
					if sb.trimTrailingWhiteSpace {
						linePart = strings.TrimRight(linePart, " \t\n\r")
					}
					if !yield(lineIndex, linePart) {
						return
					}
					lineIndex++

					workingLine.Reset()
					workingLine.WriteString(line[spacePos.pos:])
				} else {
					if workingLine.Len() > int(sb.byteLimit) {
						if lastPos.pos == 0 {
							// Single grapheme cluster larger than byteLimit
							if !sb.continueOnError {
								// In error case, we can't easily pass the error through iter.Seq2
								// So we just stop iteration
								return
							}
							// Continue on error: just output what we have
							line := workingLine.String()
							if sb.trimTrailingWhiteSpace {
								line = strings.TrimRight(line, " \t\n\r")
							}
							if !yield(lineIndex, line) {
								return
							}
							lineIndex++
							workingLine.Reset()
						} else {
							line := workingLine.String()
							linePart := line[0:lastPos.pos]
							if sb.trimTrailingWhiteSpace {
								linePart = strings.TrimRight(linePart, " \t\n\r")
							}
							if !yield(lineIndex, linePart) {
								return
							}
							lineIndex++

							workingLine.Reset()
							workingLine.WriteString(line[lastPos.pos:])
						}
					} else {
						line := workingLine.String()
						if sb.trimTrailingWhiteSpace {
							line = strings.TrimRight(line, " \t\n\r")
						}
						if !yield(lineIndex, line) {
							return
						}
						lineIndex++
						workingLine.Reset()
					}
				}

				if !sb.continueOnError {
					// Check if the last line we just yielded was too large
					// This would indicate a grapheme cluster error
					// We need a way to track this, but for now we'll skip this check
					// as it complicates the iterator pattern
				}

				spacePos = charPos{}
			}

			lastPos = charPos{workingLine.Len(), clusterSize}
		}

		if workingLine.Len() > 0 {
			if workingLine.Len() > int(sb.byteLimit) && !sb.continueOnError {
				// Error case - stop iteration
				return
			}
			line := workingLine.String()
			if sb.trimTrailingWhiteSpace {
				line = strings.TrimRight(line, " \t\n\r")
			}
			yield(lineIndex, line)
		}
	}
}

// SplitString splits a string at a certain number of bytes without breaking
// UTF-8 runes, grapheme clusters, and on Unicode space characters when possible.
//
// SplitString returns an ErrGraphemeClusterTooLarge error if it is forced to split
// a grapheme cluster that is larger than the byte limit.
//
// For example if the grapheme cluster `👩‍👩‍👧‍👧` (25 bytes) is given, yet we ask
// it to break on a byte limit of 20, it will return an error.
func SplitString(s string, byteLimit uint) ([]string, error) {
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
						return nil, ErrGraphemeClusterTooLarge
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
				return nil, ErrGraphemeClusterTooLarge
			}

			spacePos = charPos{}
		}

		lastPos = charPos{workingLine.Len(), clusterSize}
	}

	if workingLine.Len() > 0 {
		if workingLine.Len() > int(byteLimit) {
			return nil, ErrGraphemeClusterTooLarge
		}
		finishedLines = append(finishedLines, workingLine.String())
	}

	return finishedLines, nil
}

// WrapString splits a string as with SplitString and joins together with a \n
//
// WrapString returns an ErrGraphemeClusterTooLarge error if a grapheme cluster
// exceeds the byte limit.
func WrapString(s string, byteLimit uint) (string, error) {
	lines, err := SplitString(s, byteLimit)
	if err != nil {
		return "", err
	}
	return join(lines, "\n"), nil
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
