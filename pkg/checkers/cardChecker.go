package checkers

import (
	"regexp"
	"strconv"
	"strings"
)

func CreditCardNumberCheck(cardNumber string) bool {
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")

	var digits []int
	for _, digit := range cardNumber {
		num, err := strconv.Atoi(string(digit))
		if err != nil {
			return false
		}
		digits = append(digits, num)
	}

	sum := 0
	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i]
		if (len(digits)-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}

func CreditCardExpirationCheck(expiration string) bool {
	expirationPattern := regexp.MustCompile(`^(0[1-9]|1[0-2])\/(2[2-9]|[3-9][0-9])$`)
	return expirationPattern.MatchString(expiration)
}

func CreditCardCVVCheck(cvv string) bool {
	cvvPattern := regexp.MustCompile(`^\d{3,4}$`)
	return cvvPattern.MatchString(cvv)
}
