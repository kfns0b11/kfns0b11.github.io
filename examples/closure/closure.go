package main

import "fmt"

func intSeq() func() int {
	i := 0

	// the returned anonymous function references variable i
	return func() int {
		i++
		return i
	}
}

func main() {
	nextInt := intSeq()

	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	// call intSeq() generate a new enviroment, i in newINts set to 0
	newInts := intSeq()
	fmt.Println(newInts())
}
