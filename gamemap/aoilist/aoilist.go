package gamemap

import (
	"fmt"
)

// 十字链表法， 实现AOI
// 节点增加左右哨兵功能， 同时增加地标功能

type nodeType int

var (
	NormalNode    nodeType = 0 // 正常玩家节点
	LeftSentinel  nodeType = 1 // 左哨兵节点
	RightSentinel nodeType = 2 // 右哨兵节点
	LandMark      nodeType = 3 // 地标节点
)

type Direction int

var (
	XDirect Direction = 0
	YDirect Direction = 1
)

var MaxArea int = 10 // 最大视距

type Pos [2]int

// 双链表结构
type Node struct {
	EntityId int
	PreNode  *Node
	NextNode *Node
	Pos      Pos
	NodeType nodeType
}

type SceneAOIManagerByLink struct {
	SceneId       int
	Length        int
	Width         int
	XNodeList     *Node
	YNodeList     *Node
	XLandMarkList []*Node //地标节点，用于二分查找
	YLandMarkList []*Node
}

func initNode(EntityId, X, Y int, NodeType nodeType) *Node {
	node := &Node{
		EntityId: EntityId,
		Pos:      Pos{X, Y},
		NodeType: NodeType,
	}

	return node
}

func minInt(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func maxInt(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

func disInt(node1, node2 *Node) int {
	return (node1.Pos[0]-node2.Pos[0])*(node1.Pos[0]-node2.Pos[0]) + (node1.Pos[1]-node2.Pos[1])*(node1.Pos[1]-node2.Pos[1])
}

func (scene *SceneAOIManagerByLink) initLandMarkList(direct Direction) {
	var FirstNode *Node
	LandMarkList := make([]*Node, 0)

	if direct == XDirect {
		FirstNode = scene.XNodeList
	} else {
		FirstNode = scene.YNodeList
	}

	LandMarkList = append(LandMarkList, FirstNode)
	x := FirstNode.Pos[0]
	y := FirstNode.Pos[1]

	for {
		x = x + MaxArea*int(direct^YDirect)
		y = y + MaxArea*int(direct^XDirect)

		if x > scene.Length+MaxArea || y > scene.Width+MaxArea {
			break
		}

		node := initNode(0, x, y, LandMark)
		FirstNode.NextNode = node
		node.PreNode = FirstNode
		FirstNode = node

		LandMarkList = append(LandMarkList, node)
	}

	if direct == XDirect {
		scene.XLandMarkList = LandMarkList
	} else {
		scene.YLandMarkList = LandMarkList
	}
}

func (scene *SceneAOIManagerByLink) Travel() {
	fmt.Printf("Pos info:\nX: ")
	node := scene.XNodeList
	for node != nil {
		fmt.Printf("(%d, %d) ", node.Pos[0], node.Pos[1])
		node = node.NextNode
	}
	fmt.Printf("\nY: ")
	node = scene.YNodeList
	for node != nil {
		fmt.Printf("(%d, %d) ", node.Pos[0], node.Pos[1])
		node = node.NextNode
	}
	fmt.Println()
}

func (scene *SceneAOIManagerByLink) getLandMark(node *Node, direct Direction) int {
	var LandMarkList []*Node
	if direct == XDirect {
		LandMarkList = scene.XLandMarkList
	} else {
		LandMarkList = scene.YLandMarkList
	}
	start := 0
	end := len(LandMarkList) - 1
	for start <= end {
		mid := (start + end) / 2
		tmp := LandMarkList[mid]
		if tmp.Pos[direct] == node.Pos[direct] {
			return mid
		} else if tmp.Pos[direct] < node.Pos[direct] {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	return end
}

func (scene *SceneAOIManagerByLink) playerEnter(player *Player, start *Node, direct Direction) {
	// X方向节点搜索
	i := 0
	var nodeMap map[int]*Node
	if direct == XDirect {
		nodeMap = map[int]*Node{
			0: player.XLeftNode,
			1: player.XNode,
			2: player.XRightNode,
		}
	} else {
		nodeMap = map[int]*Node{
			0: player.YLeftNode,
			1: player.YNode,
			2: player.YRightNode,
		}
	}

	tmpNode := nodeMap[i]
	for i < 3 && start != nil {
		if start.Pos[direct] < tmpNode.Pos[direct] {
			// 别人进入我的视野中, i=1,2 均可
			if i > 0 && start.NodeType == NormalNode {
				otherPlayerId := start.EntityId
				otherPlayer := AllPlayerManager[otherPlayerId]

				if disInt(otherPlayer.XNode, player.XNode) < player.Area*player.Area {
					_, ok := player.Watch[otherPlayer.EntityId]
					if !ok {
						player.Watch[otherPlayer.EntityId] = true
						otherPlayer.Watched[player.EntityId] = true
						otherPlayer.enterEvent(player)
					}
				}
			}

			// 我进入别人视野中, i=0,1 均可
			if i < 2 && start.NodeType == LeftSentinel {
				otherPlayerId := start.EntityId
				otherPlayer := AllPlayerManager[otherPlayerId]
				if disInt(otherPlayer.XNode, player.XNode) < otherPlayer.Area*otherPlayer.Area {
					_, ok := player.Watched[otherPlayer.EntityId]
					if !ok {
						player.Watched[otherPlayer.EntityId] = true
						otherPlayer.Watch[player.EntityId] = true
						player.enterEvent(otherPlayer)
					}
				}
			}

			start = start.NextNode
		} else {
			tmpNode.NextNode = start
			tmpNode.PreNode = start.PreNode

			start.PreNode.NextNode = tmpNode
			start.PreNode = tmpNode

			i++
			tmpNode = nodeMap[i]
		}
	}
}

func (scene *SceneAOIManagerByLink) Enter(player *Player) {
	index := maxInt(scene.getLandMark(player.XNode, XDirect), 0)
	// 考虑到边界情况，因此搜索起点往前移动一个点
	if index > 0 {
		index = index - 1
	}
	nodeX := scene.XLandMarkList[index]
	scene.playerEnter(player, nodeX, XDirect)

	index = maxInt(scene.getLandMark(player.YNode, YDirect), 0)
	// 考虑到边界情况，因此搜索起点往前移动一个点
	if index > 0 {
		index = index - 1
	}
	nodeY := scene.YLandMarkList[index]
	scene.playerEnter(player, nodeY, YDirect)
}

func (scene *SceneAOIManagerByLink) removeNode(node *Node) {
	pre := node.PreNode
	pre.NextNode = node.NextNode
	node.NextNode.PreNode = pre
	node.NextNode = nil
	node.PreNode = nil
}

func (scene *SceneAOIManagerByLink) insertAfterNode(preNode *Node, node *Node) {
	node.NextNode = preNode.NextNode
	node.PreNode = preNode

	preNode.NextNode.PreNode = node
	preNode.NextNode = node
}

func (scene *SceneAOIManagerByLink) insertBeforeNode(nextNode *Node, node *Node) {
	node.NextNode = nextNode
	node.PreNode = nextNode.PreNode

	nextNode.PreNode.NextNode = node
	nextNode.PreNode = node
}

func (scene *SceneAOIManagerByLink) Leave(player *Player) {
	scene.removeNode(player.XNode)
	scene.removeNode(player.XLeftNode)
	scene.removeNode(player.XRightNode)

	scene.removeNode(player.YNode)
	scene.removeNode(player.YLeftNode)
	scene.removeNode(player.YRightNode)
}

func (scene *SceneAOIManagerByLink) forwardNode(node *Node, direct Direction) {
	tmp := node.NextNode
	for {
		if node.Pos[direct] < tmp.Pos[direct] {
			break
		}

		playerId := node.EntityId
		player := AllPlayerManager[playerId]
		if node.NodeType == NormalNode {
			if tmp.NodeType == LeftSentinel {
				// player进入别人视野中
				player.addNode(tmp)
			}

			if tmp.NodeType == RightSentinel {
				// player离开别人视野中
				player.leaveNode(tmp)
			}
		}

		if node.NodeType == LeftSentinel && tmp.NodeType == NormalNode {
			// 别人离开我的视野中
			player.leaveNode(tmp)
		}

		if node.NodeType == RightSentinel && tmp.NodeType == NormalNode {
			// 别人进入我的视野中
			player.addNode(tmp)
		}
		tmp = tmp.NextNode
	}
	scene.removeNode(node)
	scene.insertBeforeNode(tmp, node)
}

func (scene *SceneAOIManagerByLink) backNode(node *Node, direct Direction) {
	tmp := node.PreNode
	for {
		if node.Pos[direct] > tmp.Pos[direct] {
			break
		}

		playerId := node.EntityId
		player := AllPlayerManager[playerId]
		if node.NodeType == NormalNode {
			if tmp.NodeType == LeftSentinel {
				// node 节点离开别人视野中
				player.leaveNode(tmp)
			}

			if tmp.NodeType == RightSentinel {
				// node 节点进入别人视野中
				player.addNode(tmp)
			}
		}

		if node.NodeType == LeftSentinel && tmp.NodeType == NormalNode {
			// 别人进入我的视野中
			player.addNode(tmp)
		}

		if node.NodeType == RightSentinel && tmp.NodeType == NormalNode {
			// 别人离开我的视野中
			player.leaveNode(tmp)
		}
		tmp = tmp.PreNode
	}

	scene.removeNode(node)
	scene.insertAfterNode(tmp, node)
}

func (scene *SceneAOIManagerByLink) Move(player *Player, x, y int) {
	dir := x - player.XNode.Pos[0]
	player.XNode.Pos[0], player.XNode.Pos[1] = x, y
	player.XLeftNode.Pos[0], player.XLeftNode.Pos[1] = maxInt(player.XNode.Pos[0]-player.Area, 0), y
	player.XRightNode.Pos[0], player.XRightNode.Pos[1] = minInt(player.XNode.Pos[0]+player.Area, scene.Length), y
	if dir > 0 {
		// 注意顺序正确性
		scene.forwardNode(player.XRightNode, XDirect)
		scene.forwardNode(player.XNode, XDirect)
		scene.forwardNode(player.XLeftNode, XDirect)
	} else if dir < 0 {
		scene.backNode(player.XLeftNode, XDirect)
		scene.backNode(player.XNode, XDirect)
		scene.backNode(player.XRightNode, XDirect)
	} else {
	}

	// Y aliax

	dir = y - player.YNode.Pos[0]
	player.YNode.Pos[0], player.YNode.Pos[1] = x, y
	player.YLeftNode.Pos[0], player.YLeftNode.Pos[1] = x, maxInt(player.YNode.Pos[1]-player.Area, 0)
	player.YRightNode.Pos[0], player.YRightNode.Pos[1] = x, minInt(player.YNode.Pos[1]+player.Area, scene.Width)
	if dir > 0 {
		// 注意顺序正确性
		scene.forwardNode(player.YRightNode, YDirect)
		scene.forwardNode(player.YNode, YDirect)
		scene.forwardNode(player.YLeftNode, YDirect)
	} else if dir < 0 {
		scene.backNode(player.YLeftNode, YDirect)
		scene.backNode(player.YNode, YDirect)
		scene.backNode(player.YRightNode, YDirect)
	} else {
	}
}

func NewSceneAOIManagerByLink(SceneId int, Length int, Width int) *SceneAOIManagerByLink {
	XFirstNode := initNode(0, -MaxArea, 0, LandMark)
	YFirstNode := initNode(0, 0, -MaxArea, LandMark)

	scene := &SceneAOIManagerByLink{
		SceneId:   SceneId,
		Length:    Length,
		Width:     Width,
		XNodeList: XFirstNode,
		YNodeList: YFirstNode,
	}

	scene.initLandMarkList(XDirect)
	scene.initLandMarkList(YDirect)
	AllSceneManagerByLink[SceneId] = scene
	return scene
}

type Player struct {
	EntityId   int
	Area       int // 视距
	XNode      *Node
	XLeftNode  *Node // 左视距
	XRightNode *Node // 右视距

	YNode      *Node
	YLeftNode  *Node // 左视距
	YRightNode *Node // 右视距

	Watch   map[int]bool
	Watched map[int]bool
}

func NewPlayer(EntityId, Area int) *Player {
	if Area >= MaxArea {
		panic(fmt.Sprintf("max area is %d, your area is %d", MaxArea, Area))
	}
	player := &Player{
		EntityId: EntityId,
		Area:     Area,
		Watch:    make(map[int]bool),
		Watched:  make(map[int]bool),
	}

	AllPlayerManager[EntityId] = player
	return player
}

func (player *Player) leaveNode(node *Node) {
	// node 节点离开别人视野中
	otherPlayerId := node.EntityId
	otherPlayer := AllPlayerManager[otherPlayerId]
	_, ok := player.Watched[otherPlayerId]
	if ok {
		// 表示曾经在视野中
		delete(player.Watched, otherPlayerId)
		delete(otherPlayer.Watch, player.EntityId)
		player.leaveEvent(otherPlayer)
	}
}

func (player *Player) addNode(node *Node) {
	otherPlayerId := node.EntityId
	otherPlayer := AllPlayerManager[otherPlayerId]
	if disInt(player.XNode, otherPlayer.XNode) <= otherPlayer.Area*otherPlayer.Area {
		_, ok := player.Watched[otherPlayerId]
		if !ok {
			// 防止多次enter
			player.Watched[otherPlayerId] = true
			otherPlayer.Watch[player.EntityId] = true
			player.enterEvent(otherPlayer)
		}
	}
}

func (player *Player) enterEvent(otherPlayer *Player) {
	fmt.Printf("%d: %d Enter\n", otherPlayer.EntityId, player.EntityId)
}

func (player *Player) leaveEvent(otherPlayer *Player) {
	fmt.Printf("%d: %d Leave\n", otherPlayer.EntityId, player.EntityId)
}

func (player *Player) Enter(SceneId int, x, y int) {
	scene := AllSceneManagerByLink[SceneId]
	if scene == nil {
		return
	}

	if x < 0 || x > scene.Length {
		panic(fmt.Sprintf("x valid range is (0, %d), your x is %d", scene.Length, x))
	}

	if y < 0 || y > scene.Width {
		panic(fmt.Sprintf("y valid range is (0, %d), your y is %d", scene.Width, y))
	}

	player.XNode = initNode(player.EntityId, x, y, NormalNode)
	player.XLeftNode = initNode(player.EntityId, maxInt(x-player.Area, 0), y, LeftSentinel)
	player.XRightNode = initNode(player.EntityId, minInt(x+player.Area, scene.Length), y, RightSentinel)

	player.YNode = initNode(player.EntityId, x, y, NormalNode)
	player.YLeftNode = initNode(player.EntityId, x, maxInt(y-player.Area, 0), LeftSentinel)
	player.YRightNode = initNode(player.EntityId, x, minInt(y+player.Area, scene.Width), RightSentinel)
	scene.Enter(player)
}

func (player *Player) Leave(SceneId int) {
	scene := AllSceneManagerByLink[SceneId]
	if scene == nil {
		return
	}

	for otherPlayerId := range player.Watch {
		otherPlayer := AllPlayerManager[otherPlayerId]
		delete(otherPlayer.Watched, player.EntityId)
		otherPlayer.leaveEvent(player)
	}

	for otherPlayerId := range player.Watched {
		otherPlayer := AllPlayerManager[otherPlayerId]
		delete(otherPlayer.Watch, player.EntityId)
		player.leaveEvent(otherPlayer)
	}

	// 清空缓存
	player.Watch = make(map[int]bool)
	player.Watched = make(map[int]bool)

	scene.Leave(player)
}

func (Player *Player) Move(SceneId int, x, y int) {
	scene := AllSceneManagerByLink[SceneId]
	if scene == nil {
		return
	}

	scene.Move(Player, x, y)
}

var AllSceneManagerByLink map[int]*SceneAOIManagerByLink
var AllPlayerManager map[int]*Player

func init() {
	AllSceneManagerByLink = make(map[int]*SceneAOIManagerByLink)
	AllPlayerManager = make(map[int]*Player)
}
