package core

import "fmt"

/*
	AOI区域管理模块
*/
type AOIManager struct {
	// 区域左边界坐标
	MinX int
	// 区域右边界坐标
	MaxX int
	// X方向的格子数量
	CountX int
	// 区域上边界坐标
	MinY int
	// 区域下边界坐标
	MaxY int
	// Y方向的格子数量
	CountY int
	// 当前区域中有哪些格子
	grids map[int]*Grid
}

// 初始化
func NewAOIManager(minX, maxX, countX, minY, maxY, countY int) *AOIManager {
	m := &AOIManager{
		MinX:   minX,
		MaxX:   maxX,
		CountX: countX,
		MinY:   minY,
		MaxY:   maxY,
		CountY: countY,
		grids:  make(map[int]*Grid),
	}

	for y := 0; y < countY; y++ {
		for x := 0; x < countX; x++ {
			//计算格子ID
			gid := countX*y + x

			//初始化gid格子
			m.grids[gid] = NewGrid(gid,
				m.MinX+x*m.gridWidth(),
				m.MinX+(x+1)*m.gridWidth(),
				m.MinY+y*m.gridHeight(),
				m.MinY+(y+1)*m.gridHeight(),
			)
		}
	}
	return m
}

func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CountX
}

func (m *AOIManager) gridHeight() int {
	return (m.MaxY - m.MinY) / m.CountY
}

// 调试使用-打印当前AOI模块
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d, Max:%d, countX:%d, MinY:%d, MaxY:%d, countY:%d\n Grids in AOIManager",
		m.MinX, m.MaxX, m.CountX, m.MinY, m.MaxY, m.CountY)

	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 获取周边九宫格信息
func (m *AOIManager) GetSurroundingGridsByGid(gID int) (grids []*Grid) {
	//判断当前gID是否在AOIManager中
	if _, ok := m.grids[gID]; !ok {
		return
	}
	//将当前gID本身加入九宫格切片中
	grids = append(grids, m.grids[gID])

	//需要通过gID得到当前格子x轴的编号 - idx = id % nx
	idx := gID % m.CountX

	//左边是否有格子？
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	if idx < m.CountX-1 {
		grids = append(grids, m.grids[gID+1])
	}

	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}
	//遍历gIDsX 集合中每个格子的gID
	for _, v := range gidsX {
		//得到y轴编号
		idy := v / m.CountY
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CountX])
		}
		if idy < m.CountY-1 {
			grids = append(grids, m.grids[v+m.CountX])
		}
	}
	return
}

// 通过坐标获取周边九宫格内全部playerID
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	//得到当前玩家GID格子id
	gID := m.GetGidByPos(x, y)
	//通过GID得到周边九宫格信息
	grids := m.GetSurroundingGridsByGid(gID)

	//将九宫格信息里全部的player加入到playerIDs
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}
	return
}

// 通过坐标获取玩家所在的gID
func (m *AOIManager) GetGidByPos(x, y float32) int {
	//id = idy * nx + idx
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridHeight()
	return idy*m.CountX + idx
}

// 添加playerID到格子
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// 移除playerID从格子
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// 获取一个格子中全部playerID
func (m *AOIManager) GetPidsByGid(pID, gID int) []int {
	return m.grids[gID].GetPlayerIDs()
}

// 通过坐标将Player添加到一个格子
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.grids[gID].Add(pID)
}

// 通过坐标将Player从一个格子中移除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.grids[gID].Remove(pID)
}
