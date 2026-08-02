'use strict';

// Benchmark harness for decimal.js (v10.6.0) mirroring bench_go_test.go.
// Op keys use the form Name,precision,ns/op so they line up with the Go
// bench. Operands are pre-created once (setup) then the op runs in a tight
// loop (op(args)), matching the Go benchmarks exactly, so both sides measure
// identical work. Run: node bench.js

const { performance } = require('perf_hooks');
const Decimal = require(process.env.DECIMAL_JS_PATH || '/root/hackathon/decimal.js/decimal.js');

function setPrecision(p) {
  Decimal.config({ precision: p, rounding: 4, toExpNeg: -7, toExpPos: 21, minE: -9e15, maxE: 9e15, crypto: false });
}

const R = (s) => new Decimal(s);

function bench(name, precision, setup, op, warmup = 2000, N = 200000) {
  setPrecision(precision || 20);
  let call = setup;
  if (op) {
    const args = setup();
    call = () => op(...args);
  }
  for (let i = 0; i < warmup; i++) call(); // JIT warmup
  const t0 = performance.now();
  for (let i = 0; i < N; i++) call();
  const dt = performance.now() - t0;
  const nsOp = (dt * 1e6) / N;
  console.log(`${name},${precision || 20},${nsOp.toFixed(1)},0`);
}

bench('New', 20, () => R('3.1415926535897932384626433832795028841971'));

bench('Add', 20, () => [R('1.23456789012345'), R('9.8765432109876')], (x, y) => x.plus(y));
bench('Sub', 20, () => [R('1.23456789012345'), R('9.8765432109876')], (x, y) => x.minus(y));
bench('Mul', 20, () => [R('1.23456789012345'), R('9.8765432109876')], (x, y) => x.times(y));

bench('Div', 20, () => [R('1.23456789012345'), R('9.8765432109876')], (x, y) => x.div(y));
bench('Div', 100, () => [R('1.23456789012345'), R('9.8765432109876')], (x, y) => x.div(y));
bench('Div', 1000, () => [R('1.23456789012345'), R('9.8765432109876')], (x, y) => x.div(y));

bench('DivToInt', 20, () => [R('123456789.5'), R('0.25')], (x, y) => x.divToInt(y));
bench('Mod', 20, () => [R('123456789.123'), R('0.001')], (x, y) => x.mod(y));

bench('PowInt', 20, () => [R('9')], (x) => x.pow(75), 2000, 20000);
bench('PowFraction', 20, () => [R('2')], (x) => x.pow(0.5), 2000, 3000);

bench('Sqrt', 20, () => [R('2')], (x) => x.sqrt(), 2000, 30000);
bench('Cbrt', 20, () => [R('2')], (x) => x.cbrt(), 2000, 30000);

bench('Exp', 20, () => [R('1')], (x) => x.exp(), 2000, 20000);
bench('Ln', 20, () => [R('2')], (x) => x.ln(), 2000, 20000);
bench('LogBase10', 20, () => [R('12345.6789')], (x) => x.log(10), 2000, 20000);

bench('Sin', 20, () => [R('0.5')], (x) => x.sin(), 2000, 10000);
bench('Cos', 20, () => [R('1.2')], (x) => x.cos(), 2000, 10000);
bench('Tan', 20, () => [R('0.7')], (x) => x.tan(), 2000, 10000);
bench('Atan', 20, () => [R('1')], (x) => x.atan(), 2000, 10000);

bench('ToFixed', 20, () => [R('12345.67890123456789')], (x) => x.toFixed(8));
bench('ToExponential', 20, () => [R('0.0000123456789')], (x) => x.toExponential(6));
bench('ToPrecision', 20, () => [R('123456789.123456789')], (x) => x.toPrecision(12));
bench('Cmp', 20, () => [R('1.234567890123456789'), R('1.234567890123456788')], (x, y) => x.cmp(y));
bench('RoundTrip', 20, () => [R('12345678.901234567890123456789')], (x) => { x.plus('0.987654321098765432109'); return x.valueOf(); });