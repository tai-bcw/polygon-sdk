package ibft

import (
	"fmt"
	"testing"

	"github.com/0xPolygon/polygon-sdk/blockchain"
	"github.com/0xPolygon/polygon-sdk/consensus"
	"github.com/0xPolygon/polygon-sdk/consensus/ibft/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	i := newMockIbftFraming(t, []string{"A", "B", "C", "D"}, "A")
	fmt.Println("testMock: ", i)
	assert.NotNil(t, i)
	assert.IsType(t, &mockIbft{}, i)
}

func newMockIbftFraming(t *testing.T, accounts []string, account string) *mockIbft {
	pool := newTesterAccountPool()
	pool.add(accounts...)

	m := &mockIbft{
		t:          t,
		pool:       pool,
		blockchain: blockchain.TestBlockchain(t, pool.genesis()),
		respMsg:    []*proto.MessageReq{},
	}

	var addr *testerAccount
	if account == "" {
		// account not in validator set, create a new one that is not part
		// of the genesis
		pool.add("xx")
		addr = pool.get("xx")
	} else {
		addr = pool.get(account)
	}
	ibft := &Ibft{
		logger:           hclog.NewNullLogger(),
		config:           &consensus.Config{},
		blockchain:       m,
		validatorKey:     addr.priv,
		validatorKeyAddr: addr.Address(),
		closeCh:          make(chan struct{}),
		updateCh:         make(chan struct{}),
		operator:         &operator{},
		state:            newState(),
		epochSize:        DefaultEpochSize,
	}

	// by default set the state to (1, 0)
	ibft.state.view = proto.ViewMsg(1, 0)

	m.Ibft = ibft

	assert.NoError(t, ibft.setupSnapshot())
	assert.NoError(t, ibft.createKey())

	// set the initial validators frrom the snapshot
	ibft.state.validators = pool.ValidatorSet()

	m.Ibft.transport = m
	return m
}

// func TestFactoryBaseConsensus(t *testing.T) (consensus.Consensus, error) {

// }
