package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const TSE1 = 15
const TSE2 = 11
const v = "01011010001"
const g = "10011"

type Code struct {
	value string
}

func NewCode(code string) *Code {
	g := new(Code)
	edoc := []byte(code)
	var res []byte
	for i := len(edoc) - 1; i >= 0; i-- {
		res = append(res, edoc[i])
	}
	g.value = string(res)
	return g
}

func (code *Code) Len() int {
	return len(code.value)
}

func (code *Code) SigDigPos() int {
	return strings.LastIndexByte(code.value, '1')
}

func (code *Code) ShiftLeft(offset int) {
	code.value = strings.Repeat("0", offset) + code.value
}

func CodeSum(lhs, rhs Code) Code {
	res := ""
	for i := 0; i < lhs.Len() && i < rhs.Len(); i++ {
		char := "1"
		if rhs.value[i] == lhs.value[i] {
			char = "0"
		}
		res += char
	}
	res = strings.TrimRight(res, "0")
	return Code{value: res}
}

func CodeConcat(lhs, rhs Code) Code {
	return Code{rhs.value + lhs.value}
}

func GetDivRemain(lhs, rhs Code) Code {
	rSigDig := rhs.SigDigPos()
	for lSigDig := lhs.SigDigPos(); lSigDig >= rSigDig; lSigDig = lhs.SigDigPos() {
		temp := rhs
		temp.ShiftLeft(lSigDig - rSigDig)
		lhs = CodeSum(lhs, temp)

	}
	lhs.value += strings.Repeat("0", TSE1 - TSE2 - len(lhs.value))
	return lhs
}

func main() {
	fmt.Println("v:", v)
	fmt.Println("g:", g)

	gcode := *NewCode(g)
	vcode := *NewCode(v)
	mcode := *NewCode(v + strings.Repeat("0", TSE1 - TSE2))
	pcode := GetDivRemain(mcode, gcode)
	icode := CodeConcat(vcode, pcode)

	type Score struct {
		recog int
		unrecog int
	}
	var summary [TSE1]Score
	for i := int64(1); i < int64(math.Pow(2, TSE1)); i++ {
		emask := strconv.FormatInt(i, 2)
		emask = strings.Repeat("0", TSE1 - len(emask)) + emask
		bitness := strings.Count(emask, "1") - 1
		ecode := *NewCode(emask)
		rcode := CodeSum(ecode, icode)

		scode := GetDivRemain(rcode, gcode)
		if strings.Contains(scode.value, "1") {
			summary[bitness].recog += 1
		} else {
			summary[bitness].unrecog += 1
		}
	}

	recog := 0
	total := 0
	for pos, score := range summary {
		recog += score.recog
		total += score.recog + score.unrecog
		fmt.Println("Bit depth:", pos + 1, ", recognized:", score.recog, "/", score.recog + score.unrecog,
			", percentage:", float64(score.recog) / float64(score.unrecog + score.recog) * 100, "%")
	}

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("The detecting ability of the code:",
		float64(recog) / float64(total) * 100 , "%")
	fmt.Println("--------------------------------------------------------------")
}
