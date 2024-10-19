package cipher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptData(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test_1",
			args: args{
				data: []byte("hello world"),
			},
			want:    []byte("hello world"),
			wantErr: false,
		},
		{
			name: "test_2",
			args: args{
				data: []byte("praktikum praktikum praktikum praktikum praktikum praktikum praktikum praktikum "),
			},
			want:    []byte("praktikum praktikum praktikum praktikum praktikum praktikum praktikum praktikum "),
			wantErr: false,
		},
		{
			name: "test_3",
			args: args{
				data: []byte("yandex praktikum yandex praktikum yandex praktikum yandex praktikum yandex praktikum yandex praktikum "),
			},
			want:    []byte("yandex praktikum yandex praktikum yandex praktikum yandex praktikum yandex praktikum yandex praktikum "),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			encr, err := DecryptData(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, encr) {
				t.Errorf("EncryptData() got = %v, want %v", encr, got)
				return
			}
		})
	}
}
