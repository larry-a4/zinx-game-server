package znet

import (
	"fmt"
	"sync"

	"../ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的链接集合
	connLock    sync.RWMutex                  //保护链接集合的读写锁
}

//添加链接
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//将conn加入到connection manager中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(),
		" added to ConnManager succ: conn num = ", cm.Len())
}

//删除链接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(),
		" removed from ConnManager succ: conn num = ", cm.Len())
}

//根据ConnID获取链接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, fmt.Errorf("connection not found! connID = %d", connID)
	}
}

//得到当前链接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

//清除所有链接
func (cm *ConnManager) ClearConn() {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		//停止
		conn.Stop()
		//删除
		delete(cm.connections, connID)
	}

	fmt.Println("Cleared all connections succ! conn num = ", cm.Len())
}
