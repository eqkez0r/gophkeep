package app

import (
	"github.com/eqkez0r/gophkeep/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		logger *zap.SugaredLogger
		store  storage.Storage
		path   string
	}
	logger := zaptest.NewLogger(t).Sugar()
	tests := []struct {
		name    string
		args    args
		want    *App
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				logger: logger,
				store:  nil,
				path:   "./test_config.yaml",
			},
			want: &App{
				logger: logger,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.logger, tt.args.store, tt.args.path)
			tt.want.services = got.services
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
