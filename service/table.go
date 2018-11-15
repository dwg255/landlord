package service

import (
	"sync"
	"github.com/dwg255/landlord/common"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"time"
	"sort"
	"fmt"
)

type TableId int

const (
	GameWaitting = iota
	GameCallScore
	GamePlaying
	GameEnd
)

type Table struct {
	Lock         sync.RWMutex
	TableId      TableId
	State        int
	Creator      *Client
	TableClients map[UserId]*Client
	GameManage   *GameManage
}

type GameManage struct {
	Turn             *Client
	FirstCallScore   *Client //每局轮转
	MaxCallScore     int     //最大叫分
	MaxCallScoreTurn *Client
	LastShotClient   *Client
	Pokers           []int
	LastShotPoker    []int
	Multiple         int //加倍
}

func (table *Table) allCalled() bool {
	for _, client := range table.TableClients {
		if !client.IsCalled {
			return false
		}
	}
	return true
}

func (table *Table) gameOver(client *Client) {
	coin := table.Creator.Room.EntranceFee * table.GameManage.MaxCallScore * table.GameManage.Multiple
	table.State = GameEnd
	for _, c := range table.TableClients {
		res := []interface{}{common.ResGameOver, client.UserInfo.UserId}
		if client == c {
			res = append(res, coin*2-100)
		} else {
			res = append(res, coin)
		}
		for _, cc := range table.TableClients {
			if cc != c {
				userPokers := make([]int, 0, len(cc.HandPokers)+1)
				userPokers = append(append(userPokers, int(cc.UserInfo.UserId)), cc.HandPokers...)
				res = append(res, userPokers)
			}
		}
		c.sendMsg(res)
	}
	logs.Debug("table[%d] game over", table.TableId)
}

//叫地主结束，开始游戏
func (table *Table) callEnd() {
	//无人叫地主
	table.State = GamePlaying
	table.GameManage.FirstCallScore = table.GameManage.FirstCallScore.Next
	if table.GameManage.MaxCallScoreTurn == nil || table.GameManage.MaxCallScore == 0 {
		table.GameManage.MaxCallScoreTurn = table.Creator
		table.GameManage.MaxCallScore = 1
		//return
	}
	landLord := table.GameManage.MaxCallScoreTurn
	landLord.UserInfo.Role = RoleLandlord
	table.GameManage.Turn = landLord
	for _, poker := range table.GameManage.Pokers {
		landLord.HandPokers = append(landLord.HandPokers, poker)
	}
	res := []interface{}{common.ResShowPoker, landLord.UserInfo.UserId, table.GameManage.Pokers}
	for _, c := range table.TableClients {
		c.sendMsg(res)
	}
}

func (table *Table) joinTable(c *Client) {
	table.Lock.Lock()
	defer table.Lock.Unlock()
	if len(table.TableClients) > 2 {
		logs.Error("Player[%d] JOIN Table[%d] FULL", c.UserInfo.UserId, table.TableId)
		return
	}
	if _, ok := table.TableClients[c.UserInfo.UserId]; ok {
		logs.Error("user already in this table")
		return
	}

	c.Table = table
	c.Ready = true
	for _, client := range table.TableClients {
		if client.Next == nil {
			client.Next = c
			break
		}
	}
	table.TableClients[c.UserInfo.UserId] = c
	table.syncUser()
	if len(table.TableClients) == 3 {
		c.Next = table.Creator
		table.State = GameCallScore
		table.dealPoker()
	} else if c.Room.AllowRobot {
		go table.addRobot(c.Room) //携程添加机器人，锁不冲突
		logs.Debug("robot join ok")
	}
}

func (table *Table) addRobot(room *Room) {
	logs.Debug("robot [%v] join table", fmt.Sprintf("ROBOT-%d", len(table.TableClients)))
	rand.Seed(time.Now().UnixNano())
	if len(table.TableClients) < 3 {
		client := &Client{
			Room:       room,
			HandPokers: make([]int, 0, 21),
			UserInfo: &UserInfo{
				UserId:   UserId(rand.Intn(10000)),
				Username: fmt.Sprintf("ROBOT-%d", len(table.TableClients)),
				Coin:     10000,
			},
			IsRobot:  true,
			toRobot: make(chan []interface{}, 3),
			toServer: make(chan []interface{}, 3),
		}
		go client.runRobot()
		table.joinTable(client)
	}
}

//发牌
func (table *Table) dealPoker() {
	logs.Debug("deal poker")
	table.GameManage.Pokers = make([]int, 0)
	for i := 0; i < 54; i++ {
		table.GameManage.Pokers = append(table.GameManage.Pokers, i)
	}
	table.ShufflePokers()
	for i := 0; i < 17; i++ {
		for _, client := range table.TableClients {
			client.HandPokers = append(client.HandPokers, table.GameManage.Pokers[len(table.GameManage.Pokers)-1])
			table.GameManage.Pokers = table.GameManage.Pokers[:len(table.GameManage.Pokers)-1]
		}
	}
	response := make([]interface{}, 0, 3)
	response = append(append(append(response, common.ResDealPoker), table.GameManage.FirstCallScore.UserInfo.UserId), nil)
	for _, client := range table.TableClients {
		sort.Ints(client.HandPokers)
		response[len(response)-1] = client.HandPokers
		client.sendMsg(response)
	}
}

func (table *Table) chat(client *Client, msg string) {
	res := []interface{}{common.ResChat, client.UserInfo.UserId, msg}
	for _, c := range table.TableClients {
		c.sendMsg(res)
	}
}

//重置桌子
func (table *Table) reset() {
	table.GameManage = &GameManage{
		FirstCallScore:   table.GameManage.FirstCallScore,
		Turn:             nil,
		MaxCallScore:     0,
		MaxCallScoreTurn: nil,
		LastShotClient:   nil,
		Pokers:           table.GameManage.Pokers[:0],
		LastShotPoker:    table.GameManage.LastShotPoker[:0],
		Multiple:         1,
	}
	table.State = GameCallScore
	if table.Creator != nil {
		table.Creator.sendMsg([]interface{}{common.ResRestart})
	}
	for _, c := range table.TableClients {
		c.reset()
	}
	if len(table.TableClients) == 3 {
		table.dealPoker()
	}
}

//洗牌
func (table *Table) ShufflePokers() {
	logs.Debug("ShufflePokers")
	r := rand.New(rand.NewSource(time.Now().Unix()))
	i := len(table.GameManage.Pokers)
	for i > 0 {
		randIndex := r.Intn(i)
		table.GameManage.Pokers[i-1], table.GameManage.Pokers[randIndex] = table.GameManage.Pokers[randIndex], table.GameManage.Pokers[i-1]
		i--
	}
}

//同步房间用户信息，按链表顺序
func (table *Table) syncUser() () {
	logs.Debug("sync user")
	response := make([]interface{}, 0, 3)
	response = append(append(response, common.ResJoinTable), table.TableId)
	tableUsers := make([][2]interface{}, 0, 2)
	current := table.Creator
	for i := 0; i < len(table.TableClients); i++ {
		tableUsers = append(tableUsers, [2]interface{}{current.UserInfo.UserId, current.UserInfo.Username})
		current = current.Next
	}
	response = append(response, tableUsers)
	for _, client := range table.TableClients {
		client.sendMsg(response)
	}
}
