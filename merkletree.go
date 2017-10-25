/*
merkletree package provides effective parallel implementation of MerkleTree.

One can exchange underlaying hashing algorithm. MerkleTree uses SHA1 by default
to be compatible with BitTorrent protocol. The hash can be of course exchanged
for SHA256 for compatibility with BitCoin or any other hash registered in Go's
own crypt.Hash.
*/

package merkletree

import (
	"crypto"
	"fmt"
	"github.com/katomaso/merkletree/hash"
	"hash"
	"io"
	"log"
	"os"
)

const Size = 32

var (
	TreeHeight  uint = 40
	DefaultHash      = crypto.SHA1
	marker           = make([]byte, 0)
)

type MerkleTree struct {
	// crypto.Hash to be used as underlaying hash for data blocks
	Hash crypto.Hash
	// input channel expect to receive already hashed
	input, output chan []byte
	// cache for data written by Writer interface smaller than minimal block
	// sending `marker` or nil will clear cache and force computation with it
	cache []byte
}

func (tree *MerkleTree) New() hash.Hash {
	block := make([]byte, Size)
	// channel at the bottom of the tree used to input hashed data blocks
	input_chan := make(chan Block)
	// the only reference to middle_hash is within matching goroutine
	middle_chan = input_chan
	// by redefining this will become the root (output) channel of a tree
	output_chan := make(chan Block)

	for i := 0; i < tree_height; i++ {
		go hashPartial(middle_chan, output_chan)
		middle_chan = output_chan
		output_chan := make(chan block)
	}

	return MerkleTree{DefaultHash, input_chan, output_chan}
}

func (tree *MerkleTree) Write(data []byte) (n int, err error) {
	hasher = tree.Hash.New()
	i := 0

	if len(tree.cache) > 0 {
		hasher.Write(tree.cache)
		i = -len(tree.cache)
		tree.cache.clear()
	}

	for ; i < len(data); i += Size {
		hasher.Write(data[min(0, i) : i+Size])
		tree.input <- hasher.Sum(nil)
		hasher.Reset()
	}

	if i < len(data) {
		Copy(data[i:], tree.cache)
	}
}

// Sum appends hash of underlaying data into b and returns the hash as well
func (tree *MerkleTree) Sum(b []byte) []byte {
	input_chan <- make([]byte, 0) // end empty data to provoke hash propagation
	hash <- output_chan
}

// Sending nil through tree will terminate all running goroutines
func (tree *MerkleTree) Close() {
	tree.input <- nil
	_ <- tree.output
}

/** computePartial waits for data in `branching_factor` cycles and sends their
hash upwards in `output` channel. When nil is received the goroutine quits but
if it received some data together with the nil then it sends those data upwards.
**/
func computePartial(input chan block, output chan block) {
	for i := 0; i < branching_factor; i++ {

		if block <- channels[level]; block == nil {
			break
		}
	}

	if has_data {
		if block == nil {
			upper_channel <- hash
		}
	}

}
