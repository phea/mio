package matcher

import "testing"

// TestNewValidTemplate tests that when a valid template is passed to New,
// a Matcher is returned with the correct fields.
func TestNewValidTemplate(t *testing.T) {
	tmpl := "https://{host}/{path}"
	m := New(tmpl)
	if m.regex == nil {
		t.Errorf("regex is nil")
	}
	if m.scheme != "https" {
		t.Errorf("scheme is not https")
	}
	if len(m.idents) != 2 {
		t.Errorf("idents expected length 2, got %d", len(m.idents))
	}
	if m.idents[0] != "host" {
		t.Errorf("idents[0] is not host")
	}
	if m.idents[1] != "path" {
		t.Errorf("idents[1] is not path")
	}
}

// TestNewInvalidTemplate tests that when an invalid template is passed to New,
// a panic is raised.
func TestNewInvalidTemplate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("did not panic")
		}
	}()

	// a test table of invalid templates
	tests := []string{
		"https://{host}/{path",
		"https://{host}/{path}}",
		"https://{host}/{path}{",
		"https://{host}/{path}}{",
		"https://{host}/{path}",
	}

	for _, test := range tests {
		New(test)
	}
}

// TestIsMatch tests that IsMatch returns true when the route matches the
// template and false when it does not.
func TestIsMatch(t *testing.T) {
	// test Matcher to use for the tests
	m := New("https://{host}/{path}")

	// a test table with routes and whether they should match the template
	tests := []struct {
		route string
		match bool
	}{
		{"https://example.com/path", true},
		{"https://example.com/path/to/something", true},
		{"https://example.com/path/to/something?query=string", true},
		{"https://example.com/path/to/something?query=string#fragment", true},
		{"https://example.com/path/to/something?query=string#fragment?query=string", true},
		{"http://example.com/path/to/something", false},
		{"httpx://", false},
		{"https://host.com", false},
	}

	for _, test := range tests {
		if m.IsMatch(test.route) != test.match {
			t.Errorf("expected %v to match %v", test.route, test.match)
		}
	}
}

// TestVars tests that Vars returns the correct map of variables when the
// route matches the template.
func TestVars(t *testing.T) {
	// test Matcher to use for the tests
	m := New("https://{host}/{path1}/{path2}")

	// a test table with routes and the expected variables
	tests := []struct {
		route string
		vars  map[string]string
	}{
		{"https://example.com/path1/path2", map[string]string{"host": "example.com", "path1": "path1", "path2": "path2"}},
		{"https://example.com/path1/path2?query=string", map[string]string{"host": "example.com", "path1": "path1", "path2": "path2"}},
		{"https://example.com/path1/path2?query=string#fragment", map[string]string{"host": "example.com", "path1": "path1", "path2": "path2"}},
		{"https://example.com/path1/path2?query=string#fragment?query=string", map[string]string{"host": "example.com", "path1": "path1", "path2": "path2"}},
		{"https://example/path1/path2", map[string]string{"host": "example", "path1": "path1", "path2": "path2"}},
		{"https://example/path1/path2/path3", map[string]string{"host": "example", "path1": "path1", "path2": "path2"}},
		{"http://example.com/path1/path2", map[string]string{}},
	}

	for _, test := range tests {
		vars := m.Vars(test.route)
		if len(vars) != len(test.vars) {
			t.Errorf("expected %v variables, got %v", len(test.vars), len(vars))
		}
		for k, v := range test.vars {
			if vars[k] != v {
				t.Errorf("expected %v to be %v, got %v", k, v, vars[k])
			}
		}
	}
}
