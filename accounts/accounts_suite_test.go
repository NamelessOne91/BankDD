package accounts_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAccounts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Accounts Suite")
}
