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
