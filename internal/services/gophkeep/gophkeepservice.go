package gophkeep

import (
	"context"
	"errors"
	pb "github.com/eqkez0r/gophkeep-grpc-api/pkg"
	"github.com/eqkez0r/gophkeep/internal/services/interceptors"
	se "github.com/eqkez0r/gophkeep/internal/storage/storageerrors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
)

type GophKeepService struct {
	logger     *zap.SugaredLogger
	storage    StorageKeeperProvider
	grpcServer *grpc.Server
	host       string

	pb.UnimplementedGophKeeperServer
}

type StorageKeeperProvider interface {
	IsUserExist(context.Context, string) (bool, error)

	NewCredentials(context.Context, string, string, string, string) error
	GetCredentials(context.Context, string, string) (string, string, error)
	CredentialList(context.Context, string) ([]string, error)

	NewText(context.Context, string, string, string) error
	GetText(context.Context, string, string) (string, error)
	TextList(context.Context, string) ([]string, error)

	NewCard(context.Context, string, string, string, string, string, int32) error
	GetCard(context.Context, string, string) (string, string, string, int32, error)
	CardList(context.Context, string) ([]string, error)
}

func New(
	logger *zap.SugaredLogger,
	store StorageKeeperProvider,
	host string,
) *GophKeepService {
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.LoggerInterceptor(logger),
		interceptors.Auth(logger, store),
	))
	gs := &GophKeepService{
		logger:     logger.Named("GophKeep_Service"),
		storage:    store,
		grpcServer: grpcServer,
		host:       host,
	}
	pb.RegisterGophKeeperServer(grpcServer, gs)
	return gs
}

func (gs *GophKeepService) Run(ctx context.Context, wg *sync.WaitGroup) {
	gs.logger.Info("Starting GophKeep Service " + gs.host)

	listener, err := net.Listen("tcp", gs.host)
	if err != nil {
		gs.logger.Errorw("Failed to listen", "host", gs.host, "error", err)
		wg.Done()
		return
	}
	reflection.Register(gs.grpcServer)
	go func() {
		err = gs.grpcServer.Serve(listener)
		if err != nil {
			gs.logger.Errorw("Failed to serve", "host", gs.host, "error", err)
		}
	}()
	<-ctx.Done()
}

func (gs *GophKeepService) GracefulShutdown() {
	gs.grpcServer.GracefulStop()
}

func (gs *GophKeepService) SendCredentials(ctx context.Context, req *pb.SendCredentialsRequest) (*pb.SendCredentialsResponse, error) {
	login := ctx.Value("login").(string)
	err := gs.storage.NewCredentials(ctx, login, req.Credential.CredentialsName, req.Credential.Login, req.Credential.Password)
	if err != nil {
		gs.logger.Errorw("Failed to create credentials", "error", err)
		var state pb.SendDataState
		switch {
		case errors.Is(err, se.ErrCredentialsIsExist):
			state = pb.SendDataState_SEND_DATA_IS_EXIST
		default:
			state = pb.SendDataState_SEND_DATA_ERROR
		}
		return &pb.SendCredentialsResponse{
			State: state,
		}, err
	}
	return &pb.SendCredentialsResponse{
		State: pb.SendDataState_SEND_DATA_SUCCESS,
	}, nil
}

func (gs *GophKeepService) GetCredentials(ctx context.Context, req *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {
	login := ctx.Value("login").(string)
	log, pass, err := gs.storage.GetCredentials(ctx, login, req.CredentialName)
	if err != nil {
		gs.logger.Errorw("Failed to get credentials", "error", err)
		var state pb.GetDataState
		switch {
		case errors.Is(err, se.ErrCredentialsNotFound):
			state = pb.GetDataState_GET_DATA_IS_NOT_EXIST
		default:
			state = pb.GetDataState_GET_DATA_ERROR
		}
		return &pb.GetCredentialsResponse{
			State: state,
		}, err
	}
	return &pb.GetCredentialsResponse{
		Credentials: &pb.Credentials{
			CredentialsName: req.CredentialName,
			Login:           log,
			Password:        pass,
		},
		State: pb.GetDataState_GET_DATA_SUCCESS,
	}, nil
}

func (gs *GophKeepService) SendText(ctx context.Context, req *pb.SendTextRequest) (*pb.SendTextResponse, error) {
	login := ctx.Value("login").(string)
	err := gs.storage.NewText(ctx, login, req.Text.TextName, req.Text.Text)
	if err != nil {
		gs.logger.Errorw("Failed to create text", "error", err)
		var state pb.SendDataState
		switch {
		case errors.Is(err, se.ErrTextIsExist):
			state = pb.SendDataState_SEND_DATA_IS_EXIST
		default:
			state = pb.SendDataState_SEND_DATA_ERROR
		}
		return &pb.SendTextResponse{
			State: state,
		}, nil
	}

	return &pb.SendTextResponse{
		State: pb.SendDataState_SEND_DATA_SUCCESS,
	}, nil
}

func (gs *GophKeepService) GetText(ctx context.Context, req *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	login := ctx.Value("login").(string)
	text, err := gs.storage.GetText(ctx, login, req.TextName)
	if err != nil {
		gs.logger.Errorw("Failed to get text", "error", err)
		var state pb.GetDataState
		switch {
		case errors.Is(err, se.ErrTextNotFound):
			state = pb.GetDataState_GET_DATA_IS_NOT_EXIST
		default:
			state = pb.GetDataState_GET_DATA_ERROR
		}
		return &pb.GetTextResponse{
			State: state,
		}, err
	}
	return &pb.GetTextResponse{
		State: pb.GetDataState_GET_DATA_SUCCESS,
		Text: &pb.Text{
			Text:     text,
			TextName: req.TextName,
		},
	}, nil
}

func (gs *GophKeepService) SendCard(ctx context.Context, req *pb.SendCardRequest) (*pb.SendCardResponse, error) {
	login := ctx.Value("login").(string)
	err := gs.storage.NewCard(ctx, login, req.Card.CardName, req.Card.CardNumber, req.Card.CardHolderName, req.Card.ExpirationDate, req.Card.Cvv)
	if err != nil {
		gs.logger.Errorw("Failed to create card", "error", err)
		var state pb.SendDataState
		switch {
		case errors.Is(err, se.ErrCardIsExist):
			state = pb.SendDataState_SEND_DATA_IS_EXIST
		default:
			state = pb.SendDataState_SEND_DATA_ERROR
		}
		return &pb.SendCardResponse{
			State: state,
		}, err
	}
	return &pb.SendCardResponse{
		State: pb.SendDataState_SEND_DATA_SUCCESS,
	}, nil
}

func (gs *GophKeepService) GetCard(ctx context.Context, req *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	login := ctx.Value("login").(string)
	cardNumber, cardHolderName, expirationDate, cvv, err := gs.storage.GetCard(ctx, login, req.CardName)
	if err != nil {
		gs.logger.Errorw("Failed to get card", "error", err)
		var state pb.GetDataState
		switch {
		case errors.Is(err, se.ErrCardNotFound):
			state = pb.GetDataState_GET_DATA_IS_NOT_EXIST
		default:
			state = pb.GetDataState_GET_DATA_ERROR
		}
		return &pb.GetCardResponse{
			State: state,
		}, nil
	}
	return &pb.GetCardResponse{
		State: pb.GetDataState_GET_DATA_SUCCESS,
		Card: &pb.Card{
			CardNumber:     cardNumber,
			CardHolderName: cardHolderName,
			ExpirationDate: expirationDate,
			Cvv:            cvv,
		},
	}, nil
}

func (gs *GophKeepService) Sync(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
	login := ctx.Value("login").(string)
	texts, err := gs.storage.TextList(ctx, login)
	if err != nil {
		gs.logger.Errorw("Failed to get text list", "error", err)
		return &pb.SyncResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	credentials, err := gs.storage.CredentialList(ctx, login)
	if err != nil {
		gs.logger.Errorw("Failed to get credential list", "error", err)
		return &pb.SyncResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	cards, err := gs.storage.CardList(ctx, login)
	if err != nil {
		gs.logger.Errorw("Failed to get card list", "error", err)
		return &pb.SyncResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	return &pb.SyncResponse{
		State:       pb.GetDataState_GET_DATA_SUCCESS,
		Texts:       texts,
		Credentials: credentials,
		Cards:       cards,
	}, nil
}
