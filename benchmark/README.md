# Benchmarks

Benchmarks of this library compared to other golang implementations.

The example selectors from the original specification along with the sample data have been used to create this data.

If any library is not mentioned for a selector, that means that implementation returned an error of some kind

Test of accuracy are based off of the expected response based on the original specification and the consensus from the [json-path-comparison](https://cburgmer.github.io/json-path-comparison/)

## Command

```bash
go test -bench=. -cpu=1 -benchmem -count=1 -benchtime=10000x
```

## Libraries

- `github.com/evilmonkeyinc/jsonpath v0.7.2`
- `github.com/PaesslerAG/jsonpath v0.1.1` *uses reflection
- `github.com/bhmj/jsonslice v1.1.2` *custom parser
- `github.com/oliveagle/jsonpath v0.0.0-20180606110733-2e52cf6e6852` *uses reflection
- `github.com/spyzhov/ajson v0.7.0` *custom parser

## TL;DR

This implementation is slower than the others, but is only one of two that has a non-error response to all sample selectors or return the expected response, the other being the [spyzhov/ajson](https://github.com/spyzhov/ajson) implementation which is on average at least twice as fast but relies on its own json marshaller.

Generally the accuracy of the implementations that run are the same, with a minor deviation with how array ranges are handled with one of them when it returned an array with a single item which itself was the expected response.

## Test

```bash
goos: darwin
goarch: amd64
pkg: github.com/evilmonkeyinc/jsonpath/benchmark
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
Benchmark_Comparison/$.store.book[*].author/evilmonkeyinc         	   10000	     18869 ns/op	    6384 B/op	     185 allocs/op
Benchmark_Comparison/$.store.book[*].author/paesslerAG            	   10000	     19311 ns/op	    6944 B/op	     153 allocs/op
Benchmark_Comparison/$.store.book[*].author/bhmj                  	   10000	      7086 ns/op	    1185 B/op	      14 allocs/op
Benchmark_Comparison/$.store.book[*].author/oliveagle             	   10000	     14569 ns/op	    4784 B/op	     147 allocs/op
Benchmark_Comparison/$.store.book[*].author/spyzhov               	   10000	     14292 ns/op	    7032 B/op	     136 allocs/op
Benchmark_Comparison/$..author/evilmonkeyinc                      	   10000	     52371 ns/op	   12768 B/op	     391 allocs/op
Benchmark_Comparison/$..author/paesslerAG                         	   10000	     49396 ns/op	   20624 B/op	     630 allocs/op
Benchmark_Comparison/$..author/bhmj                               	   10000	     10053 ns/op	    1553 B/op	      27 allocs/op
Benchmark_Comparison/$..author/oliveagle                          	   10000	     12377 ns/op	    4464 B/op	     118 allocs/op
--- BENCH: Benchmark_Comparison/$..author/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..author/spyzhov                            	   10000	     18027 ns/op	    8336 B/op	     168 allocs/op
Benchmark_Comparison/$.store.*/evilmonkeyinc                      	   10000	     13030 ns/op	    4929 B/op	     130 allocs/op
Benchmark_Comparison/$.store.*/paesslerAG                         	   10000	     13948 ns/op	    6280 B/op	     126 allocs/op
Benchmark_Comparison/$.store.*/bhmj                               	   10000	      5109 ns/op	    3705 B/op	       9 allocs/op
Benchmark_Comparison/$.store.*/oliveagle                          	   10000	     12099 ns/op	    4480 B/op	     118 allocs/op
--- BENCH: Benchmark_Comparison/$.store.*/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$.store.*/spyzhov                            	   10000	     12331 ns/op	    6785 B/op	     124 allocs/op
Benchmark_Comparison/$.store..price/evilmonkeyinc                 	   10000	     45637 ns/op	   12200 B/op	     382 allocs/op
Benchmark_Comparison/$.store..price/paesslerAG                    	   10000	     42854 ns/op	   17400 B/op	     515 allocs/op
Benchmark_Comparison/$.store..price/bhmj                          	   10000	      8167 ns/op	    1192 B/op	      28 allocs/op
Benchmark_Comparison/$.store..price/oliveagle                     	   10000	     12710 ns/op	    4576 B/op	     128 allocs/op
--- BENCH: Benchmark_Comparison/$.store..price/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$.store..price/spyzhov                       	   10000	     17584 ns/op	    8256 B/op	     168 allocs/op
Benchmark_Comparison/$..book[2]/evilmonkeyinc                     	   10000	     49198 ns/op	   12864 B/op	     399 allocs/op
Benchmark_Comparison/$..book[2]/paesslerAG                        	   10000	     62442 ns/op	   20816 B/op	     643 allocs/op
Benchmark_Comparison/$..book[2]/bhmj                              	   10000	     15866 ns/op	    1257 B/op	      16 allocs/op
Benchmark_Comparison/$..book[2]/oliveagle                         	   10000	     22550 ns/op	    4560 B/op	     124 allocs/op
--- BENCH: Benchmark_Comparison/$..book[2]/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[2]/spyzhov                           	   10000	     23252 ns/op	    8296 B/op	     166 allocs/op
Benchmark_Comparison/$..book[(@.length-1)]/evilmonkeyinc          	   10000	    148979 ns/op	   13904 B/op	     470 allocs/op
Benchmark_Comparison/$..book[(@.length-1)]/paesslerAG             	   10000	     37214 ns/op	    6576 B/op	     140 allocs/op
--- BENCH: Benchmark_Comparison/$..book[(@.length-1)]/paesslerAG
    benchmark_test.go:84: unsupported
Benchmark_Comparison/$..book[(@.length-1)]/bhmj                   	   10000	       872.6 ns/op	     648 B/op	       4 allocs/op
--- BENCH: Benchmark_Comparison/$..book[(@.length-1)]/bhmj
    benchmark_test.go:102: unsupported
Benchmark_Comparison/$..book[(@.length-1)]/oliveagle              	   10000	     52143 ns/op	    4736 B/op	     143 allocs/op
--- BENCH: Benchmark_Comparison/$..book[(@.length-1)]/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[(@.length-1)]/spyzhov                	   10000	     56952 ns/op	    9248 B/op	     203 allocs/op
Benchmark_Comparison/$..book[-1:]/evilmonkeyinc                   	   10000	     83926 ns/op	   13104 B/op	     414 allocs/op
Benchmark_Comparison/$..book[-1:]/paesslerAG                      	   10000	     71105 ns/op	   21184 B/op	     654 allocs/op
Benchmark_Comparison/$..book[-1:]/bhmj                            	   10000	     14475 ns/op	    1706 B/op	      22 allocs/op
--- BENCH: Benchmark_Comparison/$..book[-1:]/bhmj
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[-1:]/oliveagle                       	   10000	     19190 ns/op	    4664 B/op	     129 allocs/op
--- BENCH: Benchmark_Comparison/$..book[-1:]/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[-1:]/spyzhov                         	   10000	     24713 ns/op	    8376 B/op	     170 allocs/op
Benchmark_Comparison/$..book[0,1]/evilmonkeyinc                   	   10000	     72378 ns/op	   13217 B/op	     417 allocs/op
Benchmark_Comparison/$..book[0,1]/paesslerAG                      	   10000	     58346 ns/op	   21312 B/op	     656 allocs/op
Benchmark_Comparison/$..book[0,1]/bhmj                            	   10000	     12148 ns/op	    2282 B/op	      23 allocs/op
--- BENCH: Benchmark_Comparison/$..book[0,1]/bhmj
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[0,1]/oliveagle                       	   10000	     14399 ns/op	    4648 B/op	     130 allocs/op
--- BENCH: Benchmark_Comparison/$..book[0,1]/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[0,1]/spyzhov                         	   10000	     45598 ns/op	    8456 B/op	     172 allocs/op
Benchmark_Comparison/$..book[:2]/evilmonkeyinc                    	   10000	     97570 ns/op	   13201 B/op	     411 allocs/op
Benchmark_Comparison/$..book[:2]/paesslerAG                       	   10000	     90545 ns/op	   21248 B/op	     656 allocs/op
Benchmark_Comparison/$..book[:2]/bhmj                             	   10000	     17076 ns/op	    2346 B/op	      24 allocs/op
--- BENCH: Benchmark_Comparison/$..book[:2]/bhmj
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[:2]/oliveagle                        	   10000	     21330 ns/op	    4648 B/op	     126 allocs/op
--- BENCH: Benchmark_Comparison/$..book[:2]/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[:2]/spyzhov                          	   10000	     36285 ns/op	    8392 B/op	     171 allocs/op
Benchmark_Comparison/$..book[?(@.isbn)]/evilmonkeyinc             	   10000	    129392 ns/op	   16185 B/op	     484 allocs/op
Benchmark_Comparison/$..book[?(@.isbn)]/paesslerAG                	   10000	     79179 ns/op	   21896 B/op	     679 allocs/op
Benchmark_Comparison/$..book[?(@.isbn)]/bhmj                      	   10000	     26555 ns/op	    2728 B/op	      30 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.isbn)]/bhmj
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[?(@.isbn)]/oliveagle                 	   10000	     35661 ns/op	    4680 B/op	     138 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.isbn)]/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[?(@.isbn)]/spyzhov                   	   10000	     28012 ns/op	    9272 B/op	     224 allocs/op
Benchmark_Comparison/$..book[?(@.price<10)]/evilmonkeyinc         	   10000	     71603 ns/op	   16073 B/op	     492 allocs/op
Benchmark_Comparison/$..book[?(@.price<10)]/paesslerAG            	   10000	     37617 ns/op	    6576 B/op	     140 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<10)]/paesslerAG
    benchmark_test.go:84: unsupported
Benchmark_Comparison/$..book[?(@.price<10)]/bhmj                  	   10000	     22856 ns/op	    2896 B/op	      43 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<10)]/bhmj
    benchmark_test.go:109: found single nested array
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[?(@.price<10)]/oliveagle             	   10000	     22944 ns/op	    4792 B/op	     146 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<10)]/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[?(@.price<10)]/spyzhov               	   10000	     35844 ns/op	   10568 B/op	     270 allocs/op
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/evilmonkeyinc         	   10000	    123683 ns/op	   17369 B/op	     556 allocs/op
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/paesslerAG            	   10000	     16665 ns/op	    6616 B/op	     140 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<$.expensive)]/paesslerAG
    benchmark_test.go:84: unsupported
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/bhmj                  	   10000	     20180 ns/op	    2992 B/op	      46 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<$.expensive)]/bhmj
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/oliveagle             	   10000	     17667 ns/op	    5080 B/op	     164 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<$.expensive)]/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/spyzhov               	   10000	     29535 ns/op	   10496 B/op	     292 allocs/op
Benchmark_Comparison/$..*/evilmonkeyinc                                    	   10000	     66367 ns/op	   17863 B/op	     496 allocs/op
Benchmark_Comparison/$..*/paesslerAG                                       	   10000	     47862 ns/op	   20440 B/op	     647 allocs/op
Benchmark_Comparison/$..*/bhmj                                             	   10000	     25386 ns/op	   31209 B/op	      69 allocs/op
Benchmark_Comparison/$..*/oliveagle                                        	   10000	     12155 ns/op	    4304 B/op	     107 allocs/op
--- BENCH: Benchmark_Comparison/$..*/oliveagle
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..*/spyzhov                                          	   10000	     24140 ns/op	   11518 B/op	     220 allocs/op
PASS
ok  	github.com/evilmonkeyinc/jsonpath/benchmark	24.627s

```
