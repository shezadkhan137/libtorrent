package libtorrent

type Choker interface {
	ChokeUnchokePeers(swarm []*peer, st swarmTally)
}

type BasicChoker struct {
}

// ChokeUnchokePeers is responsible for choking and unchoking peers
// called every RefreshInterval
func (b *BasicChoker) ChokeUnchokePeers(swarm []*peer, st swarmTally) {

	// Unchoke interested peers
	// TODO: Implement maximum unchoked peers
	// TODO: Implement optimistic unchoking algorithm

	for _, peer := range swarm {
		if peer.GetPeerInterested() && peer.GetAmChoking() {
			logger.Debug("Unchoking peer %s", peer.name)
			peer.write <- &unchokeMessage{}
			peer.SetAmChoking(false)
		}
	}
}
