package web

import (
	"fmt"
	did2 "github.com/0x51-dev/did/did"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var secure = true

// DisableHTTPS disables HTTPS for the resolver.
// NOT RECOMMENDED FOR PRODUCTION USE.
func DisableHTTPS() {
	secure = false
}

func Resolve(didURL string, u did2.DID, options did2.ResolutionOptions) did2.ResolutionResult {
	if options.Accept != "application/did+jsonutils" {
		return did2.ResolutionResult{Metadata: did2.Metadata{Error: did2.RepresentationNotSupportedError}}
	}
	path := fmt.Sprintf("%s/.well-known/did.jsonutils", decodeURI(u.MethodIDs[0]))
	if 1 < len(u.MethodIDs) {
		path = fmt.Sprintf("%s/did.jsonutils", joinMap(u.MethodIDs, "/", decodeURI))
	}

	host := "http"
	if secure {
		host += "s"
	}

	resp, err := http.Get(fmt.Sprintf("%s://%s", host, path))
	if err != nil {
		return did2.ResolutionResult{Metadata: did2.Metadata{Error: did2.NotFoundError}}
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return did2.ResolutionResult{Metadata: did2.Metadata{Error: did2.InvalidDIDError}}
	}
	document, err := did2.ParseDocument(raw)
	if err != nil {
		return did2.ResolutionResult{Metadata: did2.Metadata{Error: did2.InvalidDIDError}}
	}

	if document.ID.String() != didURL {
		return did2.ResolutionResult{Metadata: did2.Metadata{Error: did2.NotFoundError}}
	}

	return did2.ResolutionResult{
		Metadata: did2.Metadata{
			ContentType: "application/did+jsonutils",
		},
		Document: document,
	}
}

func decodeURI(s string) string {
	var uri string
	for i, s := range strings.Split(s, "%") {
		if i == 0 && s[i] != '%' {
			uri += s
			continue
		}
		v, _ := strconv.ParseInt(s[:2], 16, 32)
		uri += fmt.Sprintf("%c%s", rune(v), s[2:])
	}
	return uri
}

func joinMap(s []string, sep string, m func(s string) string) string {
	var v string
	for i, s := range s {
		if i == 0 {
			v += m(s)
			continue
		}
		v += fmt.Sprintf("%s%s", sep, m(s))
	}
	return v
}
