package storage

import (
	"context"
	"github.com/eqkez0r/gophkeep/internal/storage/memory"
	"github.com/eqkez0r/gophkeep/internal/storage/pgx"
	se "github.com/eqkez0r/gophkeep/internal/storage/storageerrors"
	"log"
)

type Storage interface {
	NewUser(context.Context, string, string) error
	IsUserExist(context.Context, string) (bool, error)
	ValidateUser(context.Context, string, string) error

	//login,credentialsname,credentiallogin, credentialspass
	NewCredentials(context.Context, string, string, string, string) error
	GetCredentials(context.Context, string, string) (string, string, error)
	CredentialList(context.Context, string) ([]string, error)

	//login, textname, text
	NewText(context.Context, string, string, string) error
	GetText(context.Context, string, string) (string, error)
	TextList(context.Context, string) ([]string, error)

	//login,cardname,cardnumber,cardholder,expirationDate,ccv
	NewCard(context.Context, string, string, string, string, string, string) error
	GetCard(context.Context, string, string) (string, string, string, string, error)
	CardList(context.Context, string) ([]string, error)
}

func New() (Storage, error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, err
	}
	log.Printf("store cfg %v", cfg)
	switch cfg.DatabaseType {
	case "memory":
		return memory.New(), nil
	case "postgres":
		return pgx.New(cfg.DatabaseURL)
	}
	return nil, se.ErrUnknownDatabaseType
}
