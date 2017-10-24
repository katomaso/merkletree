package main

import (
	"fmt"
	"log"
	"crypto/sha1"
	"os"
	"io"
)

const branching_factor = 2
const tree_height = 40

type Block []byte


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

	result := Hash(file)
	fmt.Println(result)
}


func Hash(reader io.Reader) []byte {
	block := make([]byte, 32)
	input_chan := make(chan block)  // channel at the bottom of the tree
	middle_chan = input_chan
	output_chan := make(chan block)  // (future) tree root channel

	for i := 0; i < tree_height; i++ {
		go hashPartial(middle_chan, output_chan)
		middle_chan = output_chan
		output_chan := make(chan block)
	}
	// output channel is now the last - tree root channel
	hash_reader = sha1.Reader(reader)
	for _, err := hash_reader.Read(block); err == nil; _, err = hash_reader.Read(block) {
		input_chan <- block
	}
	input_chan <- nil // signal end of computation by nil value

	return <- output_chan
}

/** computePartial waits for data in `branching_factor` cycles and sends their
hash upwards in `output` channel. When nil is received the goroutine quits but
if it received some data together with the nil then it sends those data upwards.
**/
func computePartial(input chan block, output chan block) {
	for i := 0; i < branching_factor; i++ {

		if block <- channels[level]; block == nil {
			break;
		}
	}

	if has_data {
		if block == nil {
			upper_channel <- hash
		}
	}

}