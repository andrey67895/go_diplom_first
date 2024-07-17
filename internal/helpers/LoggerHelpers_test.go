package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "positive test #1",
			want: "debug",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, Log().Level().String(), "Ошибка уровня логирования")
		})
	}
}
