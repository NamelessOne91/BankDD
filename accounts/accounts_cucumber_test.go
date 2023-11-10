package accounts

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

type accountCtxKey struct{}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Given(`^I have a bank account$`, IHaveABankAccount)
	sc.When(`^I open it with a starting sum of (\d+) euro$`, IOpenWithStartingEuro)
	sc.Then(`^it should have a balance of (\d+) euro$`, ItShouldHaveABalanceOf)
	sc.Then(`^it should be an active account$`, ItShouldBeActive)
	sc.Then(`^it should have today as opening date$`, ItShouldHaveTodayAsOpenDate)
}

func IHaveABankAccount(ctx context.Context) (context.Context, error) {
	account := NewBankAccount()
	return context.WithValue(ctx, accountCtxKey{}, account), nil
}

func IOpenWithStartingEuro(ctx context.Context, euro int) (context.Context, error) {
	account, ok := ctx.Value(accountCtxKey{}).(*BankAccount)
	if !ok {
		return ctx, errors.New("no bank account available")
	}

	if err := account.Open(Euro * Currency(euro)); err != nil {
		return ctx, fmt.Errorf("got an error opening bank account: %s", err)
	}
	return ctx, nil
}

func ItShouldHaveABalanceOf(ctx context.Context, euro int) error {
	account, ok := ctx.Value(accountCtxKey{}).(*BankAccount)
	if !ok {
		return errors.New("no bank account available")
	}

	expectedBalance := Euro * Currency(euro)
	if account.balance != expectedBalance {
		return fmt.Errorf("expected account's balance to be %d - got %d", expectedBalance, account.balance)
	}
	return nil
}

func ItShouldBeActive(ctx context.Context) error {
	account, ok := ctx.Value(accountCtxKey{}).(*BankAccount)
	if !ok {
		return errors.New("no bank account available")
	}

	if !account.active {
		return errors.New("expected the account to be flagged as active but it is not")
	}
	return nil
}

func ItShouldHaveTodayAsOpenDate(ctx context.Context) error {
	account, ok := ctx.Value(accountCtxKey{}).(*BankAccount)
	if !ok {
		return errors.New("no bank account available")
	}

	today := time.Now().Truncate(time.Minute)
	if account.openDate.Truncate(time.Minute).Compare(today) != 0 {
		return fmt.Errorf("expected account's open date to be %s - got %s", today.String(), account.openDate.Truncate(time.Minute).String())
	}
	return nil
}
