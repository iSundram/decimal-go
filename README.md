<div align="center">
<img src="assets/logo.svg" alt="decimal-go logo">

<img src="assets/tagline.svg" alt="An arbitrary-precision decimal arithmetic library for Go.">

<p>
  <img src="https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/License-MIT-green" alt="License">
  <img src="https://img.shields.io/badge/decimal.js-v10.6.0-informational" alt="decimal.js">
  <img src="https://img.shields.io/badge/Dependencies-Zero-lightgrey" alt="Dependencies">
</p>

<p>
  <img src="https://img.shields.io/badge/Test%20modules-61%20of%2061%20ported-brightgreen" alt="Tests">
  <img src="https://img.shields.io/badge/Race%20detector-Clean-success" alt="Race detector">
  <img src="https://img.shields.io/badge/Cross--validated-1518%20cases-blue" alt="Cross-validated">
  <img src="https://img.shields.io/badge/Perf-Arithmetic%20faster%20(parsing%20slower)-blueviolet" alt="Performance">
</p>

<p>
  A Go implementation inspired by the behavior and API of <code>decimal.js</code>.
</p>

</div>

---

## <img src="assets/header_about.svg" alt="About" align="absmiddle">

<img src="assets/logo2x.svg" height="24" align="absmiddle" alt="decimal-go"> is a Go implementation inspired by the behavior and API of `decimal.js`, created during the Port Mortem 2026 hackathon.

## <img src="assets/header_goals.svg" alt="Goals" align="absmiddle">

- Preserve the behavior of the original library as closely as possible.
- Provide an idiomatic Go API.
- Maintain high test compatibility.
- Produce clean, well-documented Go code.

## <img src="assets/header_install.svg" alt="Installation" align="absmiddle">

`decimal-go` has **zero external dependencies** and requires Go 1.25+.

```sh
go get github.com/iSundram/decimal-go
```

Then import it in your code:

```go
import "github.com/iSundram/decimal-go"
```

Values are arbitrary-precision and immutable; every operation returns a new
value. The package-level `New` uses the default constructor (precision 20,
half-up rounding).

## <img src="assets/header_quickstart.svg" alt="Quick Start" align="absmiddle">

```go
package main

import (
	"fmt"

	"github.com/iSundram/decimal-go"
)

func main() {
	// Money math without the float rounding surprises of float64.
	price := decimal.New("19.99")
	qty := decimal.New("3")
	total := price.Times(qty)
	fmt.Println(total) // 59.97

	// Exact value parsing from strings, integers, floats or *big.Int.
	fmt.Println(decimal.New("0.1").Plus(decimal.New("0.2"))) // 0.3

	// Custom precision via a cloned constructor.
	c := decimal.Default.Clone(&decimal.Config{Precision: decimal.I64(50)})
	pi := c.New("3.14159265358979323846264338327950288419716939937510")
	fmt.Println(pi.ToFixed(30))
}
```

Output:

```text
59.97
0.3
3.141592653589793238462643383280
```

See the [documentation site](https://iSundram.github.io/decimal-go/), the
[example tests](example_test.go), the [operations matrix](MATRIX.md),
the [decimal.js migration guide](docs/migrating-from-decimal-js.md), the
[design decisions](docs/DECISIONS.md) and the
[changelog](CHANGELOG.md) for more.

## <img src="assets/header_status.svg" alt="Status" align="absmiddle">

> **Test suite fully ported.**

The complete `decimal.js` test suite (61 modules) is ported to Go: every
module is covered by white-box tests in package `decimal`, and
`go test ./...` is green.

- Run the suite: `go test -count=1 ./...`
- Run with the race detector: `go test -race -count=1 .`
- Porting rules live in [`PORTING_TESTS.md`](PORTING_TESTS.md).

## Validation added in this fork

Beyond the ported suites, exact parity with the live `decimal.js` is proven by
additional committed tests and tooling:

- **Cross-validation** ([`xvalidate/`](xvalidate/)): a shared corpus of
  values is run through both decimal-go (`go run`) and the real
  decimal.js (`node`); the two outputs are byte-for-byte identical
  (`bash xvalidate/compare.sh`, 1518/1518 lines).
- **API parity** ([`parity.go`](parity.go)): the 42 decimal.js long method
  names and `toJSON`/`MarshalJSON` exist as aliases and are `equal`-checked
  against the primary methods — mirroring the real methods in `decimal.js`.
- **Property-based tests** ([`property_test.go`](property_test.go)):
  round-trip, commutativity,
  add/sub & mul/div inverses, sqrt/cbrt/exp-log inverses, comparisons,
  exact sqrt, modulo range.
- **Stress + race tests** ([`stress_test.go`](stress_test.go)): precision
  400–2048, integer boundaries, NaN/Infinity, and 64 concurrent
  clone-constructors under `-race`. A genuine data race inherited from
  decimal.js's module globals (`external`/`inexact`/`quadrant`) was found and
  fixed by moving those flags onto the `Constructor` (matching decimal.js's
  "one clone per context" model).
- **Regression tests** ([`regression_test.go`](regression_test.go)): the three
  bugs fixed while porting (pow10 overflow, Pow negative-base index,
  `1^±Inf`) are locked in.
- **Input matrix** ([`input_test.go`](input_test.go)): every decimal.js value
  type (number, string, `big.Int`, `Decimal`, `-0`, NaN, ±Infinity) accepted.
- **Operations matrix** ([`MATRIX.md`](MATRIX.md)) documents every op, its
  rounding semantics and edge cases.
- **Benchmarks** ([`bench/README.md`](bench/README.md) + [`results.txt`](bench/results.txt)):
  honest Go-vs-JS comparison (18 of 20 directly comparable op pairs are
  faster in Go, up to ~3× on multiplication and division); `New(string)`
  parsing is ~1.45× slower and parse+op+format round trips ~1.7× slower,
  because real-world benchmarks that re-parse inputs pay Go's allocation
  overhead while JS string parsing is heavily JIT-optimised.

## <img src="assets/header_contributing.svg" alt="Contributing" align="absmiddle">

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to get started.

## <img src="assets/header_acknowledgements.svg" alt="Acknowledgements" align="absmiddle">

<table>
  <tr>
    <td align="center">
      <a href="https://github.com/MikeMcl">
        <img src="https://avatars.githubusercontent.com/u/157787?v=4" width="60px;" alt="MikeMcl" style="border-radius: 50%;"/>
        <br />
        <sub><b>Michael Mclaughlin</b></sub>
      </a>
    </td>
    <td>
      Huge thanks to Michael Mclaughlin for the original <code>decimal.js</code> library, which heavily inspired the API, design, and behavior of this project.
    </td>
  </tr>
</table>

## <img src="assets/header_contributors.svg" alt="Contributors" align="absmiddle">

Thanks goes to these wonderful people who have contributed to this project. This section automatically updates as new developers join in!

<a href="https://github.com/iSundram/decimal-go/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=iSundram/decimal-go" alt="Contributors" />
</a>

## <img src="assets/header_license.svg" alt="License" align="absmiddle">

This project is released under the MIT License.

---

<div align="center">

Made with Go.

</div>
