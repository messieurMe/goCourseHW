package bank

import (
	"fmt"
	"strings"
	"time"
)

const (
	TopUpOp OperationType = iota
	WithdrawOp
)

type OperationType int64

type Clock interface {
	Now() time.Time
}

func NewRealTime() *RealClock {
	return &RealClock{}
}

type RealClock struct{}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

type Operation struct {
	OpTime   time.Time
	OpType   OperationType
	OpAmount int
	Balance  int
}

func (o Operation) String() string {
	var format string
	if o.OpType == TopUpOp {
		format = `%s +%d %d`
	} else {
		format = `%s -%d %d`
	}
	return fmt.Sprintf(format, o.OpTime.String()[:19], o.OpAmount, o.Balance)
}

type Account interface {
	TopUp(amount int) bool
	Withdraw(amount int) bool
	Operations() []Operation
	Statement() string
	Balance() int
}

func NewAccount(clock Clock) *AccountImpl {
	return &AccountImpl{
		amount: 0,
		clock:  clock,
		log:    []Operation{},
	}
}

type AccountImpl struct {
	amount int
	clock  Clock
	log    []Operation
}

func (a *AccountImpl) createOperation(operationType OperationType, amount int) Operation {
	return Operation{
		a.clock.Now(),
		operationType,
		amount,
		a.amount,
	}
}

func (a *AccountImpl) TopUp(amount int) bool {
	if a.amount < 0 || amount <= 0 {
		return false
	}
	a.amount += amount
	logEntry := a.createOperation(TopUpOp, amount)

	a.log = append(a.log, logEntry)
	return true
}

func (a *AccountImpl) Withdraw(amount int) bool {
	if a.amount < 0 || amount <= 0 || a.amount-amount < 0 {
		return false
	}
	a.amount -= amount
	logEntry := a.createOperation(WithdrawOp, amount)
	a.log = append(a.log, logEntry)
	return true
}

func (a *AccountImpl) Operations() []Operation {
	return a.log
}

const newLine = "\n"

func (a *AccountImpl) Statement() string {
	builder := strings.Builder{}
	for _, op := range a.log {
		builder.WriteString(op.String() + newLine)
	}
	result := builder.String()
	return strings.TrimRight(result, newLine)
}

func (a *AccountImpl) Balance() int {
	return a.amount
}
