package accounts_test

import (
	"bankdd/accounts"
	"fmt"
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	var startingCurrency = 1000*accounts.Euro + 50*accounts.Cent

	When("a user opens a new bank account", func() {
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

	When(" a user closes an existing bank account", func() {
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

	When("a user deposits money to an open account", func() {
		var account accounts.Account

		BeforeEach(func() {
			account = accounts.NewBankAccount()
			Expect(account.Open(startingCurrency)).To(Succeed())
		})

		It("fails to deposit an amount = 0", func() {
			_, err := account.Deposit(0 * accounts.Euro)
			Expect(err).To(Equal(accounts.ErrorInvalidAmount))
		})

		It("fails to deposit a negative amount", func() {
			_, err := account.Deposit(-100 * accounts.Euro)
			Expect(err).To(Equal(accounts.ErrorInvalidAmount))
		})

		DescribeTable("the account's balance is correctly updated after sequential deposits",
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
			Entry(nil, 10*accounts.Euro, []accounts.Currency{}, 10*accounts.Euro),
			Entry(nil, 10*accounts.Euro, []accounts.Currency{100 * accounts.Euro}, 110*accounts.Euro),
			Entry(nil, 10*accounts.Euro, []accounts.Currency{100 * accounts.Euro, 30 * accounts.Cent}, (110*accounts.Euro+30*accounts.Cent)),
			Entry(nil, 1000*accounts.Euro, []accounts.Currency{150 * accounts.Euro, 50 * accounts.Cent, 720 * accounts.Euro}, (1870*accounts.Euro+50*accounts.Cent)),
			Entry(nil, 1000*accounts.Euro, []accounts.Currency{150 * accounts.Euro, 50 * accounts.Cent, 720 * accounts.Euro, 3 * accounts.Cent}, (1870*accounts.Euro+53*accounts.Cent)),
		)

		DescribeTable("the account's balance is correctly updated with concurrent deposits",
			func(start accounts.Currency, deposits []accounts.Currency, expectedBalance accounts.Currency) {
				var wg sync.WaitGroup
				var account = accounts.NewBankAccount()
				Expect(account.Open(start)).To(Succeed())

				wg.Add(len(deposits))
				for _, amount := range deposits {
					go func(a accounts.Currency) {
						defer wg.Done()

						time.Sleep(100 + time.Duration(rand.Intn(500)+1)*time.Millisecond)
						_, err := account.Deposit(a)
						Expect(err).To(BeNil())
					}(amount)
				}
				wg.Wait()

				Expect(account.Balance()).To(Equal(expectedBalance))
			},
			func(start accounts.Currency, deposits []accounts.Currency, expectedBalance accounts.Currency) string {
				return fmt.Sprintf("start: %d - deposits: %v - expected balance: %d", start, deposits, expectedBalance)
			},
			Entry(nil, 10*accounts.Euro, []accounts.Currency{}, 10*accounts.Euro),
			Entry(nil, 10*accounts.Euro, []accounts.Currency{100 * accounts.Euro}, 110*accounts.Euro),
			Entry(nil, 10*accounts.Euro, []accounts.Currency{100 * accounts.Euro, 30 * accounts.Cent}, (110*accounts.Euro+30*accounts.Cent)),
			Entry(nil, 1000*accounts.Euro, []accounts.Currency{150 * accounts.Euro, 50 * accounts.Cent, 720 * accounts.Euro}, (1870*accounts.Euro+50*accounts.Cent)),
			Entry(nil, 1000*accounts.Euro, []accounts.Currency{150 * accounts.Euro, 50 * accounts.Cent, 720 * accounts.Euro, 3 * accounts.Cent}, (1870*accounts.Euro+53*accounts.Cent)),
		)
	})

	When("a user withdraws money from an open account", func() {
		var account accounts.Account

		BeforeEach(func() {
			account = accounts.NewBankAccount()
			Expect(account.Open(startingCurrency)).To(Succeed())
		})

		It("fails to withdraw an amount = 0", func() {
			_, err := account.Withdraw(0 * accounts.Euro)
			Expect(err).To(Equal(accounts.ErrorInvalidAmount))
		})

		It("fails to withdraw a negative amount", func() {
			_, err := account.Withdraw(-100 * accounts.Euro)
			Expect(err).To(Equal(accounts.ErrorInvalidAmount))
		})

		It("fails to withdraw an amount bigger than the account's balance", func() {
			_, err := account.Withdraw(2000 * accounts.Euro)
			Expect(err).To(Equal(accounts.ErrorInsufficientCredit))
		})
	})
})
