# outyet
A web server that announces whether or not a particular Go version has been tagged.



# Testing
There are two problems here if you look carefully at `main_test.go`

1. Race Condition Problem. A and B are trying to access the same resource, one
of the access is WRITE.

2. Rendezvous Problem. The need to ensure the order of executions. A cannot
happen before B.


And if you observe carefully enough, these are actually concurrency problem, and
by nature, concurrency problems are hard to solve, they are `non-deterministic`.
The key to solve concurrency problems is to making them `deterministic` again.

It turns out that, concurrency primitive of Go: `channel` is pretty good at solving
these type of problems, techniques & clever use of `channel` can therefore
enforce them to be `deterministic` again.

TO be deterministic: (making signals)
- signaling: able to ensure the order of execution by signals between more than 1 parties.
- guarding: eventually it is done with signaling also.
