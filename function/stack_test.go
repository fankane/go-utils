package function

import (
	"fmt"
	"testing"
)

func TestNewStack(t *testing.T) {
	s := NewStack()
	for i := 0; i < 10; i++ {
		s.Push(i + 1)
	}
	fmt.Println(s.Size())
	for i := 0; i < 10; i++ {
		v := s.Pop()
		fmt.Println(v, s.Size())
	}
}
