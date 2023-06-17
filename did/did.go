package did

import (
	"encoding/json"
	"fmt"
	"strings"
)

func isHexadecimal(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func isIDCharacter(c byte) bool {
	return isLowerAlphaNumeric(c) || 'A' <= c && c <= 'Z' || c == '-' || c == '.' || c == '_'
}

func isLowerAlphaNumeric(c byte) bool {
	return 'a' <= c && c <= 'z' || '0' <= c && c <= '9'
}

func isParameterCharacter(c byte) bool {
	return isIDCharacter(c) || c == ':' || c == '%'
}

type DID struct {
	Method     string
	MethodIDs  []string
	Parameters []Parameter
	Path       string
	Query      string
	Fragment   string
}

func ParseDID(didURL string) (*DID, error) {
	if didURL == "" {
		return nil, fmt.Errorf("invalid DID: %s", didURL)
	}

	if !strings.HasPrefix(didURL, "did:") {
		return nil, fmt.Errorf("invalid DID: %s", didURL)
	}

	var i = 4
	for ; i < len(didURL); i++ {
		if !isLowerAlphaNumeric(didURL[i]) {
			break
		}
	}
	var method = didURL[4:i]

	if len(didURL) <= i || didURL[i] != ':' {
		return nil, fmt.Errorf("invalid DID: %s", didURL)
	}
	i++

	var methodIDs []string
	var j = i
	for ; i < len(didURL); i++ {
		if didURL[i] == ':' {
			methodIDs = append(methodIDs, didURL[j:i])
			j = i + 1
			continue
		}
		if didURL[i] == '%' {
			// Percent-encoded character.
			if len(didURL) <= i+2 || !isHexadecimal(didURL[i+1]) || !isHexadecimal(didURL[i+2]) {
				return nil, fmt.Errorf("invalid DID: %s", didURL)
			}
			i += 2
		}
		if !isIDCharacter(didURL[i]) {
			break
		}
	}
	if i == j {
		return nil, fmt.Errorf("invalid DID: %s", didURL)
	}
	methodIDs = append(methodIDs, didURL[j:i])

	var parameters []Parameter
	for i < len(didURL) {
		if didURL[i] != ';' {
			break
		}
		i++
		for j = i; i < len(didURL); i++ {
			if !isParameterCharacter(didURL[i]) {
				break
			}
		}
		if i == j {
			return nil, fmt.Errorf("invalid DID: %s", didURL)
		}
		var k = didURL[j:i]
		if len(didURL) <= i || didURL[i] != '=' {
			return nil, fmt.Errorf("invalid DID: %s", didURL)
		}
		i++
		for j = i; i < len(didURL); i++ {
			if !isParameterCharacter(didURL[i]) {
				break
			}
		}
		var v = didURL[j:i]
		parameters = append(parameters, Parameter{
			Key:   k,
			Value: v,
		})
	}

	var path string
	if i < len(didURL) && didURL[i] == '/' {
		j = i
		for ; i < len(didURL); i++ {
			if c := didURL[i]; c == '?' || c == '#' {
				break
			}
		}
		path = didURL[j:i]
	}

	var query string
	if i < len(didURL) && didURL[i] == '?' {
		j = i + 1
		for ; i < len(didURL); i++ {
			if c := didURL[i]; c == '#' {
				break
			}
		}
		query = didURL[j:i]
	}

	var fragment string
	if i < len(didURL) && didURL[i] == '#' {
		fragment = didURL[i+1:]
		i = len(didURL)
	}

	if i != len(didURL) {
		return nil, fmt.Errorf("invalid DID: %s", didURL)
	}

	return &DID{
		Method:     method,
		MethodIDs:  methodIDs,
		Parameters: parameters,
		Path:       path,
		Query:      query,
		Fragment:   fragment,
	}, nil
}

func (d *DID) DID() string {
	var sb strings.Builder
	sb.WriteString("did:")
	sb.WriteString(d.Method)
	sb.WriteString(":")
	for i, id := range d.MethodIDs {
		if i != 0 {
			sb.WriteString(":")
		}
		sb.WriteString(id)
	}
	return sb.String()
}

func (d DID) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *DID) String() string {
	var sb strings.Builder
	sb.WriteString("did:")
	sb.WriteString(d.Method)
	sb.WriteString(":")
	for i, id := range d.MethodIDs {
		if i != 0 {
			sb.WriteString(":")
		}
		sb.WriteString(id)
	}
	for _, p := range d.Parameters {
		sb.WriteString(";")
		sb.WriteString(p.Key)
		sb.WriteString("=")
		sb.WriteString(p.Value)
	}
	if d.Path != "" {
		sb.WriteString(d.Path)
	}
	if d.Query != "" {
		sb.WriteString("?")
		sb.WriteString(d.Query)
	}
	if d.Fragment != "" {
		sb.WriteString("#")
		sb.WriteString(d.Fragment)
	}
	return sb.String()
}

func (d *DID) UnmarshalJSON(raw []byte) error {
	didURL := string(raw)
	if !strings.HasPrefix(didURL, `"`) || !strings.HasSuffix(didURL, `"`) {
		return fmt.Errorf("invalid DID: %s", didURL)
	}
	didURL = didURL[1 : len(didURL)-1]
	u, err := ParseDID(didURL)
	if err != nil {
		return err
	}
	*d = *u
	return nil
}

// Parameter is a DID parameter.
// DOCS: https://www.w3.org/TR/did-core/#did-parameters
type Parameter struct {
	Key   string
	Value string
}
