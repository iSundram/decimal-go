//go:build js && wasm

package main

import (
	"syscall/js"

	decimal "github.com/iSundram/decimal-go"
)

func main() {
	js.Global().Set("decimalGo", map[string]interface{}{
		"new":           js.FuncOf(jsNew),
		"plus":          js.FuncOf(jsPlus),
		"minus":         js.FuncOf(jsMinus),
		"times":         js.FuncOf(jsTimes),
		"div":           js.FuncOf(jsDiv),
		"mod":           js.FuncOf(jsMod),
		"pow":           js.FuncOf(jsPow),
		"sqrt":          js.FuncOf(jsSqrt),
		"cbrt":          js.FuncOf(jsCbrt),
		"exp":           js.FuncOf(jsExp),
		"ln":            js.FuncOf(jsLn),
		"log10":         js.FuncOf(jsLog10),
		"sin":           js.FuncOf(jsSin),
		"cos":           js.FuncOf(jsCos),
		"tan":           js.FuncOf(jsTan),
		"asin":          js.FuncOf(jsAsin),
		"acos":          js.FuncOf(jsAcos),
		"atan":          js.FuncOf(jsAtan),
		"sinh":          js.FuncOf(jsSinh),
		"cosh":          js.FuncOf(jsCosh),
		"tanh":          js.FuncOf(jsTanh),
		"abs":           js.FuncOf(jsAbs),
		"neg":           js.FuncOf(jsNeg),
		"trunc":         js.FuncOf(jsTrunc),
		"ceil":          js.FuncOf(jsCeil),
		"floor":         js.FuncOf(jsFloor),
		"round":         js.FuncOf(jsRound),
		"toFixed":       js.FuncOf(jsToFixed),
		"toExponential": js.FuncOf(jsToExponential),
		"toPrecision":   js.FuncOf(jsToPrecision),
		"toHex":         js.FuncOf(jsToHex),
		"toBinary":      js.FuncOf(jsToBinary),
		"toOctal":       js.FuncOf(jsToOctal),
		"cmp":           js.FuncOf(jsCmp),
		"eq":            js.FuncOf(jsEq),
		"gt":            js.FuncOf(jsGt),
		"gte":           js.FuncOf(jsGte),
		"lt":            js.FuncOf(jsLt),
		"lte":           js.FuncOf(jsLte),
		"isNaN":         js.FuncOf(jsIsNaN),
		"isInf":         js.FuncOf(jsIsInf),
		"isFinite":      js.FuncOf(jsIsFinite),
		"isInt":         js.FuncOf(jsIsInt),
		"valueOf":       js.FuncOf(jsValueOf),
		"toString":      js.FuncOf(jsToString),
	})
	select {}
}

func jsNew(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: new(value)"
	}
	return decimal.New(args[0].String()).String()
}

func jsPlus(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: plus(a, b)"
	}
	return decimal.New(args[0].String()).Plus(decimal.New(args[1].String())).String()
}
func jsMinus(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: minus(a, b)"
	}
	return decimal.New(args[0].String()).Minus(decimal.New(args[1].String())).String()
}
func jsTimes(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: times(a, b)"
	}
	return decimal.New(args[0].String()).Times(decimal.New(args[1].String())).String()
}
func jsDiv(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: div(a, b)"
	}
	return decimal.New(args[0].String()).Div(decimal.New(args[1].String())).String()
}
func jsMod(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: mod(a, b)"
	}
	return decimal.New(args[0].String()).Mod(decimal.New(args[1].String())).String()
}
func jsPow(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: pow(a, b)"
	}
	return decimal.New(args[0].String()).Pow(decimal.New(args[1].String())).String()
}
func jsSqrt(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: sqrt(a)"
	}
	return decimal.New(args[0].String()).Sqrt().String()
}
func jsCbrt(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: cbrt(a)"
	}
	return decimal.New(args[0].String()).Cbrt().String()
}
func jsExp(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: exp(a)"
	}
	return decimal.New(args[0].String()).Exp().String()
}
func jsLn(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: ln(a)"
	}
	return decimal.New(args[0].String()).Ln().String()
}
func jsLog10(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: log10(a)"
	}
	return decimal.Default.Log10(args[0].String()).String()
}
func jsSin(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: sin(a)"
	}
	return decimal.New(args[0].String()).Sin().String()
}
func jsCos(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: cos(a)"
	}
	return decimal.New(args[0].String()).Cos().String()
}
func jsTan(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: tan(a)"
	}
	return decimal.New(args[0].String()).Tan().String()
}
func jsAsin(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: asin(a)"
	}
	return decimal.New(args[0].String()).Asin().String()
}
func jsAcos(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: acos(a)"
	}
	return decimal.New(args[0].String()).Acos().String()
}
func jsAtan(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: atan(a)"
	}
	return decimal.New(args[0].String()).Atan().String()
}
func jsSinh(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: sinh(a)"
	}
	return decimal.New(args[0].String()).Sinh().String()
}
func jsCosh(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: cosh(a)"
	}
	return decimal.New(args[0].String()).Cosh().String()
}
func jsTanh(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: tanh(a)"
	}
	return decimal.New(args[0].String()).Tanh().String()
}
func jsAbs(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: abs(a)"
	}
	return decimal.New(args[0].String()).Abs().String()
}
func jsNeg(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: neg(a)"
	}
	return decimal.New(args[0].String()).Neg().String()
}
func jsTrunc(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: trunc(a)"
	}
	return decimal.New(args[0].String()).Trunc().String()
}
func jsCeil(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: ceil(a)"
	}
	return decimal.New(args[0].String()).Ceil().String()
}
func jsFloor(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: floor(a)"
	}
	return decimal.New(args[0].String()).Floor().String()
}
func jsRound(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: round(a)"
	}
	return decimal.New(args[0].String()).Round().String()
}
func jsToFixed(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 || len(args) > 2 {
		return "usage: toFixed(a, dp?)"
	}
	dp := int64(0)
	if len(args) == 2 {
		dp = int64(args[1].Int())
	}
	return decimal.New(args[0].String()).ToFixed(dp)
}
func jsToExponential(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 || len(args) > 2 {
		return "usage: toExponential(a, dp?)"
	}
	dp := int64(6)
	if len(args) == 2 {
		dp = int64(args[1].Int())
	}
	return decimal.New(args[0].String()).ToExponential(dp)
}
func jsToPrecision(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 || len(args) > 2 {
		return "usage: toPrecision(a, sd?)"
	}
	sd := int64(0)
	if len(args) == 2 {
		sd = int64(args[1].Int())
	}
	return decimal.New(args[0].String()).ToPrecision(sd)
}
func jsToHex(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 || len(args) > 2 {
		return "usage: toHex(a, sd?)"
	}
	sd := int64(0)
	if len(args) == 2 {
		sd = int64(args[1].Int())
	}
	return decimal.New(args[0].String()).ToHex(sd)
}
func jsToBinary(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 || len(args) > 2 {
		return "usage: toBinary(a, sd?)"
	}
	sd := int64(0)
	if len(args) == 2 {
		sd = int64(args[1].Int())
	}
	return decimal.New(args[0].String()).ToBinary(sd)
}
func jsToOctal(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 || len(args) > 2 {
		return "usage: toOctal(a, sd?)"
	}
	sd := int64(0)
	if len(args) == 2 {
		sd = int64(args[1].Int())
	}
	return decimal.New(args[0].String()).ToOctal(sd)
}
func jsCmp(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: cmp(a, b)"
	}
	return decimal.New(args[0].String()).Cmp(decimal.New(args[1].String()))
}
func jsEq(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: eq(a, b)"
	}
	return decimal.New(args[0].String()).Eq(decimal.New(args[1].String()))
}
func jsGt(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: gt(a, b)"
	}
	return decimal.New(args[0].String()).Gt(decimal.New(args[1].String()))
}
func jsGte(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: gte(a, b)"
	}
	return decimal.New(args[0].String()).Gte(decimal.New(args[1].String()))
}
func jsLt(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: lt(a, b)"
	}
	return decimal.New(args[0].String()).Lt(decimal.New(args[1].String()))
}
func jsLte(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "usage: lte(a, b)"
	}
	return decimal.New(args[0].String()).Lte(decimal.New(args[1].String()))
}
func jsIsNaN(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: isNaN(a)"
	}
	return decimal.New(args[0].String()).IsNaN()
}
func jsIsInf(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: isInf(a)"
	}
	return !decimal.New(args[0].String()).IsFinite()
}
func jsIsFinite(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: isFinite(a)"
	}
	return decimal.New(args[0].String()).IsFinite()
}
func jsIsInt(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: isInt(a)"
	}
	return decimal.New(args[0].String()).IsInt()
}
func jsValueOf(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: valueOf(a)"
	}
	return decimal.New(args[0].String()).ValueOf()
}
func jsToString(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "usage: toString(a)"
	}
	return decimal.New(args[0].String()).String()
}
