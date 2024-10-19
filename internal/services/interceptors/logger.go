package interceptors

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func LoggerInterceptor(
	logger *zap.SugaredLogger,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		now := time.Now()
		var resp interface{}
		defer func() {
			since := time.Since(now)
			logger.Infoln("Method", info.FullMethod,
				"Request", req,
				"Duration", since)
		}()
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
