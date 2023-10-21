package accounts_test

import (
	"bankdd/accounts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	var startingCurrency = 1000*accounts.Euro + 50*accounts.Cent

	Describe("can open a new bank account properly", func() {
		var account = accounts.NewBankAccount()
		account.Open(startingCurrency)

		It("is active", func() {
			_, err := account.Balance()
			Expect(err).To(BeNil())
		})

		It("has a balance equal to the starting amount", func() {
			b, err := account.Balance()

			Expect(err).To(BeNil())
			Expect(b).To(Equal(startingCurrency))
		})
	})

	Describe("can close an existing bank account properly", func() {
		var account accounts.Account

		BeforeEach(func() {
			account = accounts.NewBankAccount()
			account.Open(startingCurrency)
			Expect(account.Close()).To(Succeed())
		})

		It("is inactive", func() {
			_, err := account.Balance()
			Expect(err).To(Equal(accounts.ErrorInactiveAccount))
		})

		It("returns an error if we try to close it again", func() {
			Expect(account.Close()).ToNot(Succeed())
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

})
