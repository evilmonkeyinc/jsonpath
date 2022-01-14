[![codecov](https://codecov.io/gh/evilmonkeyinc/jsonpath/branch/main/graph/badge.svg?token=4PU85I7J2R)](https://codecov.io/gh/evilmonkeyinc/jsonpath)
[![main](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml)
[![develop](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml/badge.svg?branch=develop)](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/evilmonkeyinc/jsonpath.svg)](https://pkg.go.dev/github.com/evilmonkeyinc/jsonpath)

> This library is on the unstable version v0.X.X, which means there is a chance that any minor update may introduce a breaking change. Where I will endeavor to avoid this, care should be taken updating your dependency on this library until the first stable release v1.0.0 at which point any future breaking changes will result in a new major release.


## Supported standard evaluation operations

| symbol | name | supported types | example | notes |
| --- | --- | --- | --- | --- |
| == | equals | any | 1 == 1 returns true | |
| != | not equals | any | 1 != 2 returns true | |
| * | multiplication | int\|float | 2*2 returns 4 | |
| / | division | int\|float | 10/5 returns 2 | if you supply two whole numbers you will only get a whole number response, even if there is a remainder i.e. 10/4 would return 2, not 2.5. to include remainders you would need to have the numerator as a float i.e. 10.0/4 would return 2.5 |
| + | addition | int\|float | 2+2 returns 4 | |
| - | subtraction | int\|float | 2-2 returns 0 | |
| % | remainder | int\|float | 5 % 2 returns 1 | this operator will divide the numerator by the denominator and then return the remainder |
| > | greater than | | | |
| >= | greater than or equal to | | | |
| < | less than | | | |
| <= | less than or equal to  | | | |
| && | combine and | expression\|bool | x | y |
| \|\| | combine or | expression\|bool | x | y |
| () | sub-expression | expression | x | y |