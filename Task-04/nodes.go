package main

import (
	"iter"
)

// для построения дерева
type node struct {
	key   int
	value int

	parent *node

	left  *node
	right *node
}

func newNode(k, v int, p *node) *node {
	return &node{
		key:    k,
		value:  v,
		parent: p,
	}
}

func (o *node) find(key int) *node {
	switch {
	case o.key == key:
		return o
	case o.key < key && o.right != nil:
		return o.right.find(key)
	case o.key > key && o.left != nil:
		return o.left.find(key)
	}
	return nil
}

func (o *node) del(key int) *node {
	r := o.find(key)
	if r != nil {
		var n *node
		switch {
		case r.left != nil && r.right == nil:
			n = r.left
		case r.left == nil && r.right != nil:
			n = r.right
		case r.left != nil && r.right != nil:
			n = r.left
			if r.right.key > r.left.key {
				n = r.right
			}
		}

		switch {
		case r.parent.left == r:
			r.parent.left = n
		case r.parent.right == r:
			r.parent.right = n
		}

		return r
	}
	return nil
}

func (o *node) add(key, value int) *node {
	var n *node
	switch {
	case o.key == key:
		o.value = value
		return nil
	case o.right != nil && o.key < key:
		n = o.right.add(key, value)
	case o.left != nil && o.key > key:
		n = o.left.add(key, value)
	case o.left == nil && o.key > key:
		n = newNode(key, value, o)
		o.left = n
	case o.right == nil && o.key < key:
		n = newNode(key, value, o)
		o.right = n
	}

	return n
}

var preorder func(next *node)

func (n *node) Next(size int) iter.Seq[*node] {

	list := make([]*node, 0, size)

	preorder = func(next *node) {
		if next == nil {
			return
		}
		preorder(next.left)

		list = append(list, next)

		preorder(next.right)
	}

	preorder(n)

	return func(yield func(*node) bool) {
		for _, v := range list {
			if !yield(v) {
				return
			}
		}
	}
}
