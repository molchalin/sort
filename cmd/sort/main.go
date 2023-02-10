package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/molchalin/sort/internal"
)

func main() {
	for _, v := range []int{1, 2, 3, 4, 5} {
		go func(v int) {
			fmt.Println(v)
		}(v)
	}
	// <- v := 5
	fmt.Scanln()

	return

	var r io.Reader = os.Stdin
	if len(os.Args) > 1 {
		var err error
		r, err = internal.FilesSource(os.Args[1:])
		if err != nil {
			log.Fatal(err)
		}
	}
	// -r reverse
	// -n численная сортировка
	// -k 2
	err := internal.Sort(os.Stdout, r)
	if err != nil {
		log.Fatal(err)
	}
}
