package bank

var deposits = make(chan int)
var balances = make(chan int)
var withdrawRequests = make(chan int)
var withdrawResults = make(chan bool)

// Deposit ...
func Deposit(amount int) { deposits <- amount }

// Balance ...
func Balance() int { return <-balances }

// Withdraw ...
func Withdraw(amount int) bool {
	withdrawRequests <- amount
	return <-withdrawResults
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case amount := <-withdrawRequests:
			if balance < amount {
				withdrawResults <- false
			} else {
				balance -= amount
				withdrawResults <- true
			}
		}
	}
}

func init() {
	go teller()
}
