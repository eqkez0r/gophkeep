package gophkeep

import (
	"context"
	"errors"
	pb "github.com/eqkez0r/gophkeep-grpc-api/pkg"
	"github.com/eqkez0r/gophkeep/internal/services/interceptors"
	"github.com/eqkez0r/gophkeep/internal/services/servicestype"
	"github.com/eqkez0r/gophkeep/internal/storage"
	se "github.com/eqkez0r/gophkeep/internal/storage/storageerrors"
	"github.com/eqkez0r/gophkeep/pkg/cipher"
	"github.com/eqkez0r/gophkeep/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
	"sync"
)

type GophKeepService struct {
	logger     *zap.SugaredLogger
	storage    storage.Storage
	grpcServer *grpc.Server
	host       string

	pb.UnimplementedGophKeeperServer
}

func New(
	logger *zap.SugaredLogger,
	store storage.Storage,
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

func (gs *GophKeepService) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {

	err := gs.storage.ValidateUser(ctx, req.Login, req.Password)
	if err != nil {
		gs.logger.Errorw("failed to validate user", "login", req.Login, "error", err)
		return &pb.AuthResponse{
			State: pb.AuthState_AUTH_FAILED,
		}, err
	}
	token, err := jwt.CreateJWT(req.Login)
	if err != nil {
		gs.logger.Errorw("failed to generate token", "login", req.Login, "error", err)
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

func (gs *GophKeepService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	pass, err := cipher.EncryptData([]byte(req.Password))
	if err != nil {
		gs.logger.Errorw("failed to encrypt password", "error", err)
		return &pb.RegisterResponse{
			State: pb.RegisterState_REGISTER_FAILED,
		}, err
	}
	err = gs.storage.NewUser(ctx, req.Login, string(pass))
	if err != nil {
		//TODO here need check error from POSTGRES OR make union error for storage
		gs.logger.Errorw("failed to create user", "login", req.Login, "error", err)
		return &pb.RegisterResponse{
			State: pb.RegisterState_REGISTER_FAILED,
		}, err
	}
	token, err := jwt.CreateJWT(req.Login)
	if err != nil {
		gs.logger.Errorw("failed to create JWT", "login", req.Login, "error", err)
		return &pb.RegisterResponse{
			State: pb.RegisterState_REGISTER_ERROR,
		}, err
	}
	return &pb.RegisterResponse{
		State: pb.RegisterState_REGISTER_SUCCESS,
		Token: &token,
	}, nil
}

func (gs *GophKeepService) SendCredentials(ctx context.Context, req *pb.SendCredentialsRequest) (*pb.SendCredentialsResponse, error) {
	login := ctx.Value("login").(servicestype.KeyType)
	encryptedLogin, err := cipher.EncryptData([]byte(req.Credential.Login))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt credentials", "error", err)
		return &pb.SendCredentialsResponse{
			State: pb.SendDataState_SEND_DATA_ERROR,
		}, err
	}
	encryptedPassword, err := cipher.EncryptData([]byte(req.Credential.Password))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt credentials", "error", err)
		return &pb.SendCredentialsResponse{
			State: pb.SendDataState_SEND_DATA_ERROR,
		}, err
	}
	err = gs.storage.NewCredentials(ctx, string(login), req.Credential.CredentialsName, string(encryptedLogin), string(encryptedPassword))
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
	login := ctx.Value("login").(servicestype.KeyType)
	log, pass, err := gs.storage.GetCredentials(ctx, string(login), req.CredentialName)
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
	decryptedlog, err := cipher.DecryptData([]byte(log))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt credentials", "error", err)
		return &pb.GetCredentialsResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	decryptedpass, err := cipher.DecryptData([]byte(pass))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt credentials", "error", err)
		return &pb.GetCredentialsResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	return &pb.GetCredentialsResponse{
		Credentials: &pb.Credentials{
			CredentialsName: req.CredentialName,
			Login:           string(decryptedlog),
			Password:        string(decryptedpass),
		},
		State: pb.GetDataState_GET_DATA_SUCCESS,
	}, nil
}

func (gs *GophKeepService) SendText(ctx context.Context, req *pb.SendTextRequest) (*pb.SendTextResponse, error) {
	login := ctx.Value("login").(servicestype.KeyType)
	encryptText, err := cipher.EncryptData([]byte(req.Text.Text))
	if err != nil {
		gs.logger.Errorw("Failed to encrypt text", "error", err)
		return &pb.SendTextResponse{
			State: pb.SendDataState_SEND_DATA_ERROR,
		}, err
	}
	err = gs.storage.NewText(ctx, string(login), req.Text.TextName, string(encryptText))
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
	login := ctx.Value("login").(servicestype.KeyType)
	text, err := gs.storage.GetText(ctx, string(login), req.TextName)
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
	decryptedText, err := cipher.DecryptData([]byte(text))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt text", "error", err)
		return &pb.GetTextResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	return &pb.GetTextResponse{
		State: pb.GetDataState_GET_DATA_SUCCESS,
		Text: &pb.Text{
			TextName: req.TextName,
			Text:     string(decryptedText),
		},
	}, nil
}

func (gs *GophKeepService) SendCard(ctx context.Context, req *pb.SendCardRequest) (*pb.SendCardResponse, error) {
	login := ctx.Value("login").(servicestype.KeyType)
	ecnryptCardNumber, err := cipher.EncryptData([]byte(req.Card.CardNumber))
	if err != nil {
		gs.logger.Errorw("Failed to encrypt card number", "error", err)
		return &pb.SendCardResponse{
			State: pb.SendDataState_SEND_DATA_ERROR,
		}, err
	}
	encryptCardHolderName, err := cipher.EncryptData([]byte(req.Card.CardHolderName))
	if err != nil {
		gs.logger.Errorw("Failed to encrypt cardholder name", "error", err)
		return &pb.SendCardResponse{
			State: pb.SendDataState_SEND_DATA_ERROR,
		}, err
	}
	ecnryptExpirationDate, err := cipher.EncryptData([]byte(req.Card.ExpirationDate))
	if err != nil {
		gs.logger.Errorw("Failed to encrypt expiration date", "error", err)
		return &pb.SendCardResponse{
			State: pb.SendDataState_SEND_DATA_ERROR,
		}, err
	}
	cvv := strconv.Itoa(int(req.Card.Cvv))
	ecnryptCVV, err := cipher.EncryptData([]byte(cvv))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt CVV", "error", err)
		return &pb.SendCardResponse{
			State: pb.SendDataState_SEND_DATA_ERROR,
		}, err
	}
	err = gs.storage.NewCard(ctx, string(login), req.Card.CardName,
		string(ecnryptCardNumber), string(encryptCardHolderName),
		string(ecnryptExpirationDate), string(ecnryptCVV))
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
	login := ctx.Value("login").(servicestype.KeyType)
	cardNumber, cardHolderName, expirationDate, cvv, err := gs.storage.GetCard(ctx, string(login), req.CardName)
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
		}, err
	}
	decryptCardNumber, err := cipher.DecryptData([]byte(cardNumber))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt card", "error", err)
		return &pb.GetCardResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	decryptCardHolderName, err := cipher.DecryptData([]byte(cardHolderName))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt cardholder", "error", err)
		return &pb.GetCardResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	decryptExpirationDate, err := cipher.DecryptData([]byte(expirationDate))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt expiration date", "error", err)
		return &pb.GetCardResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	decryptCVV, err := cipher.DecryptData([]byte(cvv))
	if err != nil {
		gs.logger.Errorw("Failed to decrypt CVV", "error", err)
		return &pb.GetCardResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	cvvInt, err := strconv.ParseInt(string(decryptCVV), 10, 32)
	if err != nil {
		gs.logger.Errorw("Failed to parse CVV", "error", err)
		return &pb.GetCardResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	return &pb.GetCardResponse{
		State: pb.GetDataState_GET_DATA_SUCCESS,
		Card: &pb.Card{
			CardNumber:     string(decryptCardNumber),
			CardHolderName: string(decryptCardHolderName),
			ExpirationDate: string(decryptExpirationDate),
			Cvv:            int32(cvvInt),
		},
	}, nil
}

func (gs *GophKeepService) Sync(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
	login := ctx.Value("login").(servicestype.KeyType)
	texts, err := gs.storage.TextList(ctx, string(login))
	if err != nil {
		gs.logger.Errorw("Failed to get text list", "error", err)
		return &pb.SyncResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	credentials, err := gs.storage.CredentialList(ctx, string(login))
	if err != nil {
		gs.logger.Errorw("Failed to get credential list", "error", err)
		return &pb.SyncResponse{
			State: pb.GetDataState_GET_DATA_ERROR,
		}, err
	}
	cards, err := gs.storage.CardList(ctx, string(login))
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
