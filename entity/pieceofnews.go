package entity

import (
	"encoding/gob"
)

type PieceOfNews struct {
	Title  string
	Link   string
	Tags   []string
	Source string
}

func init() {
	gob.Register(PieceOfNews{})
}
