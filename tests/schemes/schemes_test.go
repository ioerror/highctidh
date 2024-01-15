package scheme

import (
	"testing"

	"codeberg.org/vula/highctidh/schemes"
)

func TestSchemes(t *testing.T) {
	t.Log("ALL Scheme names:")
	mySchemes := schemes.All()
	for _, s := range mySchemes {
		t.Logf("Scheme name: %s", s.Name())
	}
}
