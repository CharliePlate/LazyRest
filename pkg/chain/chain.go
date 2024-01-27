package chain

import (
	"net/http"
	"sync"
)

type ResponseWithError struct {
	Response *http.Response
	Error    error
}

type ChainRequest struct {
	*http.Request
	Transformer  *Transformer
	Previous     *ResponseWithError
	PipePrevious bool
	Index        int
}

type Chain struct {
	Requests  []*ChainRequest
	Responses []*ResponseWithError
}

func NewChain() Chain {
	return Chain{}
}

func (c *Chain) Add(r *http.Request, t *Transformer) {
	c.Requests = append(c.Requests, &ChainRequest{Request: r, Transformer: t})
}

func (c *Chain) AppendTransformer(i int, t *Transformer) {
	c.Requests[i].Transformer = t
}

func (c *Chain) Remove(i int) {
	c.Requests = append(c.Requests[:i], c.Requests[i+1:]...)
}

func (c *Chain) RunConcurrently() {
	wg := sync.WaitGroup{}

	c.Responses = make([]*ResponseWithError, len(c.Requests))

	for i, req := range c.Requests {
		wg.Add(1)
		go func(i int, r *http.Request) {
			defer wg.Done()
			resp, err := http.DefaultClient.Do(r)
			c.Responses[i] = &ResponseWithError{Response: resp, Error: err}
		}(i, req.Request)
	}
	wg.Wait()
}

func (c *Chain) RunSequentially() {
	c.Responses = make([]*ResponseWithError, len(c.Requests))
	for i, req := range c.Requests {
		resp, err := http.DefaultClient.Do(req.Request)
		c.Responses[i] = &ResponseWithError{Response: resp, Error: err}
	}
}
