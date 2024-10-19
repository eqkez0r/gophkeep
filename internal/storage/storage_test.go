package storage

import (
	"github.com/eqkez0r/gophkeep/internal/storage/memory"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		configPath string
	}
	mem := memory.New()

	tests := []struct {
		name    string
		args    args
		want    Storage
		wantErr bool
	}{
		{
			name: "valid_memory",
			args: args{
				configPath: "./test_memory.yaml",
			},
			want:    mem,
			wantErr: false,
		},
		{
			name: "invalid_type",
			args: args{
				configPath: "./test_invalid.yaml",
			},
			want:    nil,
			wantErr: true,
		},
		//{
		//	name: "valid_postgresql",
		//	args: args{
		//		configPath: "./test_postgresql.yaml",
		//	},
		//	want:    post,
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
