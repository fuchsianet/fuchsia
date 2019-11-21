package params

import (
	"github.com/fchnetwork/fch/common"
	"math/big"
)

// Values for FCH Genesis related to ATMOS Consensus
var (
	atmosMinDelegateNo           = 1
	atmosNetID                   = 4040
	atmosGovernanceAddress       = "0xe2b151d2eF8d7D3058E44b6481B06F71d38253c9"
	atmosTestGovernanceAddress   = "0x073b2D2Cf8D31be73BDb0109dA05dA8C85BEa279"
	atmosBlockInterval           = uint64(3)
	atmosEpochInterval           = uint64(100)
	atmosGasLimit                = uint64(126000000)
	atmosEthereumRPCProvider     = "https://mainnet.infura.io"
	atmosTestEthereumRPCProvider = "https://rinkeby.infura.io"
	atmosBlockRewards            = big.NewInt(1e+18)
)

func NewAtmosMinDelegateNo() int {
	return atmosMinDelegateNo
}

func NewAtmosNetID() int {
	return atmosNetID
}

func NewAtmosGovernanceAddress() common.Address {
	return common.HexToAddress(atmosGovernanceAddress)
}

func NewAtmosTestGovernanceAddress() common.Address {
	return common.HexToAddress(atmosTestGovernanceAddress)
}

func NewAtmosBlockInterval() uint64 {
	return atmosBlockInterval
}

func NewAtmosEpochInterval() uint64 {
	return atmosEpochInterval
}

func NewAtmosGasLimit() uint64 {
	return atmosGasLimit
}

func NewAtmosEthereumRPCProvider() string {
	return atmosEthereumRPCProvider
}

func NewAtmosTestEthereumRPCProvider() string {
	return atmosTestEthereumRPCProvider
}

func NewAtmosBlockRewards() *big.Int {
	return atmosBlockRewards
}

func NewFchPreAlloc() map[string]string {
	fchPreAlloc := map[string]string{}
	return fchPreAlloc
}
