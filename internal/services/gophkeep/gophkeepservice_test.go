package gophkeep

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
		host   string
	}
	logger := zaptest.NewLogger(t).Sugar()
	tests := []struct {
		name string
		args args
		want *GophKeepService
	}{
		{
			name: "valid",
			args: args{
				logger: logger,
				store:  nil,
				host:   "",
			},
			want: &GophKeepService{
				logger:  logger,
				storage: nil,
				host:    "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.logger, tt.args.store, tt.args.host)
			got.logger = logger
			tt.want.grpcServer = got.grpcServer
			tt.want.logger = logger

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
