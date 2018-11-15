package common

import (
	"os"
	"io"
	"encoding/json"
	"fmt"
)

var (
	Pokers       = make(map[string]*Combination, 16384)
	TypeToPokers = make(map[string][]*Combination, 38)
)

type Combination struct {
	Type  string
	Score int
	Poker string
}

func init() {
	path := "rule.json"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		write()
	}
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var jsonStrByte []byte
	for {
		buf := make([]byte, 1024)
		readNum, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		for i := 0; i < readNum; i++ {
			jsonStrByte = append(jsonStrByte, buf[i])
		}
		if 0 == readNum {
			break
		}
	}
	var rule = make(map[string][]string)
	err = json.Unmarshal(jsonStrByte, &rule)
	if err != nil {
		fmt.Printf("json unmarsha1 err:%v \n", err)
		return
	}
	for pokerType, pokers := range rule {
		for score, poker := range pokers {
			cards := SortStr(poker)
			p := &Combination{
				Type:  pokerType,
				Score: score,
				Poker: cards,
			}
			Pokers[cards] = p
			TypeToPokers[pokerType] = append(TypeToPokers[pokerType], p)
		}
	}
}
