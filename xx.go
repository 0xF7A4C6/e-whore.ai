package main

import (
	"fmt"
	"math"
)

/*
	function decrypt(a) {
		return a - Math.log10(a * Math.LN10)
	}
*/

func decrypt(a float64) float64 {
	return a - math.Log10(a * math.Ln10)
}

func main() {
	result := decrypt(0.42850214618580407)
	fmt.Println(result)
}
