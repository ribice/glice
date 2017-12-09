package main

import (
	"fmt"
	"strings"
)

func main() {
	c := fmt.Sprintf("Glice")

	a := c + c
	k := strings.Split(a, c)
	_ = k
}
