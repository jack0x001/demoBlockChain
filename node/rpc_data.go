package node

import "demoBlockChain/database"

const DefaultPort = "51239"
const EndPointGetStatus = "/status"
const EndPointGetBlocks = "/blocks"
const EndPointGetBalanceList = "/balance/list"
const EndPointPostTxAdd = "/tx/add"
const EndPointPostPeerAdd = "/peer/add"

type StatusResponse struct {
	LastBlockHash database.HashCode   `json:"last_block_hash"`
	Height        uint64              `json:"height"`
	Peers         map[string]PeerNode `json:"peers"`
}
