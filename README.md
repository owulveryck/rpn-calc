# RPN Calculator

An HP 48GX-inspired Reverse Polish Notation calculator running entirely in the browser via WebAssembly. Built with Go (no external dependencies) and vanilla JavaScript.

**[Try it live](https://owulveryck.github.io/rpn-calc/)**

## Features

- **RPN input** — operands first, then operators; no parentheses needed
- **Computation graph** — expressions are stored as a symbolic DAG, enabling exact results (e.g. `sin(π) = 0`, not `1.2e-16`) and arbitrary-precision arithmetic (256-bit floats via `math/big`)
- **Multiple bases** — DEC, HEX (`#FF`), OCT (`77oo`), BIN (`1010bb`)
- **Angle modes** — DEG, RAD, GRAD with conversion embedded in the expression graph
- **Undo** — up to 100 steps
- **Offline-first PWA** — service worker caches everything; works without a network after first load
- **Keyboard driven** — digits, operators, `s`wap, `d`up, `u`ndo, `Ctrl+Z`, arrows

## Architecture

```
┌─────────────┐       ┌──────────────────────┐       ┌───────────────┐
│  index.html │       │   Go engine (WASM)   │       │   app.js      │
│  style.css  │◄─────►│                      │◄─────►│  (vanilla JS) │
│             │       │  Stack + Expr graph   │       │  UI, keyboard │
└─────────────┘       │  BigMath (256-bit)   │       │  localStorage │
                      │  Serialize/Restore   │       └───────────────┘
                      └──────────────────────┘
```

### Engine (`engine/`)

The core is a stack machine where each value is either a concrete number (`RealValue`) or a node in a computation graph (`ExprValue`).

**Computation graph** — Every operation builds an expression tree instead of eagerly computing a float. Each node caches its result lazily (`sync.Once`). This gives two advantages:

1. **Exact symbolic results** — trigonometric functions at rational multiples of π return exact values (sin(π/6) = 0.5, not 0.4999…). The evaluator pattern-matches special angles before falling back to numeric computation.
2. **Arbitrary precision** — `bigmath.go` implements `ln`, `exp`, `pow` via Taylor/arctanh series at 256-bit precision, with cached high-precision constants for π, e, ln(2), ln(10).

**Base modes** — numbers can be entered and displayed in decimal, hex, octal, or binary. Non-decimal bases truncate to int64 for display.

**State** — the calculator state (stack, modes, history) serializes to JSON for persistence. Undo snapshots are kept in memory only (not persisted across reloads).

### WASM bridge (`cmd/wasm/`)

Go compiles to `rpn.wasm` via `GOOS=js GOARCH=wasm`. The bridge exposes functions to JavaScript:

| Function | Purpose |
|---|---|
| `rpnEnter(str)` | Parse and push a number |
| `rpnExecute(op)` | Run an operation (+, SIN, SWAP, …) |
| `rpnGetState()` | Return current state as JSON |
| `rpnUndo()` | Undo last operation |
| `rpnLast()` | Re-push last operands |
| `rpnSetAngleMode(m)` | Switch DEG/RAD/GRAD |
| `rpnSetBaseMode(m)` | Switch DEC/HEX/OCT/BIN |
| `rpnSerialize()` | Export state for localStorage |
| `rpnRestore(json)` | Restore saved state |

All functions return a JSON state object: `{stack, depth, error, angleMode, baseMode, history}`.

### Web app (`web/`)

Vanilla JS, no build step, no framework. The UI mirrors an HP 48GX layout with function, stack, and digit rows. State is auto-saved to `localStorage` (debounced 500ms). A service worker (`sw.js`) caches all assets for offline use.

## Operations

| Category | Operations |
|---|---|
| Arithmetic | `+` `-` `*` `/` |
| Power/Root | `Y^X` `SQRT` `SQ` |
| Trig | `SIN` `COS` `TAN` `ASIN` `ACOS` `ATAN` |
| Log/Exp | `LOG` `LN` `EXP` `10^X` |
| Unary | `NEG` `INV` `ABS` `FACT` |
| Constants | `PI` `E` |
| Stack | `SWAP` `DUP` `DROP` `ROT` `OVER` `DEPTH` `CLEAR` |
| Other | `UNDO` `LAST` `MIN` `MAX` |

## Build

Requires Go 1.25+.

```sh
make build   # compile rpn.wasm + copy wasm_exec.js
make test    # run tests with coverage
make run     # build + serve on http://localhost:8080
```

## Deployment

Pushes to `main` trigger a GitHub Actions workflow that runs tests, stamps the commit hash into the UI, builds the WASM binary, and deploys to GitHub Pages.

## License

MIT
