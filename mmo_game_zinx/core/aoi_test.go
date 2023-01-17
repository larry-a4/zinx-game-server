package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	//初始化AOIManager
	aoiManager := NewAOIManager(0, 250, 5, 0, 250, 5)

	//打印
	fmt.Println(aoiManager)
}

func TestAOIManagerSurroudingGridsByGid(t *testing.T) {
	//初始化AOIManager
	aoiManager := NewAOIManager(0, 250, 5, 0, 250, 5)

	for k := range aoiManager.grids {
		//得到当前gID的周边九宫格
		grids := aoiManager.GetSurroundingGridsByGid(k)
		fmt.Println("gid: ", k, "grids len = ", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Println("surrounding grids are ", gIDs)
	}
}
