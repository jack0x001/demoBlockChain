package node

import (
	"bytes"
	"context"
	"demoBlockChain/database"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"time"
)

type Node struct {
	dataDir    string
	ip         string
	port       string
	state      *database.State
	knownPeers map[string]PeerNode
}

//IsBootstrap 判断是否为引导节点
func (n Node) IsBootstrap() bool {
	return n.TcpAddress() == BootNode.TcpAddress()
}

func NewNode(dataDir string, ip string, port string, bootstrap PeerNode) *Node {
	peers := make(map[string]PeerNode)
	peers[bootstrap.TcpAddress()] = bootstrap
	return &Node{
		dataDir:    dataDir,
		ip:         ip,
		port:       port,
		knownPeers: peers,
	}
}

//节点同步
func (n *Node) sync(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 15)
	for {
		select {
		case <-ticker.C:
			n.fetchNewBlocksAndPeers()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (n Node) TcpAddress() string {
	return fmt.Sprintf("%s:%s", n.ip, n.port)
}

func getPeerStatus(url string) (StatusResponse, error) {

	response, err := http.Get(url)
	if err != nil {
		return StatusResponse{}, err
	}
	if response.StatusCode != 200 {
		return StatusResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(response.Body)

	var status StatusResponse
	err = json.NewDecoder(response.Body).Decode(&status)
	if err != nil {
		return StatusResponse{}, err
	}
	return status, nil
}

func (n *Node) addPeer(peer PeerNode) {
	if _, ok := n.knownPeers[peer.TcpAddress()]; !ok {
		log.Println("New peer added:", peer.TcpAddress())
		n.knownPeers[peer.TcpAddress()] = peer
	}
}

func (n *Node) fetchNewBlocksAndPeers() {
	for _, peer := range n.knownPeers {

		if peer.TcpAddress() == n.TcpAddress() {
			continue
		}

		url := fmt.Sprintf("http://%s%s", peer.TcpAddress(), EndPointGetStatus)
		peerStatus, err := getPeerStatus(url)
		if err != nil {
			log.Printf("Error fetching peer status: %s : %s\n", url, err)
			continue
		}

		if !n.IsBootstrap() && n.state.LatestBlock.Header.Height < peerStatus.Height {
			//TODO: 同步区块
			log.Println("TODO: 同步区块")
		}

		// 同步节点
		for _, p := range peerStatus.Peers {
			n.addPeer(p)
		}

		//将自己加入到peer的节点中
		url = fmt.Sprintf("http://%s%s", peer.TcpAddress(), EndPointPostPeerAdd)
		me := PeerNode{
			IP:          n.ip,
			Port:        n.port,
			IsBootStrap: false,
			IsActive:    true,
		}
		j, _ := json.Marshal(me)
		reader := bytes.NewReader(j)
		post, err := http.Post(url, "application/json", reader)
		if err != nil {
			log.Println("Error posting to peer:", err)
			continue
		}
		if post.StatusCode != 200 {
			log.Println("Error posting to peer:", post.StatusCode)
			continue
		}
	}
}

func (n *Node) Run() error {
	isBootStrapNode := n.TcpAddress() == BootNode.TcpAddress()

	if !isBootStrapNode {
		return runAsNormalNode(n)
	} else {
		return runAsBootStrapNode(n)
	}
}

//引导节点应只做节点发现，不做区块同步
func runAsBootStrapNode(n *Node) error {
	log.Println("Starting bootstrap node...")

	ctx := context.Background()
	go n.sync(ctx)

	router := gin.Default()

	router.GET(EndPointGetStatus, func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, StatusResponse{
			Peers: n.knownPeers,
		})
	})

	router.POST(EndPointPostPeerAdd, func(c *gin.Context) {
		var peer PeerNode
		err := c.ShouldBindJSON(&peer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		n.addPeer(peer)
		c.JSON(http.StatusOK, gin.H{"message": "Peer added"})
	})

	err := router.Run(n.TcpAddress())
	if err != nil {
		return err
	}
	return nil
}

func runAsNormalNode(n *Node) error {

	log.Println("Starting node...")

	router := gin.Default()
	s, err := database.NewStateFromDisk(n.dataDir)
	if err != nil {
		log.Fatal("Error loading state from disk: ", err)
	}
	n.state = s

	defer n.state.CloseDB()

	ctx := context.Background()
	go n.sync(ctx)

	router.GET(EndPointGetStatus, func(c *gin.Context) {

		c.IndentedJSON(http.StatusOK, StatusResponse{
			LastBlockHash: n.state.LastBlockHash,
			Height:        n.state.LatestBlock.Header.Height,
			Peers:         n.knownPeers,
		})
	})

	router.GET(EndPointGetBalanceList, func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, n.state.Balances)
	})

	router.GET(EndPointGetBlocks, func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, n.state.Blocks)
	})

	router.POST(EndPointPostPeerAdd, func(c *gin.Context) {
		var peer PeerNode
		err := c.ShouldBindJSON(&peer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		n.addPeer(peer)
		c.JSON(http.StatusOK, gin.H{"message": "Peer added"})
	})

	router.POST(EndPointPostTxAdd, func(c *gin.Context) {

		var txList []database.Tx
		if err := c.ShouldBindJSON(&txList); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for i := range txList {
			err = n.state.AddTx(txList[i])
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		hash, err := n.state.Persist()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "persist error"})
			return
		}

		//重新加载
		s, err := database.NewStateFromDisk(n.dataDir)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "persist error"})
			return
		}
		n.state = s

		c.IndentedJSON(http.StatusOK, hash)
	})

	//router.Run方法是阻塞的,直到return error
	err = router.Run(n.TcpAddress())
	if err != nil {
		return err
	}
	return nil
}
