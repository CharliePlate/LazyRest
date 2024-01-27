package chain

import (
	"bytes"
	"chain/internal/testutils"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGojqTransformer_TransformReturnsJson(t *testing.T) {
	content := `{"foo": "bar", "baz": "qux"}`
	testJson := bytes.NewBufferString(content)

	transformer := NewGojqTransformer()
	transformer.Query = "{test: .foo, test2: .baz, arbitary: \"waldo\"}"
	ctx := context.Background()

	got, err := transformer.Transform(ctx, testJson)
	assert.NoError(t, err)

	want := `{"test": "bar", "test2": "qux", "arbitary": "waldo"}`
	err = testutils.NewTestJson(got, bytes.NewBufferString(want)).Compare()
	assert.NoError(t, err)
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
