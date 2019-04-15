package service

import (
	"bytes"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"landlord/common"
	"strconv"
)

const (
	writeWait      = 1 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512

	RoleFarmer   = 0
	RoleLandlord = 1
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}//不验证origin
)

type UserId int

type UserInfo struct {
	UserId   UserId `json:"user_id"`
	Username string `json:"username"`
	Coin     int    `json:"coin"`
	Role     int
}

type Client struct {
	conn       *websocket.Conn
	UserInfo   *UserInfo
	Room       *Room
	Table      *Table
	HandPokers []int
	Ready      bool
	IsCalled   bool    //是否叫完分
	Next       *Client	//链表
	IsRobot    bool
	toRobot    chan []interface{}	//发送给robot的消息
	toServer   chan []interface{}	//robot发送给服务器
}

//重置状态
func (c *Client) reset() {
	c.UserInfo.Role = 1
	c.HandPokers = make([]int, 0, 21)
	c.Ready = false
	c.IsCalled = false
}

//发送房间内已有的牌桌信息
func (c *Client) sendRoomTables() {
	res := make([][2]int, 0)
	for _, table := range c.Room.Tables {
		if len(table.TableClients) < 3 {
			res = append(res, [2]int{int(table.TableId), len(table.TableClients)})
		}
	}
	c.sendMsg([]interface{}{common.ResTableList, res})
}

func (c *Client) sendMsg(msg []interface{}) {
	if c.IsRobot {
		c.toRobot <- msg
		return
	}
	msgByte, err := json.Marshal(msg)
	if err != nil {
		logs.Error("send msg [%v] marsha1 err:%v", string(msgByte), err)
		return
	}
	err = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		logs.Error("send msg SetWriteDeadline [%v] err:%v", string(msgByte), err)
		return
	}
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		err = c.conn.Close()
		if err != nil {
			logs.Error("close client err: %v",err)
		}
	}
	_,err = w.Write(msgByte)
	if err != nil {
		logs.Error("Write msg [%v] err: %v",string(msgByte),err)
	}
	if err := w.Close(); err != nil {
		err = c.conn.Close()
		if err != nil {
			logs.Error("close err: %v",err)
		}
	}
}

//光比客户端
func (c *Client) close() {
	if c.Table != nil {
		for _, client := range c.Table.TableClients {
			if c.Table.Creator == c && c != client {
				c.Table.Creator = client
			}
			if c == client.Next {
				client.Next = nil
			}
		}
		if len(c.Table.TableClients) != 1 {
			for _, client := range c.Table.TableClients {
				if client != client.Table.Creator {
					client.Table.Creator.Next = client
				}
			}
		}
		if len(c.Table.TableClients) == 1 {
			c.Table.Creator = nil
			delete(c.Room.Tables, c.Table.TableId)
			return
		}
		delete(c.Table.TableClients, c.UserInfo.UserId)
		if c.Table.State == GamePlaying {
			c.Table.syncUser()
			//c.Table.reset()
		}
		if c.IsRobot {
			close(c.toRobot)
			close(c.toServer)
		}
	}
}

//可能是因为版本问题，导致有些未处理的error
func (c *Client) readPump() {
	defer func() {
		//logs.Debug("readPump exit")
		c.conn.Close()
		c.close()
		if c.Room.AllowRobot {
			if c.Table != nil {
				for _, client := range c.Table.TableClients {
					client.close()
				}
			}
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logs.Error("websocket user_id[%d] unexpected close error: %v", c.UserInfo.UserId, err)
			}
			return
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		var data []interface{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			logs.Error("message unmarsha1 err, user_id[%d] err:%v", c.UserInfo.UserId, err)
		} else {
			wsRequest(data, c)
		}
	}
}

//心跳
func (c *Client) Ping() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("upgrader err:%v", err)
		return
	}
	client := &Client{conn: conn, HandPokers: make([]int, 0, 21), UserInfo: &UserInfo{}}
	var userId int
	var username string
	cookie, err := r.Cookie("userid")

	if err != nil {
		logs.Error("get cookie err: %v", err)
	} else {
		userIdStr := cookie.Value
		userId, err = strconv.Atoi(userIdStr)
	}
	cookie, err = r.Cookie("username")

	if err != nil {
		logs.Error("get cookie err: %v", err)
	} else {
		username = cookie.Value
	}

	if userId != 0 && username != "" {
		client.UserInfo.UserId = UserId(userId)
		client.UserInfo.Username = username
		go client.readPump()
		go client.Ping()
		return
	}
	logs.Error("user need login first")
	client.conn.Close()
}
