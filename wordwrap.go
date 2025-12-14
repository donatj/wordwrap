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
	continueOnError         bool
	breakGraphemeClusters   bool
	trimTrailingWhiteSpace  bool
}

// DefaultSplitBuilder is the global default SplitBuilder used by package-level Split function.
// It can be modified to change default splitting behavior globally.
var DefaultSplitBuilder = &SplitBuilder{
	continueOnError:         false,
	breakGraphemeClusters:   false,
	trimTrailingWhiteSpace:  false,
}

// SplitBuilderOption is a functional option for configuring a SplitBuilder.
type SplitBuilderOption func(*SplitBuilder)

// NewSplitBuilder creates a new SplitBuilder with the given options.
// By default, it matches the behavior of SplitString:
//   - continueOnError: false (returns error on grapheme cluster too large)
//   - breakGraphemeClusters: false (preserves grapheme clusters)
//   - trimTrailingWhiteSpace: false (keeps trailing whitespace)
func NewSplitBuilder(opts ...SplitBuilderOption) *SplitBuilder {
	sb := &SplitBuilder{
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

// Split is a package-level function that uses DefaultSplitBuilder to split a string.
// It returns an iterator that yields line content and error pairs.
func Split(s string, byteLimit uint) iter.Seq2[string, error] {
	return DefaultSplitBuilder.Split(s, byteLimit)
}

// Split returns an iterator that yields line content and error pairs.
// The iterator processes the input string according to the SplitBuilder's configuration.
// The byteLimit parameter specifies the maximum number of bytes per line.
// Each iteration returns a line (string) and an error. If there's no error for that line, error will be nil.
// If ContinueOnError is false (default), iteration stops on the first error.
func (sb *SplitBuilder) Split(s string, byteLimit uint) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		var workingLine strings.Builder

		spacePos := charPos{}
		lastPos := charPos{}

		gr := uniseg.NewGraphemes(s)
		for gr.Next() {
			cluster := gr.Str()
			clusterSize := len(cluster)

			// Check if cluster alone is too large
			if clusterSize > int(byteLimit) && !sb.breakGraphemeClusters {
				// Flush working line first if non-empty
				if workingLine.Len() > 0 {
					line := workingLine.String()
					if sb.trimTrailingWhiteSpace {
						line = strings.TrimRight(line, " \t\n\r")
					}
					if !yield(line, nil) {
						return
					}
					workingLine.Reset()
					spacePos = charPos{}
					lastPos = charPos{}
				}
				
				// Yield the oversized cluster with error
				clusterStr := cluster
				if sb.trimTrailingWhiteSpace {
					clusterStr = strings.TrimRight(clusterStr, " \t\n\r")
				}
				if !yield(clusterStr, ErrGraphemeClusterTooLarge) {
					return
				}
				if !sb.continueOnError {
					return
				}
				continue
			}

			// If breaking grapheme clusters is allowed and the cluster is too large,
			// break it down to individual runes
			if sb.breakGraphemeClusters && clusterSize > int(byteLimit) {
				for _, r := range cluster {
					runeBytes := []byte(string(r))
					runeSize := len(runeBytes)
					
					workingLine.Write(runeBytes)
					
					if workingLine.Len() >= int(byteLimit) {
						line := workingLine.String()
						if sb.trimTrailingWhiteSpace {
							line = strings.TrimRight(line, " \t\n\r")
						}
						if !yield(line, nil) {
							return
						}
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

			if workingLine.Len() >= int(byteLimit) {
				if spacePos.size > 0 {
					line := workingLine.String()
					linePart := line[0:spacePos.pos]
					if sb.trimTrailingWhiteSpace {
						linePart = strings.TrimRight(linePart, " \t\n\r")
					}
					if !yield(linePart, nil) {
						return
					}

					workingLine.Reset()
					workingLine.WriteString(line[spacePos.pos:])
				} else {
					if workingLine.Len() > int(byteLimit) {
						if lastPos.pos == 0 {
							// Single grapheme cluster larger than byteLimit
							// This should be caught earlier when breakGraphemeClusters is false
							line := workingLine.String()
							if sb.trimTrailingWhiteSpace {
								line = strings.TrimRight(line, " \t\n\r")
							}
							if !yield(line, ErrGraphemeClusterTooLarge) {
								return
							}
							if !sb.continueOnError {
								return
							}
							workingLine.Reset()
						} else {
							line := workingLine.String()
							linePart := line[0:lastPos.pos]
							if sb.trimTrailingWhiteSpace {
								linePart = strings.TrimRight(linePart, " \t\n\r")
							}
							if !yield(linePart, nil) {
								return
							}

							workingLine.Reset()
							workingLine.WriteString(line[lastPos.pos:])
						}
					} else {
						line := workingLine.String()
						if sb.trimTrailingWhiteSpace {
							line = strings.TrimRight(line, " \t\n\r")
						}
						if !yield(line, nil) {
							return
						}
						workingLine.Reset()
					}
				}

				spacePos = charPos{}
			}

			lastPos = charPos{workingLine.Len(), clusterSize}
		}

		if workingLine.Len() > 0 {
			line := workingLine.String()
			if sb.trimTrailingWhiteSpace {
				line = strings.TrimRight(line, " \t\n\r")
			}
			var err error
			if workingLine.Len() > int(byteLimit) {
				err = ErrGraphemeClusterTooLarge
			}
			yield(line, err)
		}
	}
}

// SplitString splits a string and returns the lines as a slice.
// It implements the same functionality as the global SplitString function
// while respecting the SplitBuilder's configured flags.
// It returns an error if any line encounters an error during splitting.
func (sb *SplitBuilder) SplitString(s string, byteLimit uint) ([]string, error) {
	var lines []string
	var firstErr error
	
	for line, err := range sb.Split(s, byteLimit) {
		lines = append(lines, line)
		if err != nil && firstErr == nil {
			firstErr = err
			if !sb.continueOnError {
				return lines, firstErr
			}
		}
	}
	
	return lines, firstErr
}

// SplitString splits a string at a certain number of bytes without breaking
// UTF-8 runes, grapheme clusters, and on Unicode space characters when possible.
//
// SplitString returns an ErrGraphemeClusterTooLarge error if it is forced to split
// a grapheme cluster that is larger than the byte limit.
//
// For example if the grapheme cluster `👩‍👩‍👧‍👧` (25 bytes) is given, yet we ask
// it to break on a byte limit of 20, it will return an error.
//
// This function uses the DefaultSplitBuilder with settings that match the original behavior.
func SplitString(s string, byteLimit uint) ([]string, error) {
	return DefaultSplitBuilder.SplitString(s, byteLimit)
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
