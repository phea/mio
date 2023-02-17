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

	var idents []string
	// loop through the subexpressions and add them to the idents if they
	// are not empty
	for _, name := range regex.SubexpNames() {
		if name != "" {
			idents = append(idents, name)
		}
	}

	// extract the scheme from the template
	scheme := strings.Split(tmpl, "://")[0]
	return &Matcher{regex: regex, scheme: scheme, idents: idents}
}

// IsMatch takes a route string and checks if it matches the
// service template
func (m *Matcher) IsMatch(route string) bool {
	return m.regex.MatchString(route)
}

// Vars takes a route string and returns a map of the variables
// in the route.
func (m *Matcher) Vars(route string) map[string]string {
	vars := make(map[string]string)
	// given a route string, extract the values for each of the idents
	matches := m.regex.FindStringSubmatch(route)

	// loop through subexpressions and add them to the vars if name
	// is not empty
	for i, name := range m.regex.SubexpNames() {
		if name != "" && i <= len(matches) {
			vars[name] = matches[i]
		}
	}
	return vars
}

// regex string for valid hostname including ip4 and ip6 addresses
var hostRegex = `(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z]|[A-Za-z][A-Za-z0-9\-]*[A-Za-z0-9]))`

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
			id := tmpl[start+1 : end]
			b.WriteString(`(?P`)
			b.WriteString("<" + id + ">")
			if id == "host" {
				b.WriteString(hostRegex)
			} else {
				b.WriteString(`[\w]+)`)
			}

			// add the length of the brace to the current index
			i += end - start
			// update the start and end indices
			if curIdx < len(idxs) {
				start = idxs[curIdx]
				end = idxs[curIdx+1]
				curIdx += 2
			}

		} else if tmpl[i] == '/' {
			b.WriteString("\\/")
		} else {
			// otherwise write the character to the regex
			b.WriteByte(tmpl[i])
		}
	}

	b.WriteString(".*")
	fmt.Println(b.String())
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
