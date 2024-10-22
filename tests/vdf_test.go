package vdftest

import "math/big"
import "testing"

import "github.com/sirgallo/vdf"


const N_BASE = int64(250000)
const P = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A63A3620FFFFFFFFFFFFFFFF"
const ITERATIONS = 100
const GENESIS_INPUT = "f92d7accb0004af0565148cdb657e0199574ad70dbaed4ecdfdec345943bb24a2612cfe3e45bafb59120dce8c4550640b56664072e814a55266d0441f9708650"


var gi, p *big.Int
var v *vdf.VDF


func init() {
	gi = new(big.Int)
	gi.SetString(GENESIS_INPUT, 16)

	p = new(big.Int)
	p.SetString(P, 16)

	v = vdf.NewVDF(p, N_BASE)
}


func TestVDF(t *testing.T) {
	computed, proof, computeErr := v.ComputeVDF(gi, nil)
	if computeErr != nil { t.Errorf("error computing vdf: %s\n", computeErr.Error() ) }

	isVerified := v.VerifyVDF(computed, gi, proof, nil)
	if ! isVerified { t.Error("vdf output was not verified") }

	t.Logf("input: %d\noutput: %d\nproof: %d\nisVerified: %t\n", v, computed, proof, isVerified)
}

func TestVDFSetDifficulty(t *testing.T) {
	next := gi

	for i := range make([]int, ITERATIONS) {
		computed, proof, computeErr := v.ComputeVDF(next, nil)
		if computeErr != nil { t.Errorf("error computing vdf: %s\n", computeErr.Error() )}
		
		isVerified := v.VerifyVDF(computed, next, proof, nil)
		if ! isVerified { t.Error("vdf output was not verified") }
		
		t.Logf("iteration: %d\noutput: %d\nproof: %d\nisVerified: %t\n\n", i, computed, proof, isVerified)
		next = computed
	}
}

func TestVDFDynamicDifficulty(t *testing.T) {
	next := gi

	for i := range make([]int, ITERATIONS) {
		cum := uint64(i)
		
		computed, proof, computeErr := v.ComputeVDF(next, &cum)
		if computeErr != nil { t.Errorf("error computing vdf: %s\n", computeErr.Error()) }

		isVerified := v.VerifyVDF(computed, next, proof, &cum)
		if ! isVerified { t.Error("vdf output was not verified") }

		t.Logf("iteration: %d\noutput: %d\nproof: %d\nisVerified: %t\n\n", i, computed, proof, isVerified)
		next = computed
	}

	t.Log("Done")
}