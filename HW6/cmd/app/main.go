package main

import (
	"fmt"
	"hw6/internal/fact"
	"log"
	"os"
	//"hw7/internal/fact"
)

func main() {
	slice := []int{-321, 100, -99, 98, 97, 96, 95, 94, 32}
	fmt.Println("Hello, world!")
	if err := fact.NewFactorization().Work(fact.Input{NumsOfGoroutine: 2, Numbers: slice}, os.Stdout); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finished")
}
