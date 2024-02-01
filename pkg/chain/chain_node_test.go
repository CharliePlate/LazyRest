package chain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// Mock implementation of the Condition interface
type TestCondition struct {
	Error       error
	ReturnValue bool
}

func (tc *TestCondition) Evaluate() (bool, error) {
	return tc.ReturnValue, tc.Error
}

type TestChain struct {
	Name string
}

func (tch *TestChain) Execute() (Result, error) {
	return Result{}, nil
}

func TestConditionalTraversal_FirstSuccessGetsParsed(t *testing.T) {
	rootChildren := []*ChainNode{
		{
			Chain: &TestChain{
				Name: "Success_Result_First",
			},
			Link: &TestCondition{
				ReturnValue: true,
			},
		},
		{
			Chain: &TestChain{
				Name: "Success_Result_Second",
			},
			Link: &TestCondition{
				ReturnValue: true,
			},
		},
	}

	root := &ChainNode{
		Chain: &TestChain{
			Name: "Root",
		},
		Options: &TraversalOptions{
			ParallelEvaluation: false,
			ContinueOnFailure:  false,
		},
		Childern: rootChildren,
	}

	next, err := root.Next()

	require.NoError(t, err)
	chainName := next.Chain.(*TestChain).Name
	require.Equal(t, "Success_Result_First", chainName)

	(*root).Childern[0].Link = &TestCondition{
		ReturnValue: false,
	}

	newNext, err := root.Next()
	require.NoError(t, err)

	newChainName := newNext.Chain.(*TestChain).Name
	require.Equal(t, "Success_Result_Second", newChainName)
}

func TestConditionalTraversal_StopOnErrorIfNotContinueOnFailure(t *testing.T) {
	children := []*ChainNode{
		{
			Chain: &TestChain{
				Name: "I Will Throw An Error",
			},
			Link: &TestCondition{
				Error: errors.New("im throwing!"),
			},
		},
	}

	root := &ChainNode{
		Chain: &TestChain{
			Name: "Root",
		},
		Childern: children,
		Options: &TraversalOptions{
			ParallelEvaluation: false,
			ContinueOnFailure:  false,
		},
	}

	next, err := root.Next()
	require.Error(t, err)

	chainName := next.Chain.(*TestChain).Name
	nodeError := next.Childern[0].Link.(*TestCondition).Error
	require.Equal(t, "Root", chainName)
	require.Equal(t, "im throwing!", nodeError.Error())
}
