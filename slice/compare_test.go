package slice

import (
	"fmt"
	"testing"
)

func TestStrSliContentEqual(t *testing.T) {
	sA := []string{"a", "b", "c", "b", "e"}
	sB := []string{"a", "c", "b", "e"}
	fmt.Println(StrSliContains(sA, sB))
}
