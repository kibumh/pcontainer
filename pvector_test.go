package pcontainer_test

import (
	"math/rand"
	"testing"

	"github.com/kibumh/pcontainer"
)

func TestPvector(t *testing.T) {
	to := 1025

	var pv pcontainer.PVector
	for i := 0; i < to; i++ {
		pv = pv.PushBack(i)
	}
	if pv.Len() != to {
		t.Errorf("Len() is wrong: got(%d) expected(1000)", pv.Len())
	}
	for i := 0; i < to; i++ {
		v, err := pv.At(i)
		if err != nil {
			t.Errorf("failed to At(): err(%v) at index %d", err, i)
		}
		if v != i {
			t.Errorf("wrong At() result: got(%d) expected(%d) at index %d", v, i, i)
		}
	}
	_, err := pv.At(to)
	if err == nil {
		t.Error("error is expected but nil.")
	}
}

func BenchmarkPVector_PushBack(b *testing.B) {
	var pv pcontainer.PVector
	for i := 0; i < b.N; i++ {
		pv = pv.PushBack(i)
	}
}

func BenchmarkPVector_PushBack_Transient(b *testing.B) {
	var pv pcontainer.PVector
	pv = pv.ConvertTransient()
	for i := 0; i < b.N; i++ {
		pv = pv.PushBack(i)
	}
}

func BenchmarkSlice_append(b *testing.B) {
	s := []int{}
	for i := 0; i < b.N; i++ {
		s = append(s, i)
		if len(s) >= 20000000 { // to prevent out of memory.
			s = []int{}
		}
	}
}

func BenchmarkPVector_At(b *testing.B) {
	len := 3000000
	pv := pcontainer.PVector{}
	pv = pv.ConvertTransient()
	for i := 0; i < len; i++ {
		pv = pv.PushBack(i)
	}

	for i := 0; i < b.N; i++ {
		_, _ = pv.At(rand.Intn(len))
	}
}

func BenchmarkSlice_At(b *testing.B) {
	len := 3000000
	s := []int{}
	for i := 0; i < len; i++ {
		s = append(s, i)
	}

	for i := 0; i < b.N; i++ {
		_ = s[rand.Intn(len)]
	}
}
