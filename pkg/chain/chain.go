package chain

type Result struct{}

type Chain interface {
	Execute() (Result, error)
}
