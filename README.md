# fsm
# Build on private Drone CI     [![Build Status][dci-img]][dci]
# [![Build Status][tci-img]][tci] [![Coverage Status][cov-img]][cov]

Reconstruct https://github.com/looplab/fsm

# Why do this reconstruction?
(1) Up to 44.9% improvement of reaction speed from event to callback function.

```shell
    go test -bench=. -run=none -benchtime=1s -benchmem -cpuprofile cpu.prof -memprofile mem.prof
    2019-03-21T13:22:11.392+0800	DEBUG	System init finished.
    goos: linux
    goarch: amd64
    pkg: github.com/falconray0704/fsm/fsmBenchmark
    Benchmark_falconFSM-8    	 1000000	      1361 ns/op	     387 B/op	      12 allocs/op
    Benchmark_loopLabFSM-8   	  500000	      2473 ns/op	     531 B/op	      12 allocs/op
```

# What optimizations in it?
(1) Using int as key in callback and transition indexing.

[dci-img]: https://drone.doryhub.com/api/badges/falconray0704/fsm/status.svg
[dci]: https://drone.doryhub.com/api/badges/falconray0704/fsm

[tci-img]: https://travis-ci.org/falconray0704/fsm.svg?branch=master
[tci]: https://travis-ci.org/falconray0704/fsm

[cov-img]: https://codecov.io/gh/falconray0704/fsm/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/falconray0704/fsm

