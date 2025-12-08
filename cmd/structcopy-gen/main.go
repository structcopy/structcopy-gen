package main

import (
	"fmt"
	"os"

	structcopygen "github.com/structcopy/structcopy-gen"
)

func main() {
	if err := structcopygen.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
