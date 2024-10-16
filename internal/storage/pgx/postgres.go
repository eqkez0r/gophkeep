package pgx

import (
	"context"
	"errors"
	se "github.com/eqkez0r/gophkeep/internal/storage/storageerrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxStorage struct {
	pool *pgxpool.Pool
}

func New(
	databaseURL string,
) (*pgxStorage, error) {
	const queryCreateUserTabel = `CREATE TABLE IF NOT EXISTS users(
    login VARCHAR(50) PRIMARY KEY,
    password VARCHAR(128) NOT NULL
)`
	const queryCreateUserCards = `CREATE TABLE IF NOT EXISTS user_cards(
    login VARCHAR(50) REFERENCES users(login),
    card_name VARCHAR(50) NOT NULL,
    card_number VARCHAR(16) NOT NULL,
    card_holder_name VARCHAR(50) NOT NULL,
    expiredAt VARCHAR(5),
    cvv numeric
)`
	const queryCreateUserCredentials = `CREATE TABLE IF NOT EXISTS user_credentials(
    login VARCHAR(50) REFERENCES users(login),
    credential_name VARCHAR(50) NOT NULL,
    credential_login VARCHAR(50) NOT NULL,
    credential_password VARCHAR(128) NOT NULL
)`

	const queryCreateUserTexts = `CREATE TABLE IF NOT EXISTS user_texts(
    login VARCHAR(50) REFERENCES users(login),
    text_name VARCHAR(50) NOT NULL,
    text text NOT NULL
)`

	pl, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}
	err = pl.Ping(context.Background())
	if err != nil {
		pl.Close()
		return nil, err
	}

	_, err = pl.Exec(context.Background(), queryCreateUserTabel)
	if err != nil {
		pl.Close()
		return nil, err
	}
	_, err = pl.Exec(context.Background(), queryCreateUserCredentials)
	if err != nil {
		pl.Close()
		return nil, err
	}
	_, err = pl.Exec(context.Background(), queryCreateUserCards)
	if err != nil {
		pl.Close()
		return nil, err
	}
	_, err = pl.Exec(context.Background(), queryCreateUserTexts)
	if err != nil {
		pl.Close()
		return nil, err
	}
	return &pgxStorage{
		pool: pl,
	}, nil
}

func (p *pgxStorage) NewUser(ctx context.Context, login, password string) error {
	const queryCreateNewUser = `INSERT INTO users(login, password) VALUES ($1, $2)`
	_, err := p.pool.Exec(ctx, queryCreateNewUser, login, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return se.ErrUserIsExist
			}
		}
		return err
	}
	return nil
}

func (p *pgxStorage) IsUserExist(ctx context.Context, login string) (bool, error) {
	const queryIsUserExist = `SELECT COUNT(*) FROM users WHERE login = $1`
	row, err := p.pool.Query(ctx, queryIsUserExist, login)
	if err != nil {
		return false, err
	}
	defer row.Close()
	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (p *pgxStorage) ValidateUser(ctx context.Context, login, password string) error {
	const queryValidateUser = `SELECT (login, password) FROM users WHERE login = $1`
	row, err := p.pool.Query(ctx, queryValidateUser, login)
	if err != nil {
		if err == pgx.ErrNoRows {
			return se.ErrUserNotFound
		}
		return err
	}
	defer row.Close()
	var log, pass string
	err = row.Scan(&log, &pass)
	if err != nil {
		return err
	}
	if log != login || pass != password {
		return se.ErrInvalidAuthParameters
	}
	return nil
}

func (p *pgxStorage) NewCredentials(ctx context.Context, login, credentialName, credentialLogin, credentialPassword string) error {
	const queryNewCredentials = `INSERT INTO user_credentials(login, credential_name, credential_login, credential_password) VALUES ($1, $2, $3, $4)`
	_, err := p.pool.Exec(ctx, queryNewCredentials, login, credentialName, credentialLogin, credentialPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return se.ErrCredentialsIsExist
			}
		}
		return err
	}
	return nil
}

func (p *pgxStorage) GetCredentials(ctx context.Context, login, credentialName string) (string, string, error) {
	const queryGetCredentials = `SELECT (credential_login, credential_password) FROM user_credentials WHERE login = $1 AND credential_name = $2`
	row, err := p.pool.Query(ctx, queryGetCredentials, login, credentialName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", se.ErrUserNotFound
		}
		return "", "", err
	}
	defer row.Close()
	var log, pass string
	err = row.Scan(&log, &pass)
	if err != nil {
		return "", "", err
	}
	return log, pass, nil
}

func (p *pgxStorage) CredentialList(ctx context.Context, login string) ([]string, error) {
	const queryGetCredentialsList = `SELECT (credential_name) FROM user_credentials() WHERE login = $1`
	row, err := p.pool.Query(ctx, queryGetCredentialsList, login)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []string{}, nil
		}
		return nil, err
	}
	defer row.Close()
	credentials := []string{}
	for row.Next() {
		var credName string
		err = row.Scan(&credName)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credName)
	}
	return credentials, nil
}

func (p *pgxStorage) NewText(ctx context.Context, login, textname, text string) error {
	const queryNewText = `INSERT INTO user_texts(login, text_name, text) VALUES ($1, $2, $3)`
	_, err := p.pool.Exec(ctx, queryNewText, login, textname, text)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return se.ErrTextIsExist
			}
		}
		return err
	}
	return nil
}

func (p *pgxStorage) GetText(ctx context.Context, login, textname string) (string, error) {
	const queryGetText = `SELECT (text_name, text) FROM user_texts WHERE login = $1 AND text_name = $2`
	row, err := p.pool.Query(ctx, queryGetText, login, textname)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", se.ErrTextNotFound
		}
		return "", err
	}
	defer row.Close()
	var text string
	err = row.Scan(&text)
	if err != nil {
		return "", err
	}
	return text, nil
}

func (p *pgxStorage) TextList(ctx context.Context, login string) ([]string, error) {
	const queryGetTextList = `SELECT text_name FROM user_texts WHERE login = $1`
	row, err := p.pool.Query(ctx, queryGetTextList, login)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []string{}, nil
		}
		return nil, err
	}
	defer row.Close()
	texts := []string{}
	for row.Next() {
		var text string
		err = row.Scan(&text)
		if err != nil {
			return nil, err
		}
		texts = append(texts, text)
	}
	return texts, nil
}

func (p *pgxStorage) NewCard(ctx context.Context, login, cardName, cardNumber, cardHolder, expirationTime string, ccv int32) error {
	const queryNewCard = `INSERT INTO user_cards(login, card_name, card_number, card_holder, expiration_time, ccv) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := p.pool.Exec(ctx, queryNewCard, login, cardName, cardNumber, cardHolder, expirationTime, ccv)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return se.ErrCardIsExist
			}
		}
		return err
	}
	return nil
}

func (p *pgxStorage) GetCard(ctx context.Context, login, cardName string) (string, string, string, int32, error) {
	const queryGetCard = `SELECT card_number, card_holder, expiration_time, ccv FROM user_cards WHERE login = $1 AND card_name = $2`
	row, err := p.pool.Query(ctx, queryGetCard, login, cardName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", 0, se.ErrCardNotFound
		}
		return "", "", "", 0, err
	}
	defer row.Close()
	var cardNumber, cardHolder, expirationTime string
	var ccv int32
	err = row.Scan(&cardNumber, &cardHolder, &expirationTime, &ccv)
	if err != nil {
		return "", "", "", 0, err
	}
	return cardNumber, cardHolder, expirationTime, ccv, nil
}

func (p *pgxStorage) CardList(ctx context.Context, login string) ([]string, error) {
	const queryGetCardList = `SELECT card_number FROM user_cards WHERE login = $1`
	row, err := p.pool.Query(ctx, queryGetCardList, login)
	if err != nil {
		if err != pgx.ErrNoRows {
			return []string{}, nil
		}
		return nil, err
	}
	defer row.Close()
	cards := []string{}
	for row.Next() {
		var cardNumber string
		err = row.Scan(&cardNumber)
		if err != nil {
			return nil, err
		}
		cards = append(cards, cardNumber)
	}
	return cards, nil
}
