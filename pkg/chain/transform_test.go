package chain

import (
	"bytes"
	"chain/internal/testutils"
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGojqTransformer_TransformReturnsJson(t *testing.T) {
	content := `{"foo": "bar", "baz": "qux"}`
	testJson := bytes.NewBufferString(content)

	transformer := NewGojqTransformer("{test: .foo, test2: .baz, arbitary: \"waldo\"}")
	got, err := transformer.Transform(context.Background(), testJson)
	assert.NoError(t, err)

	want := bytes.NewBufferString(`{"test": "bar", "test2": "qux", "arbitary": "waldo"}`)
	err = testutils.NewTestJson(got, want).Compare()
	assert.NoError(t, err)
}

func TestGojqTransformer_TransformerPassedSingleValues(t *testing.T) {
	testJson := bytes.NewBufferString(`{"foo": "bar"}`)

	transformer := NewGojqTransformer(".foo")
	body, err := transformer.Transform(context.Background(), testJson)
	assert.NoError(t, err)

	got, err := io.ReadAll(body)
	assert.NoError(t, err)

	want := "bar"
	assert.Equal(t, want, string(got))
}

func TestNilTransformer_TransformReturnsOriginalObject(t *testing.T) {
	content := `{"foo": "bar", "baz": "qux"}`
	testJson := bytes.NewBufferString(content)

	transformer := NewNilTransformer()

	got, err := transformer.Transform(context.Background(), testJson)
	assert.NoError(t, err)

	want := content
	err = testutils.NewTestJson(got, bytes.NewBufferString(want)).Compare()
	assert.NoError(t, err)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
