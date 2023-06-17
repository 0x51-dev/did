package did

import (
	"encoding/json"
	"github.com/0x51-dev/did/internal/jsonutils"
)

// Document is a DID document.
// DOCS: https://www.w3.org/TR/did-core/#data-model
type Document struct {
	Context            []string            `json:"@context,omitempty"`
	ID                 DID                 `json:"id"`
	AlsoKnownAs        []string            `json:"alsoKnownAs,omitempty"`
	Controller         []DID               `json:"controller,omitempty"`
	VerificationMethod VerificationMethods `json:"verificationMethod,omitempty"`
	// Verification Relationships
	Authentication       []IVerificationMethod `json:"authentication,omitempty"`
	AssertionMethod      []IVerificationMethod `json:"assertionMethod,omitempty"`
	KeyAgreement         []IVerificationMethod `json:"keyAgreement,omitempty"`
	CapabilityInvocation []IVerificationMethod `json:"capabilityInvocation,omitempty"`
	CapabilityDelegation []IVerificationMethod `json:"capabilityDelegation,omitempty"`
	// Services
	Service []Service `json:"service,omitempty"`
}

func ParseDocument(raw []byte) (*Document, error) {
	var doc Document
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (d *Document) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(*d)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	if err := jsonutils.FlattenStringOrSetOrMap(m, []jsonutils.FlattenField{
		{Name: "@context", Optional: true},
		{Name: "controller", Optional: true},
	}); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (d *Document) String() string {
	raw, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(raw)
}

func (d *Document) UnmarshalJSON(raw []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return err
	}
	if err := jsonutils.UnmarshalFields(m, []jsonutils.UnmarshalField{
		{Key: "id", Target: &d.ID},
		{Key: "alsoKnownAs", Target: &d.AlsoKnownAs, Optional: true},
		{Key: "verificationMethod", Target: &d.VerificationMethod, Optional: true},
		{Key: "service", Target: &d.Service, Optional: true},
	}); err != nil {
		return err
	}
	if err := jsonutils.UnmarshalTOrSets(m, []jsonutils.UnmarshalTOrSet[string]{
		{Key: "@context", Target: &d.Context, Optional: true},
	}); err != nil {
		return err
	}
	if err := jsonutils.UnmarshalTOrSets(m, []jsonutils.UnmarshalTOrSet[DID]{
		{Key: "controller", Target: &d.Controller, Optional: true},
	}); err != nil {
		return err
	}
	if err := unmarshalVerificationMethods(m, []unmarshalVerificationMethod{
		{key: "authentication", target: &d.Authentication, optional: true},
		{key: "assertionMethod", target: &d.AssertionMethod, optional: true},
		{key: "keyAgreement", target: &d.KeyAgreement, optional: true},
		{key: "capabilityInvocation", target: &d.CapabilityInvocation, optional: true},
		{key: "capabilityDelegation", target: &d.CapabilityDelegation, optional: true},
	}); err != nil {
		return err
	}
	return nil
}

type IVerificationMethod interface {
	Get(verificationMethods []VerificationMethod) *VerificationMethod
}

type RelativeVerificationMethod struct {
	RelativeURL string `json:"relativeURL"`
}

func (v *RelativeVerificationMethod) Get(verificationMethods []VerificationMethod) *VerificationMethod {
	for _, method := range verificationMethods {
		if method.ID == v.RelativeURL {
			return &method
		}
	}
	return nil
}

func (v *RelativeVerificationMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.RelativeURL)
}

type Service struct {
	ID              string          `json:"id"`
	Type            []string        `json:"type"`
	ServiceEndpoint ServiceEndpoint `json:"serviceEndpoint"`
}

func (s *Service) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(*s)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	if err := jsonutils.FlattenStringOrSetOrMap(m, []jsonutils.FlattenField{
		{Name: "type"},
		{Name: "serviceEndpoint"},
	}); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (s *Service) UnmarshalJSON(raw []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return err
	}
	if err := jsonutils.UnmarshalFields(m, []jsonutils.UnmarshalField{
		{Key: "id", Target: &s.ID},
	}); err != nil {
		return err
	}
	if err := jsonutils.UnmarshalTOrSets(m, []jsonutils.UnmarshalTOrSet[string]{
		{Key: "type", Target: &s.Type, Optional: true},
	}); err != nil {
		return err
	}
	if err := jsonutils.UnmarshalTOrSetOrMap(m["serviceEndpoint"], (*map[string]string)(&s.ServiceEndpoint)); err != nil {
		return err
	}
	return nil
}

type ServiceEndpoint map[string]string

type VerificationMethod struct {
	ID                 string            `json:"id"`
	Controller         string            `json:"controller"`
	Type               string            `json:"type"`
	PublicKeyJwk       map[string]string `json:"publicKeyJwk,omitempty"`
	PublicKeyMultibase string            `json:"publicKeyMultibase,omitempty"`
}

func (v *VerificationMethod) Get(_ []VerificationMethod) *VerificationMethod {
	return v
}

type VerificationMethods []VerificationMethod

// Get returns the method with the given ID.
func (v VerificationMethods) Get(id string) *VerificationMethod {
	for _, method := range v {
		if method.ID == id {
			return &method
		}
	}
	return nil
}
