package pcontainer

import (
	"fmt"
)

const (
	branchingBit    = 5
	branchingFactor = 1 << branchingBit // 32
	branchingMask   = branchingFactor - 1

	transientMask = 0x80
	lenMask       = 0x1f
)

type node struct {
	children [branchingFactor]interface{}
	status   uint8 // T00LLLLL : T (transient bit), L(len)
}

func (n *node) len() uint8 {
	return n.status & lenMask
}

func (n *node) incLen() {
	n.status++
}

func (n *node) setLen(len uint8) {
	n.status = n.status&transientMask + len
}

func (n *node) isTransient() bool {
	return n.status&transientMask == transientMask
}

func (n *node) convertPersistent(shift uint8) {
	if !n.isTransient() {
		return
	}
	if shift == 0 {
		n.status &= lenMask
		return
	}

	l := n.len()
	for i := uint8(0); i < l; i++ {
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
	if transient && n.isTransient() {
		return n
	}
	newn := &node{}
	copy(newn.children[:], n.children[:])
	newn.status = n.status
	if transient {
		newn.status |= transientMask
	}
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

func (n *node) pushBack(v interface{}, shift uint8, transient bool) (created bool, newn *node) {
	if shift == 0 {
		return n.pushBackChild(v, transient)
	}

	createdChild, newChild := n.children[n.len()-1].(*node).pushBack(v, shift-branchingBit, transient)

	if !createdChild {
		newn = n.clone(transient)
		newn.children[newn.len()-1] = newChild
		return false, newn
	}

	return n.pushBackChild(newChild, transient)
}

func (n *node) pushBackChild(child interface{}, transient bool) (created bool, newn *node) {
	if n.len() == branchingFactor {
		newn = &node{}
		newn.children[0] = child
		newn.setLen(1)
		if transient {
			newn.status |= transientMask
		}
		return true, newn
	}
	newn = n.clone(transient)
	newn.children[newn.len()] = child
	newn.incLen()
	return false, newn
}

// PVector represents radix search trie
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

// At ...
func (pv PVector) At(idx int) (interface{}, error) {
	if idx < 0 || idx >= pv.len {
		return nil, fmt.Errorf("radixtrie.At: wrong idx(%v)", idx)
	}
	return pv.root.at(idx, pv.shift), nil
}

// Update ...
func (pv PVector) Update(idx int, v interface{}) (PVector, error) {
	if idx < 0 || idx >= pv.len {
		return pv, fmt.Errorf("radixtrie.Update: wrong idx(%v)", idx)
	}
	return PVector{pv.root.update(idx, v, pv.shift, pv.transient), pv.len, pv.shift, pv.transient}, nil
}

// PushBack ...
func (pv PVector) PushBack(v interface{}) PVector {
	if pv.root == nil {
		newroot := &node{}
		newroot.children[0] = v
		newroot.setLen(1)
		if pv.transient {
			newroot.status |= transientMask
		}
		return PVector{newroot, 1, 0, pv.transient}
	}

	created, newn := pv.root.pushBack(v, pv.shift, pv.transient)
	if !created {
		return PVector{newn, pv.len + 1, pv.shift, pv.transient}
	}

	if pv.root.len() == branchingFactor {
		newroot := &node{}
		newroot.children[0] = pv.root
		newroot.children[1] = newn
		newroot.setLen(2)
		return PVector{newroot, pv.len + 1, pv.shift + branchingBit, pv.transient}
	}
	newroot := pv.root.clone(pv.transient)
	newroot.children[newroot.len()] = newn
	newroot.status++
	return PVector{newroot, pv.len + 1, pv.shift, pv.transient}
}

// ConvertTransient ...
func (pv PVector) ConvertTransient() PVector {
	newpv := pv
	newpv.transient = true
	return newpv
}

// ConvertPersistent ...
func (pv PVector) ConvertPersistent() {
	pv.transient = false
	pv.root.convertPersistent(pv.shift)
}
