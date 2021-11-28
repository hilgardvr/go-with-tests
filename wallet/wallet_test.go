package wallet

import (
	"testing"
)

func TestWallet(t *testing.T) {

	assertBalance := func(t testing.TB, w Wallet, wanted Bitcoin) {
		t.Helper()
		got := w.Balance()
		if got != wanted {
			t.Errorf("got %s want %s", got, wanted)
		}
	}

	assertError := func(t testing.TB, err error, want error) {
		t.Helper()
		if err == nil {
			t.Fatalf("wanted an error but didn't get one")
		}
		if err != want {
			t.Errorf("got %q want %q", err, want)
		}
	}
	assertNoError := func(t testing.TB, err error) {
		t.Helper()
		if err != nil {
			t.Errorf("got %q want nil", err)
		}
	}

	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		want := Bitcoin(10)
		assertBalance(t, wallet, want)
	})
	
	t.Run("withdraw", func(t *testing.T) {
		wallet := Wallet{balance: 10}
		err := wallet.Withdraw(Bitcoin(9))
		want := Bitcoin(1)
		assertBalance(t, wallet, want)
		assertNoError(t, err)
	})

	t.Run("insufficient funds", func(t *testing.T){
		startingBalance := Bitcoin(10)
		wallet := Wallet{balance: startingBalance}
		err := wallet.Withdraw(11)
		assertBalance(t, wallet, startingBalance)
		assertError(t, err, ErrInsufficientFunds)
	})

	t.Run("stringify", func(t *testing.T) {
		btc := Bitcoin(1)
		got := btc.String()
		wanted := "1 BTC"

		if got != wanted {
			t.Errorf("got %s wanted %s", got, wanted)
		}
	})
}