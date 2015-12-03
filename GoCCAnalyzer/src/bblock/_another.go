package main
import "fmt"

func main() {
	a := 1
	b := 2
	c := 3
	d := 4
	e := 5

	t := a + b
	u := a - b
	v := 0
	// #1
	if a < b { // #2
		v = t + c // #3
	} else { // #4
		w := u + c // #5
		v = w - d
	}
	x := v + e
	y := v - e

	fmt.Println(x)
	fmt.Println(y)
	// # 6

	if true { // #7
		fmt.Println("true") // #8
	}

	fmt.Printf("Done") // #9

}