[![Build Status](https://travis-ci.org/kibumh/pcontainer.png)](https://travis-ci.org/kibumh/pcontainer)

# What is this?
Persistent Containers implemented in golang inspired by Clojure programming language.


# Quick Start
```bash
go get github.com/kibumh/pcontainer
```

# GoDoc
https://godoc.org/github.com/kibumh/pcontainer

# Performance
## PushBack
- trivial PushBack API is about 30 times slower than slice version.
- transent PushBack API is about 4 times slower than slice version.
```
BenchmarkPushBack-4            	 3000000	       502 ns/op
BenchmarkTransientPushBack-4   	20000000	        65.7 ns/op
BenchmarkSlicePushBack-4       	100000000	        16.8 ns/op
```
## At
- at API is about 1.5 times slower than slice version.
```
BenchmarkAt-4                  	20000000	        89.8 ns/op
BenchmarkSliceAt-4             	20000000	        65.0 ns/op
```

# TODO
- [x] persistent vector
- [ ] persistent hash map
- [ ] advanced persistent vector using rrb-tree
