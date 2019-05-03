package ldapschemaparser

import (
	"strconv"
)

func parseOIDLength(oidLengthText string) (oid string, length int32) {
	d := []rune(oidLengthText)
	var leftPidx, rightPidx int
	for idx, ch := range d {
		switch ch {
		case '{':
			if leftPidx == 0 {
				leftPidx = idx
			}
		case '}':
			rightPidx = idx
		}
	}
	if leftPidx <= 0 {
		oid = oidLengthText
		return
	}
	oid = string(d[0:leftPidx])
	if rightPidx > leftPidx {
		lenText := string(d[leftPidx:rightPidx])
		if v, err := strconv.ParseInt(lenText, 10, 31); nil == err {
			length = int32(v)
		}
	}
	return
}
