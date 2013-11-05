package libtorrent

type Policy interface {
	ChokeUnchokePeers(swarm []*peer, st swarmTally)
	RequestBlocks(swarm []*peer, st swarmTally)
	RefreshInterval() (val int)
}

type BasicPolicy struct{}

// ChokeUnchokePeers is responsible for choking and unchoking peers
// called every RefreshInterval
func (b *BasicPolicy) ChokeUnchokePeers(swarm []*peer, st swarmTally) {
	for _, peer := range swarm {
		if peer.GetPeerInterested() && peer.GetAmChoking() {
			logger.Debug("Unchoking peer %s", peer.name)
			peer.write <- &unchokeMessage{}
			peer.SetAmChoking(false)
		}
	}
}

// RequestBlocks is responsible for sending request messages to peers
// Any game theory algs should be implemented here.
func (b *BasicPolicy) RequestBlocks(swarm []*peer, st swarmTally) {
	logger.Debug("Request Blocks called %d", st)

	// Go through st, and find the rarest and most
	// common blocks. Will request the rarest blocks first

	if len(st) > 0 {
		high := st[0]
		high_index := 0

		low := st[0]
		low_index := 0

		number_of_haves := 0

		for index, value := range st {
			if value == -1 {
				number_of_haves += 1
				continue
			}

			if value > high {
				high = value
				high_index = index
			} else if value < low {
				low = value
				low_index = low_index
			}

		}

		// If high and low ==-1  and number_of_haves == len(st)
		// We have got all of the pieces and should be called tor.State(SEEDING)
		// or we return with some value

		// Else we can implement some kind of state machiene to control beginning
		// and endgame algs.

		if number_of_haves == 0 {
			b.beginningGame()
		} else {
			b.middleGame()
		}

		logger.Debug("%d %d %d %d %d", high, high_index, low, low_index, number_of_haves)
	}
}

func (b *BasicPolicy) beginningGame() {
	logger.Debug("In the beginning of the game")

	// Choose by most common pieces
	// for _, peer := range swarm {
	// 	peer.RequestPiece(index)
	// }
}

func (b *BasicPolicy) middleGame() {
	logger.Debug("Middle of the game")

	// Choose by least common

}

//RefreshInterval returns the interval between
//how often the policys Refresh method is called.
//i.e how often choking/unchoking of peers should
//occur.
func (b *BasicPolicy) RefreshInterval() (val int) {
	return 5
}
