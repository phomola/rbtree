package rbtree

import "fmt"

// Comparable is a totally orderable type.
type Comparable interface {
	Compare(interface{}) int
}

type color byte

const (
	black color = iota
	red
)

type direction byte

const (
	exact direction = iota
	left
	right
)

type node struct {
	key    Comparable
	value  interface{}
	color  color
	parent *node
	tree   *Tree
	left   *node
	right  *node
}

func (n *node) check() bool {
	if n.left != nil {
		if n.left.key.Compare(n.key) >= 0 {
			return false
		}
		if n.left.parent != n {
			return false
		}
		if !n.left.check() {
			return false
		}
	}
	if n.right != nil {
		if n.key.Compare(n.right.key) >= 0 {
			return false
		}
		if n.right.parent != n {
			return false
		}
		if !n.right.check() {
			return false
		}
	}
	return true
}

func (n *node) depth() int {
	var ld, rd int
	if n.left != nil {
		ld = n.left.depth()
	}
	if n.right != nil {
		rd = n.right.depth()
	}
	if ld > rd {
		return ld + 1
	} else {
		return rd + 1
	}
}

func (n *node) size() int {
	r := 1
	if n.left != nil {
		r += n.left.size()
	}
	if n.right != nil {
		r += n.right.size()
	}
	return r
}

func (n *node) keys() []Comparable {
	var ks []Comparable
	if n.left != nil {
		ks = n.left.keys()
	}
	ks = append(ks, n.key)
	if n.right != nil {
		ks = append(ks, n.right.keys()...)
	}
	return ks
}

func (n *node) str() string {
	var s string
	if n.left != nil {
		s += "(" + n.left.str() + ") "
	}
	s += fmt.Sprintf("%v:%v", n.key, n.value)
	if n.color == black {
		s += "/B"
	} else {
		s += "/R"
	}
	if n.right != nil {
		s += " (" + n.right.str() + ")"
	}
	return s
}

func (n *node) find(key Comparable) (*node, direction) {
	c := key.Compare(n.key)
	switch {
	case c == 0:
		return n, exact
	case c < 0:
		if n.left == nil {
			return n, left
		} else {
			return n.left.find(key)
		}
	case c > 0:
		if n.right == nil {
			return n, right
		} else {
			return n.right.find(key)
		}
	}
	panic("bad red-black node")
}

func (n *node) rotateRight() {
	p := n.parent
	pp := p.parent
	a, b, c := n.left, n.right, p.right
	if pp != nil {
		switch p.dir() {
		case left:
			pp.left = n
		case right:
			pp.right = n
		default:
			panic("bad red-black node")
		}
	} else {
		n.tree.root = n
	}
	n.parent, p.parent = pp, n
	n.left, n.right = a, p
	p.left, p.right = b, c
	if a != nil {
		a.parent = n
	}
	if b != nil {
		b.parent = p
	}
	if c != nil {
		c.parent = p
	}
}

func (n *node) rotateLeft() {
	p := n.parent
	pp := p.parent
	a, b, c := p.left, n.left, n.right
	if pp != nil {
		switch p.dir() {
		case left:
			pp.left = n
		case right:
			pp.right = n
		default:
			panic("bad red-black node")
		}
	} else {
		n.tree.root = n
	}
	n.parent, p.parent = pp, n
	n.left, n.right = p, c
	p.left, p.right = a, b
	if c != nil {
		c.parent = n
	}
	if a != nil {
		a.parent = p
	}
	if b != nil {
		b.parent = p
	}
}

func (n *node) rotate() {
	switch n.dir() {
	case right:
		n.rotateLeft()
	case left:
		n.rotateRight()
	}
}

func (n *node) dir() direction {
	p := n.parent
	switch {
	case p.left == n:
		return left
	case p.right == n:
		return right
	}
	panic("bad red-black node")
}

func (n *node) brother() *node {
	p := n.parent
	switch {
	case p.left == n:
		return p.right
	case p.right == n:
		return p.left
	}
	panic("bad red-black node")
}

func (n *node) ensureInvariants() {
	p := n.parent
	if p == nil {
		n.color = black
		return
	}
	if p.color == black {
		return
	}
	pp := p.parent
	if pp != nil && pp.color == black {
		u := p.brother()
		if u != nil && u.color == red {
			p.color, pp.color, u.color = black, red, black
			pp.ensureInvariants()
		} else {
			if n.dir() == p.dir() {
				p.rotate()
				p.color, pp.color = black, red
			} else {
				n.rotate()
				n.rotate()
				n.color, pp.color = black, red
			}
		}
	}
}

// Tree is a generic red-black tree.
type Tree struct {
	root *node
}

// New creates a new red-black tree.
func New() *Tree { return new(Tree) }

// Depth returns the depth of the tree.
func (t *Tree) Depth() int {
	if t.root == nil {
		return 0
	}
	return t.root.depth()
}

// Size returns the size of the tree.
func (t *Tree) Size() int {
	if t.root == nil {
		return 0
	}
	return t.root.size()
}

// Keys returns the keys of the items in the tree.
func (t *Tree) Keys() []Comparable {
	if t.root == nil {
		return nil
	}
	return t.root.keys()
}

// Insert inserts a new key-value pair into the tree or replaces the value for an existing key.
func (t *Tree) Insert(key Comparable, value interface{}) (interface{}, bool) {
	if t.root == nil {
		t.root = &node{key: key, value: value, color: black, tree: t}
		return nil, false
	}
	n, dir := t.root.find(key)
	switch dir {
	case exact:
		oldValue := n.value
		n.value = value
		return oldValue, true
	case left:
		l := &node{key: key, value: value, color: red, parent: n, tree: t}
		n.left = l
		l.ensureInvariants()
	case right:
		l := &node{key: key, value: value, color: red, parent: n, tree: t}
		n.right = l
		l.ensureInvariants()
	}
	return nil, false
}

// Get returns the value for the given key or nil if the key can't be found.
func (t *Tree) Get(key Comparable) (interface{}, bool) {
	if t.root == nil {
		return nil, false
	}
	n, dir := t.root.find(key)
	if dir == exact {
		return n.value, true
	}
	return nil, false
}

// String returns the textual representation of the tree.
func (t *Tree) String() string {
	if t.root == nil {
		return "-"
	}
	return t.root.str()
}

// Check verifies that the keys in the tree are ordered correctly.
func (t *Tree) Check() bool {
	if t.root == nil {
		return true
	}
	return t.root.check()
}
