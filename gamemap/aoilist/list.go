// 利用十字链表方式管理玩家视野
package gamemap

import (
	"math"
	"github.com/felix/felixgo/global"
	. "github.com/felix/felixgo/container/list"
	. "github.com/felix/felixgo/container"
)

type PosNode struct {
	pID uint64
	pos int
}

func newNode(pID uint64, pos int) *PosNode {
	return &PosNode{
		pID: pID,
		pos: pos,
	}
}

type PosList struct {
	*LinkList
}

func newPosList() *PosList {
	return &PosList{}
}

// 插入时保持有序链表
func (l *PosList) Add(pID uint64, pos int) {
	node := newNode(pID, pos)
	// 找到链表中比pos大的节点，然后插入到此节点的前面
	listHnode := l.Head
	listTnode := l.Tail
	if listHnode == nil {
		l.PushFront(node)
		return
	}
	for listHnode != nil && listTnode != nil {
		hnode := listHnode.Data.(*PosNode)
		if hnode.pos >= pos {
			l.PushBefore(listHnode, node)
			return
		}
		tnode := listTnode.Data.(*PosNode)
		if tnode.pos < pos {
			l.PushAfter(listTnode, node)
			return
		}
		listHnode = listHnode.Next
		listTnode = listTnode.Prev
	}
}

// 从链表中删除
func (l *PosList) Remove(pID uint64, pos int) {
	listNode := l.Head
	listTnode := l.Tail
	if listNode == nil {
		return
	}
	for listNode != nil && listTnode != nil {
		hnode := listNode.Data.(*PosNode)
		if hnode.pID == pID && hnode.pos == pos {
			l.Pop(listNode)
			return
		}

		tnode := listTnode.Data.(*PosNode)
		if tnode.pID == pID && tnode.pos == pos {
			l.Pop(listTnode)
			return
		}
		listNode = listNode.Next
		listTnode = listTnode.Prev
	}
}

// 获取玩家屏幕内所有周围玩家
func (l *PosList) GetScreenPlayers(pos int) SetUInt64 {
	listNode := l.Head
	if listNode == nil {
		return nil
	}
	var players SetUInt64
	for listNode != nil {
		node := listNode.Data.(*PosNode)
		dis := math.Abs(float64(node.pos - pos))
		if dis <= global.SCREEN_SIGHT {
			players.Add(node.pID)
		}
		if node.pos > pos && dis <= global.SCREEN_SIGHT {
			break
		}
		listNode = listNode.Next
	}
	return players
}

// 怪物或者玩家可以看到的玩家
func (l *PosList) GetViewPlayers(pos int) SetUInt64 {
	listNode := l.Head
	if listNode == nil {
		return nil
	}
	var players SetUInt64
	for listNode != nil {
		node := listNode.Data.(*PosNode)
		dis := math.Abs(float64(node.pos - pos))
		if dis <= global.SIGHT {
			players.Add(node.pID)
		}
		if node.pos > pos && dis > global.SIGHT {
			break
		}
		listNode = listNode.Next
	}
	return players
}
