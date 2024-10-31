package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLuhnValid(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "negative test #1",
			args: args{
				number: 1,
			},
			want: false,
		},
		{
			name: "positive test #1",
			args: args{
				number: 4207452,
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equalf(t, test.want, LuhnValid(test.args.number), "LuhnValid(%v) is not valid", test.args.number)
		})
	}
}
