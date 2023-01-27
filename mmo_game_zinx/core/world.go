package core

import "sync"

/*
	当前游戏世界的总管理模块
*/
type WorldManager struct {
	//AOIManager 当前世界地图AOI的管理模块
	AoiMgr *AOIManager
	//当前全部在线的玩家集合
	Players map[int32]*Player
	// 保护玩家集合的锁
	pLock sync.RWMutex
}

var WorldMgrObj *WorldManager

func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_COUNT_X, AOI_MIN_Y, AOI_MAX_Y, AOI_COUNT_Y),
		Players: make(map[int32]*Player),
	}
}

func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	//将player添加在AOIManager中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

func (wm *WorldManager) RemovePlayer(pid int32) {
	player, ok := wm.Players[pid]
	if ok {
		wm.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Z)

		wm.pLock.Lock()
		delete(wm.Players, pid)
		defer wm.pLock.Unlock()
	}
}

func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}

func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)
	for _, v := range wm.Players {
		players = append(players, v)
	}
	return players
}
