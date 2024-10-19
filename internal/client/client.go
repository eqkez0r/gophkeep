package client

import (
	"context"
	pb "github.com/eqkez0r/gophkeep-grpc-api/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn       *grpc.ClientConn
	gophkeepcl pb.GophKeeperClient
}

func New() *GRPCClient {
	return &GRPCClient{}
}

func (gc *GRPCClient) Login(
	host, login, password string,
) (string, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	grpccl := pb.NewGophKeeperClient(conn)
	res, err := grpccl.Auth(context.Background(), &pb.AuthRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	gc.conn = conn
	gc.gophkeepcl = grpccl
	return *res.Token, nil
}

func (gc *GRPCClient) Register(
	host, login, password string,
) (string, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}

	grpccl := pb.NewGophKeeperClient(conn)
	res, err := grpccl.Register(context.Background(), &pb.RegisterRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	gc.conn = conn
	gc.gophkeepcl = grpccl
	return *res.Token, nil
}

func (gc *GRPCClient) SendCredentials(token, credentialsName, credentialsLogin, credentialsPassword string) error {
	_, err := gc.gophkeepcl.SendCredentials(context.Background(), &pb.SendCredentialsRequest{
		Token: token,
		Credential: &pb.Credentials{
			CredentialsName: credentialsName,
			Login:           credentialsLogin,
			Password:        credentialsPassword,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (gc *GRPCClient) GetCredentials(token, credentialsName string) (string, string, error) {
	res, err := gc.gophkeepcl.GetCredentials(context.Background(), &pb.GetCredentialsRequest{
		Token:          token,
		CredentialName: credentialsName,
	})
	if err != nil {
		return "", "", err
	}
	return res.Credentials.Login, res.Credentials.Password, nil
}

func (gc *GRPCClient) SendCard(token, cardName, cardNumber, cardHolderName, expirationTime string, cvv int32) error {
	_, err := gc.gophkeepcl.SendCard(context.Background(), &pb.SendCardRequest{
		Card: &pb.Card{
			CardName:       cardName,
			CardNumber:     cardNumber,
			CardHolderName: cardHolderName,
			ExpirationDate: expirationTime,
			Cvv:            cvv,
		},
		Token: token,
	})
	if err != nil {
		return err
	}
	return nil
}

func (gc *GRPCClient) GetCard(token, cardName string) (string, string, string, int32, error) {
	res, err := gc.gophkeepcl.GetCard(context.Background(), &pb.GetCardRequest{
		Token:    token,
		CardName: cardName,
	})
	if err != nil {
		return "", "", "", 0, err
	}
	return res.Card.CardNumber, res.Card.CardHolderName, res.Card.ExpirationDate, res.Card.Cvv, nil
}

func (gc *GRPCClient) SendText(token, textName, text string) error {
	_, err := gc.gophkeepcl.SendText(context.Background(), &pb.SendTextRequest{
		Token: token,
		Text: &pb.Text{
			TextName: textName,
			Text:     text,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (gc *GRPCClient) GetText(token, textName string) (string, error) {
	res, err := gc.gophkeepcl.GetText(context.Background(), &pb.GetTextRequest{
		Token:    token,
		TextName: textName,
	})
	if err != nil {
		return "", err
	}
	return res.Text.Text, nil
}

func (gc *GRPCClient) Sync(token string) ([]string, []string, []string, error) {
	res, err := gc.gophkeepcl.Sync(context.Background(), &pb.SyncRequest{
		Token: token,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return res.Credentials, res.Cards, res.Texts, nil
}

func (gc *GRPCClient) Logout() error {
	err := gc.conn.Close()
	if err != nil {
		return err
	}
	gc.conn = nil
	gc.gophkeepcl = nil
	return nil
}

func (gc *GRPCClient) CheckConnection() bool {
	return gc.conn != nil || gc.gophkeepcl != nil
}
