package accounts_test

import (
	"bankdd/accounts"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	var startingCurrency = 1000*accounts.Euro + 50*accounts.Cent

	Describe("can open a new bank account properly", func() {
		var account = accounts.NewBankAccount()
		Expect(account.Open(startingCurrency)).To(Succeed())

		It("is active", func() {
			_, err := account.Balance()
			Expect(err).To(BeNil())
		})

		It("has a balance equal to the starting amount", func() {
			b, err := account.Balance()

			Expect(err).To(BeNil())
			Expect(b).To(Equal(startingCurrency))
		})

		It("returns an error if we try to open an active account", func() {
			Expect(account.Open(startingCurrency)).To(Equal(accounts.ErrorAccountIsOpen))
		})
	})

	Describe("can close an existing bank account properly", func() {
		var account accounts.Account

		BeforeEach(func() {
			account = accounts.NewBankAccount()
			Expect(account.Open(startingCurrency)).To(Succeed())
			Expect(account.Close()).To(Succeed())
		})

		It("is inactive", func() {
			_, err := account.Balance()
			Expect(err).To(Equal(accounts.ErrorInactiveAccount))
		})

		It("returns an error if we try to close it again", func() {
			err := account.Close()
			Expect(err).To(Equal(accounts.ErrorInactiveAccount))
		})

		It("returns an error if we try to deposit money", func() {
			_, err := account.Deposit(100 * accounts.Euro)
			Expect(err).To(Equal(accounts.ErrorInactiveAccount))
		})

		It("returns an error if we try to withdraw money", func() {
			_, err := account.Withdraw(100 * accounts.Euro)
			Expect(err).To(Equal(accounts.ErrorInactiveAccount))
		})

		It("returns an error if we try to transfer money", func() {
			_, err := account.Transfer(nil, 100*accounts.Euro)
			Expect(err).To(Equal(accounts.ErrorInactiveAccount))
		})
	})

	Describe("can deposit money to an open account", func() {

		DescribeTable("correctly updates the account's balance after a deposit",
			func(start accounts.Currency, deposits []accounts.Currency, expectedBalance accounts.Currency) {
				var account = accounts.NewBankAccount()
				Expect(account.Open(start)).To(Succeed())

				for _, amount := range deposits {
					_, err := account.Deposit(amount)
					Expect(err).To(BeNil())
				}
				Expect(account.Balance()).To(Equal(expectedBalance))
			},
			func(start accounts.Currency, deposits []accounts.Currency, expectedBalance accounts.Currency) string {
				return fmt.Sprintf("start: %d - deposits: %v - expected balance: %d", start, deposits, expectedBalance)
			},
			Entry(nil, 0*accounts.Euro, []accounts.Currency{}, 0*accounts.Euro),
			Entry(nil, 0*accounts.Euro, []accounts.Currency{100 * accounts.Euro}, 100*accounts.Euro),
			Entry(nil, 0*accounts.Euro, []accounts.Currency{100 * accounts.Euro, 30 * accounts.Cent}, (100*accounts.Euro+30*accounts.Cent)),
			Entry(nil, 1000*accounts.Euro, []accounts.Currency{150 * accounts.Euro, 50 * accounts.Cent, 720 * accounts.Euro}, (1870*accounts.Euro+50*accounts.Cent)),
			Entry(nil, 1000*accounts.Euro, []accounts.Currency{150 * accounts.Euro, 50 * accounts.Cent, 720 * accounts.Euro, 3 * accounts.Cent}, (1870*accounts.Euro+53*accounts.Cent)),
		)
	})
})
