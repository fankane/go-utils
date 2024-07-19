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

func TestSubOfUTF8(t *testing.T) {
	type args struct {
		s     string
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				s:     "我爱中国good",
				start: -1,
				end:   3,
			},
			want: "",
		},
		{
			args: args{
				s:     "我爱中国good",
				start: 1,
				end:   9,
			},
			want: "",
		},
		{
			args: args{
				s:     "我爱中国good",
				start: 5,
				end:   3,
			},
			want: "",
		},
		{
			args: args{
				s:     "我爱中国good",
				start: 0,
				end:   3,
			},
			want: "我爱中",
		},
		{
			args: args{
				s:     "我爱中国good",
				start: 0,
				end:   0,
			},
			want: "",
		},
		{
			args: args{
				s:     "我爱中国good",
				start: 0,
				end:   1,
			},
			want: "我",
		},
		{
			args: args{
				s:     "我爱中国good",
				start: 7,
				end:   8,
			},
			want: "d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubOfUTF8(tt.args.s, tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("SubOfUTF8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsChinese(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no chinese",
			args: args{s: "hello"},
			want: false,
		},
		{
			name: "has chinese",
			args: args{s: "hello, 你好"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsChinese(tt.args.s); got != tt.want {
				t.Errorf("ContainsChinese() = %v, want %v", got, tt.want)
			}
		})
	}
}
