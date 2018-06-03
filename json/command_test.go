package json

import (
	"errors"
	"testing"

	"github.com/caalberts/localroast"
	"github.com/stretchr/testify/assert"
)

type mockReader struct {
	result []byte
	called bool
	args   []string
	err    error
}

func (r *mockReader) Read(args []string) ([]byte, error) {
	r.called = true
	r.args = args
	return r.result, r.err
}

type mockParser struct {
	result []localroast.Schema
	called bool
	args   []byte
	err    error
}

func (p *mockParser) Parse(bytes []byte) ([]localroast.Schema, error) {
	p.called = true
	p.args = bytes
	return p.result, p.err
}

func TestExecuteJSONCommand(t *testing.T) {
	r := &mockReader{result: []byte("content")}
	p := &mockParser{result: []localroast.Schema{}}

	cmd := &Command{r, p}
	args := []string{"filename"}
	schema, err := cmd.Execute(args)

	assert.Nil(t, err)
	assert.Equal(t, p.result, schema)

	assert.True(t, r.called)
	assert.Equal(t, args, r.args)

	assert.True(t, p.called)
	assert.Equal(t, r.result, p.args)
}

func TestReadError(t *testing.T) {
	r := &mockReader{err: errors.New("Failed to read file")}
	p := &mockParser{}

	cmd := &Command{r, p}
	args := []string{"fakefile"}
	schema, err := cmd.Execute(args)

	assert.NotNil(t, err)
	assert.Equal(t, "Failed to read file", err.Error())
	assert.Nil(t, schema)

	assert.True(t, r.called)
	assert.False(t, p.called)
}

func TestParseError(t *testing.T) {
	r := &mockReader{result: []byte("content")}
	p := &mockParser{err: errors.New("Failed to parse schema")}

	cmd := &Command{r, p}
	args := []string{"filename"}
	schema, err := cmd.Execute(args)

	assert.NotNil(t, err)
	assert.Equal(t, "Failed to parse schema", err.Error())
	assert.Nil(t, schema)

	assert.True(t, r.called)
	assert.True(t, p.called)
}