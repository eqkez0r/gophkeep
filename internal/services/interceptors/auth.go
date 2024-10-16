package interceptors

import (
	"context"
	goph_keeper_v1 "github.com/eqkez0r/gophkeep-grpc-api/pkg"
	"github.com/eqkez0r/gophkeep/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type UserStorageProvider interface {
	IsUserExist(context.Context, string) (bool, error)
}

func Auth(
	logger *zap.SugaredLogger,
	storage UserStorageProvider,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		namedLogger := logger.Named("Authentication")
		namedLogger.Debug(info.FullMethod)
		token := getToken(req)
		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "Token is empty")
		}
		login, t, err := jwt.JWTPayload(token)
		if err != nil {
			namedLogger.Error(info.FullMethod + ": " + err.Error())
			return nil, err
		}
		ok, err := storage.IsUserExist(ctx, login)
		if err != nil {
			namedLogger.Error(info.FullMethod + ": " + err.Error())
			return nil, err
		}
		if !ok {
			namedLogger.Error(info.FullMethod + ": user is not exist")
			return nil, status.Errorf(codes.Aborted, "user is not exist")
		}
		tOk := t.After(time.Now())
		if tOk {
			namedLogger.Info(info.FullMethod + ": token expired")
			return nil, status.Errorf(codes.Aborted, "token expired")
		}
		ctx = context.WithValue(ctx, "login", login)
		return handler(ctx, req)
	}
}

func getToken(req any) string {
	switch v := req.(type) {
	case *goph_keeper_v1.SendCardRequest:
		return v.Token
	case *goph_keeper_v1.GetCardRequest:
		return v.Token
	case *goph_keeper_v1.SendCredentialsRequest:
		return v.Token
	case *goph_keeper_v1.GetCredentialsRequest:
		return v.Token
	case *goph_keeper_v1.SendTextRequest:
		return v.Token
	case *goph_keeper_v1.GetTextRequest:
		return v.Token
	case *goph_keeper_v1.SyncRequest:
		return v.Token
	}
	return ""
}
