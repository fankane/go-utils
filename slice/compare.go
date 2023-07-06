package slice

// StrSliContentEqual 比较两个字符串切片内容是否相同，忽略顺序
func StrSliContentEqual(strA, strB []string) bool {
	if len(strA) != len(strB) {
		return false
	}
	if len(strA) == 0 {
		return true
	}
	mapA := StrSliToMapCnt(strA)
	for _, val := range strB {
		v, ok := mapA[val]
		if !ok {
			return false
		}
		if v == 1 { //只有一个
			delete(mapA, val)
		} else {
			mapA[val] = v - 1
		}
	}
	return true
}

// StrSliContains more 是否包含less里面的全部数据
func StrSliContains(more, less []string) bool {
	if len(more) < len(less) {
		return false
	}
	if len(more) == len(less) {
		return StrSliContentEqual(more, less)
	}
	mapMore := StrSliToMapCnt(more)
	for _, val := range less {
		v, ok := mapMore[val]
		if !ok {
			return false
		}
		if v == 1 { //只有一个
			delete(mapMore, val)
		} else {
			mapMore[val] = v - 1
		}
	}
	return true
}
