package magicTheGathering

import (
	mtg "github.com/MagicTheGathering/mtg-sdk-go"
)

const (
	version = "0.1.0"
)

func Version() string {
	return version
}

func GetMTG() {
	ca := mtg.CardArtist
	println(ca)
}

func GetSet(name string) ([]*mtg.Set, error) {
	q := mtg.NewSetQuery()
	q.Where(mtg.SetName, name) // "murders at karlov manor")
	sets, _, err := q.PageS(0, 10)
	return sets, err
}
