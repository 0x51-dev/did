package did_test

import (
	"fmt"
	"github.com/0x51-dev/did/did"
	"testing"
)

func ExampleDID() {
	u, _ := did.ParseDID("did:example:test:21tDAKCERh95uGgKbJNHYp;service=agent;foo:bar=high/some/path?foo=bar#key1")
	fmt.Println(u.DID())
	fmt.Println(u.String())
	// Output:
	// did:example:test:21tDAKCERh95uGgKbJNHYp
	// did:example:test:21tDAKCERh95uGgKbJNHYp;service=agent;foo:bar=high/some/path?foo=bar#key1
}

func TestParseDID_examples(t *testing.T) {
	// SOURCE: https://www.w3.org/TR/did-core/#did-url-syntax
	for i, test := range []string{
		"did:example:123456/path",
		"did:example:123456?versionId=1",
		"did:example:123#public-key-0",                                  // A unique verification method in a DID Document.
		"did:example:123#agent",                                         // A unique service in a DID Document.
		"did:example:123?service=agent&relativeRef=/credentials#degree", // A resource external to a DID Document.
		"did:example:123?versionTime=2021-05-10T17:00:00Z",              // A DID URL with a 'versionTime' DID parameter.
		"did:example:123?service=files&relativeRef=/resume.pdf",         // A DID URL with a 'service' and a 'relativeRef' DID parameter.
	} {
		t.Run(fmt.Sprintf("Example%d", i+2), func(t *testing.T) {
			if _, err := did.ParseDID(test); err != nil {
				t.Errorf("ParseDID(%q): %v", test, err)
			}
		})
	}
}
