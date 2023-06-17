package did

import (
	"encoding/json"
	"fmt"
)

func unmarshalVerificationMethods(m map[string]json.RawMessage, methods []unmarshalVerificationMethod) error {
	for _, method := range methods {
		raw, ok := m[method.key]
		if !ok && !method.optional {
			return fmt.Errorf("missing field: %s", method.key)
		}
		if len(raw) == 0 {
			continue // Ignore optional empty fields.
		}
		var a []json.RawMessage
		if err := json.Unmarshal(raw, &a); err != nil {
			return err
		}

		var verificationMethods = make([]IVerificationMethod, len(a))
		for i, raw := range a {
			var vm VerificationMethod
			if err := json.Unmarshal(raw, &vm); err == nil {
				verificationMethods[i] = &vm
				continue
			}
			var relativeURL string
			if err := json.Unmarshal(raw, &relativeURL); err != nil {
				return err
			}
			verificationMethods[i] = &RelativeVerificationMethod{RelativeURL: relativeURL}
		}
		*method.target = verificationMethods
	}
	return nil
}

type unmarshalVerificationMethod struct {
	key      string
	target   *[]IVerificationMethod
	optional bool
}
