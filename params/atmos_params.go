package params

import (
	"github.com/fchnetwork/fch/common"
	"math/big"
)

// Values for FCH Genesis related to ATMOS Consensus
var (
	atmosMinDelegateNo           = 3
	atmosNetID                   = 538
	atmosGovernanceAddress       = "0x073b2D2Cf8D31be73BDb0109dA05dA8C85BEa279"
	atmosTestGovernanceAddress   = "0x02c362540efc9FA5592621C9212D0bF776732050"
	atmosBlockInterval           = uint64(3)
	atmosEpochInterval           = uint64(100)
	atmosGasLimit                = uint64(126000000)
	// TODO: This is for initial launch we need to switch to mainnet later
	atmosEthereumRPCProvider     = "https://rinkeby.infura.io"
	atmosTestEthereumRPCProvider = "https://rinkeby.infura.io"
	atmosBlockRewards            = new(big.Int).Mul(big.NewInt(888),big.NewInt(1e+18))
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
