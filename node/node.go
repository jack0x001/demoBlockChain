package node

import (
	"demoBlockChain/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const Port = "51239"

func Run(dataDir string) error {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {

		state, err := database.NewStateFromDisk(dataDir)
		if err != nil {
			log.Fatal("Error loading state from disk: ", err)
		}
		defer state.CloseDB()

		c.IndentedJSON(http.StatusOK, gin.H{
			"lastBlockHash": state.LastBlockHash,
			"height":        state.LatestBlock.Header.Height,
		})
	})

	router.GET("/balance/list", func(c *gin.Context) {
		state, err := database.NewStateFromDisk(dataDir)
		if err != nil {
			log.Fatal("Error loading state from disk: ", err)
		}
		defer state.CloseDB()

		c.IndentedJSON(http.StatusOK, state.Balances)
	})

	router.GET("/blocks", func(c *gin.Context) {
		state, err := database.NewStateFromDisk(dataDir)
		if err != nil {
			log.Fatal("Error loading state from disk: ", err)
		}
		defer state.CloseDB()
		c.IndentedJSON(http.StatusOK, state.Blocks)
	})

	router.POST("/tx/add", func(c *gin.Context) {
		state, err := database.NewStateFromDisk(dataDir)
		if err != nil {
			log.Fatal("Error loading state from disk: ", err)
		}
		defer state.CloseDB()

		var txList []database.Tx

		if err := c.ShouldBindJSON(&txList); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for i := range txList {
			err = state.AddTx(txList[i])
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		hash, err := state.Persist()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "persist error"})
			return
		}
		c.IndentedJSON(http.StatusOK, hash)
	})

	//router.Run方法是阻塞的,直到return error
	err := router.Run(":" + Port)
	if err != nil {
		return err
	}
	return nil
}
