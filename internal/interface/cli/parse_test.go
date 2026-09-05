package cli

import "testing"

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want options
	}{
		{
			name: "short template",
			args: []string{
				"dev",
				"-t",
				"example",
			},
			want: options{
				template: "example",
				target:   "dev",
			},
		},
		{
			name: "template before target",
			args: []string{
				"-t",
				"example",
				"dev",
			},
			want: options{
				template: "example",
				target:   "dev",
			},
		},
		{
			name: "long template",
			args: []string{
				"dev",
				"--template=example",
			},
			want: options{
				template: "example",
				target:   "dev",
			},
		},
		{
			name: "open",
			args: []string{
				"dev",
				"--template=example",
				"-o",
			},
			want: options{
				template: "example",
				target:   "dev",
				open:     true,
			},
		},
		{
			name: "list",
			args: []string{
				"--list",
			},
			want: options{
				list: true,
			},
		},
		{
			name: "help",
			args: []string{
				"-h",
			},
			want: options{
				help: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf(
					"expected %#v, got %#v",
					tt.want,
					got,
				)
			}
		})
	}
}

func TestParseArgsErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "unknown option",
			args: []string{
				"--wat",
			},
			want: "unknown option: --wat",
		},
		{
			name: "extra target",
			args: []string{
				"one",
				"two",
			},
			want: "unexpected argument: two",
		},
		{
			name: "missing short template value",
			args: []string{
				"-t",
			},
			want: "-t requires a value",
		},
		{
			name: "missing long template value",
			args: []string{
				"--template",
			},
			want: "--template requires a value",
		},
		{
			name: "empty template assignment",
			args: []string{
				"--template=",
			},
			want: "--template requires a value",
		},
		{
			name: "blank template",
			args: []string{
				"-t",
				"   ",
			},
			want: "-t requires a value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseArgs(tt.args)

			if err == nil {
				t.Fatal("expected an error")
			}

			if err.Error() != tt.want {
				t.Fatalf(
					"expected %q, got %q",
					tt.want,
					err,
				)
			}
		})
	}
}
