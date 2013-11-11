package libtorrent

import (
	"errors"
	"github.com/torrance/libtorrent/metainfo"
	"sync"
	"time"
)

type PieceBlock struct {
	pieceIndex  int
	offset      int
	length      int
	timeRequest time.Time
}

type PiecePicker interface {
	NewRequests(st swarmTally) (pieceBlocks []PieceBlock)
	PieceReceived(msg pieceMessage) (isComplete bool, err error)
}

type BasicPiecePicker struct {
	meta           *metainfo.Metainfo
	activeRequests map[int][]PieceBlock
	totalLength    int
	pieceLength    int
	mutex          sync.RWMutex
}

func NewPiecePicker(meta *metainfo.Metainfo, totalLength int) (piecePicker *BasicPiecePicker, err error) {
	piecePicker = &BasicPiecePicker{
		meta:           meta,
		activeRequests: make(map[int][]PieceBlock, meta.PieceCount),
		pieceLength:    int(meta.PieceLength),
		totalLength:    totalLength,
	}
	return
}

func (b *BasicPiecePicker) NewRequests(st swarmTally) (pieceBlocks []PieceBlock) {
	// Identify what pieces we need

	b.mutex.Lock()
	needs := st.GetNeeds()

	logger.Debug("Still Need: %d", needs)

	if needs == nil {
		return
	}

	for _, pieceIndex := range needs {
		// Pieces we still need
		if pbs, ok := b.activeRequests[pieceIndex]; !ok && len(b.activeRequests) < 10 {
			// This piece has not been requested
			tmpPieceBlocks := b.pieceToPieceBlocks(pieceIndex)
			b.activeRequests[pieceIndex] = tmpPieceBlocks
			pieceBlocks = append(pieceBlocks, tmpPieceBlocks...)
		} else {
			for i, pb := range pbs {
				// Piece has been requested, check how long ago it was
				// And re request if necessary
				// TODO: Add some randomness(?)
				if thirtySecondsAgo := time.Now().Add(-time.Second * 30); pb.timeRequest.Before(thirtySecondsAgo) {
					pb.timeRequest = time.Now()
					pbs[i] = pb
					pieceBlocks = append(pieceBlocks, pb)
				}
			}
		}
	}

	b.mutex.Unlock()
	return
}

func (b *BasicPiecePicker) pieceToPieceBlocks(pieceIndex int) (pieceBlocks []PieceBlock) {
	var pieceSize int
	s1 := b.pieceLength
	s2 := b.totalLength - b.pieceLength*pieceIndex
	if s1 > s2 {
		pieceSize = s2
	} else {
		pieceSize = s1
	}

	blockSize := int(0x4000)
	offset := 0

	for ; offset+blockSize < pieceSize; offset += blockSize {
		pc := PieceBlock{
			pieceIndex:  pieceIndex,
			offset:      offset,
			length:      blockSize,
			timeRequest: time.Now(),
		}
		pieceBlocks = append(pieceBlocks, pc)
	}

	pc := PieceBlock{
		pieceIndex:  pieceIndex,
		offset:      offset,
		length:      pieceSize - offset,
		timeRequest: time.Now(),
	}

	pieceBlocks = append(pieceBlocks, pc)

	return
}

func (b *BasicPiecePicker) PieceReceived(msg pieceMessage) (isComplete bool, err error) {
	// Do something
	b.mutex.Lock()
	pieceIndex := int(msg.pieceIndex)
	pbs, ok := b.activeRequests[pieceIndex]
	var newPbs []PieceBlock
	isComplete = false
	if ok {
		for _, pb := range pbs {
			if pb.offset == int(msg.blockOffset) && pb.length == len(msg.data) {
				continue
			}
			newPbs = append(newPbs, pb)
		}
		if len(newPbs) == 0 {
			delete(b.activeRequests, pieceIndex)
			isComplete = true
			b.mutex.Unlock()
			return
		}
		b.activeRequests[pieceIndex] = newPbs
		b.mutex.Unlock()
		return
	}
	err = errors.New("Could not identify the piece")
	b.mutex.Unlock()
	return
}
