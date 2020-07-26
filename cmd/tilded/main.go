package main

import (
	"fmt"

	"github.com/mikecamilleri/tilde/internal/state"
)

func main() {
	fmt.Println("Hello World!")
	s := state.NewState()
	_ = &s
}
