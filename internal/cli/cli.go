package cli

import (
	"bufio"
	"fmt"
	"github.com/eqkez0r/gophkeep/internal/client"
	"github.com/eqkez0r/gophkeep/pkg/checkers"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	statusAuth         = "authorized"
	statusUnauthorized = "unauthorized"

	usage = `
	auth %s
	choose option:
	1. auth
	2. register
	3. add new credentials
	4. get credentials
	5. add new card
	6. get card
	7. add new text
	8. get text
	9. data list
	10. logout
	e/q. exit
`
)

type cli struct {
	token  string
	host   string
	status string

	client *client.GRPCClient
	reader *bufio.Reader

	credentials, cards, texts []string
}

func New() *cli {
	c := &cli{
		status: statusUnauthorized,
		client: client.New(),
	}
	c.reader = bufio.NewReader(os.Stdin)

	return c
}

func (c *cli) Run() {
	for {
		fmt.Printf(usage, c.status)
		input, err := c.reader.ReadString('\n')
		if err != nil {
			fmt.Printf("input err %v", err)
			continue
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			err = c.login()
			if err != nil {
				continue
			}
		case "2":
			err = c.register()
			if err != nil {
				continue
			}
		case "3":
			err = c.sendCredentials()
			if err != nil {
				continue
			}
		case "4":
			err = c.getCredentials()
			if err != nil {
				continue
			}
		case "5":
			err = c.sendCard()
			if err != nil {
				continue
			}
		case "6":
			err = c.getCard()
			if err != nil {
				continue
			}
		case "7":
			err = c.sendText()
			if err != nil {
				continue
			}
		case "8":
			err = c.getText()
			if err != nil {
				continue
			}
		case "9":
			err = c.dataList()
			if err != nil {
				continue
			}
		case "10":
			c.logout()
		case "e", "q":
			if c.status == statusAuth {
				c.client.Logout()
			}
			return
		}

	}
}

func (c *cli) login() error {
	fmt.Println("enter host please:")
	host, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	host = strings.TrimSpace(host)

	fmt.Println("enter login please:")
	login, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	login = strings.TrimSpace(login)

	fmt.Println("enter password please:")
	password, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	password = strings.TrimSpace(password)
	token, err := c.client.Login(host, login, password)
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}

	c.credentials, c.cards, c.texts, err = c.client.Sync(token)
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}

	c.status = statusAuth
	c.token = token
	c.host = host

	return nil
}

func (c *cli) register() error {
	fmt.Println("registration...")
	fmt.Println("enter host please:")
	host, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	host = strings.TrimSpace(host)

	fmt.Println("enter login please: ")
	login, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	login = strings.TrimSpace(login)

	fmt.Println("enter password please: ")
	password, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	password = strings.TrimSpace(password)

	c.token, err = c.client.Register(host, login, password)
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}

	c.status = statusAuth

	return nil
}

func (c *cli) sendCard() error {
	fmt.Println("enter send card name please:")
	cardName, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	cardName = strings.TrimSpace(cardName)

enterNumber:
	fmt.Println("enter send card number please:")
	cardNumber, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	cardNumber = strings.TrimSpace(cardNumber)
	if !checkers.CreditCardNumberCheck(cardNumber) {
		fmt.Println("card number is invalid. try again")
		goto enterNumber
	}

	fmt.Println("enter send cardholder name please:")
	cardHolderName, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	cardHolderName = strings.TrimSpace(cardHolderName)

enterExpiration:
	fmt.Println("enter send card expiration time please:")
	cardExpiration, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	cardExpiration = strings.TrimSpace(cardExpiration)
	if !checkers.CreditCardExpirationCheck(cardExpiration) {
		fmt.Println("card expiration time is invalid. try again")
		goto enterExpiration
	}

enterCVV:
	fmt.Println("enter send card cvv token please:")
	cvv, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	cvv = strings.TrimSpace(cvv)
	if !checkers.CreditCardCVVCheck(cvv) {
		fmt.Println("card cvv token is invalid. try again")
		goto enterCVV
	}

	cvvInt, err := strconv.Atoi(cvv)
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}

	err = c.client.SendCard(c.token, cardName, cardNumber, cardHolderName, cardExpiration, int32(cvvInt))
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}
	c.cards = append(c.cards, cardName)
	return nil
}

func (c *cli) getCard() error {
	if len(c.cards) > 0 {
		fmt.Println("available cards: ", strings.Join(c.cards, ","))
	} else {
		fmt.Println("no available cards")
		return nil
	}
	fmt.Println("enter get card name please:")
	cardName, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	cardName = strings.TrimSpace(cardName)
	cardNumber, cardHolderName, cardExpirationTime, cvv, err := c.client.GetCard(c.token, cardName)
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}
	log.Println("card name: ", cardName)
	log.Println("card number: ", cardNumber)
	log.Println("card holder: ", cardHolderName)
	log.Println("card expiration time: ", cardExpirationTime)
	log.Println("card cvv: ", cvv)
	return nil
}

func (c *cli) sendCredentials() error {
	fmt.Println("enter send credentials name please:")
	credentialName, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	credentialName = strings.TrimSpace(credentialName)

	fmt.Println("enter send credentials login please:")
	login, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	login = strings.TrimSpace(login)

	fmt.Println("enter send credentials password please:")
	password, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	password = strings.TrimSpace(password)

	err = c.client.SendCredentials(c.token, credentialName, login, password)
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}
	c.credentials = append(c.credentials, credentialName)
	return nil
}

func (c *cli) getCredentials() error {
	if len(c.credentials) > 0 {
		fmt.Println("available credentials: ", strings.Join(c.credentials, ","))
	} else {
		fmt.Println("no available credentials")
		return nil
	}
	fmt.Println("enter get credentials name please:")
	credentialName, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	credentialName = strings.TrimSpace(credentialName)
	log, pass, err := c.client.GetCredentials(c.token, credentialName)
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}
	fmt.Println("credentials name: ", credentialName)
	fmt.Println("credentials login: ", log)
	fmt.Println("credentials password: ", pass)
	return nil
}

func (c *cli) sendText() error {
	fmt.Println("enter send text name please:")
	textName, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	textName = strings.TrimSpace(textName)

	fmt.Println("enter send text number please:")
	text, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	text = strings.TrimSpace(text)

	err = c.client.SendText(c.token, textName, text)
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}
	c.texts = append(c.texts, textName)
	return nil
}

func (c *cli) getText() error {
	if len(c.texts) > 0 {
		fmt.Println("available texts: ", strings.Join(c.texts, ","))
	} else {
		fmt.Println("no available texts")
		return nil
	}

	fmt.Println("enter get text name please:")
	textName, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("input err %v", err)
		return err
	}
	textName = strings.TrimSpace(textName)

	text, err := c.client.GetText(c.token, textName)
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}

	fmt.Println("text name: ", textName)
	fmt.Println("text: ", text)

	return nil
}

func (c *cli) dataList() error {
	credentials, cards, text, err := c.client.Sync(c.token)
	if err != nil {
		fmt.Printf("client err %v", err)
		return err
	}

	fmt.Println("after sync")
	fmt.Println("credentials: ", credentials)
	fmt.Println("cards: ", cards)
	fmt.Println("text: ", text)
	c.credentials = credentials
	c.cards = cards
	c.texts = text
	return nil
}

func (c *cli) logout() {
	c.resetAuth()
}

//func (c *cli) checkAuth() bool {
//	return c.client.CheckConnection()
//}

func (c *cli) resetAuth() {
	c.cards = nil
	c.texts = nil
	c.credentials = nil
	c.token = ""
	err := c.client.Logout()
	if err != nil {
		fmt.Printf("client err %v", err)
	}
	c.client = nil
	c.status = statusUnauthorized
	c.host = ""
}
