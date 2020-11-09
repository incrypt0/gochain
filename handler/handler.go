package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/incrypt0/gochain/blockchain"
	"github.com/incrypt0/gochain/models"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	name   string
	bchain blockchain.BlockChain
}

func New() *Handler {
	return &Handler{name: "gochain"}
}

func (h *Handler) Register(e *echo.Echo) {
	log.Println("Hi from handler !!")
	e.GET("/", h.getBlockChain)
	e.POST("/", h.createBlock)
}

func (h *Handler) createBlock(c echo.Context) error {
	var m models.Message

	if err := c.Bind(&m); err != nil {
		log.Println(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{"success": false})
	}

	newBlock, err := blockchain.GenerateBlock(h.bchain.Blocks[len(h.bchain.Blocks)-1], m.BPM)

	if err != nil {
		log.Println(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{"success": false})
	}

	if newBlock.IsBlockValid(h.bchain.Blocks[len(h.bchain.Blocks)-1]) {
		h.bchain.Blocks = append(h.bchain.Blocks, newBlock)
	}

	return c.JSON(http.StatusOK, newBlock)
}

func (h *Handler) getBlockChain(c echo.Context) error {
	return c.JSON(http.StatusOK, h.bchain.Blocks)
}

func (h *Handler) GenesisBlock() {
	t := time.Now()

	genesisBlock := blockchain.Block{0, t.String(), 0, "", ""}

	spew.Dump(genesisBlock)
	h.bchain.Blocks = append(h.bchain.Blocks, genesisBlock)
}
