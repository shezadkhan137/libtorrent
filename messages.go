package libtorrent

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/torrance/libtorrent/bitfield"
	"io"
	"io/ioutil"
)

const (
	Choke = uint8(iota)
	Unchoke
	Interested
	Uninterested
	Have
	Bitfield
	Request
	Piece
	Cancel
)

type binaryDumper interface {
	BinaryDump(w io.Writer) error
}

type handshake struct {
	protocol []byte
	infoHash []byte
	peerId   []byte
}

func newHandshake(infoHash []byte) (hs *handshake) {
	hs = &handshake{
		protocol: []byte("BitTorrent protocol"),
		infoHash: infoHash,
		peerId:   PeerId,
	}
	return
}

func parseHandshake(r io.Reader) (hs *handshake, err error) {
	buf := make([]byte, 20)
	hs = new(handshake)

	// Name length
	_, err = r.Read(buf[0:1])
	if err != nil {
		return
	} else if int(buf[0]) != 19 {
		err = errors.New("Handshake halted: name length was not 19")
		return
	}

	// Protocol
	_, err = r.Read(buf[0:19])
	if err != nil {
		return
	} else if !bytes.Equal(buf[0:19], []byte("BitTorrent protocol")) {
		err = errors.New(fmt.Sprintf("Handshake halted: incompatible protocol: %s", buf[0:19]))
	}
	hs.protocol = append(hs.protocol, buf[0:19]...)

	// Skip reserved bytes
	_, err = r.Read(buf[0:8])
	if err != nil {
		return
	}

	// Info Hash
	_, err = r.Read(buf)
	if err != nil {
		return
	}
	hs.infoHash = append(hs.infoHash, buf...)

	// PeerID
	_, err = r.Read(buf)
	if err != nil {
		return
	}
	hs.peerId = append(hs.peerId, buf...)

	return
}

func (hs *handshake) BinaryDump(w io.Writer) error {
	mw := &monadWriter{w: w}
	mw.Write(uint8(19))       // Name length
	mw.Write(hs.protocol)     // Protocol name
	mw.Write(make([]byte, 8)) // Reserved 8 bytes
	mw.Write(hs.infoHash)     // InfoHash
	mw.Write(hs.peerId)       // PeerId
	return mw.err
}

func (hs *handshake) String() string {
	return fmt.Sprintf("[Handshake Protocol: %s infoHash: %x peerId: %s]", hs.protocol, hs.infoHash, hs.peerId)

}

func parsePeerMessage(r io.Reader) (msg interface{}, err error) {
	// Read message length (4 bytes)
	var length uint32
	err = binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return
	} else if length == 0 {
		// Keepalive message
		return
	} else if length > 783870696 {
		// Set limit at 2^17. Might need to revise this later
		err = errors.New(fmt.Sprintf("Message size too long: %d", length))
		return
	}

	// Read message id (1 byte)
	var id uint8
	err = binary.Read(r, binary.BigEndian, &id)
	if err != nil {
		return
	} else if id > Cancel {
		// Return error on unknown messages
		discard := make([]byte, length-1)
		_, err = r.Read(discard)
		if err != nil {
			return
		}
		err = unknownMessage{id: id, length: length}
		return
	}

	// Read payload (arbitrary size)
	payload := make([]byte, length-1)
	if length-1 > 0 {
		if _, err = r.Read(payload); err != nil {
			return
		}
	}
	payloadReader := bytes.NewReader(payload)

	switch id {
	case Choke:
		return parseChokeMessage(payloadReader)
	case Unchoke:
		return parseUnchokeMessage(payloadReader)
	case Interested:
		return parseInterestedMessage(payloadReader)
	case Have:
		return parseHaveMessage(payloadReader)
	case Bitfield:
		return parseBitfieldMessage(payloadReader)
	case Request:
		return parseRequestMessage(payloadReader)
	case Piece:
		return parsePieceMessage(payloadReader)
	}

	return
}

type chokeMessage struct{}

func parseChokeMessage(r io.Reader) (msg *chokeMessage, err error) {
	msg = new(chokeMessage)
	return
}

func (msg *chokeMessage) BinaryDump(w io.Writer) error {
	mw := monadWriter{w: w}
	mw.Write(uint32(1))
	mw.Write(Choke)
	return mw.err
}

type unchokeMessage struct{}

func parseUnchokeMessage(r io.Reader) (msg *unchokeMessage, err error) {
	msg = new(unchokeMessage)
	return
}

func (msg *unchokeMessage) BinaryDump(w io.Writer) error {
	mw := monadWriter{w: w}
	mw.Write(uint32(1))
	mw.Write(Unchoke)
	return mw.err
}

type interestedMessage struct{}

func parseInterestedMessage(r io.Reader) (msg *interestedMessage, err error) {
	msg = new(interestedMessage)
	return
}

func (msg *interestedMessage) BinaryDump(w io.Writer) error {
	mw := monadWriter{w: w}
	mw.Write(uint32(1))
	mw.Write(Interested)
	return mw.err
}

type haveMessage struct {
	pieceIndex uint32
}

func parseHaveMessage(r io.Reader) (msg *haveMessage, err error) {
	msg = new(haveMessage)
	mw := monadReader{r: r}
	mw.Read(&msg.pieceIndex)
	return msg, mw.err
}

func (msg *haveMessage) BinaryDump(w io.Writer) error {
	mw := monadWriter{w: w}
	mw.Write(uint32(5))
	mw.Write(Have)
	mw.Write(msg.pieceIndex)
	return mw.err
}

type bitfieldMessage struct {
	bitf *bitfield.Bitfield
}

func parseBitfieldMessage(r io.Reader) (msg *bitfieldMessage, err error) {
	bitf, err := bitfield.ParseBitfield(r)
	msg = &bitfieldMessage{bitf: bitf}
	return
}

func (msg *bitfieldMessage) BinaryDump(w io.Writer) error {
	length := uint32(msg.bitf.ByteLength() + 1)
	mw := monadWriter{w: w}
	mw.Write(length)
	mw.Write(Bitfield)
	mw.Write(msg.bitf.Bytes())
	return mw.err
}

func (msg *bitfieldMessage) String() string {
	return "Bitfield message"
}

type requestMessage struct {
	pieceIndex  uint32
	blockOffset uint32
	blockLength uint32
}

func parseRequestMessage(r io.Reader) (msg *requestMessage, err error) {
	msg = new(requestMessage)
	mr := &monadReader{r: r}
	mr.Read(&msg.pieceIndex)
	mr.Read(&msg.blockOffset)
	mr.Read(&msg.blockLength)
	return msg, mr.err
}

func (msg requestMessage) BinaryDump(w io.Writer) (err error) {
	mw := &monadWriter{w: w}
	mw.Write(uint32(13)) // Length: status + 12 byte payload
	mw.Write(Request)    // Message id
	mw.Write(msg.pieceIndex)
	mw.Write(msg.blockOffset)
	mw.Write(msg.blockLength)
	return mw.err
}

type pieceMessage struct {
	pieceIndex  uint32
	blockOffset uint32
	data        []byte
}

func parsePieceMessage(r io.Reader) (msg *pieceMessage, err error) {
	msg = new(pieceMessage)
	mr := &monadReader{r: r}
	mr.Read(&msg.pieceIndex)
	mr.Read(&msg.blockOffset)
	if err = mr.err; err != nil {
		return
	}
	msg.data, err = ioutil.ReadAll(r)
	return
}

func (msg *pieceMessage) BinaryDump(w io.Writer) error {
	length := uint32(len(msg.data) + 9)
	mw := monadWriter{w: w}
	mw.Write(length)
	mw.Write(Piece)
	mw.Write(msg.pieceIndex)
	mw.Write(msg.blockOffset)
	mw.Write(msg.data)
	return mw.err
}

type keepAliveMessage struct{}

func (msg *keepAliveMessage) BinaryDump(w io.Writer) error {
	mw := monadWriter{w: w}
	mw.Write(uint32(0))
	return mw.err
}

type unknownMessage struct {
	id     uint8
	length uint32
}

func (e unknownMessage) Error() string {
	return fmt.Sprintf("Unknown message id: %d, length: %d", e.id, e.length)
}
