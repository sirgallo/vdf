# vdf

### a verifiable delay function module implemented in go


## design

This `vdf` implementation uses the repeated squaring method along with `Wesolowski` proofs. Check out [design](./docs/design.md) for a more in depth explanation.


## usage

```go
package main

import "math/big"

import "github.com/sirgallo/vdf"


const N_BASE = int64(<num-of-base-iterations>)
const P = "<string-representation-of-1024-bit-prime>"
const GENESIS_INPUT = "<string-representation-of-1024-bit-input>"


func main() {
  p = new(big.Int)
  p.SetString(P, 16)

  gIn := new(big.Int)
  gIn.SetString(GENESIS_INPUT, 16)

  // initialize
  v := vdf.NewVDF(p, N_BASE)


  //========================  using set difficulty

  computed, proof, computeErr := v.ComputeVDF(gIn, nil)
  if computeErr != nil { panic(computeErr.Error()) }

  isVerified := v.VerifyVDF(computed, gIn, proof, nil)
  if ! isVerified { t.Error("vdf output was not verified") }


  //========================  using dynamically increasing difficulty

  level := uint64(1)
  
  computed, proof, computeErr := v.ComputeVDF(gIn, &level)
  if computeErr != nil { panic(computeErr.Error()) }

  isVerified := v.VerifyVDF(computed, gIn, proof, &level)
  if ! isVerified { t.Error("vdf output was not verified") }
}
```


## tests

```bash
go test -v ./tests
```


## godoc

For in depth definitions of types and functions, `godoc` can generate documentation from the formatted function comments. If `godoc` is not installed, it can be installed with the following:
```bash
go install golang.org/x/tools/cmd/godoc
```

To run the `godoc` server and view definitions for the package:
```bash
godoc -http=:6060
```

Then, in your browser, navigate to:
```
http://localhost:6060/pkg/github.com/sirgallo/vdf/
```