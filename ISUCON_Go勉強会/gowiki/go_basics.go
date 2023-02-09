package main

import (
	"fmt"
)

func main() {
	// how to define variables
	var a int
	a = 5
	var b float32 = 3.14
	c := "I'm a perfect human"
	fmt.Println(a, b, c)

	// basic control syntaxes
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	d := "human"
	if d != "dog" {
		fmt.Println("I'm not a human")
	}

	// functions
	fmt.Println(add(1, 2))
	e := "cat"
	if isHuman(e) {
		fmt.Printf("%s is not a human\n", e)
	}

	// struct and methods
	taro := Dog{Kind: "Shiba", Name: "Taro", Age: 7}
	taro.explain()
	if taro.youngerThan(10) {
		fmt.Println("Taro is still young!")
	}

	// pointers
	i, j := 42, 2701

	p := &i         // point to i
	fmt.Println(*p) // read i through the pointer
	*p = 21         // set i through the pointer
	fmt.Println(i)  // see the new value of i
	p = &j          // point to j
	*p = *p / 37    // divide j through the pointer
	fmt.Println(j)  // see the new value of j

	// slices
	q := []int{2, 3, 5, 7, 11, 13}
	fmt.Println(q)
	for i, v := range q {
		fmt.Printf("index %d, value %d\n", i, v)
	}
}

func add(a int, b int) int {
	return a + b
}

func isHuman(s string) bool {
	return s == "human"
}

type Dog struct {
	Kind string
	Name string
	Age  int
}

func (d Dog) explain() {
	fmt.Printf("I'm %s, named %s and %d years old.\n", d.Kind, d.Name, d.Age)
}

func (d Dog) youngerThan(age int) bool {
	return d.Age < age
}
