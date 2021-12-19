package main

import "fmt"

func f(n int) {
	for i := 0; i < 10; i++ {
		fmt.Println(n, ":", i)
	}
}

func a() {
	go f(0)
	go f(5)
	var input string
	fmt.Scanln(&input)
}
