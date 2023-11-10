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

type accountCtx struct {
	account *BankAccount
	err     error
}

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
	sc.Given(`^I have a new bank account$`, IHaveANewBankAccount)
	sc.Given(`^I have an active bank account$`, IHaveAnActiveBankAccount)
	sc.When(`^I open it with a starting sum of (-?\d+) euro$`, IOpenWithStartingEuro)
	sc.Then(`^it should have a balance of (\d+) euro$`, ItShouldHaveABalanceOf)
	sc.Then(`^it should be an active account$`, ItShouldBeActive)
	sc.Then(`^it should have today as opening date$`, ItShouldHaveTodayAsOpenDate)
	sc.Then(`^I should get an invalid amount error$`, IShouldGetInvalidAmountError)
	sc.Then(`^I should get an error because the account is already open$`, IShouldGetAlreadyOpenError)
}

func IHaveANewBankAccount(ctx context.Context) (context.Context, error) {
	account := NewBankAccount()
	return context.WithValue(ctx, accountCtxKey{}, &accountCtx{account: account}), nil
}

func IHaveAnActiveBankAccount(ctx context.Context) (context.Context, error) {
	account := NewBankAccount()
	account.active = true
	return context.WithValue(ctx, accountCtxKey{}, &accountCtx{account: account}), nil
}

func IOpenWithStartingEuro(ctx context.Context, euro int) (context.Context, error) {
	accountCtx, ok := ctx.Value(accountCtxKey{}).(*accountCtx)
	if !ok {
		return ctx, errors.New("no bank account available")
	}

	if err := accountCtx.account.Open(Euro * Currency(euro)); err != nil {
		if euro > 0 && !accountCtx.account.active {
			return ctx, err
		}
		accountCtx.err = err
	}
	return ctx, nil
}

func ItShouldHaveABalanceOf(ctx context.Context, euro int) error {
	accountCtx, ok := ctx.Value(accountCtxKey{}).(*accountCtx)
	if !ok {
		return errors.New("no bank account available")
	}
	if accountCtx.err != nil {
		return fmt.Errorf("error in previous step: %s", accountCtx.err)
	}

	expectedBalance := Euro * Currency(euro)
	if accountCtx.account.balance != expectedBalance {
		return fmt.Errorf("expected account's balance to be %d - got %d", expectedBalance, accountCtx.account.balance)
	}
	return nil
}

func ItShouldBeActive(ctx context.Context) error {
	accountCtx, ok := ctx.Value(accountCtxKey{}).(*accountCtx)
	if !ok {
		return errors.New("no bank account available")
	}
	if accountCtx.err != nil {
		return fmt.Errorf("error in previous step: %s", accountCtx.err)
	}

	if !accountCtx.account.active {
		return errors.New("expected the account to be flagged as active but it is not")
	}
	return nil
}

func ItShouldHaveTodayAsOpenDate(ctx context.Context) error {
	accountCtx, ok := ctx.Value(accountCtxKey{}).(*accountCtx)
	if !ok {
		return errors.New("no bank account available")
	}
	if accountCtx.err != nil {
		return fmt.Errorf("error in previous step: %s", accountCtx.err)
	}

	today := time.Now().Truncate(time.Minute)
	if accountCtx.account.openDate.Truncate(time.Minute).Compare(today) != 0 {
		return fmt.Errorf(
			"expected account's open date to be %s - got %s",
			today.String(),
			accountCtx.account.openDate.Truncate(time.Minute).String(),
		)
	}
	return nil
}

func IShouldGetInvalidAmountError(ctx context.Context) error {
	accountCtx, ok := ctx.Value(accountCtxKey{}).(*accountCtx)
	if !ok {
		return errors.New("no bank account available")
	}
	if accountCtx.err != ErrorInvalidAmount {
		return fmt.Errorf("expected error in previous step to be %s - got %s", ErrorInvalidAmount, accountCtx.err)
	}

	return nil
}

func IShouldGetAlreadyOpenError(ctx context.Context) error {
	accountCtx, ok := ctx.Value(accountCtxKey{}).(*accountCtx)
	if !ok {
		return errors.New("no bank account available")
	}
	if accountCtx.err != ErrorAccountIsOpen {
		return fmt.Errorf("expected error in previous step to be %s - got %s", ErrorAccountIsOpen, accountCtx.err)
	}

	return nil
}
