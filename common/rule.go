package common

import (
	"sort"
	"github.com/astaxie/beego/logs"
)

func SortStr(pokers string) (sortPokers string) {
	runeArr := make([]int, 0)
	for _, s := range pokers {
		runeArr = append(runeArr, int(s))
	}
	sort.Ints(runeArr)
	res := make([]byte, 0)
	for _, v := range runeArr {
		res = append(res, byte(v))
	}
	return string(res)
}

// 出的牌是否在手牌中存在
func IsContains(parent, child string) (result bool) {
	for _, childCard := range child {
		inHand := false
		for i, parentCard := range parent {
			if childCard == parentCard {
				inHand = true
				tmp := []byte(parent)
				copy(tmp[i:], tmp[i+1:])
				tmp = tmp[:len(tmp)-1]
				parent = string(tmp)
				break
			}
		}
		if !inHand {
			return
		}
	}
	return true
}

//将牌编号转换为扑克牌
func ToPokers(num []int) (string) {
	totalCards := "A234567890JQK"
	res := make([]byte, 0)
	for _, poker := range num {
		if poker == 52 {
			res = append(res, 'W')
		} else if poker == 53 {
			res = append(res, 'w')
		} else {
			res = append(res, totalCards[poker%13])
		}
	}
	return string(res)
}

//将牌转换为编号
func ToPoker(card byte) (poker []int) {
	if card == 'W' {
		return []int{52}
	}
	if card == 'w' {
		return []int{53}
	}
	cards := "A234567890JQK"
	for i, c := range []byte(cards) {
		if c == card {
			return []int{i, i + 13, i + 13*2, i + 13*3}
		}
	}
	return []int{54}
}

//将机器人要出的牌转换为编号
func pokersInHand(num []int, findPokers string) (pokers []int) {
	var isInResPokers = func(poker int) bool {
		for _, p := range pokers {
			if p == poker {
				return true
			}
		}
		return false
	}

	for _, poker := range findPokers {
		poker := ToPoker(byte(poker))
	out:
		for _,pItem := range poker {
			for _, n := range num {
				if pItem == n && !isInResPokers(n) {
					pokers = append(pokers, pItem)
					break out
				}
			}
		}
	}
	return
}

//获得牌型和大小
func pokersValue(pokers string) (cardType string, score int) {
	if combination, ok := Pokers[SortStr(pokers)]; ok {
		cardType = combination.Type
		score = combination.Score
	}
	return
}

//比较牌大小,并返回是否翻倍
func ComparePoker(baseNum, comparedNum []int) (int, bool) {
	logs.Debug("comparedNum %v  %v", baseNum, comparedNum)
	if len(baseNum) == 0 || len(comparedNum) == 0 {
		if len(baseNum) == 0 && len(comparedNum) == 0 {
			return 0, false
		} else {
			if len(baseNum) != 0 {
				return -1, false
			} else {
				comparedType, _ := pokersValue(ToPokers(comparedNum))
				if comparedType == "rocket" || comparedType == "bomb" {
					return 1, true
				}
				return 1, false
			}
		}
	}
	baseType, baseScore := pokersValue(ToPokers(baseNum))
	comparedType, comparedScore := pokersValue(ToPokers(comparedNum))
	logs.Debug("compare poker %v, %v, %v, %v", baseType, baseScore, comparedType, comparedScore)
	if baseType == comparedType {
		return comparedScore - baseScore, false
	}
	if comparedType == "rocket" {
		return 1, true
	}
	if baseType == "rocket" {
		return -1, false
	}
	if comparedType == "bomb" {
		return 1, true
	}
	return 0, false
}

//查找手牌中是否有比被比较牌型大的牌
func CardsAbove(handsNum, lastShotNum []int) (aboveNum []int) {
	handCards := ToPokers(handsNum)
	turnCards := ToPokers(lastShotNum)
	cardType, cardScore := pokersValue(turnCards)
	logs.Debug("CardsAbove handsNum %v ,lastShotNum %v, handCards %v,cardType %v,turnCards %v",
		handsNum, lastShotNum, handCards, cardType, turnCards)
	if len(cardType) == 0 {
		return
	}
	for _, combination := range TypeToPokers[cardType] {
		if combination.Score > cardScore && IsContains(handCards, combination.Poker) {
			aboveNum = pokersInHand(handsNum, combination.Poker)
			return
		}
	}
	if cardType != "boom" && cardType != "rocket" {
		for _, combination := range TypeToPokers["boom"] {
			if IsContains(handCards, combination.Poker) {
				aboveNum = pokersInHand(handsNum, combination.Poker)
				return
			}
		}
	} else if IsContains(handCards, "Ww") {
		aboveNum = pokersInHand(handsNum, "Ww")
		return
	}
	return
}
