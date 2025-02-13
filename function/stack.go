package function

import (
	"container/list"
	"sync"
)

type Stack struct {
	l    *sync.Mutex
	list *list.List
}

func NewStack() *Stack {
	return &Stack{l: &sync.Mutex{}, list: list.New()}
}

func (s *Stack) Push(val any) {
	s.l.Lock()
	defer s.l.Unlock()
	if s == nil || s.list == nil || val == nil {
		return
	}
	s.list.PushBack(val)
}

func (s *Stack) Pop() any {
	s.l.Lock()
	defer s.l.Unlock()
	e := s.list.Back()
	if e == nil {
		return nil
	}
	s.list.Remove(e)
	return e.Value
}

func (s *Stack) Peek() any {
	s.l.Lock()
	defer s.l.Unlock()
	e := s.list.Front()
	if e == nil {
		return nil
	}
	return e.Value
}

func (s *Stack) Size() int {
	if s == nil || s.list == nil {
		return -1
	}
	return s.list.Len()
}
