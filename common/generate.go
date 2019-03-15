package common

import (
	"os"
	"github.com/astaxie/beego/logs"
	"strconv"
	"encoding/json"
)

func write() {
	path := "rule.json"
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			generate()
		} else {
			logs.Error("rule.json err:%v", err)
			return
		}
	}
}

//生成连续num个的单牌的顺子
func generateSeq(num int, seq []string) (res []string) {
	for i, _ := range seq {
		if i+num > 12 {
			break
		}
		var sec string
		for j := i; j < i+num; j++ {
			sec += seq[j]
		}
		res = append(res, sec)
	}
	return
}

//生成num个不同单的组合
func combination(seq []string, num int) (comb []string) {
	if num == 0 {
		panic("generate err , combination count can not be 0")
	}
	if len(seq) < num {
		logs.Error("seq: %v,num:%d",seq,num)
		return
		//panic("generate err , seq length less than num")
	}
	if num == 1 {
		return seq
	}
	if len(seq) == num {
		allSingle := ""
		for _,single := range seq{
			allSingle += single
		}
		return []string{allSingle}
	}
	noFirst := combination(seq[1:],num)
	hasFirst := []string(nil)
	for _,comb := range combination(seq[1:],num-1) {
		hasFirst = append(hasFirst, string(seq[0])+comb)
	}
	comb = append(comb, noFirst...)
	comb = append(comb, hasFirst...)
	return
}

func generate() {
	CARDS := "34567890JQKA2"
	RULE := map[string][]string{}
	RULE["single"] = []string{}
	RULE["pair"] = []string{}
	RULE["trio"] = []string{}
	RULE["bomb"] = []string{}
	for _, c := range CARDS {
		card := string(c)
		RULE["single"] = append(RULE["single"], card)
		RULE["pair"] = append(RULE["pair"], card+card)
		RULE["trio"] = append(RULE["trio"], card+card+card)
		RULE["bomb"] = append(RULE["bomb"], card+card+card+card)
	}
	for _, num := range []int{5, 6, 7, 8, 9, 10, 11, 12} {
		RULE["seq_single"+strconv.Itoa(num)] = generateSeq(num, RULE["single"])
	}
	for _, num := range []int{3, 4, 5, 6, 7, 8, 9, 10} {
		RULE["seq_pair"+strconv.Itoa(num)] = generateSeq(num, RULE["pair"])
	}
	for _, num := range []int{2, 3, 4, 5, 6} {
		RULE["seq_trio"+strconv.Itoa(num)] = generateSeq(num, RULE["trio"])
	}
	RULE["single"] = append(RULE["single"], "w")
	RULE["single"] = append(RULE["single"], "W")
	RULE["rocket"] = append(RULE["rocket"], "Ww")

	RULE["trio_single"] = make([]string, 0)
	RULE["trio_pair"] = make([]string, 0)

	for _, t := range RULE["trio"] {
		for _, s := range RULE["single"] {
			if s[0] != t[0] {
				RULE["trio_single"] = append(RULE["trio_single"], t+s)
			}
		}
		for _, p := range RULE["pair"] {
			if p[0] != t[0] {
				RULE["trio_pair"] = append(RULE["trio_pair"], t+p)
			}
		}
	}
	for _, num := range []int{2, 3, 4, 5} {
		seqTrioSingle := []string(nil)
		seqTrioPair := []string(nil)
		for _, seqTrio := range RULE["seq_trio"+strconv.Itoa(num)] {
			seq := make([]string, len(RULE["single"]))
			copy(seq, RULE["single"])
			for i := 0; i < len(seqTrio); i = i + 3 {
				for k, v := range seq {
					if v[0] == seqTrio[i] {
						copy(seq[k:], seq[k+1:])
						seq = seq[:len(seq)-1 ]
						break
					}
				}
			}
			for _, singleCombination := range combination(seq, len(seqTrio)/3) {
				seqTrioSingle = append(seqTrioSingle, seqTrio+singleCombination)
				var hasJoker bool
				for _, single := range singleCombination {
					if single == 'w' || single == 'W' {
						hasJoker = true
					}
				}
				if !hasJoker {
					seqTrioPair = append(seqTrioPair, seqTrio+singleCombination+singleCombination)
				}
			}
		}
		RULE["seq_trio_single"+strconv.Itoa(num)] = seqTrioSingle
		RULE["seq_trio_pair"+strconv.Itoa(num)] = seqTrioPair
	}

	RULE["bomb_single"] = []string(nil)
	RULE["bomb_pair"] = []string(nil)
	for _, b := range RULE["bomb"] {
		seq := make([]string, len(RULE["single"]))
		copy(seq, RULE["single"])
		for i, single := range seq {
			if single[0] == b[0] {
				copy(seq[i:], seq[i+1:])
				seq = seq[:len(seq)-1]
			}
		}
		for _,comb := range combination(seq,2){
			RULE["bomb_single"] = append(RULE["bomb_single"],b+comb)
			if comb[0] != 'w' && comb[0] != 'W' && comb[1] != 'w' && comb[1] != 'W' {
				RULE["bomb_pair"] = append(RULE["bomb_pair"],b+comb+comb)
			}
		}
	}

	res,err := json.Marshal(RULE)
	if err != nil {
		panic("json marsha1 RULE err :" + err.Error())
	}
	file, err := os.Create("rule.json")
	defer func(){
		err = file.Close()
		if err != nil {
			logs.Error("generate err: %v",err)
		}
	}()
	if err != nil {
		panic("create rule.json err:" + err.Error())
	}
	_,err = file.Write(res)
	if err != nil {
		panic("create rule.json err:" + err.Error())
	}
}
