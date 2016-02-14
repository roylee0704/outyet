# outyet

A web server that announces whether or not a particular Go version has been tagged.


## Testing

There are two problems here if you look carefully at `main_test.go`

1. Race Condition Problem. A and B are trying to access the same resource, one
of the access is WRITE.

2. Rendezvous Problem. The need to ensure the order of executions. A cannot
happen before B.

---
## Analysis

If you study carefully enough, these are actually concurrency problem, how do I
know? Easy, there is `go` routine in it.

Concurrency problems are hard to solve by nature, and it always boils down to
one evil root-cause: they are `non-deterministic`.

Now you know the root cause, to solve any kind of concurrency problems, it easy,
you just need to make them `deterministic` again.

---
## Solution

It turns out that, concurrency primitive of Go: `channel` is pretty good at solving
these type of problems. Just like `sleep()`, **it blocks**.

General techniques in solving concurrency problem:

- Signal. *guarantees order of execution by communications*
- Guard. *guarantees number of simultaneous access & order, too*

The two Problems described above is solved using technique#1, signaling.
