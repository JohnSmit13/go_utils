/*Generated Test_main
Generated Test_unpack
Generated Test_isNumber*/
package main

import "testing"

/*func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}*/

func Test_unpack(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
        //false
        {"",args{`\1`},`1`,false},
        {"",args{`\0`},`0`,false},
        {"",args{`\13`},`111`,false},
        {"",args{`\21`},`2`,false},
        {"",args{`\\`},`\`,false},
        {"",args{`asd`},`asd`,false},
        {"",args{`\a`},`a`,false},
        {"",args{`\5a`},`5a`,false},
        {"",args{`\d5`},`ddddd`,false},
        {"",args{`\\\\`},`\\`,false},
        {"",args{`a0`},``,false},
        {"",args{`a\0`},`a0`,false},
        //true
        {"",args{`\\\`},``,true},
        {"",args{`0`},``,true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unpack(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unpack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isNumber(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
        {"",args{'0'},true},
        {"",args{'1'},true},
        {"",args{'2'},true},
        {"",args{'3'},true},
        {"",args{'4'},true},
        {"",args{'5'},true},
        {"",args{'6'},true},
        {"",args{'7'},true},
        {"",args{'8'},true},
        {"",args{'9'},true},
        {"",args{'a'},false},
        {"",args{'r'},false},
        {"",args{'y'},false},
        {"",args{','},false},
        {"",args{'\''},false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNumber(tt.args.r); got != tt.want {
				t.Errorf("isNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
