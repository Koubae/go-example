/*

@credit: https://github.com/MuxN4/pow
*/
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

// ------------------------------------------------------------
// MAIN
type Blockchain struct {
	Blocks []*Block
}

// Creates a blockchain with a genesis block
func NewBlockchain() *Blockchain {
	fmt.Printf("🚀 Starting Mining Process...")

	genesisBlock := NewBlock(0, "Genesis Block", "", 2)
	powInstance := NewProofOfWork(*genesisBlock)
	success, _, _ := powInstance.Mine()

	if !success {
		log.Fatal("Failed to mine genesis block")
	}

	return &Blockchain{Blocks: []*Block{genesisBlock}}
}

// converts duration to human-readable format
func formatDuration(d time.Duration) string {
	switch {
	case d < time.Millisecond:
		return fmt.Sprintf("%dµs", d.Microseconds())
	case d < time.Second:
		return fmt.Sprintf("%dms", d.Milliseconds())
	default:
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

// Where mining and adding a new block to the blockchain takes place
func (bc *Blockchain) AddBlock(data string, difficulty int) *Block {
	previousHash := ""
	if len(bc.Blocks) > 0 {
		previousHash = bc.Blocks[len(bc.Blocks)-1].Hash
	}

	newBlock := NewBlock(len(bc.Blocks), data, previousHash, difficulty)
	powInstance := NewProofOfWork(*newBlock)

	success, nonce, duration := powInstance.Mine()
	if !success {
		fmt.Println("❌ Failed to mine block")
		return nil
	}

	bc.Blocks = append(bc.Blocks, newBlock)

	fmt.Printf("⛏️  Block mined!")
	fmt.Printf("    %s:   %d\n", "Difficulty", difficulty)
	fmt.Printf("    %s:        %d\n", "Nonce", nonce)
	fmt.Printf("    %s:         %s\n", "Time", formatDuration(duration))
	fmt.Printf("    %s:         %s\n", "Hash", newBlock.Hash)

	// Handle previous hash display
	prevHashDisplay := previousHash
	if prevHashDisplay == "" {
		prevHashDisplay = "[Genesis Block]"
	}
	fmt.Printf("    %s:    %s\n\n", "Prev Hash", prevHashDisplay)

	return newBlock
}

// ------------------------------------------------------------
// POW

type ProofOfWork struct {
	block    Block
	target   string
	maxNonce int
}

// This function initializes a new Proof of Work instance with some ground rules
func NewProofOfWork(block Block) *ProofOfWork {
	target := strings.Repeat("0", block.Difficulty)
	maxNonce := calculateMaxNonce(block.Difficulty)

	return &ProofOfWork{
		block:    block,
		target:   target,
		maxNonce: maxNonce,
	}
}

// determines how many attempts we'll tolerate
func calculateMaxNonce(difficulty int) int {
	return math.MaxInt64
	//return 1_000_000 * difficulty
}

// Mine attempts to find a valid hash meeting the difficulty criteria
func (pow *ProofOfWork) Mine() (bool, int, time.Duration) {
	startTime := time.Now()

	for pow.block.Nonce < pow.maxNonce {
		hash := pow.block.CalculateHash()

		if strings.HasPrefix(hash, pow.target) {
			pow.block.Hash = hash
			return true, pow.block.Nonce, time.Since(startTime)
		}

		pow.block.Nonce++
	}

	return false, -1, time.Since(startTime)
}

// checks if the block's hash meets difficulty requirements
func (pow *ProofOfWork) Validate() bool {
	return strings.HasPrefix(pow.block.Hash, pow.target)
}

// ------------------------------------------------------------
// BLOCK

// Represents a single block in the blockchain
type Block struct {
	Index        int
	Timestamp    int64
	Data         string
	PreviousHash string
	Hash         string
	Nonce        int
	Difficulty   int
}

// NewBlock creates a new block with given parameters
func NewBlock(index int, data string, previousHash string, difficulty int) *Block {
	return &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Data:         data,
		PreviousHash: previousHash,
		Nonce:        0,
		Difficulty:   difficulty,
	}
}

// Here CalculateHash generates a unique hash for the block
func (b *Block) CalculateHash() string {
	record := strconv.Itoa(b.Index) +
		strconv.FormatInt(b.Timestamp, 10) +
		b.Data +
		b.PreviousHash +
		strconv.Itoa(b.Nonce)

	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	blockchain := NewBlockchain()

	// A simulation with different difficulties
	//difficulties := []int{2, 4, 5, 6}
	difficulties := []int{7, 8}
	for _, diff := range difficulties {
		blockchain.AddBlock(fmt.Sprintf("Transaction at difficulty %d", diff), diff)
	}

	fmt.Println("🎉 Mining Complete!")
}
