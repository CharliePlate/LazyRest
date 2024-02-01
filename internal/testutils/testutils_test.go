package testutils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare_IdenticalJSON(t *testing.T) {
	got := bytes.NewReader([]byte(`{"name": "John", "age": 30}`))
	want := bytes.NewReader([]byte(`{"name": "John", "age": 30}`))

	testJson := NewTestJson(got, want)
	err := testJson.Compare()

	assert.NoError(t, err, "Compare should not produce an error for identical JSON")
}

func TestCompare_DifferentJSON(t *testing.T) {
	got := bytes.NewReader([]byte(`{"name": "John", "age": 30}`))
	want := bytes.NewReader([]byte(`{"name": "Jane", "age": 25}`))

	testJson := NewTestJson(got, want)
	err := testJson.Compare()

	assert.Error(t, err, "Compare should produce an error for different JSON")
}

func TestCompare_InvalidGotJSON(t *testing.T) {
	got := bytes.NewReader([]byte(`{"name": "John", "age": `)) // Invalid JSON
	want := bytes.NewReader([]byte(`{"name": "John", "age": 30}`))

	testJson := NewTestJson(got, want)
	err := testJson.Compare()

	assert.Error(t, err, "Compare should produce an error for invalid 'got' JSON")
}

func TestCompare_InvalidWantJSON(t *testing.T) {
	got := bytes.NewReader([]byte(`{"name": "John", "age": 30}`))
	want := bytes.NewReader([]byte(`{"name": "John", "age": `)) // Invalid JSON

	testJson := NewTestJson(got, want)
	err := testJson.Compare()

	assert.Error(t, err, "Compare should produce an error for invalid 'want' JSON")
}
