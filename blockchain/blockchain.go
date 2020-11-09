package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"
)

type BlockChain struct {
	Blocks []Block
}

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

func (block Block) hash() (string, error) {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()

	if _, err := h.Write([]byte(record)); err != nil {
		log.Println(err)

		return "", err
	}

	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed), nil
}

func GenerateBlock(oldBlock Block, bpm int) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = bpm
	newBlock.PrevHash = oldBlock.Hash
	hash, err := newBlock.hash()

	if err != nil {
		return newBlock, err
	}

	newBlock.Hash = hash

	return newBlock, err
}

func (block Block) IsBlockValid(oldBlock Block) bool {
	if oldBlock.Index+1 != block.Index {
		return false
	}

	if oldBlock.Hash != block.PrevHash {
		return false
	}

	hash, err := block.hash()

	if err != nil {
		return false
	}

	if hash != block.Hash {
		return false
	}

	return true
}

func (bchain BlockChain) ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > len(bchain.Blocks) {
		bchain.Blocks = newBlocks
	}
}
