<div align="center">
<img src="assets/logo.svg" alt="decimal-go logo">

<img src="assets/tagline.svg" alt="An arbitrary-precision decimal arithmetic library for Go.">

<p>
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/License-MIT-green" alt="License">
  <img src="https://img.shields.io/badge/Status-Work%20In%20Progress-orange" alt="Status">
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

## <img src="assets/header_status.svg" alt="Status" align="absmiddle">

> **Test suite fully ported.**

The complete `decimal.js` test suite is ported to Go: all 62 test modules are
covered by white-box tests in package `decimal`, and `go test ./...` is green.

- Run the suite: `go test -count=1 ./...`
- Porting rules live in [`PORTING_TESTS.md`](PORTING_TESTS.md).

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
