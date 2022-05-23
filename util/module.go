package util

import (
	"bytes"
	"github.com/psilva261/sparkle/js"
	"github.com/psilva261/sparkle/require"
)

type Util struct {
	runtime *js.Runtime
}

func (u *Util) format(f rune, val js.Value, w *bytes.Buffer) bool {
	switch f {
	case 's':
		w.WriteString(val.String())
	case 'd':
		w.WriteString(val.ToNumber().String())
	case 'j':
		if json, ok := u.runtime.Get("JSON").(*js.Object); ok {
			if stringify, ok := js.AssertFunction(json.Get("stringify")); ok {
				res, err := stringify(json, val)
				if err != nil {
					panic(err)
				}
				w.WriteString(res.String())
			}
		}
	case '%':
		w.WriteByte('%')
		return false
	default:
		w.WriteByte('%')
		w.WriteRune(f)
		return false
	}
	return true
}

func (u *Util) Format(b *bytes.Buffer, f string, args ...js.Value) {
	pct := false
	argNum := 0
	for _, chr := range f {
		if pct {
			if argNum < len(args) {
				if u.format(chr, args[argNum], b) {
					argNum++
				}
			} else {
				b.WriteByte('%')
				b.WriteRune(chr)
			}
			pct = false
		} else {
			if chr == '%' {
				pct = true
			} else {
				b.WriteRune(chr)
			}
		}
	}

	for _, arg := range args[argNum:] {
		b.WriteByte(' ')
		b.WriteString(arg.String())
	}
}

func (u *Util) js_format(call js.FunctionCall) js.Value {
	var b bytes.Buffer
	var fmt string

	if arg := call.Argument(0); !js.IsUndefined(arg) {
		fmt = arg.String()
	}

	var args []js.Value
	if len(call.Arguments) > 0 {
		args = call.Arguments[1:]
	}
	u.Format(&b, fmt, args...)

	return u.runtime.ToValue(b.String())
}

func Require(runtime *js.Runtime, module *js.Object) {
	u := &Util{
		runtime: runtime,
	}
	obj := module.Get("exports").(*js.Object)
	obj.Set("format", u.js_format)
}

func New(runtime *js.Runtime) *Util {
	return &Util{
		runtime: runtime,
	}
}

func init() {
	require.RegisterNativeModule("util", Require)
}
