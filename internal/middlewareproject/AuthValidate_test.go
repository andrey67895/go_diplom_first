package middlewareproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_itAuth(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test #1",
			args: args{
				paths: []string{"/api/user/register", "/api/user/login"},
			},
			want: false,
		},
		{
			name: "positive test #2",
			args: args{
				paths: []string{"/api/user/orders", "/api/user/balance", "/api/user/balance/withdraw", "/api/user/withdrawals"},
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, path := range test.args.paths {
				assert.Equalf(t, test.want, itCheckAuth(path), "Ошибка в методе: %s", path)
			}
		})
	}
}
