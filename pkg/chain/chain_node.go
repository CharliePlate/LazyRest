package chain

import (
	"errors"
	"sync"
)

type ChainNode struct {
	EvalError  error
	Chain      Chain
	Link       Condition
	Default    *ChainNode
	Options    *TraversalOptions
	Childern   []*ChainNode
	EvalResult bool
}

type TraversalOptions struct {
	ParallelEvaluation bool
	ContinueOnFailure  bool
}

func NewChainNode(chain Chain, link Condition, children []*ChainNode, options *TraversalOptions) *ChainNode {
	return &ChainNode{
		Chain:    chain,
		Link:     link,
		Options:  options,
		Childern: children,
	}
}

func (cn *ChainNode) parallelEvaluation() (*ChainNode, error) {
	wg := sync.WaitGroup{}
	for _, child := range cn.Childern {
		wg.Add(1)
		go func(c *ChainNode) {
			defer wg.Done()
			c.EvalResult, c.EvalError = c.Link.Evaluate()
		}(child)
	}

	wg.Wait()

	for _, child := range cn.Childern {
		if child.EvalError != nil {
			if !cn.Options.ContinueOnFailure {
				return cn, errors.New("one of the parsed expressions failed before a success was found")
			}
		}
		if child.EvalResult {
			return child, nil
		}
	}

	return cn, nil
}

func (cn *ChainNode) syncEvaluation() (*ChainNode, error) {
	for _, child := range cn.Childern {
		child.EvalResult, child.EvalError = child.Link.Evaluate()
		if child.EvalError != nil {
			if !cn.Options.ContinueOnFailure {
				return cn, errors.New("one of the parsed expressions failed before a success was found")
			}
		}
		if child.EvalResult {
			return child, nil
		}
	}

	cn.EvalError = errors.New("none of the parsed expressions evaluated to true")
	return cn, nil
}

func (cn *ChainNode) Next() (*ChainNode, error) {
	var next *ChainNode
	var err error
	if cn.Options.ParallelEvaluation {
		next, err = cn.parallelEvaluation()
	} else {
		next, err = cn.syncEvaluation()
	}

	if err != nil {
		return cn, err
	}

	return next, nil
}
