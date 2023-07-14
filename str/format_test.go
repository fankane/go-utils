package str

import (
	"fmt"
	"testing"
)

func TestFormatFileSize(t *testing.T) {
	fmt.Println(FormatFileSize(13293363585771110499))
	fmt.Println(FormatFileSize(234))
	fmt.Println(FormatFileSize(2340))
	fmt.Println(FormatFileSize(23423534))
	fmt.Println(FormatFileSize(2342353444))
	fmt.Println(FormatFileSize(234235344499))
}

func TestGetStrIndex(t *testing.T) {
	fmt.Println(len("hi中国"))
	fmt.Println(LenOfUTF8("hi中国"))
	res := GetStrIndex("hi 你好", 1)
	fmt.Println(res)

	sli := SliceOfChar("hello 你好")
	fmt.Println(sli)
}
