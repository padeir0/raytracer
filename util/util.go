package util

import (
	"fmt"
	"math/rand"
	"os"
)

func Println(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func Print(a ...interface{}) {
	fmt.Fprint(os.Stderr, a...)
}

func Printf(f string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, f, a...)
}

func Bar(processed, total int) {
	bar := makebar(processed, total)
	fmt.Fprintf(os.Stderr, "\033[1A\033[K%v %v / %v                       \n", bar, processed, total)
}

const backgroundGreen = "\u001b[42m"
const reset = "\u001b[0m"

func makebar(processed, total int) string {
	bars := (processed * 20 / total)
	output := "|" + backgroundGreen
	for i := 0; i < bars; i++ {
		output += " "
	}
	output += reset
	for i := bars; i < 20; i++ {
		output += " "
	}
	return output + "|"
}

func Random(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}
