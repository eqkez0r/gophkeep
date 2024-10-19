package memory

import (
	"context"
	"github.com/eqkez0r/gophkeep/internal/models"
	se "github.com/eqkez0r/gophkeep/internal/storage/storageerrors"
	"github.com/eqkez0r/gophkeep/pkg/cipher"
	"sync"
)

type memoryStorage struct {
	mu sync.RWMutex
	//user db
	users map[string]string
	//user --- cred_name --- cred
	credentials map[string]map[string]models.Credentials
	//user --- text_name --- text
	texts map[string]map[string]models.Text
	//user -- card_name --- card
	cards map[string]map[string]models.Card
	// user --- attachment_name --- attachment
	attachments map[string]map[string]models.Attachment
}

func New() *memoryStorage {
	return &memoryStorage{
		users:       make(map[string]string),
		credentials: make(map[string]map[string]models.Credentials),
		texts:       make(map[string]map[string]models.Text),
		cards:       make(map[string]map[string]models.Card),
		attachments: make(map[string]map[string]models.Attachment),
	}
}

func (m *memoryStorage) NewUser(ctx context.Context, login, password string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.users[login]
	if ok {
		return se.ErrUserIsExist
	}
	m.users[login] = password
	return nil
}

func (m *memoryStorage) IsUserExist(ctx context.Context, login string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.users[login]
	return ok, nil
}

func (m *memoryStorage) ValidateUser(ctx context.Context, login, password string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pass, ok := m.users[login]
	if !ok {
		return se.ErrUserNotFound
	}
	depass, err := cipher.DecryptData([]byte(pass))
	if err != nil {
		return err
	}
	if string(depass) != password {
		return se.ErrInvalidAuthParameters
	}
	return nil
}

func (m *memoryStorage) NewCredentials(ctx context.Context, login, credentialName, credentialLogin, credentialPassword string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mp, ok := m.credentials[login]
	if !ok {
		mp = map[string]models.Credentials{}
	}
	mp[credentialName] = models.Credentials{
		CredentialsName: credentialName,
		Login:           credentialLogin,
		Password:        credentialPassword,
	}
	m.credentials[login] = mp
	return nil
}

func (m *memoryStorage) GetCredentials(ctx context.Context, login, credentialName string) (string, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mp, ok := m.credentials[login]
	if !ok {
		return "", "", se.ErrCredentialsNotFound
	}
	credential, ok := mp[credentialName]
	if !ok {
		return "", "", se.ErrCredentialsNotFound
	}
	return credential.Login, credential.Password, nil
}

func (m *memoryStorage) CredentialList(ctx context.Context, login string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mp, ok := m.credentials[login]
	if !ok {
		return []string{}, nil
	}
	credentials := make([]string, 0, len(mp))
	for k := range mp {
		credentials = append(credentials, k)
	}
	return credentials, nil
}

func (m *memoryStorage) NewText(ctx context.Context, login, textname, text string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mp, ok := m.texts[login]
	if !ok {
		mp = map[string]models.Text{}
	}
	mp[textname] = models.Text{
		TextName: textname,
		Text:     text,
	}
	m.texts[login] = mp
	return nil
}

func (m *memoryStorage) GetText(ctx context.Context, login, textname string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mp, ok := m.texts[login]
	if !ok {
		return "", se.ErrTextNotFound
	}
	text, ok := mp[textname]
	if !ok {
		return "", se.ErrTextNotFound
	}
	return text.Text, nil
}

func (m *memoryStorage) TextList(ctx context.Context, login string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mp, ok := m.texts[login]
	if !ok {
		return []string{}, nil
	}
	texts := make([]string, 0, len(mp))
	for k := range mp {
		texts = append(texts, k)
	}
	return texts, nil
}

func (m *memoryStorage) NewCard(ctx context.Context, login, cardName, cardNumber, cardholder, expirationDate, cvv string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mp, ok := m.cards[login]
	if !ok {
		mp = map[string]models.Card{}
	}
	mp[cardName] = models.Card{
		CardName:       cardName,
		CardNumber:     cardNumber,
		CardHolderName: cardholder,
		ExpirationDate: expirationDate,
		CVV:            cvv,
	}
	m.cards[login] = mp
	return nil
}

func (m *memoryStorage) GetCard(ctx context.Context, login, cardname string) (string, string, string, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mp, ok := m.cards[login]
	if !ok {
		return "", "", "", "", se.ErrCardNotFound
	}
	card, ok := mp[cardname]
	if !ok {
		return "", "", "", "", se.ErrCardNotFound
	}
	return card.CardNumber, card.CardHolderName, card.ExpirationDate, card.CVV, nil
}

func (m *memoryStorage) CardList(ctx context.Context, login string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mp, ok := m.cards[login]
	if !ok {
		return []string{}, nil
	}
	cards := make([]string, 0, len(mp))
	for k := range mp {
		cards = append(cards, k)
	}
	return cards, nil
}
