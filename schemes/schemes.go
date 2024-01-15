package schemes

import (
	"errors"
	"strings"

	"codeberg.org/vula/highctidh/interfaces"
)

var allSchemes map[string]interfaces.CTIDH

func init() {
	allSchemes = make(map[string]interfaces.CTIDH)
}

// ByName returns the NIKE scheme by string name.
func ByName(name string) (interfaces.CTIDH, error) {
	v, ok := allSchemes[strings.ToLower(name)]
	if !ok {
		return nil, errors.New("ByName lookup failed.")
	}
	return v, nil
}

// All returns all CTIDH schemes supported.
func All() map[string]interfaces.CTIDH {
	a := allSchemes
	return a
}

func Register(name string, scheme interfaces.CTIDH) error {
	_, err := ByName(name)
	if err == nil {
		return errors.New("Register failed, already registered.")
	}
	allSchemes[name] = scheme
	return nil
}
