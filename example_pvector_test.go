package pcontainer_test

import (
	"fmt"

	"github.com/kibumh/pcontainer"
)

func ExamplePVector() {
	var pv0 pcontainer.PVector

	pv1 := pv0.PushBack(0)
	for i := 1; i < 5; i++ {
		pv1 = pv1.PushBack(i)
	}
	pv2, _ := pv1.Update(0, 1000)

	fmt.Println("pv0:", pv0)
	fmt.Println("pv1:", pv1)
	fmt.Println("pv2:", pv2)

	// Output:
	// pv0: []
	// pv1: [0 1 2 3 4]
	// pv2: [1000 1 2 3 4]
}

func ExamplePVector_ConvertTransient() {
	var pv0 pcontainer.PVector
	pv1 := pv0.PushBack(0)
	pv1 = pv1.PushBack(1)
	pv2 := pv1.ConvertTransient()
	pv2 = pv2.PushBack(2)
	pv2 = pv2.PushBack(3)
	pv2, _ = pv2.Update(0, 1000)
	pv2.ConvertPersistent()

	fmt.Println("pv0:", pv0)
	fmt.Println("pv1:", pv1)
	fmt.Println("pv2:", pv2)

	// Output:
	// pv0: []
	// pv1: [0 1]
	// pv2: [1000 1 2 3]
}
