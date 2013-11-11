package libtorrent

import (
	"errors"
	"fmt"
	"github.com/torrance/libtorrent/bitfield"
)

type swarmTally []int

//InitializeWithBitfield use just after swarm tally is initialized/
//Use to set the initial state of indexes which we have
func (st swarmTally) InitializeWithBitfield(bitf *bitfield.Bitfield) (err error) {

	if len(st) != bitf.Length() {
		err = errors.New(fmt.Sprintf("addBitfield: Supplied bitfield incorrect size, want %d, got %d", len(st), bitf.Length()))
		return
	}

	for i := 0; i < len(st); i++ {
		if bitf.Get(i) {
			st[i] = -1
		} else {
			st[i] = 0
		}
	}
	return

}

//AddBitfield used to add a new bitfield from a peer into the swarm
func (st swarmTally) AddBitfield(bitf *bitfield.Bitfield) (err error) {
	if len(st) != bitf.Length() {
		err = errors.New(fmt.Sprintf("addBitfield: Supplied bitfield incorrect size, want %d, got %d", len(st), bitf.Length()))
		return
	}

	for i := 0; i < len(st); i++ {
		if st[i] == -1 {
			// We have this piece.
			continue
		}
		if bitf.Get(i) {
			st[i]++
		}
	}
	return
}

//RemoveBitfield used when a peer leaves the swarm TODO: Check this is called in cleanup
func (st swarmTally) RemoveBitfield(bitf *bitfield.Bitfield) (err error) {
	if len(st) != bitf.Length() {
		err = errors.New(fmt.Sprintf("removeBitfield: Supplied bitfield incorrect size, want %d, got %d", len(st), bitf.Length()))
		return
	}

	for i := 0; i < len(st); i++ {
		if st[i] <= 0 {
			// We either have this piece, or something's gone wrong. Either way, leave as is.
			continue
		}
		if bitf.Get(i) {
			st[i]--
		}
	}
	return
}

//Zero zeros all values not set to -1
func (st swarmTally) Zero() {
	for i := 0; i < len(st); i++ {
		if st[i] != -1 {
			st[i] = 0
		}
	}
}

//AddIndex Adds individual index to the swarm (ie when peer sends us a HAVE message)
func (st swarmTally) AddIndex(i int) {
	if i < len(st) {
		val := st[i]
		if val != -1 {
			st[i] = val + 1
		}
	}
}

//SetWeHave indicates we now have this piece
func (st swarmTally) SetWeHave(i int) {
	if i < len(st) {
		st[i] = -1
	}
}

//GetNeeds returns a slice of the pieceIndexes which
//we still need. Nil is returned if we have all
func (st swarmTally) GetNeeds() (needs []int) {
	// This function will return nil, if we have all
	for i := 0; i < len(st); i++ {
		if st[i] != -1 {
			needs = append(needs, i)
		}
	}
	return
}

func (st swarmTally) GetMostPopularIndex() (highIndex int, highValue int) {

	for i := 0; i < len(st); i++ {
		if st[i] >= highValue && st[i] != -1 {
			highValue = st[i]
			highIndex = i
		}
	}

	return
}

func (st swarmTally) GetRarestIndex() (lowestIndex int, lowestValue int) {

	for i := 0; i < len(st); i++ {
		if st[i] <= lowestValue && st[i] != -1 {
			lowestValue = st[i]
			lowestIndex = i
		}
	}

	return
}
