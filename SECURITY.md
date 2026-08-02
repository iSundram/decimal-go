# Security Policy

## Reporting a vulnerability

Please do **not** open a public issue for security problems. Report them
privately by emailing the maintainer or by opening a GitHub security advisory
at https://github.com/iSundram/decimal-go/security/advisories/new .

Include:

- a description of the issue,
- a minimal reproduction case (input string + operation),
- the affected version(s).

## What we care about

- A `[DecimalError]` panic is the documented contract for invalid input; it is
  not a security issue.
- A genuine bug is: a crash, hang, or wrong/unexpected result on **valid**
  input, memory unsafety, or unbounded resource growth triggered by attacker-
  controlled input.

## Supported versions

The latest tagged release is supported. There are no release branches yet.