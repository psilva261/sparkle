package console

import (
	"log"

	"github.com/psilva261/sparkle/js"
	"github.com/psilva261/sparkle/require"
	_ "github.com/psilva261/sparkle/util"
)

type Console struct {
	runtime *js.Runtime
	util    *js.Object
	printer Printer
}

type Printer interface {
	Log(string)
	Warn(string)
	Error(string)
}

type PrinterFunc func(s string)

func (p PrinterFunc) Log(s string) { p(s) }

func (p PrinterFunc) Warn(s string) { p(s) }

func (p PrinterFunc) Error(s string) { p(s) }

var defaultPrinter Printer = PrinterFunc(func(s string) { log.Print(s) })

func (c *Console) log(p func(string)) func(js.FunctionCall) js.Value {
	return func(call js.FunctionCall) js.Value {
		if format, ok := js.AssertFunction(c.util.Get("format")); ok {
			ret, err := format(c.util, call.Arguments...)
			if err != nil {
				panic(err)
			}

			p(ret.String())
		} else {
			panic(c.runtime.NewTypeError("util.format is not a function"))
		}

		return nil
	}
}

func Require(runtime *js.Runtime, module *js.Object) {
	requireWithPrinter(defaultPrinter)(runtime, module)
}

func RequireWithPrinter(printer Printer) require.ModuleLoader {
	return requireWithPrinter(printer)
}

func requireWithPrinter(printer Printer) require.ModuleLoader {
	return func(runtime *js.Runtime, module *js.Object) {
		c := &Console{
			runtime: runtime,
			printer: printer,
		}

		c.util = require.Require(runtime, "util").(*js.Object)

		o := module.Get("exports").(*js.Object)
		o.Set("log", c.log(c.printer.Log))
		o.Set("error", c.log(c.printer.Error))
		o.Set("warn", c.log(c.printer.Warn))
	}
}

func Enable(runtime *js.Runtime) {
	runtime.Set("console", require.Require(runtime, "console"))
}

func init() {
	require.RegisterNativeModule("console", Require)
}
