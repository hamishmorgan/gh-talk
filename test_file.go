package main

import "fmt"

// TestFunction demonstrates some code for review testing
func TestFunction() {
	// TODO: This could be improved
	x := 1
	y := 2
	z := x + y
	
	// This is a comment that might get feedback
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	
	// Another section reviewers might comment on
	if z > 0 {
		fmt.Println("Positive result")
	}
	
	// End of test function
	fmt.Println("Result:", z)
}

