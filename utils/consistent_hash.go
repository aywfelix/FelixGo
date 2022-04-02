// 一致性hash算法，初始化n个虚拟节点，然后将真实的服务器节点包装成虚拟节点,每个真实节点对应100个虚拟节点
// 根据虚拟节点获取提供服务的真实节点
// ServerID_Index
package utils

import (
	"fmt"
	"sort"

	"github.com/aywfelix/felixgo/encrypt"
)

// 定义虚拟节点
type VirtualNode struct {
	serverID int
	index    int
}

func NewVNode(serverID int, index int) *VirtualNode {
	return &VirtualNode{
		serverID: serverID,
		index:    index,
	}
}

func (v *VirtualNode) GetHash() uint32 {
	str := fmt.Sprintf("%d-%d", v.serverID, v.index)
	return encrypt.CRC32(str)
}

func (v *VirtualNode) Index() int {
	return v.index
}

func (v *VirtualNode) ServerID() int {
	return v.serverID
}

//----------------------------------------------------------------------------------
const VNODE_NUMBER = 100

type SortKeys []uint32

func (s SortKeys) Len() int {
	return len(s)
}
func (s SortKeys) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s SortKeys) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHash struct {
	virtualNodes map[uint32]*VirtualNode // key是虚拟节点hash值
	serverMap    map[int]bool
	keys         SortKeys
}

func NewConsistentHash() *ConsistentHash {
	consistentHash := &ConsistentHash{
		serverMap:    make(map[int]bool, 0),
		virtualNodes: make(map[uint32]*VirtualNode, 0),
	}

	return consistentHash
}

func (c *ConsistentHash) Insert(serverID int) {
	if _, ok := c.serverMap[serverID]; ok {
		return
	}
	c.serverMap[serverID] = true

	for i := 0; i < VNODE_NUMBER; i++ {
		vnode := NewVNode(serverID, i)
		hash := vnode.GetHash()
		if _, ok := c.virtualNodes[hash]; !ok {
			c.virtualNodes[hash] = vnode
			c.keys = append(c.keys, hash)
		}
	}
	sort.Sort(c.keys)
}

func (c *ConsistentHash) Delete(serverID int) {
	if _, ok := c.serverMap[serverID]; !ok {
		return
	}
	delete(c.serverMap, serverID)

	for i := 0; i < VNODE_NUMBER; i++ {
		vnode := NewVNode(serverID, i)
		hash := vnode.GetHash()
		if _, ok := c.virtualNodes[hash]; ok {
			delete(c.virtualNodes, hash)
		}
	}

	keys := SortKeys{}
	for key, _ := range c.virtualNodes {
		keys = append(keys, key)
	}
	sort.Sort(keys)
	c.keys = keys
}

// 获取某个服务
func (c *ConsistentHash) Get(key string) int {
	hash := encrypt.CRC32(key)
	vnode, ok := c.virtualNodes[hash]
	if ok {
		return vnode.ServerID()
	}
	// 查找最接近的服务节点
	index := c.search(hash)
	node := c.virtualNodes[c.keys[index]]
	return node.ServerID()
}

func (c *ConsistentHash) search(hash uint32) int {
	i := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})
	if i < len(c.keys) {
		if i == len(c.keys)-1 {
			return 0
		} else {
			return i
		}
	}
	return len(c.keys) - 1
}

//--------------------------------------------------------------------------------------------------------------
