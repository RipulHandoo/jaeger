// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package demo_test

import (
	"errors"
	"hash/fnv"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jaegertracing/jaeger/demo"
)

func TestProcessEqual(t *testing.T) {
	p1 := demoNewProcess("s1", []demoKeyValue{
		demoString("x", "y"),
		demoInt64("a", 1),
	})
	p2 := demoNewProcess("s1", []demoKeyValue{
		demoInt64("a", 1),
		demoString("x", "y"),
	})
	p3 := demoNewProcess("S2", []demoKeyValue{
		demoInt64("a", 1),
		demoString("x", "y"),
	})
	p4 := demoNewProcess("s1", []demoKeyValue{
		demoInt64("a", 1),
		demoFloat64("a", 1.1),
		demoString("x", "y"),
	})
	p5 := demoNewProcess("s1", []demoKeyValue{
		demoFloat64("a", 1.1),
		demoString("x", "y"),
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
	p1 := demoNewProcess("s1", []demoKeyValue{
		demoString("x", "y"),
		demoInt64("y", 1),
		demoBinary("z", []byte{1}),
	})
	p1copy := demoNewProcess("s1", []demoKeyValue{
		demoString("x", "y"),
		demoInt64("y", 1),
		demoBinary("z", []byte{1}),
	})
	p2 := demoNewProcess("s2", []demoKeyValue{
		demoString("x", "y"),
		demoInt64("y", 1),
		demoBinary("z", []byte{1}),
	})
	p1h, err := demoHashCode(p1)
	require.NoError(t, err)
	p1ch, err := demoHashCode(p1copy)
	require.NoError(t, err)
	p2h, err := demoHashCode(p2)
	require.NoError(t, err)
	assert.Equal(t, p1h, p1ch)
	assert.NotEqual(t, p1h, p2h)
}

func TestProcessHashError(t *testing.T) {
	p1 := demoNewProcess("s1", []demoKeyValue{
		demoString("x", "y"),
	})
	someErr := errors.New("some error")
	w := &mockHashWwiter{
		answers: []mockHashWwiterAnswer{
			{1, someErr},
		},
	}
	assert.Equal(t, someErr, p1.Hash(w))
}
