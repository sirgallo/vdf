package vdf

import "crypto/rand"
import "math"
import "math/big"


// NewVDF
//	initialize the vdf
func NewVDF(p *big.Int, nBase int64) *VDF {
	return &VDF{ p: p, nBase: nBase }
}

// ComputeVDF
//	generates both an output using repeated squaring and a proof that the computation was actually calculated.
//	the proof is generated by taking a random prime number 
func (vdf *VDF) ComputeVDF(input *big.Int, totSeq *uint64) (*big.Int, *VDFProof, error) {
	output := vdf.generateOutput(input, totSeq)
	proof, genProofErr := vdf.generateProof(input, totSeq)
	if genProofErr != nil { return nil, nil, genProofErr }

	return output, proof, nil
}

// VerifyVDF
//	verifies the output of a vdf by taking the proof and performing a partial computation to get the expected output.
//	if the partial computation produces the correct output when multiplied by the proof output and mod P, the vdf output has been verified.
//	this should be orders of magnitude faster than computing the vdf output
func (vdf *VDF) VerifyVDF(output *big.Int, input *big.Int, proof *VDFProof, totSeq *uint64) bool {
	r := vdf.calculateR(proof.L, totSeq)

	twoToNs := new(big.Int).Exp(big.NewInt(2), vdf.calculateNs(totSeq), nil)
	difftwoToNsR := new(big.Int).Sub(twoToNs, r)

	partialComputation := new(big.Int).Exp(input, difftwoToNsR, vdf.p)
	expOutput := new(big.Int).Mod(new(big.Int).Mul(partialComputation, proof.Y), vdf.p)

	return output.Cmp(expOutput) == 0
}

// generateOutput
//	perform the repeated squaring for the input and the total iterations.
//	introduces the "delay" since the calculations are computationally expensive
func (vdf *VDF) generateOutput(input *big.Int, totSeq *uint64) *big.Int {
	exponent := new(big.Int).Exp(big.NewInt(2), vdf.calculateNs(totSeq), nil)
	return new(big.Int).Exp(input, exponent, vdf.p)
}

// generateProof
//	create the proof for the verification process once the vdf has been computed.
func (vdf *VDF) generateProof(input *big.Int, totSeq *uint64) (*VDFProof, error) {
	l, err := calculateL()
	if err != nil { return nil, err }

	r := vdf.calculateR(l, totSeq)
	y := new(big.Int).Exp(input, r, vdf.p)

	return &VDFProof{ Y: y, L: l }, nil
}

// calculateL
//	generate a random prime number in a sample space that is 0 <= L <= P to create the security parameter.
func calculateL() (*big.Int, error) {
	generateRandomPrime128 := func() (*big.Int, error) {
		bitSize := 128
		for {
			randomNumber, genErr := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(bitSize)))
			if genErr != nil { return nil, genErr }
			
			randomNumber.SetBit(randomNumber, bitSize - 1, 1)
			randomNumber.SetBit(randomNumber, 0, 1)

			if randomNumber.ProbablyPrime(20) { return randomNumber, nil }
		}
	}

	l, err := generateRandomPrime128()
	if err != nil { return nil, err }

	return l, nil
}

// calculateNs
//	Ns is a dynamically adjusting N, where total successful iterations by a single writer will result in N growing by:
//	N * log2(total successful)
//	this introduces a form of "rate limiting".
func (vdf *VDF) calculateNs(totSeq *uint64) *big.Int {
	Ns := vdf.nBase
	if totSeq != nil && *totSeq > 1 {
		growthFactor := math.Log2(float64(*totSeq))
		Ns = int64((float64(vdf.nBase) * growthFactor))
	}

	return big.NewInt(Ns)
}

// calculateR
//	r is random nonce. It is just 2^N mod l
func (vdf *VDF) calculateR(l *big.Int, totSeq *uint64) *big.Int {
	return new(big.Int).Exp(big.NewInt(2), vdf.calculateNs(totSeq), l)
}


/*
NOTE:
	P = 1024-bit prime number -> preselected and publicly known
	L = The security parameter of the challenge group: 128-bit prime number -> randomly generated
	N_BASE = base number of iterations

	calculateL() -> L:

	calculateNs(totalSequential) -> Ns:
		if totalSeqential > 1:
			growthFactor = log2(totalSequential)
			return N_BASE * growthFactor
		else:
			return N_BASE

	calculateR() -> r:
		return 2^Ns mod L

	generateOutput(input) -> output:
		return input^(2^Ns) mod P

	generateProof(input) -> (y, L):
		r = 2^Ns mod L
		y = input^r mod P

		return (y, L)

	VerifyVDF(output, input, y) -> boolean:
		partial_compute = input^(2^Ns - r) mod P
		exp_output = (partial_compute * y) mod P

		return exp_output == output

	
	repeated squaring is used for generating output, and the Wesolowski proof is used for generating + verifying proofs.
*/