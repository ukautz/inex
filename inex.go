// Package inex provides a framework and glue code to create and execute arbitrary deep boolean filter.
package inex

type (

	// Inex is a node in a doubly linked list. It provides an include, exclude builder
	// to compose arbitrary deep filter expressions.
	Inex struct {
		matcher Matcher
		parent  *Inex
		child   *Inex
	}
)

// NewRoot returns a new root Inex instance matching the provided Matcher.
func NewRoot(matcher ...Matcher) *Inex {
	var m Matcher
	if l := len(matcher); l == 1 {
		m = matcher[0]
	} else if l > 1 {
		m = And(matcher...)
	}
	return NewInex(m, nil)
}

// NewInex returns a new Inex instance, which has the provided Inex instance set as parent.
func NewInex(matcher Matcher, parent *Inex) *Inex {
	return &Inex{
		matcher: matcher,
		parent:  parent,
	}
}

// Child set new child from matcher with current instance as parent.
// WILL overwrite existing child.
func (i *Inex) Child(matcher Matcher) *Inex {
	i.child = NewInex(matcher, i)
	return i.child
}

// Exclude sets new child with inverted result of given matcher. Same as `instance.Child(inex.Not(matcher))`.
// Think `find . ! -name "*.go"`.
// WILL overwrite existing child.
// Returns pointer to newly created child so that follow up include can be chained.
func (i *Inex) Exclude(matcher Matcher) *Inex {
	return i.Child(Not(matcher))
}

// Include sets new child with result of given matcher. Same as `instance.Child(matcher)`.
// Think `find . -name "*.go"`
// WILL overwrite existing child.
// Consider `inex.All` instead executing on a node which was also create using Include.
// Returns pointer to newly created child so that follow up exclude can be chained.
func (i *Inex) Include(matcher Matcher) *Inex {
	return i.Child(matcher)
}

// IsRoot returns bool whether this node is root
func (i *Inex) IsRoot() bool {
	return i.parent == nil
}

// Root returns the upper most parent (which has no parent) in the (doubly linked) list.
func (i *Inex) Root() *Inex {
	if i.parent == nil {
		return i
	}
	return i.parent.Root()
}

// IsEnd returns bool whether this node is last (child)
func (i *Inex) IsEnd() bool {
	return i.child == nil
}

// End returns the bottom most child (which has no child) in the (doubly linked) list. Last element in the include / exclude chain
func (i *Inex) End() *Inex {
	if i.child == nil {
		return i
	}
	return i.child.End()
}

// Match returns bool if given word is matched by this instance and all following includes and excludes (childs).
// Best called as `instance.Root().Match(word)`.
func (i *Inex) Match(str string) (res bool) {
	if i.matcher != nil {
		res = i.matcher.Match(str)
	}
	if i.child != nil && (res || i.matcher == nil) {
		return i.child.Match(str)
	}
	return res
}
