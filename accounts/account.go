package accounts

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrorInactiveAccount = errors.New("operation failed: the account is inactive")
var ErrorInsufficientCredit = errors.New("operation failed: the accoun't balance is insufficient")

type Currency int64

const (
	Cent Currency = 100
	Euro Currency = 10000
)

type Account interface {
	Open(Currency) Account
	Close() error
	Balance() (Currency, error)
	Deposit(Currency) (Currency, error)
	Withdraw(Currency) (Currency, error)
	Transfer(Account, Currency) (Currency, error)
}

type BankAccount struct {
	m         sync.RWMutex
	id        uuid.UUID
	balance   Currency
	openDate  time.Time
	closeDate time.Time
	active    bool
}

func NewBankAccount() *BankAccount {
	return &BankAccount{
		id: uuid.New(),
	}
}

func (a *BankAccount) Open(startingAmount Currency) Account {
	a.balance = startingAmount
	a.openDate = time.Now()
	a.active = true
	return a
}

func (a *BankAccount) Close() error {
	a.m.Lock()
	defer a.m.Unlock()

	if !a.active {
		return ErrorInactiveAccount
	}
	a.active = false
	a.closeDate = time.Now()
	a.balance = 0

	return nil
}

func (a *BankAccount) Balance() (Currency, error) {
	a.m.RLock()
	defer a.m.RUnlock()

	if !a.active {
		return 0, ErrorInactiveAccount
	}
	return a.balance, nil
}

func (a *BankAccount) Deposit(amount Currency) (Currency, error) {
	a.m.Lock()
	defer a.m.Unlock()

	if !a.active {
		return 0, ErrorInactiveAccount
	}
	a.balance += amount
	return a.balance, nil
}

func (a *BankAccount) Withdraw(amount Currency) (Currency, error) {
	a.m.Lock()
	defer a.m.Unlock()

	if !a.active {
		return 0, ErrorInactiveAccount
	}
	if a.balance < amount {
		return 0, ErrorInsufficientCredit
	}
	a.balance -= amount
	return a.balance, nil
}

func (a *BankAccount) Transfer(target Account, amount Currency) (Currency, error) {
	if _, err := a.Withdraw(amount); err != nil {
		return 0, err
	}
	if _, err := target.Deposit(amount); err != nil {
		// should never error
		a.Deposit(amount)
		return 0, err
	}

	return a.balance, nil
}
