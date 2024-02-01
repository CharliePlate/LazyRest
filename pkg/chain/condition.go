package chain

type Condition interface {
	Evaluate() (bool, error)
}

type TodoCondition struct {
	Result bool
}

func (t *TodoCondition) Evaluate() (bool, error) {
	return t.Result, nil
}
