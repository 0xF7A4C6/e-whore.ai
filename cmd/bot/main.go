package main

import (
	"fmt"

	"github.com/0xF7A4C6/e-whore.ai/internal/vercel"
)

func main() {
	v, err := vercel.NewVercelClient()
	if err != nil {
		panic(err)
	}

	text, err := v.GetPrompt("hi!!!")
	if err != nil {
		panic(err)
	}

	fmt.Println(text)
}
