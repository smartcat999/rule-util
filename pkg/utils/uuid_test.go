package utils

import (
	"fmt"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	u1 := GenerateUUID()
	u2 := GenerateUUID()

	fmt.Println(u1, u2)

	if u1 == u2 {
		t.Errorf("Expected not equals but got %s & %s", u1, u2)
	}
}
