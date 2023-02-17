package service

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/phea/mio/internal/matcher"
)

// TestJSONTemplate tests the JSON templates to see if they match
// given routes.
func TestJSONTemplate(t *testing.T) {
	var jsonMatchers []*matcher.Matcher
	for _, tmpl := range jsonTemplates {
		jsonMatchers = append(jsonMatchers, matcher.New(tmpl))
	}

	tests := []struct {
		route string
		match bool
	}{
		{"json://localhost", true},
		{"json://localhost:8080", true},
		{"json://localhost:8080/abc/def", true},
		{"json://localhost:8080/abc/def/ghi", true},
		{"json://test@localhost/abc/def", true},
		{"json://test@localhost:8080/abc/def", true},
		{"json://test:password@localhost/abc/def", true},
		{"json://test:password@localhost:8080/abc/def", true},
	}

	for _, test := range tests {
		var matched bool
		for _, m := range jsonMatchers {
			if m.IsMatch(test.route) {
				matched = true
				break
			}
		}

		if matched != test.match {
			t.Errorf("expected %s to match %t", test.route, test.match)
		}
	}
}

// TestJSONInit tests the JSON service initialization.
func TestJSONInit(t *testing.T) {
	svc := defaultJSONService()
	m := matcher.New("json://{user}:{pass}@{host}:{port}")
	mVars := m.Vars("json://test:pass@localhost:8080/abc/def")

	err := svc.Init("json://test:pass@localhost:8080/abc/def", mVars, SetTLS(false))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if svc.isTLS {
		t.Errorf("expected isTLS to be false, got %t", svc.isTLS)
	}

	vars := svc.vars
	if vars["user"] != "test" {
		t.Errorf("expected user to be test, got %s", vars["user"])
	}

	if vars["pass"] != "pass" {
		t.Errorf("expected pass to be pass, got %s", vars["pass"])
	}

	if vars["host"] != "localhost" {
		t.Errorf("expected host to be localhost, got %s", vars["host"])
	}

	if vars["port"] != "8080" {
		t.Errorf("expected port to be 8080, got %s", vars["port"])
	}
}

// TestJSONSend tests the JSON service send.
func TestJSONSend(t *testing.T) {
	// create a test server to handle the request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("intercepted request")
		if r.Method != "POST" {
			t.Errorf("expected method to be POST, got %s", r.Method)
		}

		if r.URL.Path != "/abc/def" {
			t.Errorf("expected path to be /abc/def, got %s", r.URL.Path)
		}

		// check the headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		// unmarshal the body and check if title and body are correct
		var data map[string]string
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if data["title"] != "title" {
			t.Errorf("expected title to be title, got %s", data["title"])
		}

		if data["body"] != "message" {
			t.Errorf("expected message to be message, got %s", data["message"])
		}
	}))
	defer ts.Close()

	route := "json://" + strings.TrimPrefix(ts.URL, "http://") + "/abc/def"
	svc := defaultJSONService()
	svc.Init(route, nil, SetTLS(false))
	err := svc.Send("title", "message")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
