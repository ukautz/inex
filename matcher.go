package inex

import (
	"regexp"
)

type (

	// Matcher is a boolean filter, either accepting or not accepting given input string
	Matcher interface {

		// Match returns true if this Matcher accepts the input string
		Match(str string) bool
	}

	// StringMatcher matches only a single string
	//	m := inex.StringMatcher("foo")
	//	m.Match("foo") // == true
	//	m.Match("bar") // == false
	// Consider also StringMatcherP
	StringMatcher  string

	// StringsMatcher matches a slice of strings
	//	m := inex.StringsMatcher([]string{"foo", "bar"})
	//	m.Match("foo") // == true
	//	m.Match("bar") // == true
	//	m.Match("baz") // == false
	StringsMatcher []string

	// RegexpMatcher matches a regular expression
	//	m := inex.RegexpMatcher(regexp.MustCompile(`^(foo|bar)$)
	//	m.Match("foo") // == true
	//	m.Match("bar") // == true
	//	m.Match("baz") // == false
	// Consider also NewRegexpMatcher
	RegexpMatcher  struct{ *regexp.Regexp }

	// FuncMatcher is function which returns bool for given string
	FuncMatcher    func(word string) bool
)

// Match returns true if input equals self
func (m StringMatcher) Match(str string) bool {
	return string(m) == str
}

// NewStringMatcher create new *StringMatcher instance from given string
func NewStringMatcher(str string) *StringMatcher {
	m := StringMatcher(str)
	return &m
}

// StringMatcherP is an alias for NewStringMatcher
func StringMatcherP(str string) *StringMatcher {
	return NewStringMatcher(str)
}

// Match returns true if any of self equal the given string
func (m StringsMatcher) Match(str string) bool {
	for _, s := range m {
		if s == str {
			return true
		}
	}
	return false
}

// Match returns true if self regular expression matches given string
func (m *RegexpMatcher) Match(str string) bool {
	return m.MatchString(str)
}

// Match returns true if call of self with given string returns true
func (m FuncMatcher) Match(str string) bool {
	return m(str)
}

// And returns Matcher which matches only if ALL provided Matchers match.
func And(matchers ...Matcher) Matcher {
	empty := len(matchers) == 0
	return FuncMatcher(func(word string) bool {
		if empty {
			return false
		}
		for _, matcher := range matchers {
			if !matcher.Match(word) {
				return false
			}
		}
		return true
	})
}

// Or returns Matcher witch matches if ANY provided Matcher matches.
func Or(matchers ...Matcher) Matcher {
	return FuncMatcher(func(word string) bool {
		for _, matcher := range matchers {
			if matcher.Match(word) {
				return true
			}
		}
		return false
	})
}

// IfMost returns Matcher which matches if the absolute majority of provided Matchers match.
func IfMost(matchers ...Matcher) Matcher {
	return FuncMatcher(func(word string) bool {
		matched := 0
		for _, matcher := range matchers {
			if matcher.Match(word) {
				matched++
			}
		}
		return matched > len(matchers)-matched
	})
}

// Nothing returns a Matcher which NEVER matches to ANYTHING.
func Nothing() Matcher {
	return FuncMatcher(func(word string) bool {
		return false
	})
}

// Everything returns a Matcher which ALWAYS matches to ANYTHING.
func Everything() Matcher {
	return FuncMatcher(func(word string) bool {
		return true
	})
}

// Not returns inverted Match result of the provided Matcher
func Not(matcher Matcher) Matcher {
	return FuncMatcher(func(word string) bool {
		return !matcher.Match(word)
	})
}
