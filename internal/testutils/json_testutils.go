package testutils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/go-cmp/cmp"
)

type TestJson struct {
	ParsedGot  interface{}
	ParsedWant interface{}
	Got        io.Reader
	Want       io.Reader
}

func NewTestJson(actual, expected io.Reader) *TestJson {
	return &TestJson{
		Got:  actual,
		Want: expected,
	}
}

func (t *TestJson) Compare() error {
	err := t.Parse()
	if err != nil {
		return err
	}

	if diff := cmp.Diff(t.ParsedWant, t.ParsedGot); diff != "" {
		return fmt.Errorf("JSON mismatch (-want +got):\n%s", diff)
	}

	return nil
}

func (t *TestJson) Parse() error {
	err := json.NewDecoder(t.Got).Decode(&t.ParsedGot)
	if err != nil {
		return fmt.Errorf("error parsing got JSON: %w", err)
	}

	err = json.NewDecoder(t.Want).Decode(&t.ParsedWant)
	if err != nil {
		return fmt.Errorf("error parsing want JSON: %w", err)
	}

	return nil
}
