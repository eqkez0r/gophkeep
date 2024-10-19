package checkers

import "testing"

func TestCreditCardNumberCheck(t *testing.T) {
	type args struct {
		cardNumber string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid_card_1",
			args: args{
				cardNumber: "4111111111111111",
			},
			want: true,
		},
		{
			name: "valid_card_2",
			args: args{
				cardNumber: "4627-1001-0165-4724",
			},
			want: true,
		},
		{
			name: "valid_card_3",
			args: args{
				cardNumber: "5467 9298 5807 4128",
			},
			want: true,
		},
		{
			name: "invalid_card_1",
			args: args{
				cardNumber: "1111/2222/3333/4444",
			},
			want: false,
		},
		{
			name: "invalid_card_2",
			args: args{
				cardNumber: "1223334444",
			},
			want: false,
		},
		{
			name: "invalid_card_3",
			args: args{
				cardNumber: "5467_9298_5807_4128",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreditCardNumberCheck(tt.args.cardNumber); got != tt.want {
				t.Errorf("CreditCardNumberCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreditCardExpirationCheck(t *testing.T) {
	type args struct {
		expiration string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid_expiration_date_1",
			args: args{
				expiration: "01/23",
			},
			want: true,
		},
		{
			name: "valid_expiration_date_2",
			args: args{
				expiration: "10/30",
			},
			want: true,
		},
		{
			name: "invalid_expiration_date_1",
			args: args{
				expiration: "1/23",
			},
			want: false,
		},
		{
			name: "invalid_expiration_date_2",
			args: args{
				expiration: "13/30",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreditCardExpirationCheck(tt.args.expiration); got != tt.want {
				t.Errorf("CreditCardExpirationCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreditCardCVVCheck(t *testing.T) {
	type args struct {
		cvv string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid_cvv_1",
			args: args{
				cvv: "412",
			},
			want: true,
		},
		{
			name: "valid_cvv_2",
			args: args{
				cvv: "128",
			},
			want: true,
		},
		{
			name: "valid_cvv_4",
			args: args{
				cvv: "1294",
			},
			want: true,
		},
		{
			name: "invalid_cvv_1",
			args: args{
				cvv: "12",
			},
			want: false,
		},
		{
			name: "invalid_cvv_2",
			args: args{
				cvv: "1",
			},
			want: false,
		},
		{
			name: "invalid_cvv_3",
			args: args{
				cvv: "12945",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreditCardCVVCheck(tt.args.cvv); got != tt.want {
				t.Errorf("CreditCardCVVCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
