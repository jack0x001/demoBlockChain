package node

import (
	"demoBlockChain/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

const Port = "51239"

func Run(dataDir string) error {
	state, err := database.NewStateFromDisk(dataDir)
	if err != nil {
		return err
	}
	defer state.CloseDB()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, "{hi demoBlockChain}")
	})

	router.GET("/balance/list", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, state.Balances)
	})

	router.GET("/blocks", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, state.Blocks)
	})

	router.POST("/tx/add", func(c *gin.Context) {
		var tx database.Tx
		if err := c.ShouldBindJSON(&tx); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := state.AddTx(tx)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := state.Persist()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "persist error"})
			return
		}
		c.IndentedJSON(http.StatusOK, hash)
	})

	//router.Run方法是阻塞的,直到return error
	err = router.Run(":" + Port)
	if err != nil {
		return err
	}
	return nil
}
