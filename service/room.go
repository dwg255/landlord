package service

import (
	"sync"
	"github.com/astaxie/beego/logs"
)

var (
	roomManager = RoomManager{
		Rooms: map[int]*Room{
			1: {
				RoomId:      1,
				AllowRobot:  true,
				EntranceFee: 200,
				Tables:      make(map[TableId]*Table),
			},
			2: {
				RoomId:      2,
				AllowRobot:  false,
				EntranceFee: 200,
				Tables:      make(map[TableId]*Table),
			},
		},
	}
)

type RoomId int

type RoomManager struct {
	Lock       sync.RWMutex
	Rooms      map[int]*Room
	TableIdInc TableId
}

type Room struct {
	RoomId      RoomId
	Lock        sync.RWMutex
	AllowRobot  bool
	Tables      map[TableId]*Table
	EntranceFee int
}

//新建牌桌
func (r *Room) newTable(client *Client) (table *Table) {
	roomManager.Lock.Lock()
	defer roomManager.Lock.Unlock()

	r.Lock.Lock()
	defer r.Lock.Unlock()
	roomManager.TableIdInc = roomManager.TableIdInc + 1
	table = &Table{
		TableId:      roomManager.TableIdInc,
		Creator:      client,
		TableClients: make(map[UserId]*Client, 3),
		GameManage: &GameManage{
			FirstCallScore: client,
			Multiple:       1,
			LastShotPoker:  make([]int, 0),
			Pokers:         make([]int, 0, 54),
		},
	}
	r.Tables[table.TableId] = table
	logs.Debug("create new table ok! allow robot :%v", r.AllowRobot)
	return
}

//func init()  {
//	go func() {		//压测
//		time.Sleep(time.Second * 3)
//		for i:=0;i<1;i++{
//			client := &Client{
//				Room:       roomManager.Rooms[1],
//				HandPokers: make([]int, 0, 21),
//				UserInfo: &UserInfo{
//					UserId:   UserId(rand.Intn(10000)),
//					Username: "ROBOT-0",
//					Coin:     10000,
//				},
//				IsRobot:  true,
//				toRobot: make(chan []interface{}, 3),
//				toServer: make(chan []interface{}, 3),
//			}
//			go client.runRobot()
//			table := client.Room.newTable(client)
//			table.joinTable(client)
//		}
//	}()
//}