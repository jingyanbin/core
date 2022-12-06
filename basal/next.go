package basal

import (
	"github.com/jingyanbin/core/internal"
)

type NextNumber = internal.NextNumber

func NewNextNumber(s string) *NextNumber {
	return internal.NewNextNumber(s)
}
