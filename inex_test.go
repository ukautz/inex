package inex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRoot(t *testing.T) {
	root := NewRoot(StringMatcher("foo"))
	assert.NotNil(t, root)
	assert.Nil(t, root.parent, "has no parent")
	assert.Nil(t, root.child, "has no child")
	assert.Equal(t, StringMatcher("foo"), root.matcher, "has no child")
}

func TestNewRoot_Multi(t *testing.T) {
	root := NewRoot(StringsMatcher{"foo", "bar"}, StringsMatcher{"bar", "baz"})
	assert.NotNil(t, root)
	assert.False(t, root.Match("foo"))
	assert.True(t, root.Match("bar"))
	assert.False(t, root.Match("baz"))
}

func TestInex_IsRoot(t *testing.T) {
	root := NewRoot()
	assert.True(t, root.IsRoot())
	child := root.Child(StringMatcher("foo"))
	assert.False(t, child.IsRoot())
}

func TestInex_IsEnd(t *testing.T) {
	root := NewRoot()
	assert.True(t, root.IsEnd())
	child := root.Child(StringMatcher("foo"))
	assert.False(t, root.IsEnd())
	assert.True(t, child.IsEnd())
}

func TestInex_Match(t *testing.T) {
	root := NewRoot(StringsMatcher(_testWords))
	for _, word1 := range _testWords {
		assert.True(t, root.Match(word1), "word \"%s\" is in list of words")
	}
	for _, word1 := range _testWords[0:2] {
		for _, word2 := range _testWords[0:2] {
			assert.False(t, root.Match(word1+word2), "word combinations \"%s%s\" not in list", word1, word2)
		}
	}
}

func TestInex_NewChild(t *testing.T) {
	r1 := true
	r2 := true
	r3 := true
	root := NewRoot(FuncMatcher(func(word string) bool {
		return r1
	}))
	child1 := root.Child(FuncMatcher(func(word string) bool {
		return r2
	}))
	child2 := child1.Child(FuncMatcher(func(word string) bool {
		return r3
	}))

	assert.True(t, root.Match("foo"))
	assert.True(t, child1.Match("foo"))
	assert.True(t, child2.Match("foo"))
	r3 = false
	assert.False(t, root.Match("foo"))
	assert.False(t, child1.Match("foo"))
	assert.False(t, child2.Match("foo"))
	r3 = true
	r2 = false
	assert.False(t, root.Match("foo"))
	assert.False(t, child1.Match("foo"))
	assert.True(t, child2.Match("foo"))
	r2 = true
	r1 = false
	assert.False(t, root.Match("foo"))
	assert.True(t, child1.Match("foo"))
	assert.True(t, child2.Match("foo"))
	r2 = false
	assert.False(t, root.Match("foo"))
	assert.False(t, child1.Match("foo"))
	assert.True(t, child2.Match("foo"))
	r3 = false
	assert.False(t, root.Match("foo"))
	assert.False(t, child1.Match("foo"))
	assert.False(t, child2.Match("foo"))
}

func TestInex_NewChild_Chain(t *testing.T) {
	root := NewRoot(StringsMatcher{"foo", "bar", "baz"})
	assert.True(t, root.Match("foo"))
	assert.True(t, root.Match("bar"))
	assert.True(t, root.Match("baz"))

	child1 := root.Child(StringsMatcher{"foo", "bar"})
	assert.True(t, root.Match("foo"))
	assert.True(t, root.Match("bar"))
	assert.False(t, root.Match("baz"))
	assert.True(t, child1.Match("foo"))
	assert.True(t, child1.Match("bar"))
	assert.False(t, child1.Match("baz"))

	child2 := child1.Child(StringsMatcher{"foo"})
	assert.True(t, root.Match("foo"))
	assert.False(t, root.Match("bar"))
	assert.False(t, root.Match("baz"))
	assert.True(t, child1.Match("foo"))
	assert.False(t, child1.Match("bar"))
	assert.False(t, child1.Match("baz"))
	assert.True(t, child2.Match("foo"))
	assert.False(t, child2.Match("bar"))
	assert.False(t, child2.Match("baz"))

	assert.False(t, root.Match("bla"))
	assert.False(t, child1.Match("bla"))
	assert.False(t, child2.Match("bla"))
}

func TestInex_Root(t *testing.T) {
	end := NewRoot().Child(StringMatcher("foo"))
	assert.NotNil(t, end.parent)
	assert.Nil(t, end.Root().parent)
}

func TestInex_End(t *testing.T) {
	root := NewRoot()
	child := root.Child(StringMatcher("foo"))
	assert.Nil(t, root.parent)
	assert.NotNil(t, child.parent)
	assert.NotNil(t, root.End().parent)
}

func TestInex_Exclude(t *testing.T) {
	root := NewRoot(StringsMatcher{"foo", "bar", "baz"}).
		Exclude(StringMatcherP("bar")).
		Root()
	assert.True(t, root.Match("foo"))
	assert.False(t, root.Match("bar"))
	assert.True(t, root.Match("baz"))
}

func TestInex_Include(t *testing.T) {
	root := NewRoot(StringsMatcher{"foo", "bar", "baz"}).
		Include(StringMatcherP("bar")).
		Root()
	assert.False(t, root.Match("foo"))
	assert.True(t, root.Match("bar"))
	assert.False(t, root.Match("baz"))
}
