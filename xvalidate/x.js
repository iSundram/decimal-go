'use strict';

// x.js mirrors xvalid.go: reads xvalidate/cases.txt, applies the same unary
// ops, conversions and binary ops with the same default config, and prints
// identical "op\tinput\tresult" lines. Diff against the Go output.
//
// Usage:
//   node xvalidate/x.js | sort > js.out
const fs = require('fs');
const path = require('path');

const Decimal = require(process.env.DECIMAL_JS_PATH || '/root/hackathon/decimal.js/decimal.js');
Decimal.config({ precision: 20, rounding: 4, toExpNeg: -7, toExpPos: 21, minE: -9e15, maxE: 9e15, crypto: false });

const cases = fs.readFileSync(path.join(__dirname, 'cases.txt'), 'utf8').trim().split(/\s+/);

function fmtf(v) { return String(v); } // JS Number::toString, shortest form

function result(v) { return v.valueOf(); }

for (const v of cases) {
  const x = new Decimal(v);
  console.log(`abs\t${v}\t${result(x.abs())}`);
  console.log(`neg\t${v}\t${result(x.neg())}`);
  console.log(`trunc\t${v}\t${result(x.trunc())}`);
  console.log(`ceil\t${v}\t${result(x.ceil())}`);
  console.log(`floor\t${v}\t${result(x.floor())}`);
  console.log(`round\t${v}\t${result(x.round())}`);
  console.log(`sqrt\t${v}\t${result(x.sqrt())}`);
  console.log(`cbrt\t${v}\t${result(x.cbrt())}`);
  console.log(`exp\t${v}\t${result(x.exp())}`);
  console.log(`ln\t${v}\t${result(x.ln())}`);
}

for (const v of cases) {
  const x = new Decimal(v);
  console.log(`dp\t${v}\t${fmtf(x.dp())}`);
  console.log(`sd\t${v}\t${fmtf(x.sd())}`);
  console.log(`toFixed2\t${v}\t${x.toFixed(2)}`);
  console.log(`toExp6\t${v}\t${x.toExponential(6)}`);
  console.log(`toPrec8\t${v}\t${x.toPrecision(8)}`);
}

const ops = [
  ['add', (a, b) => a.plus(b)],
  ['sub', (a, b) => a.minus(b)],
  ['mul', (a, b) => a.times(b)],
  ['div', (a, b) => a.div(b)],
  ['divToInt', (a, b) => a.divToInt(b)],
  ['mod', (a, b) => a.mod(b)],
  ['pow', (a, b) => a.pow(b)],
];
for (let i = 0; i < cases.length; i++) {
  const a = cases[i];
  const b = cases[(i + 1) % cases.length];
  const xa = new Decimal(a);
  const xb = new Decimal(b);
  for (const [name, fn] of ops) {
    console.log(`${name}\t${a}|${b}\t${result(fn(xa, xb))}`);
  }
  console.log(`cmp\t${a}|${b}\t${String(xa.cmp(xb))}`);
}