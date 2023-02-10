/*
 * BSD 3-Clause License
 *
 * Copyright (c) 2023, Phea Duch <phea.duch@gmail.com>
 * All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 *
 */

package matcher

import (
	"fmt"
	"regexp"
	"strings"
)

type Matcher struct {
	regex  *regexp.Regexp
	scheme string
	idents []string
	vars   map[string]string
}

// NewMatcher takes a route string and returns a Matcher.
func New(tmpl string) *Matcher {
	regex, err := newTmplRegex(tmpl)
	if err != nil {
		panic(err)
	}

	idents := regex.SubexpNames()

	return &Matcher{regex: regex, idents: idents}
}

// IsMatch takes a route string and checks if it matches the
// service template
func (m *Matcher) IsMatch(route string) bool {
	return m.regex.MatchString(route)
}

// Vars takes a route string and returns a map of the variables
// in the route.
func (m *Matcher) Vars(route string) map[string]string {
	matches := m.regex.FindStringSubmatch(route)
	vars := make(map[string]string)
	for i, match := range matches {
		if i != 0 {
			vars[m.idents[i]] = match
		}
	}
	return vars
}

// newTmplRegex takes a template string and returns a regexp.Regexp
// used to match the route to a service template.
func newTmplRegex(tmpl string) (*regexp.Regexp, error) {
	// check if the route is well formed
	idxs, err := branceIndeces(tmpl)
	if err != nil {
		return nil, err
	}

	// build the regex string
	var b strings.Builder
	b.WriteRune('^')
	// intialize variables for the first start and end indices
	start, end := -1, -1
	curIdx := 0
	if len(idxs) > 0 {
		start = idxs[curIdx]
		end = idxs[curIdx+1]
		curIdx += 2
	}

	for i := 0; i < len(tmpl); i++ {
		if i == start {
			// if the current index is the start of a brace
			// write the regex for the variable
			b.WriteString(`(?P<`)
			b.WriteString(tmpl[start+1 : end])
			b.WriteString(`>[\w]+)`)
			// update the start and end indices
			if curIdx < len(idxs) {
				start = idxs[curIdx]
				end = idxs[curIdx+1]
				curIdx += 2
			}

			// add the length of the brace to the current index
			i += end - start
		} else if tmpl[i] == '/' {
			b.WriteString("\\/")
		} else {
			// otherwise write the character to the regex
			b.WriteByte(tmpl[i])
		}
	}
	return regexp.Compile(b.String())
}

// branceIndices returns the indices of the opening and closing braces
// in the string along with checking if the braces are balanced.
func branceIndeces(str string) ([]int, error) {
	idxs := []int{}
	var level int
	for i, c := range str {
		if c == '{' {
			idxs = append(idxs, i)
			level++

		} else if c == '}' {
			idxs = append(idxs, i)
			level--
			if level < 0 {
				return nil, fmt.Errorf("unbalanced braces")
			}
		}
	}

	if len(idxs)%2 != 0 || level != 0 {
		return nil, fmt.Errorf("unbalanced braces")
	}
	return idxs, nil
}
