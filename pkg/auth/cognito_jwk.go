package auth

import (
	"encoding/json"

	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

// CognitoJWK parses the JWK set embedded in the code and returns it
func CognitoJWK(key string) (*jose.JSONWebKeySet, error) {
	if key == "" {
		key = jwkKey
	}

	result := &jose.JSONWebKeySet{}
	if err := json.Unmarshal([]byte(key), result); err != nil {
		return nil, errors.Wrap(err, "unable to decode the JWK set")
	}

	injected := jose.JSONWebKey{}
	if err := json.Unmarshal([]byte(injectedKey), &injected); err != nil {
		return nil, errors.Wrap(err, "unable to decode injected JWK")
	}

	result.Keys = append(result.Keys, injected)

	return result, nil
}
