package timezone

import (
	"strconv"
	"time"
)

type TimeZone int64

func (tz TimeZone) Name() string {
	z := int(tz.Zone())
	if z > 0 {
		return "E" + strconv.Itoa(z)
	} else if z < 0 {
		return "W" + strconv.Itoa(z)
	} else {
		return "ZERO"
	}
}

func (tz TimeZone) Zone() int32 {
	return int32(tz) / 3600
}

func (tz TimeZone) Offset() int64 {
	return int64(tz)
}

func (tz TimeZone) String() string {
	return tz.Name()
}

var _, offset = time.Now().Zone()
var local = TimeZone(offset)

func Local() TimeZone {
	return local
}

const E12 TimeZone = 3600 * 12
const E11 TimeZone = 3600 * 11
const E10 TimeZone = 3600 * 10
const E9 TimeZone = 3600 * 9
const E8 TimeZone = 3600 * 8
const E7 TimeZone = 3600 * 7
const E6 TimeZone = 3600 * 6
const E5 TimeZone = 3600 * 5
const E4 TimeZone = 3600 * 4
const E3 TimeZone = 3600 * 3
const E2 TimeZone = 3600 * 2
const E1 TimeZone = 3600 * 1
const ZERO TimeZone = 3600 * 0
const W1 TimeZone = 3600 * -1
const W2 TimeZone = 3600 * -2
const W3 TimeZone = 3600 * -3
const W4 TimeZone = 3600 * -4
const W5 TimeZone = 3600 * -5
const W6 TimeZone = 3600 * -6
const W7 TimeZone = 3600 * -7
const W8 TimeZone = 3600 * -8
const W9 TimeZone = 3600 * -9
const W10 TimeZone = 3600 * -10
const W11 TimeZone = 3600 * -11
const W12 TimeZone = 3600 * -12
