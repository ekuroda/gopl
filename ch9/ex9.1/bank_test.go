package bank

import "testing"

func TestWithdraw(t *testing.T) {

	Deposit(1000)
	if ok := Withdraw(500); !ok {
		t.Errorf("Withdraw(%d) = %t, want %t", 500, ok, true)
	}
	if b := Balance(); b != 500 {
		t.Errorf("Balance() = %d, want %d", b, 500)
	}
	if ok := Withdraw(501); ok {
		t.Errorf("Withdraw(%d) = %t, want %t", 501, ok, false)
	}
}
