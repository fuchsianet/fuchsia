// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fchnetwork/fch/common"
	"github.com/fchnetwork/fch/common/hexutil"
	"github.com/fchnetwork/fch/core"
	"github.com/fchnetwork/fch/core/rawdb"
	"github.com/fchnetwork/fch/core/state"
	"github.com/fchnetwork/fch/core/types"
	"github.com/fchnetwork/fch/internal/ethapi"
	"github.com/fchnetwork/fch/rlp"
	"github.com/fchnetwork/fch/rpc"
	"github.com/fchnetwork/fch/trie"
)

// PublicEthereumAPI provides an API to access Ethereum full node-related
// information.
type PublicEthereumAPI struct {
	e *Ethereum
}

// NewPublicEthereumAPI creates a new Ethereum protocol API for full nodes.
func NewPublicEthereumAPI(e *Ethereum) *PublicEthereumAPI {
	return &PublicEthereumAPI{e}
}

// Etherbase is the address that mining rewards will be send to
func (api *PublicEthereumAPI) Etherbase() (common.Address, error) {
	return api.e.Etherbase()
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (api *PublicEthereumAPI) Coinbase() (common.Address, error) {
	return api.Etherbase()
}

// Hashrate returns the POW hashrate
func (api *PublicEthereumAPI) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(api.e.Miner().HashRate())
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (api *PublicEthereumAPI) ChainId() hexutil.Uint64 {
	chainID := new(big.Int)
	if config := api.e.blockchain.Config(); config.IsEIP155(api.e.blockchain.CurrentBlock().Number()) {
		chainID = config.ChainID
	}
	return (hexutil.Uint64)(chainID.Uint64())
}

// PublicMinerAPI provides an API to control the miner.
// It offers only methods that operate on data that pose no security risk when it is publicly accessible.
type PublicMinerAPI struct {
	e *Ethereum
}

// NewPublicMinerAPI create a new PublicMinerAPI instance.
func NewPublicMinerAPI(e *Ethereum) *PublicMinerAPI {
	return &PublicMinerAPI{e}
}

// Mining returns an indication if this node is currently mining.
func (api *PublicMinerAPI) Mining() bool {
	return api.e.IsMining()
}

// PrivateMinerAPI provides private RPC methods to control the miner.
// These methods can be abused by external users and must be considered insecure for use by untrusted users.
type PrivateMinerAPI struct {
	e *Ethereum
}

// NewPrivateMinerAPI create a new RPC service which controls the miner of this node.
func NewPrivateMinerAPI(e *Ethereum) *PrivateMinerAPI {
	return &PrivateMinerAPI{e: e}
}

// Start starts the miner with the given number of threads. If threads is nil,
// the number of workers started is equal to the number of logical CPUs that are
// usable by this process. If mining is already running, this method adjust the
// number of threads allowed to use and updates the minimum price required by the
// transaction pool.
func (api *PrivateMinerAPI) Start(threads *int) error {
	if threads == nil {
		return api.e.StartMining(runtime.NumCPU())
	}
	return api.e.StartMining(*threads)
}

// Stop terminates the miner, both at the consensus engine level as well as at
// the block creation level.
func (api *PrivateMinerAPI) Stop() {
	api.e.StopMining()
}

// SetExtra sets the extra data string that is included when this miner mines a block.
func (api *PrivateMinerAPI) SetExtra(extra string) (bool, error) {
	if err := api.e.Miner().SetExtra([]byte(extra)); err != nil {
		return false, err
	}
	return true, nil
}

// SetGasPrice sets the minimum accepted gas price for the miner.
func (api *PrivateMinerAPI) SetGasPrice(gasPrice hexutil.Big) bool {
	api.e.lock.Lock()
	api.e.gasPrice = (*big.Int)(&gasPrice)
	api.e.lock.Unlock()

	api.e.txPool.SetGasPrice((*big.Int)(&gasPrice))
	return true
}

// SetEtherbase sets the etherbase of the miner
func (api *PrivateMinerAPI) SetEtherbase(etherbase common.Address) bool {
	api.e.SetEtherbase(etherbase)
	return true
}

// SetRecommitInterval updates the interval for miner sealing work recommitting.
func (api *PrivateMinerAPI) SetRecommitInterval(interval int) {
	api.e.Miner().SetRecommitInterval(time.Duration(interval) * time.Millisecond)
}

// GetHashrate returns the current hashrate of the miner.
func (api *PrivateMinerAPI) GetHashrate() uint64 {
	return api.e.miner.HashRate()
}

// PrivateAdminAPI is the collection of Ethereum full node-related APIs
// exposed over the private admin endpoint.
type PrivateAdminAPI struct {
	eth *Ethereum
}

// NewPrivateAdminAPI creates a new API definition for the full node private
// admin methods of the Ethereum service.
func NewPrivateAdminAPI(eth *Ethereum) *PrivateAdminAPI {
	return &PrivateAdminAPI{eth: eth}
}

// ExportChain exports the current blockchain into a local file.
func (api *PrivateAdminAPI) ExportChain(file string) (bool, error) {
	if _, err := os.Stat(file); err == nil {
		// File already exists. Allowing overwrite could be a DoS vecotor,
		// since the 'file' may point to arbitrary paths on the drive
		return false, errors.New("location would overwrite an existing file")
	}
	// Make sure we can create the file to export into
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer out.Close()

	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	// Export the blockchain
	if err := api.eth.BlockChain().Export(writer); err != nil {
		return false, err
	}
	return true, nil
}

func hasAllBlocks(chain *core.BlockChain, bs []*types.Block) bool {
	for _, b := range bs {
		if !chain.HasBlock(b.Hash(), b.NumberU64()) {
			return false
		}
	}

	return true
}

// ImportChain imports a blockchain from a local file.
func (api *PrivateAdminAPI) ImportChain(file string) (bool, error) {
	// Make sure the can access the file to import
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()

	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}

	// Run actual the import in pre-configured batches
	stream := rlp.NewStream(reader, 0)

	blocks, index := make([]*types.Block, 0, 2500), 0
	for batch := 0; ; batch++ {
		// Load a batch of blocks from the input file
		for len(blocks) < cap(blocks) {
			block := new(types.Block)
			if err := stream.Decode(block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}

		if hasAllBlocks(api.eth.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		// Import the batch and reset the buffer
		if _, err := api.eth.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}

// PublicDebugAPI is the collection of Ethereum full node APIs exposed
// over the public debugging endpoint.
type PublicDebugAPI struct {
	eth *Ethereum
}

// NewPublicDebugAPI creates a new API definition for the full node-
// related public debug methods of the Ethereum service.
func NewPublicDebugAPI(eth *Ethereum) *PublicDebugAPI {
	return &PublicDebugAPI{eth: eth}
}

// DumpBlock retrieves the entire state of the database at a given block.
func (api *PublicDebugAPI) DumpBlock(blockNr rpc.BlockNumber) (state.Dump, error) {
	if blockNr == rpc.PendingBlockNumber {
		// If we're dumping the pending state, we need to request
		// both the pending block as well as the pending state from
		// the miner and operate on those
		_, stateDb := api.eth.miner.Pending()
		return stateDb.RawDump(false, false, true), nil
	}
	var block *types.Block
	if blockNr == rpc.LatestBlockNumber {
		block = api.eth.blockchain.CurrentBlock()
	} else {
		block = api.eth.blockchain.GetBlockByNumber(uint64(blockNr))
	}
	if block == nil {
		return state.Dump{}, fmt.Errorf("block #%d not found", blockNr)
	}
	stateDb, err := api.eth.BlockChain().StateAt(block.Root())
	if err != nil {
		return state.Dump{}, err
	}
	return stateDb.RawDump(false, false, true), nil
}

// PrivateDebugAPI is the collection of Ethereum full node APIs exposed over
// the private debugging endpoint.
type PrivateDebugAPI struct {
	eth *Ethereum
}

// NewPrivateDebugAPI creates a new API definition for the full node-related
// private debug methods of the Ethereum service.
func NewPrivateDebugAPI(eth *Ethereum) *PrivateDebugAPI {
	return &PrivateDebugAPI{eth: eth}
}

// Preimage is a debug API function that returns the preimage for a sha3 hash, if known.
func (api *PrivateDebugAPI) Preimage(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	if preimage := rawdb.ReadPreimage(api.eth.ChainDb(), hash); preimage != nil {
		return preimage, nil
	}
	return nil, errors.New("unknown preimage")
}

// BadBlockArgs represents the entries in the list returned when bad blocks are queried.
type BadBlockArgs struct {
	Hash  common.Hash            `json:"hash"`
	Block map[string]interface{} `json:"block"`
	RLP   string                 `json:"rlp"`
}

// GetBadBlocks returns a list of the last 'bad blocks' that the client has seen on the network
// and returns them as a JSON list of block-hashes
func (api *PrivateDebugAPI) GetBadBlocks(ctx context.Context) ([]*BadBlockArgs, error) {
	blocks := api.eth.BlockChain().BadBlocks()
	results := make([]*BadBlockArgs, len(blocks))

	var err error
	for i, block := range blocks {
		results[i] = &BadBlockArgs{
			Hash: block.Hash(),
		}
		if rlpBytes, err := rlp.EncodeToBytes(block); err != nil {
			results[i].RLP = err.Error() // Hacky, but hey, it works
		} else {
			results[i].RLP = fmt.Sprintf("0x%x", rlpBytes)
		}
		if results[i].Block, err = ethapi.RPCMarshalBlock(block, true, true); err != nil {
			results[i].Block = map[string]interface{}{"error": err.Error()}
		}
	}
	return results, nil
}

// AccountRangeResult returns a mapping from the hash of an account addresses
// to its preimage. It will return the JSON null if no preimage is found.
// Since a query can return a limited amount of results, a "next" field is
// also present for paging.
type AccountRangeResult struct {
	Accounts map[common.Hash]*common.Address `json:"accounts"`
	Next     common.Hash                     `json:"next"`
}

func accountRange(st state.Trie, start *common.Hash, maxResults int) (AccountRangeResult, error) {
	if start == nil {
		start = &common.Hash{0}
	}
	it := trie.NewIterator(st.NodeIterator(start.Bytes()))
	result := AccountRangeResult{Accounts: make(map[common.Hash]*common.Address), Next: common.Hash{}}

	if maxResults > AccountRangeMaxResults {
		maxResults = AccountRangeMaxResults
	}

	for i := 0; i < maxResults && it.Next(); i++ {
		if preimage := st.GetKey(it.Key); preimage != nil {
			addr := &common.Address{}
			addr.SetBytes(preimage)
			result.Accounts[common.BytesToHash(it.Key)] = addr
		} else {
			result.Accounts[common.BytesToHash(it.Key)] = nil
		}
	}

	if it.Next() {
		result.Next = common.BytesToHash(it.Key)
	}

	return result, nil
}

// AccountRangeMaxResults is the maximum number of results to be returned per call
const AccountRangeMaxResults = 256

// AccountRange enumerates all accounts in the latest state
func (api *PrivateDebugAPI) AccountRange(ctx context.Context, start *common.Hash, maxResults int) (AccountRangeResult, error) {
	var statedb *state.StateDB
	var err error
	block := api.eth.blockchain.CurrentBlock()

	if len(block.Transactions()) == 0 {
		statedb, err = api.computeStateDB(block, defaultTraceReexec)
		if err != nil {
			return AccountRangeResult{}, err
		}
	} else {
		_, _, statedb, err = api.computeTxEnv(block.Hash(), len(block.Transactions())-1, 0)
		if err != nil {
			return AccountRangeResult{}, err
		}
	}

	trie, err := statedb.Database().OpenTrie(block.Header().Root)
	if err != nil {
		return AccountRangeResult{}, err
	}

	return accountRange(trie, start, maxResults)
}

// StorageRangeResult is the result of a debug_storageRangeAt API call.
type StorageRangeResult struct {
	Storage storageMap   `json:"storage"`
	NextKey *common.Hash `json:"nextKey"` // nil if Storage includes the last key in the trie.
}

type storageMap map[common.Hash]storageEntry

type storageEntry struct {
	Key   *common.Hash `json:"key"`
	Value common.Hash  `json:"value"`
}

// StorageRangeAt returns the storage at the given block height and transaction index.
func (api *PrivateDebugAPI) StorageRangeAt(ctx context.Context, blockHash common.Hash, txIndex int, contractAddress common.Address, keyStart hexutil.Bytes, maxResult int) (StorageRangeResult, error) {
	_, _, statedb, err := api.computeTxEnv(blockHash, txIndex, 0)
	if err != nil {
		return StorageRangeResult{}, err
	}
	st := statedb.StorageTrie(contractAddress)
	if st == nil {
		return StorageRangeResult{}, fmt.Errorf("account %x doesn't exist", contractAddress)
	}
	return storageRangeAt(st, keyStart, maxResult)
}

func storageRangeAt(st state.Trie, start []byte, maxResult int) (StorageRangeResult, error) {
	it := trie.NewIterator(st.NodeIterator(start))
	result := StorageRangeResult{Storage: storageMap{}}
	for i := 0; i < maxResult && it.Next(); i++ {
		_, content, _, err := rlp.Split(it.Value)
		if err != nil {
			return StorageRangeResult{}, err
		}
		e := storageEntry{Value: common.BytesToHash(content)}
		if preimage := st.GetKey(it.Key); preimage != nil {
			preimage := common.BytesToHash(preimage)
			e.Key = &preimage
		}
		result.Storage[common.BytesToHash(it.Key)] = e
	}
	// Add the 'next key' so clients can continue downloading.
	if it.Next() {
		next := common.BytesToHash(it.Key)
		result.NextKey = &next
	}
	return result, nil
}

// GetModifiedAccountsByNumber returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateDebugAPI) GetModifiedAccountsByNumber(startNum uint64, endNum *uint64) ([]common.Address, error) {
	var startBlock, endBlock *types.Block

	startBlock = api.eth.blockchain.GetBlockByNumber(startNum)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startNum)
	}

	if endNum == nil {
		endBlock = startBlock
		startBlock = api.eth.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.eth.blockchain.GetBlockByNumber(*endNum)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %d not found", *endNum)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

// GetModifiedAccountsByHash returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateDebugAPI) GetModifiedAccountsByHash(startHash common.Hash, endHash *common.Hash) ([]common.Address, error) {
	var startBlock, endBlock *types.Block
	startBlock = api.eth.blockchain.GetBlockByHash(startHash)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startHash)
	}

	if endHash == nil {
		endBlock = startBlock
		startBlock = api.eth.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.eth.blockchain.GetBlockByHash(*endHash)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %x not found", *endHash)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

func (api *PrivateDebugAPI) getModifiedAccounts(startBlock, endBlock *types.Block) ([]common.Address, error) {
	if startBlock.Number().Uint64() >= endBlock.Number().Uint64() {
		return nil, fmt.Errorf("start block height (%d) must be less than end block height (%d)", startBlock.Number().Uint64(), endBlock.Number().Uint64())
	}
	triedb := api.eth.BlockChain().StateCache().TrieDB()

	oldTrie, err := trie.NewSecure(startBlock.Root(), triedb)
	if err != nil {
		return nil, err
	}
	newTrie, err := trie.NewSecure(endBlock.Root(), triedb)
	if err != nil {
		return nil, err
	}
	diff, _ := trie.NewDifferenceIterator(oldTrie.NodeIterator([]byte{}), newTrie.NodeIterator([]byte{}))
	iter := trie.NewIterator(diff)

	var dirty []common.Address
	for iter.Next() {
		key := newTrie.GetKey(iter.Key)
		if key == nil {
			return nil, fmt.Errorf("no preimage found for hash %x", iter.Key)
		}
		dirty = append(dirty, common.BytesToAddress(key))
	}
	return dirty, nil
}
