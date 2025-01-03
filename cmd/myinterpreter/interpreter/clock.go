package interpreter

import "time"

type Clock struct{}

func (c Clock) arity() int {
	return 0
}

func (c Clock) call(_interpreter *Interpreter, _arguments []interface{}) interface{} {
	return float64(time.Now().UnixMilli() / 1000)
}

func (c Clock) String() string {
	return "<native fn>"
}
