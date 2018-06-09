package inex

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

var _testWords = []string{
	"foo",
	"bar",
	"baz",
	"bam",
	"zoing",
	"zing",
}
var _testWordsMatch = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(_testWords, "|")))

func TestStringMatcher_Match(t *testing.T) {
	for _, word1 := range _testWords {
		for _, word2 := range _testWords {
			assert.Equal(t, word1 == word2, StringMatcher(word1).Match(word2), "%v = (\"%s\" == \"%s\")", word1 == word2, word1, word2)
		}
	}
}

func TestStringsMatcher_Match(t *testing.T) {
	matcher := StringsMatcher(_testWords)
	for _, word := range _testWords {
		assert.True(t, matcher.Match(word), "word \"%s\" of self must match", word)
	}
	assert.False(t, matcher.Match("not"), "not on the list must not match")
	assert.False(t, matcher.Match("unseen"), "unseen on the list must not match")
	assert.False(t, matcher.Match("missing"), "missing on the list must not match")
}

func TestRegexpMatcher_Match(t *testing.T) {
	matcher := RegexpMatcher{_testWordsMatch}
	for _, word := range _testWords {
		assert.True(t, matcher.Match(word), "word \"%s\" of self must match", word)
	}
	for _, word1 := range _testWords[0:2] {
		for _, word2 := range _testWords[0:2] {
			assert.False(t, matcher.Match(word1+word2), "word \"%s\" of self must not match", word1+word2)
			assert.False(t, matcher.Match(word1+" "+word2), "word \"%s %s\" of self must not match", word1, word2)
		}
	}
}

func TestFuncMatcher_Match(t *testing.T) {
	ok := make(map[string]bool)
	for _, word := range _testWords {
		ok[word] = true
	}
	fn := func(word string) bool {
		return ok[word]
	}
	matcher := FuncMatcher(fn)
	for _, word1 := range _testWords {
		assert.True(t, matcher.Match(word1), "\"%s\" must match", word1)
		for _, word2 := range _testWords {
			assert.True(t, matcher.Match(word2), "\"%s\" must match", word2)
			assert.False(t, matcher.Match(word1+word2), "\"%s\" must not match", word1+word2)
			assert.False(t, matcher.Match(word1+" "+word2), "\"%s\" must not match", word1+" "+word2)
		}
	}
}

func TestStringMatcherP(t *testing.T) {
	sm := StringMatcher("foo")
	smp := StringMatcherP("foo")
	assert.Equal(t, "foo", string(sm))
	assert.Equal(t, sm, *smp)
}

func TestAnd(t *testing.T) {
	sm1 := StringsMatcher([]string{"a", "b", "c"})
	sm2 := StringsMatcher([]string{"b", "c", "d"})
	sm3 := StringsMatcher([]string{"c", "b", "e"})
	am := And(sm1, sm2, sm3)
	assert.False(t, am.Match("a"))
	assert.True(t, am.Match("b"))
	assert.True(t, am.Match("c"))
	assert.False(t, am.Match("d"))
	assert.False(t, am.Match("e"))
	assert.False(t, And().Match("a"))
}

func TestOr(t *testing.T) {
	sm1 := StringMatcher("b")
	sm2 := StringMatcher("d")
	sm3 := StringMatcher("f")
	am := Or(sm1, sm2, sm3)
	assert.False(t, am.Match("a"))
	assert.True(t, am.Match("b"))
	assert.False(t, am.Match("c"))
	assert.True(t, am.Match("d"))
	assert.False(t, am.Match("e"))
	assert.True(t, am.Match("f"))
}

func TestIfMost(t *testing.T) {
	sm1 := StringsMatcher{"a", "b", "c", "d"}
	sm2 := StringsMatcher{"c", "d", "e", "f"}
	sm3 := StringsMatcher{"d", "e", "f", "g"}
	im := IfMost(sm1, sm2, sm3)
	assert.False(t, im.Match("a"))
	assert.False(t, im.Match("b"))
	assert.True(t, im.Match("c"))
	assert.True(t, im.Match("d"))
	assert.True(t, im.Match("e"))
	assert.True(t, im.Match("f"))
	assert.False(t, im.Match("g"))
}

func TestNothing(t *testing.T) {
	nm := Nothing()
	assert.False(t, nm.Match("a"))
	assert.False(t, nm.Match("b"))
	assert.False(t, nm.Match("c"))
}

func TestEverything(t *testing.T) {
	nm := Everything()
	assert.True(t, nm.Match("a"))
	assert.True(t, nm.Match("b"))
	assert.True(t, nm.Match("c"))
}