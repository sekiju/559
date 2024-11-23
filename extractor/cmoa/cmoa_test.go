package cmoa

import (
	"os"
	"testing"
)

func TestProvider(t *testing.T) {
	session := os.Getenv("CMOA_COOKIE")
	if session == "" {
		t.Skip("Authorization token required for tests")
	}

	//provider := New(session)

}
