package function

import "strings"

// Set 基于空结构体实现 set
type Set map[string]struct{}

func NewSet() Set {
	return make(Set)
}

// Add 添加元素到 set
func (s Set) Add(element string) {
	s[element] = struct{}{}
}

// Remove 从 set 中移除元素
func (s Set) Remove(element string) {
	delete(s, element)
}

// Contains 检查 set 中是否包含指定元素
func (s Set) Contains(element string) bool {
	_, exists := s[element]
	return exists
}

// Size 返回 set 大小
func (s Set) Size() int {
	return len(s)
}

// String implements fmt.Stringer
func (s Set) String() string {
	format := "("
	for element := range s {
		format += element + " "
	}
	format = strings.TrimRight(format, " ") + ")"
	return format
}
