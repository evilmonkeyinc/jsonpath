# Benchmarks

Benchmarks of this library compared to other golang implementations.

The example selectors from the original specification along with the sample data have been used to create this data.

If any library is not mentioned for a selector, that means that implementation returned an error of some kind

Test of accuracy are based off of the expected response based on the original specification and the consensus from the [json-path-comparison](https://cburgmer.github.io/json-path-comparison/)

## Command

```bash
go test -bench=. -cpu=1 -benchmem -count=1 -benchtime=1000x
```

## Libraries

- `github.com/evilmonkeyinc/jsonpath v0.7.0`
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
Benchmark_Comparison/$.store.book[*].author/evilmonkeyinc         	   10000	     21165 ns/op	    6496 B/op	     188 allocs/op
Benchmark_Comparison/$.store.book[*].author/paesslerAG            	   10000	     16595 ns/op	    6944 B/op	     153 allocs/op
Benchmark_Comparison/$.store.book[*].author/bhmj                  	   10000	      4880 ns/op	    1185 B/op	      14 allocs/op
Benchmark_Comparison/$.store.book[*].author/oliveagle             	   10000	     17070 ns/op	    4784 B/op	     147 allocs/op
Benchmark_Comparison/$.store.book[*].author/spyzhov               	   10000	     16361 ns/op	    7032 B/op	     136 allocs/op
Benchmark_Comparison/$..author/evilmonkeyinc                      	   10000	     56820 ns/op	   16688 B/op	     458 allocs/op
Benchmark_Comparison/$..author/paesslerAG                         	   10000	    113776 ns/op	   20624 B/op	     630 allocs/op
Benchmark_Comparison/$..author/bhmj                               	   10000	     33157 ns/op	    1553 B/op	      27 allocs/op
Benchmark_Comparison/$..author/oliveagle                          	   10000	     37042 ns/op	    4464 B/op	     118 allocs/op
--- BENCH: Benchmark_Comparison/$..author/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..author/spyzhov                            	   10000	     51169 ns/op	    8336 B/op	     168 allocs/op
Benchmark_Comparison/$.store.*/evilmonkeyinc                      	   10000	     43641 ns/op	    4929 B/op	     130 allocs/op
Benchmark_Comparison/$.store.*/paesslerAG                         	   10000	     31748 ns/op	    6280 B/op	     126 allocs/op
Benchmark_Comparison/$.store.*/bhmj                               	   10000	      5370 ns/op	    3705 B/op	       9 allocs/op
Benchmark_Comparison/$.store.*/oliveagle                          	   10000	     20300 ns/op	    4480 B/op	     118 allocs/op
--- BENCH: Benchmark_Comparison/$.store.*/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$.store.*/spyzhov                            	   10000	     12892 ns/op	    6785 B/op	     124 allocs/op
Benchmark_Comparison/$.store..price/evilmonkeyinc                 	   10000	     50521 ns/op	   15672 B/op	     443 allocs/op
Benchmark_Comparison/$.store..price/paesslerAG                    	   10000	     60435 ns/op	   17400 B/op	     515 allocs/op
Benchmark_Comparison/$.store..price/bhmj                          	   10000	     12666 ns/op	    1192 B/op	      28 allocs/op
Benchmark_Comparison/$.store..price/oliveagle                     	   10000	     13522 ns/op	    4576 B/op	     128 allocs/op
--- BENCH: Benchmark_Comparison/$.store..price/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$.store..price/spyzhov                       	   10000	     18554 ns/op	    8256 B/op	     168 allocs/op
Benchmark_Comparison/$..book[2]/evilmonkeyinc                     	   10000	     54938 ns/op	   16960 B/op	     471 allocs/op
Benchmark_Comparison/$..book[2]/paesslerAG                        	   10000	     50810 ns/op	   20816 B/op	     643 allocs/op
Benchmark_Comparison/$..book[2]/bhmj                              	   10000	     10774 ns/op	    1257 B/op	      16 allocs/op
Benchmark_Comparison/$..book[2]/oliveagle                         	   10000	     12317 ns/op	    4560 B/op	     124 allocs/op
--- BENCH: Benchmark_Comparison/$..book[2]/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[2]/spyzhov                           	   10000	     18695 ns/op	    8296 B/op	     166 allocs/op
Benchmark_Comparison/$..book[(@.length-1)]/evilmonkeyinc          	   10000	    129960 ns/op	   18000 B/op	     542 allocs/op
Benchmark_Comparison/$..book[(@.length-1)]/paesslerAG             	   10000	     18974 ns/op	    6576 B/op	     140 allocs/op
--- BENCH: Benchmark_Comparison/$..book[(@.length-1)]/paesslerAG
    benchmark_test.go:84: unsupported
    benchmark_test.go:84: unsupported
Benchmark_Comparison/$..book[(@.length-1)]/bhmj                   	   10000	       622.8 ns/op	     648 B/op	       4 allocs/op
--- BENCH: Benchmark_Comparison/$..book[(@.length-1)]/bhmj
    benchmark_test.go:102: unsupported
    benchmark_test.go:102: unsupported
Benchmark_Comparison/$..book[(@.length-1)]/oliveagle              	   10000	     14319 ns/op	    4736 B/op	     143 allocs/op
--- BENCH: Benchmark_Comparison/$..book[(@.length-1)]/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[(@.length-1)]/spyzhov                	   10000	     24926 ns/op	    9248 B/op	     203 allocs/op
Benchmark_Comparison/$..book[-1:]/evilmonkeyinc                   	   10000	    127852 ns/op	   17200 B/op	     486 allocs/op
Benchmark_Comparison/$..book[-1:]/paesslerAG                      	   10000	    134288 ns/op	   21184 B/op	     654 allocs/op
Benchmark_Comparison/$..book[-1:]/bhmj                            	   10000	     29420 ns/op	    1706 B/op	      22 allocs/op
--- BENCH: Benchmark_Comparison/$..book[-1:]/bhmj
    benchmark_test.go:109: found single nested array
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[-1:]/oliveagle                       	   10000	     20008 ns/op	    4664 B/op	     129 allocs/op
--- BENCH: Benchmark_Comparison/$..book[-1:]/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[-1:]/spyzhov                         	   10000	     63175 ns/op	    8376 B/op	     170 allocs/op
Benchmark_Comparison/$..book[0,1]/evilmonkeyinc                   	   10000	    116278 ns/op	   17297 B/op	     489 allocs/op
Benchmark_Comparison/$..book[0,1]/paesslerAG                      	   10000	     80451 ns/op	   21312 B/op	     656 allocs/op
Benchmark_Comparison/$..book[0,1]/bhmj                            	   10000	     11671 ns/op	    2282 B/op	      23 allocs/op
--- BENCH: Benchmark_Comparison/$..book[0,1]/bhmj
    benchmark_test.go:109: found single nested array
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[0,1]/oliveagle                       	   10000	     15306 ns/op	    4648 B/op	     130 allocs/op
--- BENCH: Benchmark_Comparison/$..book[0,1]/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[0,1]/spyzhov                         	   10000	     22789 ns/op	    8456 B/op	     172 allocs/op
Benchmark_Comparison/$..book[:2]/evilmonkeyinc                    	   10000	     92287 ns/op	   17281 B/op	     483 allocs/op
Benchmark_Comparison/$..book[:2]/paesslerAG                       	   10000	     89657 ns/op	   21248 B/op	     656 allocs/op
Benchmark_Comparison/$..book[:2]/bhmj                             	   10000	     26255 ns/op	    2346 B/op	      24 allocs/op
--- BENCH: Benchmark_Comparison/$..book[:2]/bhmj
    benchmark_test.go:109: found single nested array
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[:2]/oliveagle                        	   10000	     14178 ns/op	    4648 B/op	     126 allocs/op
--- BENCH: Benchmark_Comparison/$..book[:2]/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[:2]/spyzhov                          	   10000	     37200 ns/op	    8392 B/op	     171 allocs/op
Benchmark_Comparison/$..book[?(@.isbn)]/evilmonkeyinc             	   10000	     69058 ns/op	   20265 B/op	     556 allocs/op
Benchmark_Comparison/$..book[?(@.isbn)]/paesslerAG                	   10000	     55980 ns/op	   21896 B/op	     679 allocs/op
Benchmark_Comparison/$..book[?(@.isbn)]/bhmj                      	   10000	     15960 ns/op	    2728 B/op	      30 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.isbn)]/bhmj
    benchmark_test.go:109: found single nested array
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[?(@.isbn)]/oliveagle                 	   10000	     12832 ns/op	    4680 B/op	     138 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.isbn)]/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[?(@.isbn)]/spyzhov                   	   10000	     22150 ns/op	    9272 B/op	     224 allocs/op
Benchmark_Comparison/$..book[?(@.price<10)]/evilmonkeyinc         	   10000	     65729 ns/op	   20153 B/op	     564 allocs/op
Benchmark_Comparison/$..book[?(@.price<10)]/paesslerAG            	   10000	     15610 ns/op	    6576 B/op	     140 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<10)]/paesslerAG
    benchmark_test.go:84: unsupported
    benchmark_test.go:84: unsupported
Benchmark_Comparison/$..book[?(@.price<10)]/bhmj                  	   10000	     17494 ns/op	    2896 B/op	      43 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<10)]/bhmj
    benchmark_test.go:109: found single nested array
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[?(@.price<10)]/oliveagle             	   10000	     13628 ns/op	    4792 B/op	     146 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<10)]/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[?(@.price<10)]/spyzhov               	   10000	     25878 ns/op	   10568 B/op	     270 allocs/op
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/evilmonkeyinc         	   10000	     73693 ns/op	   21449 B/op	     628 allocs/op
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/paesslerAG            	   10000	     15426 ns/op	    6616 B/op	     140 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<$.expensive)]/paesslerAG
    benchmark_test.go:84: unsupported
    benchmark_test.go:84: unsupported
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/bhmj                  	   10000	     20282 ns/op	    2992 B/op	      46 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<$.expensive)]/bhmj
    benchmark_test.go:109: found single nested array
    benchmark_test.go:109: found single nested array
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/oliveagle             	   10000	     14472 ns/op	    5080 B/op	     164 allocs/op
--- BENCH: Benchmark_Comparison/$..book[?(@.price<$.expensive)]/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..book[?(@.price<$.expensive)]/spyzhov               	   10000	     32897 ns/op	   10496 B/op	     292 allocs/op
Benchmark_Comparison/$..*/evilmonkeyinc                                    	   10000	     72445 ns/op	   20199 B/op	     546 allocs/op
Benchmark_Comparison/$..*/paesslerAG                                       	   10000	     46927 ns/op	   20440 B/op	     647 allocs/op
Benchmark_Comparison/$..*/bhmj                                             	   10000	     25616 ns/op	   31209 B/op	      69 allocs/op
Benchmark_Comparison/$..*/oliveagle                                        	   10000	     12465 ns/op	    4304 B/op	     107 allocs/op
--- BENCH: Benchmark_Comparison/$..*/oliveagle
    benchmark_test.go:137: unsupported
    benchmark_test.go:137: unsupported
Benchmark_Comparison/$..*/spyzhov                                          	   10000	     22968 ns/op	   11519 B/op	     220 allocs/op
PASS
ok  	github.com/evilmonkeyinc/jsonpath/benchmark	25.365s
```
