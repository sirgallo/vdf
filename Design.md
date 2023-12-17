# VDF Design


## Overview

`Verifiable Delay Functions` are a relatively new cryptographic primitive that take a predictable amount of time to compute, even on multi-core systems, but are easily verifiable by viewers of the result.

For computing the `vdf`, repeated squaring is chosen. No additional hashing is done on the VDF to apply randomness for security. When a entity accumulates successful computations, the difficulty of solving the next vdf increases. This is done to give other entities in a distributed network a higher chance to submit the next output in the sequence and introduces a built in rate limiting feature. Once a different entity has successfully generated the next VDF in the sequence, the entity being rate limited has the total accumulated successful computations set to $0$


## Computation

### Constants

$$
\begin{align}
  &P = \text{a large prime number that defines the group}
  &N_{base} = \text{the base number of squaring operations to use} \\
\end{align}
$$

### Variables

$$
\begin{align}
  &v_{curr} = \text{the current version output} \\
  &v_{next} = \text{the next version output from computing the vdf for the current version} \\
  &s = \text{the total number of successful sequential writes} \\
  &N_{s} = \text{N dynamically adjusted after s sequential successful writes} \\
  &L = \text{a randomly security parameter of the challenge group (a 128-bit prime)} \\
  &y= \text{the output from generating a proof based on a partial computation of the next computed}
\end{align}
$$

### Generating Output

$N_{s}$ is calculated by taking $N_{base}$, where after each sequential successul write to the state machine causes a gradual growth in the difficulty of the vdf.

$$
\begin{align}
  &N_{s} = N_{base}\times{log_{2}{s}}
\end{align}
$$

Where the computation for the `vdf` is as follows. $s$, the total successful sequential writes, is used to increase the value of $N_{base}$, where $N_{base}$ is multiplied by the logarithm of the number of successfuly sequential writes, so that each attempt to write multiple proposals to the state machine will take increasingly longer.

$$
\begin{align}
  &VDF(P, v_{curr}, s) = v_{curr}^{2^{N_{s}}}\mod{P}
\end{align}
$$

The value for $N_{base}$ is a major determinant in the total time to solve the `vdf`, so for longer delays a larger value should be selected to increase the total number of iterations required to compute the vdf output.


### Generating Proof

The [Wesolowski Proof](https://eprint.iacr.org/2018/623.pdf) has been selected to generate the proof for the computed output from the `vdf`. The proof generates a partial computation of the computed output, so any entity verifying the output only needs to compute a portion of the output based on the input instead of recomputing the entire output again.

$$
\begin{align}
  &proof(y, L) = v_{curr}^{2^{N_{s}}\mod{L}}\mod{P}
\end{align}
$$


### Verification

$$
\begin{align}
  &isVerified = (((v_{curr}^{2^{N_{s} - 2^{N_{s}}\mod{L}}}\mod{P}) * y)\mod{P} = v_{next})
\end{align}
$$