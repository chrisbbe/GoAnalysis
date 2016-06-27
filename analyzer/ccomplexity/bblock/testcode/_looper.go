package man

func main() { // BB #0 ending.
	sum := 0
	for i := 0; i < 10; i++ { // BB #1 ending.
		sum += i // BB #2 ending.
	}
	fmt.Println(sum) // BB #3 ending.
}
