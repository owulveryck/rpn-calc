//go:build js && wasm

// Point d'entrée WebAssembly de la calculatrice RPN.
// Expose les fonctions de la calculatrice comme fonctions globales JavaScript
// (rpnEnter, rpnExecute, rpnGetState, rpnUndo, rpnLast, rpnSerialize, rpnRestore,
// rpnSetAngleMode, rpnSetBaseMode).
package main

import (
	"syscall/js"

	"github.com/owulveryck/rpn-calc/engine"
)

func main() {
	calc := engine.NewCalculator()

	js.Global().Set("rpnEnter", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return calc.GetStateJSON()
		}
		return calc.Enter(args[0].String())
	}))

	js.Global().Set("rpnExecute", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return calc.GetStateJSON()
		}
		return calc.Execute(args[0].String())
	}))

	js.Global().Set("rpnGetState", js.FuncOf(func(this js.Value, args []js.Value) any {
		return calc.GetStateJSON()
	}))

	js.Global().Set("rpnUndo", js.FuncOf(func(this js.Value, args []js.Value) any {
		return calc.Undo()
	}))

	js.Global().Set("rpnLast", js.FuncOf(func(this js.Value, args []js.Value) any {
		return calc.Last()
	}))

	js.Global().Set("rpnSerialize", js.FuncOf(func(this js.Value, args []js.Value) any {
		return calc.Serialize()
	}))

	js.Global().Set("rpnRestore", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return ""
		}
		err := calc.Restore(args[0].String())
		if err != nil {
			return err.Error()
		}
		return ""
	}))

	js.Global().Set("rpnSetAngleMode", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return calc.GetStateJSON()
		}
		return calc.SetAngleMode(args[0].String())
	}))

	js.Global().Set("rpnSetBaseMode", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return calc.GetStateJSON()
		}
		return calc.SetBaseMode(args[0].String())
	}))

	<-make(chan bool)
}
