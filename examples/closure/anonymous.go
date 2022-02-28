package main

import "fmt"

func testanonymousfunc() {
	// execute locally
	func() {
		fmt.Println("Anonymous function execute locally!")
	}()

	// assign to a variable
	var anonyfunc func() = func() {
		fmt.Println("Anoymous function assign to a variable")
	}

	anonyfunc()

	// change value of anonyfunc
	anonyfunc = somefunc
	anonyfunc()
}

func somefunc() {
	fmt.Println("some standard func")
}
