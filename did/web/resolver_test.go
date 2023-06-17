package web

import (
	"fmt"
	did2 "github.com/0x51-dev/did/did"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	for _, test := range []struct {
		methodID string
		expected string
	}{
		{"example.com:user:alice", "example.com:user:alice"},
		{"example.com:path:some%2Bsubpath", "example.com:path:some+subpath"},
		{"localhost%3A8443", "localhost:8443"},
	} {
		if u := decodeURI(test.methodID); u != test.expected {
			t.Error(u, test.expected)
		}
	}
}

func TestResolve(t *testing.T) {
	DisableHTTPS()

	var s *httptest.Server
	s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/did.jsonutils":
			u := didFromServer(s)
			document := did2.Document{
				Context: []string{"https://www.w3.org/ns/did/v1"},
				ID:      u,
				Authentication: []did2.IVerificationMethod{
					&did2.RelativeVerificationMethod{RelativeURL: fmt.Sprintf("%s#owner", u.String())},
				},
			}
			raw, err := document.MarshalJSON()
			if err != nil {
				t.Error(err)
			}
			if _, err := w.Write(raw); err != nil {
				t.Error(err)
			}
		}
	}))

	u := didFromServer(s)
	result := Resolve(u.String(), u, did2.ResolutionOptions{Accept: "application/did+jsonutils"})
	if result.Metadata.Error != "" {
		t.Error(result.Metadata.Error)
	}
	if result.Metadata.ContentType != "application/did+jsonutils" {
		t.Error(result.Metadata.ContentType)
	}
	if result.Document.ID.String() != u.String() {
		t.Error(result.Document.ID.String(), u.String())
	}
}

func didFromServer(s *httptest.Server) did2.DID {
	url := strings.TrimPrefix(s.URL, "http://")
	url = strings.ReplaceAll(url, ":", "%3A")

	u, _ := did2.ParseDID(fmt.Sprintf("did:web:%s", url))
	return *u
}
