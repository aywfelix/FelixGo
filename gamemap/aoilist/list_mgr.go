// 利用十字链表方式管理玩家视野
package gamemap

import (
	. "github.com/aywfelix/felixgo/container"
)

type PosListManager struct {
	xList *PosList
	yList *PosList
}

func NewListManager() *PosListManager {
	return &PosListManager{
		xList: newPosList(),
		yList: newPosList(),
	}
}

// 这里的xy 代表的是地图格子坐标，不再是玩家真是像素坐标
func (pm *PosListManager) AddPlayerByPos(pID uint64, x, y int) {
	pm.xList.Add(pID, x)
	pm.yList.Add(pID, y)
}

func (pm *PosListManager) RemovePlayerByPos(pID uint64, x, y int) {
	pm.xList.Remove(pID, x)
	pm.yList.Remove(pID, y)
}

func (pm *PosListManager) MovePlayer(pID uint64, x_ori, y_ori, x, y int) {
	pm.RemovePlayerByPos(pID, x_ori, y_ori)
	pm.AddPlayerByPos(pID, x, y)
}

func (pm *PosListManager) GetScreenSurround(x, y int) SetUInt64 {
	xplayers := pm.xList.GetScreenPlayers(x)
	yPlayers := pm.yList.GetScreenPlayers(y)

	var players SetUInt64
	for pID, _ := range xplayers {
		players.Add(pID)
	}
	for pID, _ := range yPlayers {
		players.Add(pID)
	}
	return players
}

func (pm *PosListManager) GetViewSurround(x, y int) SetUInt64 {
	xplayers := pm.xList.GetViewPlayers(x)
	yPlayers := pm.yList.GetViewPlayers(y)

	var players SetUInt64
	for pID, _ := range xplayers {
		players.Add(pID)
	}
	for pID, _ := range yPlayers {
		players.Add(pID)
	}
	return players
}
