package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeHash(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{
				value: "TEST_HASH",
			},
			want: "7265c6cf1d6ad0ae64424fdd35ee7a88b8a1d4e237eeb6424e6b4c68a14ed441",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash := EncodeHash(test.args.value)
			assert.Equal(t, test.want, hash)
		})
	}
}
