// 利用九宫格方式管理玩家视野
package gamemap

import (
	"fmt"
	"sync"
)

// 地图上一个虚拟格子类
type Grid struct {
	GID  int // 格子id
	minX int // 左边界
	maxX int // 有边界
	minY int // 下边界
	maxY int // 上边界

	playerMap map[uint64]bool // 格子内保存的玩家
	pLocker   sync.RWMutex    // 玩家id保护锁
}

func NewGrid(gid int, minX int, maxX int, minY int, maxY int) *Grid {
	grid := &Grid{
		GID:       gid,
		minX:      minX,
		maxX:      maxX,
		minY:      minY,
		maxY:      maxY,
		playerMap: make(map[uint64]bool),
	}
	return grid
}

func (g *Grid) Add(pID uint64) {
	g.pLocker.Lock()
	defer g.pLocker.Unlock()
	if _, ok := g.playerMap[pID]; ok {
		return
	}
	g.playerMap[pID] = true
}

func (g *Grid) Remove(pID uint64) {
	g.pLocker.Lock()
	defer g.pLocker.Unlock()
	if _, ok := g.playerMap[pID]; ok {
		delete(g.playerMap, pID)
	}
}

func (g *Grid) GetPlayerIDs() []uint64 {
	g.pLocker.Lock()
	defer g.pLocker.Unlock()

	var playerIDs []uint64
	for pID, _ := range g.playerMap {
		playerIDs = append(playerIDs, pID)
	}
	return playerIDs
}

// 打印字符串信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid gid=%d, minX=%d, maxX=%d, minY=%d, maxY=%d",
		g.GID, g.minX, g.maxX, g.minY, g.maxY)
}
