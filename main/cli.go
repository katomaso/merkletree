package main

import (
	"fmt"
	"github.com/katomaso/merkletree"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <filename>", os.Args[0])
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tree = merkletree.New()
	defer tree.Close()

	io.Copy(file, tree)
	fmt.Println(tree.Sum(nil))
}
