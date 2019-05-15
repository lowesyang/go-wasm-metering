package toolkit

import (
	"regexp"
	"strings"
	"unicode"
)

var immediates = ReadImmediates()

type queue struct {
	i   int
	str []string
}

func (q *queue) head() string {
	if q.i >= len(q.str) {
		return ""
	}
	return q.str[q.i]
}

func (q *queue) shift() string {
	if q.i >= len(q.str) {
		return ""
	}
	e := q.str[q.i]
	q.i += 1
	return e
}

func (q *queue) length() int {
	return len(q.str) - q.i
}

func Text2Json(text string) (res []JSON) {
	reg := regexp.MustCompile(`\s|\n`)
	textArr := &queue{str: reg.Split(text, -1)}
	for textArr.length() > 0 {
		textOp := textArr.shift()
		jsonOp := make(JSON)

		opArr := strings.Split(textOp, ".") // [type, name]
		typ := opArr[0]
		name := typ
		if len(opArr) > 1 {
			name = opArr[1]
			jsonOp["ReturnType"] = typ
		}

		jsonOp["Name"] = name

		key := name
		if name == "const" {
			key = jsonOp["ReturnType"].(string)
		}
		immediate, exist := immediates[key]

		if exist {
			jsonOp["Immediates"] = immediataryParser(immediate.(string), textArr)
		}

		res = append(res, jsonOp)
	}

	return
}

func immediataryParser(typ string, txt *queue) interface{} {
	json := make(JSON)
	switch typ {
	case "br_table":
		var dests []string
		for {
			dest := txt.head()
			if !isNumber(dest) {
				break
			}
			txt.shift()
			dests = append(dests, dest)
		}
		return dests

	case "call_indirect":
		json["Index"] = txt.shift()
		json["Reserved"] = 0
		return json
	case "memory_immediate":
		json["Flags"] = txt.shift()
		json["Offset"] = txt.shift()
		return json
	default:
		return txt.shift()
	}
}

func isNumber(s string) bool {
	for i, digit := range s {
		if !unicode.IsNumber(digit) {
			if i != 0 || i == 0 && digit != '-' {
				return false
			}
		}
	}
	return true
}
