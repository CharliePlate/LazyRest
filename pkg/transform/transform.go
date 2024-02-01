package transform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"

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

func NewGojqTransformer(query string) *GojqTransformer {
	return &GojqTransformer{
		Query: query,
	}
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

	var buffer bytes.Buffer
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if buffer.Len() > 0 {
			return nil, errors.New("gojq query returned more than one result")
		}

		if err, ok := v.(error); ok {
			return nil, err
		}

		// handle strings because jq will return it surrounded in ""
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		if jsonData[0] == '"' && jsonData[len(jsonData)-1] == '"' {
			var resultString string
			if err := json.Unmarshal(jsonData, &resultString); err != nil {
				return nil, err
			}
			buffer.WriteString(resultString)
		} else {
			buffer.Write(jsonData)
		}
	}

	return &buffer, nil
}

type NilTransformer struct{}

func NewNilTransformer() *NilTransformer {
	return &NilTransformer{}
}

func (nt *NilTransformer) Transform(ctx context.Context, body io.Reader) (io.Reader, error) {
	return body, nil
}

type CustomScriptTransformer struct {
	Script string
	Args   []string
}

func NewCustomScriptTransformer(script string, args ...string) *CustomScriptTransformer {
	return &CustomScriptTransformer{Script: script, Args: args}
}

func (cs *CustomScriptTransformer) setupCommand(ctx context.Context) (*exec.Cmd, io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, cs.Script, cs.Args...)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	stderrPipe, err := cmd.StderrPipe() // Capture stderr
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return cmd, stdinPipe, stdoutPipe, stderrPipe, nil
}

func (cs *CustomScriptTransformer) Transform(ctx context.Context, body io.Reader) (io.Reader, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	cmd, stdinPipe, stdoutPipe, stderrPipe, err := cs.setupCommand(ctx)
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	_, err = stdinPipe.Write(b)
	if err != nil {
		return nil, err
	}
	stdinPipe.Close()

	stdout, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return nil, err
	}
	stdoutPipe.Close()

	stderr, err := io.ReadAll(stderrPipe)
	if err != nil {
		return nil, err
	}
	stderrPipe.Close()

	if len(stderr) > 0 {
		return nil, fmt.Errorf("command stderr: %s", string(stderr))
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return bytes.NewReader(stdout), nil
}
