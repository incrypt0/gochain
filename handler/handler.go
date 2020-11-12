package handler

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/incrypt0/gochain/blockchain"
)

type Handler struct {
	chain    *blockchain.BlockChain
	bcServer chan []*blockchain.Block
	mutex    *sync.Mutex
}

func (h *Handler) handleConn(conn net.Conn) {
	defer conn.Close()

	_, _ = io.WriteString(conn, "Enter a new message:")

	go h.scanBlockCreation(conn)

	go func() {
		for {
			time.Sleep(30 * time.Second)
			h.mutex.Lock()
			output, err := json.Marshal(h.chain)

			if err != nil {
				log.Panic(err)
			}

			h.mutex.Unlock()

			_, _ = io.WriteString(conn, string(output))
		}
	}()

	for range h.bcServer {
		spew.Dump(h.chain)
	}

	log.Println("4")
}

func (h *Handler) scanBlockCreation(conn io.ReadWriter) {
	blocks := h.chain.Blocks
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		data := scanner.Text()
		newBlock, err := blockchain.GenerateBlock(blocks[len(blocks)-1], data)

		if err != nil {
			log.Println(err)

			continue
		}

		if newBlock.IsBlockValid(blocks[len(blocks)-1]) {
			newBlockChain := append(h.chain.Blocks, newBlock)

			h.chain.ReplaceChain(newBlockChain)
		}

		h.bcServer <- h.chain.Blocks

		_, err = io.WriteString(conn, "\nEnter a new message:")

		if err != nil {
			log.Println(err)

			continue
		}
	}
}

func New(chain *blockchain.BlockChain) {
	server, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	bcServer := make(chan []*blockchain.Block)

	var mutex = &sync.Mutex{}

	h := Handler{chain: chain, bcServer: bcServer, mutex: mutex}

	if err != nil {
		log.Panic(err)
	}

	defer server.Close()

	for {
		conn, err := server.Accept()

		if err != nil {
			log.Panic(err)
		}

		go h.handleConn(conn)
	}
}
