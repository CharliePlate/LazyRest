package chain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/itchyny/gojq"
)

type Transformer interface {
	Transform(ctx context.Context, body io.Reader) (io.Reader, error)
}

type TransformerType string

const (
	NIL_TRANSFORMER  TransformerType = "nil"
	GOJQ_TRANSFORMER TransformerType = "gojq"
)

type GojqTransformer struct {
	Payload *io.Reader
	Query   string
}

func NewGojqTransformer() *GojqTransformer {
	return &GojqTransformer{}
}

// Transform executes a gojq query on the given JSON input and ensures only a single result is produced.
// It returns an error if the query produces more than one result or if any other processing error occurs.
func (jq *GojqTransformer) Transform(ctx context.Context, body io.Reader) (io.Reader, error) {
	query, err := gojq.Parse(jq.Query)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	if err := json.NewDecoder(body).Decode(&jsonData); err != nil {
		return nil, err
	}

	iter := query.RunWithContext(ctx, jsonData)

	var output bytes.Buffer
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if output.Len() > 0 {
			return nil, errors.New("gojq query returned more than one result")
		}

		if err, ok := v.(error); ok {
			return nil, err
		}

		jsonOutput, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		output.Write(jsonOutput)
	}

	return &output, nil
}

type NilTransformer struct{}

func NewNilTransformer() *NilTransformer {
	return &NilTransformer{}
}

func (nt *NilTransformer) Transform(ctx context.Context, body io.Reader) (io.Reader, error) {
	return body, nil
}
