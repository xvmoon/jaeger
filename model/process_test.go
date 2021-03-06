// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package model_test

import (
	"hash/fnv"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"errors"

	"github.com/stretchr/testify/require"
	"github.com/uber/jaeger/model"
)

func TestProcessEqual(t *testing.T) {
	p1 := model.NewProcess("s1", []model.KeyValue{
		model.String("x", "y"),
		model.Int64("a", 1),
	})
	p2 := model.NewProcess("s1", []model.KeyValue{
		model.Int64("a", 1),
		model.String("x", "y"),
	})
	p3 := model.NewProcess("S2", []model.KeyValue{
		model.Int64("a", 1),
		model.String("x", "y"),
	})
	p4 := model.NewProcess("s1", []model.KeyValue{
		model.Int64("a", 1),
		model.Float64("a", 1.1),
		model.String("x", "y"),
	})
	p5 := model.NewProcess("s1", []model.KeyValue{
		model.Float64("a", 1.1),
		model.String("x", "y"),
	})
	assert.Equal(t, p1, p2)
	assert.True(t, p1.Equal(p2))
	assert.False(t, p1.Equal(p3))
	assert.False(t, p1.Equal(p4))
	assert.False(t, p1.Equal(p5))
}

func Hash(w io.Writer) {
	w.Write([]byte("hello"))
}

func TestX(t *testing.T) {
	h := fnv.New64a()
	Hash(h)
}

func TestProcessHash(t *testing.T) {
	p1 := model.NewProcess("s1", []model.KeyValue{
		model.String("x", "y"),
		model.Int64("y", 1),
		model.Binary("z", []byte{1}),
	})
	p1copy := model.NewProcess("s1", []model.KeyValue{
		model.String("x", "y"),
		model.Int64("y", 1),
		model.Binary("z", []byte{1}),
	})
	p2 := model.NewProcess("s2", []model.KeyValue{
		model.String("x", "y"),
		model.Int64("y", 1),
		model.Binary("z", []byte{1}),
	})
	p1h, err := model.HashCode(p1)
	require.NoError(t, err)
	p1ch, err := model.HashCode(p1copy)
	require.NoError(t, err)
	p2h, err := model.HashCode(p2)
	require.NoError(t, err)
	assert.Equal(t, p1h, p1ch)
	assert.NotEqual(t, p1h, p2h)
}

func TestProcessHashError(t *testing.T) {
	p1 := model.NewProcess("s1", []model.KeyValue{
		model.String("x", "y"),
	})
	someErr := errors.New("some error")
	w := &mockHashWwiter{
		answers: []mockHashWwiterAnswer{
			{1, someErr},
		},
	}
	assert.Equal(t, someErr, p1.Hash(w))
}
