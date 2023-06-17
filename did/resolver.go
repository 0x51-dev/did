package did

type ConditionWeightedThreshold struct {
	Condition VerificationMethod
	Weight    int
}

type DefaultResolver struct {
	Registry Registry
}

func (r DefaultResolver) Resolve(didURL string, options ResolutionOptions) ResolutionResult {
	u, err := ParseDID(didURL)
	if err != nil {
		return ResolutionResult{
			Metadata: Metadata{
				Error: InvalidDIDError,
			},
		}
	}
	resolver, ok := r.Registry[u.Method]
	if !ok {
		return ResolutionResult{
			Metadata: Metadata{
				Error: UnsupportedDidMethodError,
			},
		}
	}
	return resolver(didURL, *u, r, options)
}

type Error string

const (
	// InvalidDIDError - the supplied DID to the DID resolution function does not conform to valid syntax.
	InvalidDIDError Error = "invalidDid"
	// NotFoundError - The DID resolver was unable to find the DID document resulting from this resolution request.
	NotFoundError Error = "notFound"
	// RepresentationNotSupportedError - This error code is returned if the representation requested via the accept
	// input metadata property is not supported by the DID method and/or DID resolver implementation.
	RepresentationNotSupportedError Error = "representationNotSupported"
	// UnsupportedDidMethodError - This error code is returned if the DID method specified in the DID is not supported.
	UnsupportedDidMethodError Error = "unsupportedDidMethod"
)

// Metadata are metadata about the DID resolution process.
// DOCS: https://www.w3.org/TR/did-core/#did-resolution-metadata
type Metadata struct {
	// The Media Type of the returned didDocument.
	ContentType string
	// The error code from the resolution process.
	// This property is REQUIRED when there is an error in the resolution process
	Error Error
}

type Registry map[string]Resolver

// ResolutionOptions are options for resolving a DID.
// DOCS: https://www.w3.org/TR/did-core/#did-resolution-options
type ResolutionOptions struct {
	// The Media Type of the caller's preferred representation of the DID document.
	Accept string
}

type ResolutionResult struct {
	Metadata Metadata  `json:"metadata"`
	Document *Document `json:"didDocument"`
}

type Resolvable interface {
	// Resolve returns the DID document in its abstract form (a map).
	Resolve(didURL string, options ResolutionOptions) ResolutionResult
}

type Resolver func(didURL string, did DID, resolver Resolvable, options ResolutionOptions) ResolutionResult
