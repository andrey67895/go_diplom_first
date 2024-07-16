package model

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWithdrawnBalanceModelDecodePositive(t *testing.T) {
	type want struct {
		code      int
		errorText *string
	}
	type args struct {
		order       string
		processedAT *time.Time
		withdrawn   float64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				order:     "4207452",
				withdrawn: 100,
			},
			want: want{
				code:      200,
				errorText: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tModel := WithdrawnBalanceModel{Order: &test.args.order, Withdrawn: &test.args.withdrawn}
			body := io.NopCloser(bytes.NewReader(tModel.Marshal()))
			got, err := WithdrawnBalanceModelDecode(w, body)
			res := w.Result()
			assert.Equal(t, nil, err)
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, *tModel.Order, *got.Order)
			assert.Equal(t, *tModel.Withdrawn, *got.Withdrawn)
		})
	}
}

func TestWithdrawnBalanceModelDecodeNegative(t *testing.T) {
	type want struct {
		code      int
		errorText *string
	}
	type args struct {
		order       string
		processedAT *time.Time
		withdrawn   float64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				order:     "1",
				withdrawn: 100,
			},
			want: want{
				code:      422,
				errorText: helpers.GetAdrressString("Неверный формат номера заказа!"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tModel := WithdrawnBalanceModel{Order: &test.args.order, Withdrawn: &test.args.withdrawn}
			body := io.NopCloser(bytes.NewReader(tModel.Marshal()))
			got, err := WithdrawnBalanceModelDecode(w, body)
			res := w.Result()
			assert.Equal(t, nil, err)
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, *tModel.Order, *got.Order)
			assert.Equal(t, *tModel.Withdrawn, *got.Withdrawn)
		})
	}
}

func TestWithdrawnBalanceModel_Marshal(t *testing.T) {
	type fields struct {
		Login       string
		Order       string
		ProcessedAT time.Time
		Withdrawn   float64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "positive test #1",
			fields: fields{
				Login:       "TEST",
				Order:       "",
				ProcessedAT: time.Now(),
				Withdrawn:   100,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tModel := WithdrawnBalanceModel{
				Login:       &test.fields.Login,
				Order:       &test.fields.Order,
				ProcessedAT: &test.fields.ProcessedAT,
				Withdrawn:   &test.fields.Withdrawn,
			}
			marshal, err := json.Marshal(tModel)
			assert.Equal(t, nil, err)
			assert.Equalf(t, marshal, tModel.Marshal(), "Marshal() is error")
		})
	}
}
