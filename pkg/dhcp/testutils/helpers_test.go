package testutils

import "testing"

func TestBool2String(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test for True",
			args: args{
				b: true,
			},
			want: "True",
		},
		{
			name: "Test for False",
			args: args{
				b: false,
			},
			want: "False",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bool2String(tt.args.b); got != tt.want {
				t.Errorf("Bool2String() = %v, want %v", got, tt.want)
			}
		})
	}
}
