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

func TestGojqTransformReturnsJson(t *testing.T) {
	testJson := bytes.NewBufferString(`{"foo": "bar", "baz": "qux"}`)

	transformer := NewGojq()
	transformer.Query = "{test: .foo, test2: .baz, arbitary: \"waldo\"}"
	ctx := context.Background()

	output, err := transformer.Transform(ctx, testJson)
	if err != nil {
		t.Errorf("Error transforming JSON: %s", err)
	}

	outputStr, err := io.ReadAll(output)
	if err != nil {
		t.Errorf("Error reading output: %s", err)
	}

	test := testutils.TestJson{
		Actual:   string(outputStr),
		Expected: `{"test":"bar", "test2":"qux", "arbitary":"waldo"}`,
	}

	err = test.Compare()
	assert.NoError(t, err)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
