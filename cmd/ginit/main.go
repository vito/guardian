package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Pid 1 Running")

	ignore := make([]byte, 1)
	os.Stdin.Read(ignore)
}
