package pcontainer

import (
	"fmt"
	"strings"
)

const (
	branchingBit    = 5
	branchingFactor = 1 << branchingBit // 32
	branchingMask   = branchingFactor - 1

	// Following masks are for node.status.
	transientMask = 0x80
	lenMask       = 0x3f
)

type node struct {
	children [branchingFactor]interface{}
	status   uint8 // T0LLLLLL : T (transient bit), L(len)
}

func newNode(transient bool) *node {
	n := &node{}
	if transient {
		n.status |= transientMask
	}
	return n
}

func (n *node) len() uint8 {
	return n.status & lenMask
}

func (n *node) incLen() {
	n.status++
}

func (n *node) isTransient() bool {
	return (n.status & transientMask) == transientMask
}

func (n *node) convertPersistent(shift uint8) {
	if !n.isTransient() {
		return
	}
	n.status &= lenMask
	if shift == 0 {
		return
	}
	for i := uint8(0); i < n.len(); i++ {
		n.children[i].(*node).convertPersistent(shift - branchingBit)
	}
}

func (n *node) at(idx int, shift uint8) interface{} {
	for ; shift > 0; shift -= branchingBit {
		n = n.children[(idx>>shift)&branchingMask].(*node)
	}
	return n.children[idx&branchingMask]
}

func (n *node) clone(transient bool) *node {
	// In transient mode, we update in-place. So no clone!
	if transient && n.isTransient() {
		return n
	}
	newn := newNode(transient)
	copy(newn.children[:], n.children[:])
	newn.status |= n.status & lenMask
	return newn
}

func (n *node) update(idx int, v interface{}, shift uint8, transient bool) *node {
	newn := n.clone(transient)
	if shift == 0 {
		newn.children[idx&branchingMask] = v
		return newn
	}
	cidx := (idx >> shift) & branchingMask
	newn.children[cidx] = newn.children[cidx].(*node).update(idx, v, shift-branchingBit, transient)
	return newn
}

func (n *node) pushBack(v interface{}, shift uint8, transient bool) (newn *node, overflowed bool) {
	if shift == 0 {
		return n.pushBackChild(v, transient)
	}
	newChild, overflowed := n.children[n.len()-1].(*node).pushBack(v, shift-branchingBit, transient)
	if overflowed {
		return n.pushBackChild(newChild, transient)
	}
	newn = n.clone(transient)
	newn.children[newn.len()-1] = newChild
	return newn, false
}

func (n *node) pushBackChild(child interface{}, transient bool) (newn *node, overflowed bool) {
	overflowed = n.len() == branchingFactor
	if overflowed {
		newn = newNode(transient)
	} else {
		newn = n.clone(transient)
	}

	newn.children[newn.len()] = child
	newn.incLen()
	return newn, overflowed
}

// PVector is a persistent vector.
type PVector struct {
	root      *node
	len       int
	shift     uint8
	transient bool
}

// Len returns the number of values
func (pv PVector) Len() int {
	return pv.len
}

// At returns a value at a given index.
func (pv *PVector) At(idx int) (interface{}, error) {
	if idx < 0 || idx >= pv.len {
		return nil, fmt.Errorf("wrong index: idx(%d), len(%d)", idx, pv.len)
	}
	return pv.root.at(idx, pv.shift), nil
}

// Update updates a value at a given index.
func (pv PVector) Update(idx int, v interface{}) (PVector, error) {
	if idx < 0 || idx >= pv.len {
		return pv, fmt.Errorf("wrong index: idx(%d), len(%d)", idx, pv.len)
	}
	return PVector{pv.root.update(idx, v, pv.shift, pv.transient), pv.len, pv.shift, pv.transient}, nil
}

// PushBack adds a value at the end.
func (pv PVector) PushBack(v interface{}) PVector {
	if pv.root == nil {
		pv.root = newNode(pv.transient)
	}
	newn, overflowed := pv.root.pushBack(v, pv.shift, pv.transient)
	if !overflowed {
		return PVector{newn, pv.len + 1, pv.shift, pv.transient}
	}

	newroot := newNode(pv.transient)
	newroot.children[0] = pv.root
	newroot.incLen()
	newroot.children[1] = newn
	newroot.incLen()
	return PVector{newroot, pv.len + 1, pv.shift + branchingBit, pv.transient}
}

// ConvertTransient makes a given pvector transient.
// (A transient vector updates in-place. This is for a performance optimization.)
func (pv PVector) ConvertTransient() PVector {
	pv.transient = true
	return pv
}

// ConvertPersistent converts a transient vector to a persistent one.
func (pv *PVector) ConvertPersistent() {
	pv.transient = false
	pv.root.convertPersistent(pv.shift)
}

func (pv PVector) String() string {
	var b strings.Builder
	fmt.Fprint(&b, "[")
	sep := ""
	for i := 0; i < pv.Len(); i++ {
		v, _ := pv.At(i)
		fmt.Fprintf(&b, "%s%v", sep, v)
		sep = " "
	}
	fmt.Fprint(&b, "]")
	return b.String()
}
