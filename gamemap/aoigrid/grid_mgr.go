package gamemap

import (
	"fmt"
)

// 地图视野管理器
type GridManager struct {
	minX int // 场景地图左边界
	maxX int // 场景地图右边界
	minY int // 场景地图下边界
	maxY int // 场景地图上边界

	cntsx int // 横向格子数量
	cntsy int // 纵向格子数量

	grids map[int]*Grid // 保存地图虚拟格子数
}

func NewAOIManager(minX int, maxX int, minY int, maxY int, cntsx int, cntsy int) *GridManager {
	aoiMgr := &GridManager{
		minX:  minX,
		maxX:  maxX,
		minY:  minY,
		maxY:  maxY,
		cntsx: cntsx,
		cntsy: cntsy,
		grids: make(map[int]*Grid, cntsx*cntsy),
	}
	// 初始化地图视野所有格子
	for x := 0; x < cntsx; x++ {
		for y := 0; y < cntsy; y++ {
			// 计算格子编号
			gID := y*cntsx + x
			grid := NewGrid(gID,
				aoiMgr.minX+x*aoiMgr.gridLength(),
				aoiMgr.minX+(x+1)*aoiMgr.gridLength(),
				aoiMgr.minY+y*aoiMgr.gridWidth(),
				aoiMgr.minY+(y+1)*aoiMgr.gridWidth())
			aoiMgr.grids[gID] = grid
		}
	}
	return aoiMgr
}

//得到每个格子在x轴方向的长度
func (gm *GridManager) gridLength() int {
	return (gm.maxX - gm.minX) / gm.cntsx
}

//得到每个格子在x轴方向的宽度
func (gm *GridManager) gridWidth() int {
	return (gm.maxY - gm.minY) / gm.cntsy
}

// 打印字符串信息
func (gm *GridManager) String() string {
	s := fmt.Sprintf("AoiManager: minX=%d, maxX=%d, minY=%d, maxY=%d, cntsx=%d, cntsy=%d, gridWidth=%d, gridLength=%d\n",
		gm.minX, gm.maxX, gm.minY, gm.maxY, gm.cntsx, gm.cntsy, gm.gridWidth(), gm.gridLength())
	for _, grid := range gm.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据格子id获取周边九宫格信息
func (gm *GridManager) getSurroundGridsByGridID(gID int) []*Grid {
	if _, ok := gm.grids[gID]; !ok {
		return nil
	}

	var retGrids []*Grid
	retGrids = append(retGrids, gm.grids[gID])
	// 根据gID, 得到格子所在的坐标
	x, y := gID%gm.cntsx, gID/gm.cntsx
	// 新建一个临时存储周围格子的数组
	surroundGID := make([]int, 0)
	// 新建8个方向向量: 左上: (-1, -1), 左中: (-1, 0), 左下: (-1,1), 中上: (0,-1), 中下: (0,1), 右上:(1, -1)
	// 右中: (1, 0), 右下: (1, 1), 分别将这8个方向的方向向量按顺序写入x, y的分量数组
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	// 根据8个方向向量, 得到周围点的相对坐标, 挑选出没有越界的坐标, 将坐标转换为gID
	for i := 0; i < 8; i++ {
		newX := x + dx[i]
		newY := y + dy[i]

		if newX >= 0 && newX < gm.cntsx && newY >= 0 && newY < gm.cntsy {
			surroundGID = append(surroundGID, newY*gm.cntsx+newX)
		}
	}
	// 根据没有越界的gID, 得到格子信息
	for _, gID := range surroundGID {
		retGrids = append(retGrids, gm.grids[gID])
	}
	return retGrids
}

// 地图坐标=》视野编号
func (gm *GridManager) GetGIDbyPos(x, y float32) int {
	gx := (int(x) - gm.minX) / gm.gridLength()
	gy := (int(y) - gm.minY) / gm.gridWidth()

	return gy*gm.cntsx + gx
}

// 根据地图左边获取九宫格内玩家
func (gm *GridManager) GetPIDsByPos(x, y float32) []uint64 {
	gID := gm.GetGIDbyPos(x, y)
	grids := gm.getSurroundGridsByGridID(gID)
	var playerIDs []uint64
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}
	return playerIDs
}

// 将玩家加入到某个格子中
func (gm *GridManager) AddPlayerByPos(pID uint64, x, y float32) {
	gID := gm.GetGIDbyPos(x, y)
	if _, ok := gm.grids[gID]; !ok {
		return
	}
	grid := gm.grids[gID]
	grid.Add(pID)
}

// 将玩家从某个格子中移除
func (gm *GridManager) RemovePlayerByPos(pID uint64, x, y float32) {
	gID := gm.GetGIDbyPos(x, y)
	if _, ok := gm.grids[gID]; !ok {
		return
	}
	grid := gm.grids[gID]
	grid.Remove(pID)
}

// 将玩家加入到某个格子中
func (gm *GridManager) AddPlayerByGridID(pID uint64, gID int) {
	if _, ok := gm.grids[gID]; !ok {
		return
	}
	grid := gm.grids[gID]
	grid.Add(pID)
}

// 将玩家从某个格子中移除
func (gm *GridManager) RemovePlayerGridID(pID uint64, gID int) {
	if _, ok := gm.grids[gID]; !ok {
		return
	}
	grid := gm.grids[gID]
	grid.Remove(pID)
}
