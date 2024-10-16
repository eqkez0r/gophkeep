package auth

import (
	"context"
	pb "github.com/eqkez0r/gophkeep-grpc-api/pkg"
	"github.com/eqkez0r/gophkeep/internal/services/interceptors"
	"github.com/eqkez0r/gophkeep/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
)

type AuthService struct {
	logger     *zap.SugaredLogger
	store      UserStorageProvider
	grpcServer *grpc.Server
	host       string

	pb.UnimplementedAuthServiceServer
}

type UserStorageProvider interface {
	NewUser(context.Context, string, string) error
	ValidateUser(context.Context, string, string) error
}

func New(
	logger *zap.SugaredLogger,
	storage UserStorageProvider,
	host string,
) *AuthService {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggerInterceptor(logger),
		))
	as := &AuthService{
		logger:     logger.Named("Auth_Service"),
		store:      storage,
		host:       host,
		grpcServer: grpcServer,
	}
	pb.RegisterAuthServiceServer(grpcServer, as)
	return as
}

func (as *AuthService) Run(ctx context.Context, wg *sync.WaitGroup) {
	as.logger.Info("Starting Auth Service at " + as.host)

	listener, err := net.Listen("tcp", as.host)
	if err != nil {
		as.logger.Errorw("failed to listen", "host", as.host, "error", err)
		wg.Done()
		return
	}
	reflection.Register(as.grpcServer)
	go func() {
		err = as.grpcServer.Serve(listener)
		if err != nil {
			as.logger.Errorw("failed to serve", "host", as.host, "error", err)
		}
	}()
	<-ctx.Done()
}

func (as *AuthService) GracefulShutdown() {
	as.grpcServer.GracefulStop()
}

func (as *AuthService) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	err := as.store.ValidateUser(ctx, req.Login, req.Password)
	if err != nil {
		as.logger.Errorw("failed to validate user", "login", req.Login, "error", err)
		return &pb.AuthResponse{
			State: pb.AuthState_AUTH_FAILED,
		}, err
	}
	token, err := jwt.CreateJWT(req.Login)
	if err != nil {
		as.logger.Errorw("failed to generate token", "login", req.Login, "error", err)
		return &pb.AuthResponse{
			State: pb.AuthState_AUTH_ERROR,
			Token: nil,
		}, err
	}
	return &pb.AuthResponse{
		State: pb.AuthState_AUTH_SUCCESS,
		Token: &token,
	}, nil
}

func (as *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := as.store.NewUser(ctx, req.Login, req.Password)
	if err != nil {
		//TODO here need check error from POSTGRES OR make union error for storage
		as.logger.Errorw("failed to create user", "login", req.Login, "password", req.Password, "error", err)
		return &pb.RegisterResponse{
			State: pb.RegisterState_REGISTER_FAILED,
		}, err
	}
	token, err := jwt.CreateJWT(req.Login)
	if err != nil {
		as.logger.Errorw("failed to create JWT", "login", req.Login, "error", err)
		return &pb.RegisterResponse{
			State: pb.RegisterState_REGISTER_ERROR,
		}, err
	}
	return &pb.RegisterResponse{
		State: pb.RegisterState_REGISTER_SUCCESS,
		Token: &token,
	}, nil
}
