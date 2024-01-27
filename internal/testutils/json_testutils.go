package testutils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type TestJson struct {
	parsedActual   interface{}
	parsedExpected interface{}
	Actual         string
	Expected       string
}

func NewTestJson(actual, expected string) *TestJson {
	return &TestJson{
		Actual:   actual,
		Expected: expected,
	}
}

func (t *TestJson) Parse() error {
	if err := json.Unmarshal([]byte(t.Actual), &t.parsedActual); err != nil {
		return fmt.Errorf("error parsing actual JSON: %w", err)
	}

	if err := json.Unmarshal([]byte(t.Expected), &t.parsedExpected); err != nil {
		return fmt.Errorf("error parsing expected JSON: %w", err)
	}

	return nil
}

func (t *TestJson) Compare() error {
	if err := t.Parse(); err != nil {
		return err
	}

	if !jsonEqual(t.parsedActual, t.parsedExpected) {
		pretty, err := t.prettyPrint()
		if err != nil {
			pretty = err.Error()
		}

		return fmt.Errorf("actual JSON does not match expected JSON. %s", pretty)
	}

	return nil
}

func (t *TestJson) prettyPrint() (string, error) {
	if t.Actual == "" || t.Expected == "" {
		err := t.Parse()
		if err != nil {
			return "", fmt.Errorf("error parsing JSON: %w", err)
		}
	}

	actualBytes, err := json.MarshalIndent(t.parsedActual, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling actual JSON: %w", err)
	}

	expectedBytes, err := json.MarshalIndent(t.parsedExpected, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling expected JSON: %w", err)
	}

	return fmt.Sprintf("Actual:\n%s\nExpected:\n%s", string(actualBytes), string(expectedBytes)), nil
}

func jsonEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
