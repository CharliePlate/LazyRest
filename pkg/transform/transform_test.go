package transform

import (
	"bytes"
	"chain/internal/testutils"
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type TransformTester struct {
	Transformer           Transformer
	Content               string
	Want                  string
	ShouldFailOnTransform bool
	ShouldFailOnCompare   bool
}

func TestGojqTransformer_TransformReturnsJson(t *testing.T) {
	tester := TransformTester{
		Content:     `{"foo": "bar", "baz": "qux"}`,
		Transformer: NewGojqTransformer(`{test: .foo, test2: .baz, arbitary: "waldo"}`),
		Want:        `{"test": "bar", "test2": "qux", "arbitary": "waldo"}`,
	}

	tester.Test(t)
}

func TestGojqTransformer_TransformerReturnsSingleString(t *testing.T) {
	tester := TransformTester{
		Content:     `{"foo": "bar"}`,
		Transformer: NewGojqTransformer(`.foo`),
		Want:        `bar`,
	}

	result, err := tester.Transformer.Transform(context.Background(), bytes.NewBufferString(tester.Content))
	require.NoError(t, err)

	got, err := io.ReadAll(result)
	require.NoError(t, err)
	want := tester.Want
	require.Equal(t, want, string(got))
}

func TestNilTransformer_TransformReturnsOriginalObject(t *testing.T) {
	tester := TransformTester{
		Content:     `{"foo": "bar", "baz": "qux"}`,
		Transformer: NewNilTransformer(),
		Want:        `{"foo": "bar", "baz": "qux"}`,
	}

	tester.Test(t)
}

func TestCustomScriptTransformer_DataInAndOut(t *testing.T) {
	content := `{"foo": "bar", "baz": "qux"}`
	testJson := bytes.NewBufferString(content)

	transformer := NewCustomScriptTransformer("jq", `{test: .foo, test2: .baz, arbitary: "waldo"}`)
	got, err := transformer.Transform(context.Background(), testJson)
	require.NoError(t, err)

	want := bytes.NewBufferString(`{"test": "bar", "test2": "qux", "arbitary": "waldo"}`)
	err = testutils.NewTestJson(got, want).Compare()
	require.NoError(t, err)
}

func TestTester(t *testing.T) {
	tester := TransformTester{
		Content:     `{"foo": "bar", "baz": "qux"}`,
		Transformer: NewNilTransformer(),
		Want:        `{"foo": "bar", "baz": "qux"}`,
	}
	tester.Test(t)

	tester1 := TransformTester{
		Content:             `{"foo": "bar", "baz": "qux"}`,
		Transformer:         NewNilTransformer(),
		Want:                `{"foo": "qux", "baz": "bar"}`,
		ShouldFailOnCompare: true,
	}
	tester1.Test(t)
}

func (tt *TransformTester) Test(t *testing.T) {
	testJson := bytes.NewBufferString(tt.Content)
	got, err := tt.Transformer.Transform(context.Background(), testJson)
	if tt.ShouldFailOnTransform {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)

	want := bytes.NewBufferString(tt.Want)
	err = testutils.NewTestJson(got, want).Compare()
	if tt.ShouldFailOnCompare {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
