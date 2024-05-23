package assert

import "testing"

func TestIsArrayOrSlice(t *testing.T) {
	type args struct {
		i interface{}
	}
	array := []int{1, 2, 3}
	slice1 := make([]int, 0)
	slice1 = append(slice1, 1)
	intV := 1
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "not array or slice",
			args: args{i: intV},
			want: false,
		},
		{
			name: "is array",
			args: args{i: array},
			want: true,
		},
		{
			name: "is slice",
			args: args{i: slice1},
			want: true,
		},
		{
			name: "is slice ptr",
			args: args{i: &slice1},
			want: true,
		},
		{
			name: "is array ptr",
			args: args{i: &array},
			want: true,
		},
		{
			name: "is int ptr",
			args: args{i: &intV},
			want: false,
		},
		{
			name: "is nil",
			args: args{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsArrayOrSlice(tt.args.i); got != tt.want {
				t.Errorf("IsArrayOrSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	type args struct {
		i interface{}
	}
	intV, fv := 1, 1.2
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is int",
			args: args{i: intV},
			want: true,
		},
		{
			name: "is int prt",
			args: args{i: &intV},
			want: true,
		},
		{
			name: "is float ",
			args: args{i: fv},
			want: true,
		},
		{
			name: "is float prt",
			args: args{i: &fv},
			want: true,
		},
		{
			name: "not number",
			args: args{i: "a"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumeric(tt.args.i); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}
