package libtorrent

type PeerRequester interface {
	RequestFromPeers(swarm []*peer, pieceBlocks []PieceBlock)
}

type BasicPeerRequester struct {
}

func NewBasicPeerRequester() *BasicPeerRequester {
	return new(BasicPeerRequester)
}

func (b *BasicPeerRequester) RequestFromPeers(swarm []*peer, pieceBlocks []PieceBlock) {
	for _, pb := range pieceBlocks {
		for _, peer := range swarm {
			if peer.CheckHasPiece(pb.pieceIndex) {
				if peer.GetPeerChoking() {
					// Peer is choking us so send interested message
					if !peer.GetAmInterested() {
						peer.SetAmInterested(true)
						peer.write <- &interestedMessage{}

					}
				} else {
					// Peer is not chokeing us, so send request

					peer.write <- &requestMessage{
						pieceIndex:  uint32(pb.pieceIndex),
						blockOffset: uint32(pb.offset),
						blockLength: uint32(pb.length),
					}
				}
			}
		}
	}
}
