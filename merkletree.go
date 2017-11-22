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
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"log"
)

const BlockSize = 32

var (
	TreeHeight      = 3
	DefaultHash     = crypto.SHA1
	BranchingFactor = 2
	marker          = make([]byte, 0)
	_               = sha1.Size
	_               = sha256.Size
	_               = sha512.Size
)

type Block []byte

type MerkleTree struct {
	// input channel expect to receive already hashed
	channels []chan Block
	// crypto.Hash to be used as underlaying hash for data blocks
	hash crypto.Hash
	// cached holds number of bytes already in hasher (less than BlockSize)
	cached int
}

func New() hash.Hash {
	// channel at the bottom of the tree used to input hashed data blocks
	channels := make([]chan Block, TreeHeight)
	for i, _ := range channels {
		channels[i] = make(chan Block)
	}
	tree := MerkleTree{channels, DefaultHash, 0}
	tree.Start()
	return tree
}

func (tree *MerkleTree) Start() {
	for i := 0; i < TreeHeight-1; i++ {
		go tree.compute(tree.channels[i], tree.channels[i+1])
	}
}

func (tree MerkleTree) Write(data []byte) (n int, err error) {
	var (
		i      int = 0
		imax       = max(0, len(data)-BlockSize)
		hasher     = tree.hash.New()
	)

	if tree.cached > 0 {
		i = -tree.cached
		tree.cached = 0
	}
	log.Printf("Write called with Slice %d/%d\n", len(data), cap(data))
	for ; i <= imax; i += BlockSize {
		log.Printf("Writing <%d, %d> data\n", max(0, i), min(imax, i+BlockSize))
		hasher.Write(data[max(0, i):min(imax, i+BlockSize)])
		tree.channels[0] <- hasher.Sum(nil)
		hasher.Reset()
	}

	if i < len(data) {
		log.Printf("Caching %d of data", len(data)-i)
		hasher.Write(data[i:])
		tree.cached = len(data) - i
	}

	return i, nil
}

// Sum appends hash of underlaying data into b and returns the hash as well
func (tree MerkleTree) Sum(b []byte) []byte {
	log.Println("Writing marker (empty cache)!")
	tree.Write(marker) // ensure emptying cache
	log.Println("Sending marker (propagate hash)!")
	tree.channels[0] <- marker // send 'marker' to provoke hash propagation
	log.Println("Receiving hash!")
	return <-tree.channels[TreeHeight-1]
}

func (tree MerkleTree) Reset() {
	tree.Close()
	tree.Start()
}

// Sending nil through tree will terminate all running goroutines
func (tree *MerkleTree) Close() {
	tree.channels[0] <- nil
	<-tree.channels[TreeHeight-1]
}

func (tree MerkleTree) Size() int {
	return tree.hash.Size()
}

func (tree MerkleTree) BlockSize() int {
	return BlockSize
}

/** computePartial waits for data in `branching_factor` cycles and sends their
hash upwards in `output` channel. When nil is received the goroutine quits but
if it received some data together with the nil then it sends those data upwards.
**/
func (tree *MerkleTree) compute(input, output chan Block) {
	var hasher hash.Hash = tree.hash.New()

Always:
	for {
		for i := 0; i < BranchingFactor; i++ {
			block := <-input

			switch {
			case block == nil:
				break Always
			case len(block) == 0:
				break
			default:
				if i == 0 {
					hasher.Reset()
				}
				hasher.Write(block)
			}
		}
		output <- hasher.Sum(nil)
	}
	output <- hasher.Sum(nil)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
