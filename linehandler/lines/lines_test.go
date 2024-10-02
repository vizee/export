package lines

import (
	"fmt"
	"testing"
)

func Test_nextString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 bool
	}{
		{args: args{s: " "}, want: "", want1: "", want2: false},
		{args: args{s: "a b c"}, want: "a", want1: " b c", want2: true},
		{args: args{s: " a b"}, want: "a", want1: " b", want2: true},
		{args: args{s: "\"a\\tb\"c"}, want: "a\tb", want1: "c", want2: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := nextString(tt.args.s)
			if got != tt.want {
				t.Errorf("nextString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("nextString() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("nextString() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_parseCommand(t *testing.T) {
	p := LineProcessor{
		Cmds: map[string]*Command{
			"LIST": {
				Args: 0,
				Handle: func(args []string) ([]byte, error) {
					fmt.Println("LIST", args)
					return nil, nil
				},
			},
			"GET": {
				Args: 1,
				Handle: func(args []string) ([]byte, error) {
					fmt.Println("GET", args)
					return nil, nil
				},
			},
			"SET": {
				Args: 2,
				Handle: func(args []string) ([]byte, error) {
					fmt.Println("SET", args)
					return nil, nil
				},
			},
		},
	}
	p.Default = p.Cmds["GET"]
	_, err := p.ExecuteLine("LIST")
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.ExecuteLine("GET A")
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.ExecuteLine("SET A 1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.ExecuteLine("DEFAULT")
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.ExecuteLine("LIST 1")
	if err != ErrInvalidArgs {
		t.Fatal(err)
	}
	_, err = p.ExecuteLine("GET")
	if err != ErrInvalidArgs {
		t.Fatal(err)
	}
	_, err = p.ExecuteLine("GET ")
	if err != ErrInvalidLine {
		t.Fatal(err)
	}
}
