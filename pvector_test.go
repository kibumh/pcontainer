package pcontainer_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/kibumh/pcontainer"
)

func ExamplePushBack() {
	pv0 := pcontainer.PVector{}
	fmt.Println("pv0's len:", pv0.Len())
	for i := 0; i < pv0.Len(); i++ {
		fmt.Println(pv0.At(i))
	}
	fmt.Println()

	pv1 := pv0.PushBack(0)
	fmt.Println("pv1's len:", pv1.Len())
	for i := 0; i < pv1.Len(); i++ {
		fmt.Println(pv1.At(i))
	}
	fmt.Println()

	pv2 := pv1.PushBack(1)
	fmt.Println("pv2's len:", pv2.Len())
	for i := 0; i < pv2.Len(); i++ {
		fmt.Println(pv2.At(i))
	}
	fmt.Println()

	pv33 := pv2
	for i := 3; i <= 33; i++ {
		pv33 = pv33.PushBack(i - 1)
	}
	fmt.Println("pv33's len:", pv33.Len())
	fmt.Println(pv33.At(31))
	fmt.Println(pv33.At(32))
	fmt.Println()

	pv65 := pv33
	for i := 34; i <= 65; i++ {
		pv65 = pv65.PushBack(i - 1)
	}
	fmt.Println("pv65's len:", pv65.Len())
	fmt.Println(pv65.At(63))
	fmt.Println(pv65.At(64))
	fmt.Println()

	pv1025 := pv65
	for i := 66; i <= 1025; i++ {
		pv1025 = pv1025.PushBack(i - 1)
	}
	fmt.Println("pv1025's len:", pv1025.Len())
	fmt.Println(pv1025.At(1023))
	fmt.Println(pv1025.At(1024))

	// Output:
	// pv0's len: 0
	//
	// pv1's len: 1
	// 0 <nil>
	//
	// pv2's len: 2
	// 0 <nil>
	// 1 <nil>
	//
	// pv33's len: 33
	// 31 <nil>
	// 32 <nil>
	//
	// pv65's len: 65
	// 63 <nil>
	// 64 <nil>
	//
	// pv1025's len: 1025
	// 1023 <nil>
	// 1024 <nil>
}

func ExampleUpdate() {
	pv0 := pcontainer.PVector{}
	pv1 := pv0.PushBack(0)
	pv2 := pv1.PushBack(1)
	pv2_2, _ := pv2.Update(1, 2)

	fmt.Println("pv2's len:", pv2.Len())
	fmt.Println(pv2.At(0))
	fmt.Println(pv2.At(1))
	fmt.Println()

	fmt.Println("pv2_2's len:", pv2_2.Len())
	fmt.Println(pv2_2.At(0))
	fmt.Println(pv2_2.At(1))

	// Output:
	// pv2's len: 2
	// 0 <nil>
	// 1 <nil>
	//
	// pv2_2's len: 2
	// 0 <nil>
	// 2 <nil>
}

func TestTransient(t *testing.T) {
	// pv0 {}
	// pv1 {0}
	// pv2 {0, 1}
	// pv3 {0, 10}
	// pv4 {0, 100, 2}
	// pv5 {0, 1000, 2}
	pv0 := pcontainer.PVector{}
	pv1 := pv0.PushBack(0)
	pv2 := pv1.PushBack(1)
	pv3, _ := pv2.Update(1, 10)
	pv4 := pv3.ConvertTransient()
	pv4 = pv4.PushBack(2)
	pv4, _ = pv4.Update(1, 100)
	pv4.ConvertPersistent()
	pv5, _ := pv4.Update(1, 1000)

	var cases = []struct {
		pv     pcontainer.PVector
		idx    int
		wanted int
	}{
		{pv1, 0, 0},
		{pv2, 0, 0},
		{pv2, 1, 1},
		{pv3, 0, 0},
		{pv3, 1, 10},
		{pv4, 0, 0},
		{pv4, 1, 100},
		{pv4, 2, 2},
		{pv5, 0, 0},
		{pv5, 1, 1000},
		{pv5, 2, 2},
	}

	for _, c := range cases {
		got, err := c.pv.At(c.idx)
		if err != nil || got != c.wanted {
			t.Errorf("At failed: pv(%v), idx(%v), got(%v), wanted(%v)", c.pv, c.idx, got, c.wanted)
		}
	}

}

func BenchmarkPushBack(b *testing.B) {
	pv := pcontainer.PVector{}
	for i := 0; i < b.N; i++ {
		pv = pv.PushBack(i)
	}
}

func BenchmarkTransientPushBack(b *testing.B) {
	pv := pcontainer.PVector{}
	pv = pv.ConvertTransient()
	for i := 0; i < b.N; i++ {
		pv = pv.PushBack(i)
	}
}

func BenchmarkSlicePushBack(b *testing.B) {
	s := []int{}
	for i := 0; i < b.N; i++ {
		s = append(s, i)
		if len(s) >= 20000000 { // to prevent out of memory.
			s = []int{}
		}
	}
}

func BenchmarkAt(b *testing.B) {
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

func BenchmarkSliceAt(b *testing.B) {
	len := 3000000
	s := []int{}
	for i := 0; i < len; i++ {
		s = append(s, i)
	}

	for i := 0; i < b.N; i++ {
		_ = s[rand.Intn(len)]
	}
}
