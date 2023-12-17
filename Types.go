package vdf

import "math/big"


type VDF struct {
	p *big.Int
	nBase int64
}

type VDFProof struct {
	Y *big.Int
	L *big.Int
}