package node

import "fmt"

//BootNode is a node that is used to bootstrap the network.
//关于引导节点: https://ethereum.stackexchange.com/questions/8948/difference-between-full-node-and-bootstrap-nodes-in-ethereum
var BootNode = NewPeerNode("127.0.0.1", DefaultPort, true, true)

type PeerNode struct {
	IP          string `json:"ip"`
	Port        string `json:"port"`
	IsBootStrap bool   `json:"is_boot_strap"` // BootStrap node: 引导节点,用于发现网络上的其它节点
	IsActive    bool   `json:"is_active"`
}

func NewPeerNode(ip string, port string, isBootStrap bool, isActive bool) *PeerNode {
	return &PeerNode{
		IP:          ip,
		Port:        port,
		IsBootStrap: isBootStrap,
		IsActive:    isActive,
	}
}

func (n PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%s", n.IP, n.Port)
}
