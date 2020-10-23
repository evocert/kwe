package console

import (
	"log"

	"github.com/lnksnk/ecma/es6"
	"github.com/lnksnk/nodejs/require"
	_ "github.com/lnksnk/nodejs/util"
)

type Console struct {
	runtime *es6.Runtime
	util    *es6.Object
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

func (c *Console) log(p func(string)) func(es6.FunctionCall) es6.Value {
	return func(call es6.FunctionCall) es6.Value {
		if format, ok := es6.AssertFunction(c.util.Get("format")); ok {
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

func Require(runtime *es6.Runtime, module *es6.Object) {
	requireWithPrinter(defaultPrinter)(runtime, module)
}

func RequireWithPrinter(printer Printer) require.ModuleLoader {
	return requireWithPrinter(printer)
}

func requireWithPrinter(printer Printer) require.ModuleLoader {
	return func(runtime *es6.Runtime, module *es6.Object) {
		c := &Console{
			runtime: runtime,
			printer: printer,
		}

		c.util = require.Require(runtime, "util").(*es6.Object)

		o := module.Get("exports").(*es6.Object)
		o.Set("log", c.log(c.printer.Log))
		o.Set("error", c.log(c.printer.Error))
		o.Set("warn", c.log(c.printer.Warn))
	}
}

func Enable(runtime *es6.Runtime) {
	runtime.Set("console", require.Require(runtime, "console"))
}

func init() {
	require.RegisterNativeModule("console", Require)
}
