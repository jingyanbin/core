package basal

import "github.com/jingyanbin/core/internal"

type OnceSuccess = internal.OnceSuccess

type Waiter = internal.Waiter

type Integer interface {
	internal.Integer
}

type Float interface {
	internal.Float
}

type Number interface {
	internal.Number
}

type LinkListNode = internal.LinkListNode

type LinkList = internal.LinkList

func NewLinkList() *LinkList

//func NewLinkListNode() *LinkListNode

var Compress = internal.Compress
