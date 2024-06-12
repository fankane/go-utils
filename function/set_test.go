package function

import (
	"fmt"
	"testing"
)

func TestNewSet(t *testing.T) {
	s := NewSet()
	s.Add("1")
	s.Add("1")
	s.Add("2")
	s.Add("3")
	fmt.Println(s.String())
	fmt.Println(s.Size())
	fmt.Println(s.Contains("2"))
	fmt.Println(s.Contains("a"))
	s.Remove("1")
	fmt.Println(s)
}
