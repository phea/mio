package service

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/phea/mio/internal/matcher"
)

// TestXMLTemplate tests the XML templates to see if they match
// given routes.
func TestXMLTemplate(t *testing.T) {
	var xmlMatchers []*matcher.Matcher
	for _, tmpl := range xmlTemplates {
		xmlMatchers = append(xmlMatchers, matcher.New(tmpl))
	}

	tests := []struct {
		route string
		match bool
	}{
		{"xml://localhost", true},
		{"xml://localhost:8080", true},
		{"xml://localhost:8080/abc/def", true},
		{"xml://localhost:8080/abc/def/ghi", true},
		{"xml://test@localhost/abc/def", true},
		{"xml://test@localhost:8080/abc/def", true},
		{"xml://test:password@localhost/abc/def", true},
		{"xml://test:password@localhost:8080/abc/def", true},
	}

	for _, test := range tests {
		var matched bool
		for _, m := range xmlMatchers {
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

// TestXMLInit tests the XML service initialization.
func TestXMLInit(t *testing.T) {
	svc := defaultXMLService()
	m := matcher.New("xml://{user}:{pass}@{host}:{port}")
	mVars, err := m.Vars("xml://test:pass@localhost:8080/abc/def")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = svc.Init("xml://test:pass@localhost:8080/abc/def", mVars, SetTLS(false))
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

// TestXMLSend tests the XML service send.
func TestXMLSend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method to be POST, got %s", r.Method)
		}

		// check headers
		if r.Header.Get("Content-Type") != "application/xml" {
			t.Errorf("expected content type to be application/xml, got %s", r.Header.Get("Content-Type"))
		}

		// unmarshal body
		var data xmlPayload
		err := xml.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if data.Title != "test" {
			t.Errorf("expected title to be test, got %s", data.Title)
		}

		if data.Body != "message" {
			t.Errorf("expected body to be test body, got %s", data.Body)
		}
	}))
	defer ts.Close()

	route := "xml://" + strings.TrimPrefix(ts.URL, "http://") + "/abc/def"
	svc := defaultXMLService()
	svc.Init(route, nil, SetTLS(false))
	err := svc.Send("test", "message")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
