package dag_test

import (
	"errors"
	"fmt"

	"github.com/drone/dag"
)

func ExampleRunner() {
	var r dag.Runner

	r.AddVertex("one", func() error {
		fmt.Println("one and two will run in parallel before three")
		return nil
	})
	r.AddVertex("two", func() error {
		fmt.Println("one and two will run in parallel before three")
		return nil
	})
	r.AddVertex("three", func() error {
		fmt.Println("three will run before four")
		return errors.New("three is broken")
	})
	r.AddVertex("four", func() error {
		fmt.Println("four will never run")
		return nil
	})

	r.AddEdge("one", "three")
	r.AddEdge("two", "three")

	r.AddEdge("three", "four")

	fmt.Printf("the runner terminated with: %v\n", r.Run())
	// Output:
	// one and two will run in parallel before three
	// one and two will run in parallel before three
	// three will run before four
	// the runner terminated with: three is broken
}

func ExampleRunnerWithFailure() {
	var r dag.Runner

	r.AddVertex("one", func() error {
		fmt.Println("one and two will run in parallel before three")
		return nil
	})
	r.AddVertex("two", func() error {
		fmt.Println("one and two will run in parallel before three")
		return nil
	})
	r.AddVertex("three", func() error {
		fmt.Println("three will run before four")
		return errors.New("three is broken")
	})
	r.AddVertex("four", func() error {
		fmt.Println("four will never run")
		return nil
	})
	r.AddVertex("five", func() error {
		fmt.Println("five ran as a dep of one")
		return nil
	})
	r.AddVertex("six", func() error {
		fmt.Println("six will run as a dep of five.")
		return nil
	})

	r.AddEdge("one", "three")
	r.AddEdge("one", "five")
	r.AddEdge("five", "six")
	r.AddEdge("two", "three")

  // -------->
	// 1 -> 2 -> 3 -> 4
	//   -> 5 -> 6

	r.AddEdge("three", "four")

	fmt.Printf("the runner terminated with: %v\n", r.RunOnFailure())
	
	// Output:
	// one and two will run in parallel before three
	// one and two will run in parallel before three
	// five ran as a dep of one
	// six will run as a dep of five.
	// three will run before four
	// the runner terminated with: three is broken
}
