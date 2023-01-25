/*Generated Test_main
Generated Test_countColumns
Generated Test_sorter_CmpString
Generated Test_sorter_CmpInt
Generated Test_sorter_Sort
Generated Test_toInt*/
package main

import (
	"reflect"
	"testing"
)

/*func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}*/

func Test_countColumns(t *testing.T) {
	type args struct {
		s string
		c int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
        {"1",args{"abc 123",2} ,4, true},
        {"2",args{"abcc 123",2},5, true},
        {"3",args{"abc  123",2},5, true},
        {"4",args{"abc  123",1},0, true},
        {"5",args{"abc 123",3},0, false},
        {"6",args{"abc 123 5",2},4,true},
        {"6",args{"фыв 123 5",3},11,true},
        {"6",args{"фыв пра 5",2},7,true},
        {"6",args{"фыв  пра 5",2},8,true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := countColumns(tt.args.s, tt.args.c)
			if got != tt.want {
				t.Errorf("countColumns() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("countColumns() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_sorter_CmpString(t *testing.T) {
	type fields struct {
		num  bool
		uniq bool
		rev  bool
		col  int
	}
	type args struct {
		x string
		y string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sorter{
				num:  tt.fields.num,
				uniq: tt.fields.uniq,
				rev:  tt.fields.rev,
				col:  tt.fields.col,
			}
			if got := s.CmpString(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("sorter.CmpString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sorter_CmpInt(t *testing.T) {
	type fields struct {
		num  bool
		uniq bool
		rev  bool
		col  int
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sorter{
				num:  tt.fields.num,
				uniq: tt.fields.uniq,
				rev:  tt.fields.rev,
				col:  tt.fields.col,
			}
			if got := s.CmpInt(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("sorter.CmpInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sorter_Sort(t *testing.T) {
	type fields struct {
		num  bool
		uniq bool
		rev  bool
		col  int
	}
	type args struct {
		strs []columnString
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []columnString
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sorter{
				num:  tt.fields.num,
				uniq: tt.fields.uniq,
				rev:  tt.fields.rev,
				col:  tt.fields.col,
			}
			if got := s.Sort(tt.args.strs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sorter.Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toInt(t *testing.T) {
	type args struct {
		x string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		// TODO: Add test cases.
        {"1",args{"123"},123,true},
        {"2",args{"123d"},123,true},
        {"3",args{"1023"},1023,true},
        {"4",args{"12b3"},12,true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := toInt(tt.args.x)
			if got != tt.want {
				t.Errorf("toInt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("toInt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
