package data

import (
	"testing"
)

func TestGenerateReqID(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Logf(GenerateReqID())
	}
}
