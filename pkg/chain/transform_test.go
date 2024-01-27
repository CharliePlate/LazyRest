package chain

import (
	"bytes"
	"chain/internal/testutils"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGojqTransformReturnsJson(t *testing.T) {
	testJson := bytes.NewBufferString(`{"foo": "bar", "baz": "qux"}`)

	transformer := NewGojqTransformer()
	transformer.Query = "{test: .foo, test2: .baz, arbitary: \"waldo\"}"
	ctx := context.Background()

	transformed, err := transformer.Transform(ctx, testJson)
	if err != nil {
		t.Errorf("Error transforming JSON: %s", err)
	}

	compare := testutils.NewTestJson(testutils.ReaderToString(transformed), `{"test":"bar", "test2":"qux", "arbitary":"waldo"}`)
	err = compare.Compare()
	assert.NoError(t, err)
}

func TestNilTransformerReturnsOriginalObject(t *testing.T) {
	testJson := bytes.NewBufferString(`{"foo": "bar", "baz": "qux"}`)

	transformer := NewNilTransformer()

	transformed, err := transformer.Transform(context.Background(), testJson)
	if err != nil {
		t.Errorf("Error transforming JSON: %s", err)
	}

	compare := testutils.NewTestJson(testutils.ReaderToString(transformed), `{"foo":"bar", "baz":"qux"}`)
	err = compare.Compare()
	assert.NoError(t, err)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
