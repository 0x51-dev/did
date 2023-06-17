package did

import (
	_ "embed"
	"fmt"
)

var (
	//go:embed testdata/example1.json
	example1 []byte
	//go:embed testdata/example9.json
	example9 []byte
	//go:embed testdata/example10.json
	example10 []byte
	//go:embed testdata/example11.json
	example11 []byte
	//go:embed testdata/example13.json
	example13 []byte
	//go:embed testdata/example14.json
	example14 []byte
	//go:embed testdata/example20.json
	example20 []byte
)

func ExampleDID_controller() {
	doc, _ := ParseDocument(example11)
	fmt.Println(doc)
	// Output:
	// {
	// 	"@context": "https://www.w3.org/ns/did/v1",
	// 	"controller": "did:example:bcehfew7h32f32h7af3",
	// 	"id": "did:example:123456789abcdefghi"
	// }
}

func ExampleDID_subject() {
	doc, _ := ParseDocument(example10)
	fmt.Println(doc)
	// Output:
	// {
	// 	"id": "did:example:123456789abcdefghijk"
	// }
}

// SOURCE: https://www.w3.org/TR/did-core/#example-a-simple-did-document
func ExampleDocument() {
	doc, _ := ParseDocument(example1)
	fmt.Println(doc)
	// Output:
	// {
	// 	"@context": [
	// 		"https://www.w3.org/ns/did/v1",
	// 		"https://w3id.org/security/suites/ed25519-2020/v1"
	// 	],
	// 	"authentication": [
	// 		{
	// 			"controller": "did:example:123456789abcdefghi",
	// 			"id": "did:example:123456789abcdefghi#keys-1",
	// 			"publicKeyMultibase": "zH3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV",
	// 			"type": "Ed25519VerificationKey2020"
	// 		}
	// 	],
	// 	"id": "did:example:123456789abcdefghi"
	// }
}

func ExampleDocument_publicKey() {
	doc, _ := ParseDocument(example13)
	fmt.Println(doc)
	// Output:
	// {
	// 	"@context": [
	// 		"https://www.w3.org/ns/did/v1",
	// 		"https://w3id.org/security/suites/jws-2020/v1",
	// 		"https://w3id.org/security/suites/ed25519-2020/v1"
	// 	],
	// 	"id": "did:example:123456789abcdefghi",
	// 	"verificationMethod": [
	// 		{
	// 			"controller": "did:example:123",
	// 			"id": "did:example:123#_Qq0UL2Fq651Q0Fjd6TvnYE-faHiOpRlPVQcY_-tA4A",
	// 			"publicKeyJwk": {
	// 				"crv": "Ed25519",
	// 				"kid": "_Qq0UL2Fq651Q0Fjd6TvnYE-faHiOpRlPVQcY_-tA4A",
	// 				"kty": "OKP",
	// 				"x": "VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ"
	// 			},
	// 			"type": "JsonWebKey2020"
	// 		},
	// 		{
	// 			"controller": "did:example:pqrstuvwxyz0987654321",
	// 			"id": "did:example:123456789abcdefghi#keys-1",
	// 			"publicKeyMultibase": "zH3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV",
	// 			"type": "Ed25519VerificationKey2020"
	// 		}
	// 	]
	// }
}

func ExampleDocument_references() {
	doc, _ := ParseDocument(example14)
	fmt.Println(doc)
	// Output:
	// {
	// 	"authentication": [
	// 		"did:example:123456789abcdefghi#keys-1",
	// 		{
	// 			"controller": "did:example:123456789abcdefghi",
	// 			"id": "did:example:123456789abcdefghi#keys-2",
	// 			"publicKeyMultibase": "zH3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV",
	// 			"type": "Ed25519VerificationKey2020"
	// 		}
	// 	],
	// 	"id": "did:example:123456789abcdefghi"
	// }
}

func ExampleDocument_relativeURL() {
	doc, _ := ParseDocument(example9)
	fmt.Println(doc)
	// Output:
	// {
	// 	"@context": [
	// 		"https://www.w3.org/ns/did/v1",
	// 		"https://w3id.org/security/suites/ed25519-2020/v1"
	// 	],
	// 	"authentication": [
	// 		"#key-1"
	// 	],
	// 	"id": "did:example:123456789abcdefghi",
	// 	"verificationMethod": [
	// 		{
	// 			"controller": "did:example:123456789abcdefghi",
	// 			"id": "did:example:123456789abcdefghi#key-1",
	// 			"publicKeyMultibase": "zH3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV",
	// 			"type": "Ed25519VerificationKey2020"
	// 		}
	// 	]
	// }
}

func ExampleDocument_service() {
	doc, _ := ParseDocument(example20)
	fmt.Println(doc)
	// Output:
	// {
	// 	"id": "did:example:123456789abcdefghi",
	// 	"service": [
	// 		{
	// 			"id": "did:example:123#linked-domain",
	// 			"serviceEndpoint": "https://bar.example.com",
	// 			"type": "LinkedDomains"
	// 		}
	// 	]
	// }
}
