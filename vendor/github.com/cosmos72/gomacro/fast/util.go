/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017 Massimiliano Ghilardi
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Lesser General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Lesser General Public License for more details.
 *
 *     You should have received a copy of the GNU Lesser General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/lgpl>.
 *
 *
 * identifier.go
 *
 *  Created on Apr 01, 2017
 *      Author Massimiliano Ghilardi
 */

package fast

import (
	"go/ast"
	"go/constant"
	r "reflect"

	. "github.com/cosmos72/gomacro/base"
	xr "github.com/cosmos72/gomacro/xreflect"
)

func eFalse(*Env) bool {
	return false
}

func eTrue(*Env) bool {
	return true
}

func eNone(*Env) {
}

func eNil(*Env) r.Value {
	return Nil
}

func eXVNone(*Env) (r.Value, []r.Value) {
	return None, nil
}

func nop() {
}

var valueOfNopFunc = r.ValueOf(nop)

func asIdent(node ast.Expr) *ast.Ident {
	ident, _ := node.(*ast.Ident)
	return ident
}

func (e *Expr) TryAsPred() (value bool, fun func(*Env) bool, err bool) {
	if e.Untyped() {
		untyp := e.Value.(UntypedLit)
		if untyp.Kind != r.Bool {
			return false, nil, true
		}
		return constant.BoolVal(untyp.Obj), nil, false
	}
	if e.Type.Kind() != r.Bool {
		return false, nil, true
	}
	if e.Const() {
		v := r.ValueOf(e.Value)
		return v.Bool(), nil, false
	}
	switch fun := e.Fun.(type) {
	case func(*Env) bool:
		return false, fun, false
	case func(*Env) (r.Value, []r.Value):
		e.CheckX1()
		return false, func(env *Env) bool {
			ret, _ := fun(env)
			return ret.Bool()
		}, false
	default:
		fun1 := e.AsX1()
		return false, func(env *Env) bool {
			return fun1(env).Bool()
		}, false
	}
}

func (c *Comp) invalidPred(node ast.Expr, x *Expr) Stmt {
	return c.badPred("invalid", node, x)
}

func (c *Comp) badPred(reason string, node ast.Expr, x *Expr) Stmt {
	var t xr.Type = nil
	if x.NumOut() != 0 {
		t = x.Out(0)
	}
	c.Errorf("%s boolean predicate, expecting <bool> expression, found <%v>: %v",
		reason, t, node)
	return nil
}

func (e *Expr) AsX() func(*Env) {
	if e == nil || e.Const() {
		return nil
	}
	return funAsX(e.Fun)
}

func funAsX(any I) func(*Env) {
	switch fun := any.(type) {
	case nil:
	case func(*Env):
		return fun
	case func(*Env) r.Value:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) (r.Value, []r.Value):
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) bool:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) int:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) int8:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) int16:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) int32:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) int64:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) uint:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) uint8:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) uint16:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) uint32:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) uint64:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) uintptr:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) float32:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) float64:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) complex64:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) complex128:
		return func(env *Env) {
			fun(env)
		}
	case func(*Env) string:
		return func(env *Env) {
			fun(env)
		}
	default:
		Errorf("unsupported function type, cannot convert to func(*Env): %v <%v>", any, r.TypeOf(any))
	}
	return nil
}

// CheckX1() panics if given expression cannot be used in single-value context,
// for example because it returns no value at all.
// It just prints a warning if expression returns multiple values.
func (e *Expr) CheckX1() {
	if e != nil && e.Const() {
		return
	}
	if e == nil || e.NumOut() == 0 {
		Errorf("expression returns no values, cannot convert to func(env *Env) r.Value")
		return
	} else if e.NumOut() > 1 {
		Warnf("expression returns %d values, using only the first one: %v", e.NumOut(), e.Types)
	}
}

func (e *Expr) AsX1() func(*Env) r.Value {
	if e == nil {
		return eNil
	}
	if e.Const() {
		return valueAsX1(e.Value, e.Type, OptDefaults)
	}
	e.CheckX1()
	return funAsX1(e.Fun, e.Type)
}

func (e *Expr) AsXV(opts CompileOptions) func(*Env) (r.Value, []r.Value) {
	if e == nil {
		return eXVNone
	}
	if e.Const() {
		return valueAsXV(e.Value, e.Type, opts)
	}
	return funAsXV(e.Fun, e.Type)
}

func valueAsX1(any I, t xr.Type, opts CompileOptions) func(*Env) r.Value {
	convertuntyped := opts&OptKeepUntyped == 0
	untyp, untyped := any.(UntypedLit)
	if untyped && convertuntyped {
		if t == nil || t.ReflectType() == rtypeOfUntypedLit {
			t = untyp.DefaultType()
		}
		// Debugf("late conversion of untyped constant %v <%v> to <%v>", untyp, r.TypeOf(untyp), t)
		any = untyp.ConstTo(t)
	}
	v := r.ValueOf(any)
	if t != nil {
		rtype := t.ReflectType()
		if !v.IsValid() {
			v = r.Zero(rtype)
		} else if convertuntyped || !untyped {
			v = v.Convert(rtype)
		}
	}
	return func(*Env) r.Value {
		return v
	}
}

func valueAsXV(any I, t xr.Type, opts CompileOptions) func(*Env) (r.Value, []r.Value) {
	convertuntyped := opts&OptKeepUntyped == 0
	untyp, untyped := any.(UntypedLit)
	if convertuntyped {
		if untyped {
			if t == nil || t.ReflectType() == rtypeOfUntypedLit {
				t = untyp.DefaultType()
				// Debugf("valueAsXV: late conversion of untyped constant %v <%v> to its default type <%v>", untyp, r.TypeOf(untyp), t)
			} else {
				// Debugf("valueAsXV: late conversion of untyped constant %v <%v> to <%v>", untyp, r.TypeOf(untyp), t.ReflectType())
			}
			any = untyp.ConstTo(t)
		}
	}
	v := r.ValueOf(any)
	if t != nil {
		rtype := t.ReflectType()
		if ValueType(v) == nil {
			v = r.Zero(rtype)
		} else if convertuntyped || !untyped {
			v = v.Convert(rtype)
		}
	}
	return func(*Env) (r.Value, []r.Value) {
		return v, nil
	}
}

func funAsX1(fun I, t xr.Type) func(*Env) r.Value {
	// Debugf("funAsX1() %v -> %v", TypeOf(fun), t)
	var rt r.Type
	if t != nil {
		rt = t.ReflectType()
	}
	switch fun := fun.(type) {
	case nil:
	case func(*Env):
		if fun == nil {
			break
		}
		return func(env *Env) r.Value {
			fun(env)
			return None
		}
	case func(*Env) r.Value:
		return fun
	case func(*Env) (r.Value, []r.Value):
		return func(env *Env) r.Value {
			ret, _ := fun(env)
			return ret
		}
	case func(*Env) bool:
		if rt == nil || rt == TypeOfBool {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int:
		if rt == nil || rt == TypeOfInt {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int8:
		if rt == nil || rt == TypeOfInt8 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int16:
		if rt == nil || rt == TypeOfInt16 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int32:
		if rt == nil || rt == TypeOfInt32 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int64:
		if rt == nil || rt == TypeOfInt64 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint:
		if rt == nil || rt == TypeOfUint {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint8:
		if rt == nil || rt == TypeOfUint8 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint16:
		if rt == nil || rt == TypeOfUint16 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint32:
		if rt == nil || rt == TypeOfUint32 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint64:
		if rt == nil || rt == TypeOfUint64 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uintptr:
		if rt == nil || rt == TypeOfUintptr {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) float32:
		if rt == nil || rt == TypeOfFloat32 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) float64:
		if rt == nil || rt == TypeOfFloat64 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) complex64:
		if rt == nil || rt == TypeOfComplex64 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) complex128:
		if rt == nil || rt == TypeOfComplex128 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) string:
		if rt == nil || rt == TypeOfString {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *bool:
		if rt == nil || rt == TypeOfPtrBool {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *int:
		if rt == nil || rt == TypeOfPtrInt {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *int8:
		if rt == nil || rt == TypeOfPtrInt8 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *int16:
		if rt == nil || rt == TypeOfPtrInt16 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *int32:
		if rt == nil || rt == TypeOfPtrInt32 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *int64:
		if rt == nil || rt == TypeOfPtrInt64 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *uint:
		if rt == nil || rt == TypeOfPtrUint {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *uint8:
		if rt == nil || rt == TypeOfPtrUint8 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *uint16:
		if rt == nil || rt == TypeOfPtrUint16 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *uint32:
		if rt == nil || rt == TypeOfPtrUint32 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *uint64:
		if rt == nil || rt == TypeOfPtrUint64 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *uintptr:
		if rt == nil || rt == TypeOfPtrUintptr {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *float32:
		if rt == nil || rt == TypeOfPtrFloat32 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *float64:
		if rt == nil || rt == TypeOfPtrFloat64 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) *complex64:
		if rt == nil || rt == TypeOfPtrComplex64 {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env))
			}
		} else {
			return func(env *Env) r.Value {
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	default:
		Errorf("unsupported expression type, cannot convert to func(*Env) r.Value: %v <%v>", fun, r.TypeOf(fun))
	}
	return nil
}

func funAsXV(fun I, t xr.Type) func(*Env) (r.Value, []r.Value) {
	// Debugf("funAsXV() %v -> %v", TypeOf(fun), t)
	var rt r.Type
	if t != nil {
		rt = t.ReflectType()
	}
	switch fun := fun.(type) {
	case nil:
	case func(*Env):
		if fun == nil {
			break
		}
		return func(env *Env) (r.Value, []r.Value) {
			fun(env)
			return None, nil
		}
	case func(*Env) r.Value:
		return func(env *Env) (r.Value, []r.Value) {
			return fun(env), nil
		}
	case func(*Env) (r.Value, []r.Value):
		return fun
	case func(*Env) bool:
		if rt == nil || rt == TypeOfBool {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) int:
		if rt == nil || rt == TypeOfInt {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) int8:
		if rt == nil || rt == TypeOfInt8 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) int16:
		if rt == nil || rt == TypeOfInt16 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) int32:
		if rt == nil || rt == TypeOfInt32 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) int64:
		if rt == nil || rt == TypeOfInt64 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) uint:
		if rt == nil || rt == TypeOfUint {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) uint8:
		if rt == nil || rt == TypeOfUint8 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) uint16:
		if rt == nil || rt == TypeOfUint16 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) uint32:
		if rt == nil || rt == TypeOfUint32 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) uint64:
		if rt == nil || rt == TypeOfUint64 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) uintptr:
		if rt == nil || rt == TypeOfUintptr {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) float32:
		if rt == nil || rt == TypeOfFloat32 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) float64:
		if rt == nil || rt == TypeOfFloat64 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) complex64:
		if rt == nil || rt == TypeOfComplex64 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) complex128:
		if rt == nil || rt == TypeOfComplex128 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) string:
		if rt == nil || rt == TypeOfString {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *bool:
		if rt == nil || rt == TypeOfPtrBool {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *int:
		if rt == nil || rt == TypeOfPtrInt {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *int8:
		if rt == nil || rt == TypeOfPtrInt8 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *int16:
		if rt == nil || rt == TypeOfPtrInt16 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *int32:
		if rt == nil || rt == TypeOfPtrInt32 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *int64:
		if rt == nil || rt == TypeOfPtrInt64 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *uint:
		if rt == nil || rt == TypeOfPtrUint {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *uint8:
		if rt == nil || rt == TypeOfPtrUint8 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *uint16:
		if rt == nil || rt == TypeOfPtrUint16 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *uint32:
		if rt == nil || rt == TypeOfPtrUint32 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *uint64:
		if rt == nil || rt == TypeOfPtrUint64 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *uintptr:
		if rt == nil || rt == TypeOfPtrUintptr {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *float32:
		if rt == nil || rt == TypeOfPtrFloat32 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *float64:
		if rt == nil || rt == TypeOfPtrFloat64 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	case func(*Env) *complex64:
		if rt == nil || rt == TypeOfPtrComplex64 {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)), nil
			}
		} else {
			return func(env *Env) (r.Value, []r.Value) {
				return r.ValueOf(fun(env)).Convert(rt), nil
			}
		}
	default:
		Errorf("unsupported expression, cannot convert to func(*Env) (r.Value, []r.Value) : %v <%v>",
			fun, r.TypeOf(fun))
	}
	return nil
}

func (e *Expr) AsStmt() Stmt {
	if e == nil || e.Const() {
		return nil
	}
	return funAsStmt(e.Fun)
}

func funAsStmt(fun I) Stmt {
	var ret func(env *Env) (Stmt, *Env)

	switch fun := fun.(type) {
	case nil:
	case func(*Env):
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) r.Value:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) (r.Value, []r.Value):
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) bool:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) int:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) int8:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) int16:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) int32:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) int64:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) uint:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) uint8:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) uint16:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) uint32:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) uint64:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) uintptr:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) float32:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) float64:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) complex64:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) complex128:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case func(*Env) string:
		ret = func(env *Env) (Stmt, *Env) {
			fun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	default:

		Errorf("unsupported expression type, cannot convert to Stmt : %v <%v>",
			fun, r.TypeOf(fun))
	}
	return ret
}

// funTypeOut returns the first return type of given function
func funTypeOut(fun I) r.Type {
	rt := r.TypeOf(fun)
	if rt == nil || rt.Kind() != r.Func || rt.NumOut() == 0 {
		return nil
	}
	return rt.Out(0)
}

// funTypeOuts returns the return types of given function
func funTypeOuts(fun I) []r.Type {
	rt := r.TypeOf(fun)
	if rt == nil || rt.Kind() != r.Func {
		return []r.Type{rt}
	}
	n := rt.NumOut()
	rts := make([]r.Type, n)
	for i := 0; i < n; i++ {
		rts[i] = rt.Out(i)
	}
	return rts
}

// exprList merges together a list of expressions,
// and returns an expression that evaluates each one
func exprList(list []*Expr, opts CompileOptions) *Expr {
	// skip constant expressions (except the last one)
	var n int
	for i, ni := 0, len(list)-1; i <= ni; i++ {
		// preserve the last expression even if constant
		// because it will be returned to the user
		if i == ni || !list[i].Const() {
			list[n] = list[i]
			n++
		}
	}
	switch n {
	case 0:
		return nil
	case 1:
		return list[0]
	}
	list = list[:n]

	funs := make([]func(*Env), n-1)
	for i := range funs {
		funs[i] = list[i].AsX()
	}
	return &Expr{
		Lit:   Lit{Type: list[n-1].Type},
		Types: list[n-1].Types,
		Fun:   funList(funs, list[n-1], opts),
	}
}

// funList merges together a list of functions,
// and returns a function that evaluates each one
func funList(funs []func(*Env), last *Expr, opts CompileOptions) I {
	var rt r.Type
	if last.Type != nil {
		// keep untyped constants only if requested
		if opts != OptKeepUntyped && last.Untyped() {
			last.ConstTo(last.DefaultType())
		}
		rt = last.Type.ReflectType()
	}
	switch fun := last.WithFun().(type) {
	case nil:
		return func(env *Env) {
			for _, f := range funs {
				f(env)
			}
		}
	case func(*Env):
		return func(env *Env) {
			for _, f := range funs {
				f(env)
			}
			fun(env)
		}
	case func(*Env) r.Value:
		return func(env *Env) r.Value {
			for _, f := range funs {
				f(env)
			}
			return fun(env)
		}
	case func(*Env) (r.Value, []r.Value):
		return func(env *Env) (r.Value, []r.Value) {
			for _, f := range funs {
				f(env)
			}
			return fun(env)
		}
	case func(*Env) bool:
		if rt == nil || rt == TypeOfBool {
			return func(env *Env) bool {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int:
		if rt == nil || rt == TypeOfInt {
			return func(env *Env) int {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int8:
		if rt == nil || rt == TypeOfInt8 {
			return func(env *Env) int8 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int16:
		if rt == nil || rt == TypeOfInt16 {
			return func(env *Env) int16 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int32:
		if rt == nil || rt == TypeOfInt32 {
			return func(env *Env) int32 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) int64:
		if rt == nil || rt == TypeOfInt64 {
			return func(env *Env) int64 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint:
		if rt == nil || rt == TypeOfUint {
			return func(env *Env) uint {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint8:
		if rt == nil || rt == TypeOfUint8 {
			return func(env *Env) uint8 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint16:
		if rt == nil || rt == TypeOfUint16 {
			return func(env *Env) uint16 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint32:
		if rt == nil || rt == TypeOfUint32 {
			return func(env *Env) uint32 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uint64:
		if rt == nil || rt == TypeOfUint64 {
			return func(env *Env) uint64 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) uintptr:
		if rt == nil || rt == TypeOfUintptr {
			return func(env *Env) uintptr {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) float32:
		if rt == nil || rt == TypeOfFloat32 {
			return func(env *Env) float32 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) float64:
		if rt == nil || rt == TypeOfFloat64 {
			return func(env *Env) float64 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) complex64:
		if rt == nil || rt == TypeOfComplex64 {
			return func(env *Env) complex64 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) complex128:
		if rt == nil || rt == TypeOfComplex128 {
			return func(env *Env) complex128 {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	case func(*Env) string:
		if rt == nil || rt == TypeOfString {
			return func(env *Env) string {
				for _, f := range funs {
					f(env)
				}
				return fun(env)
			}
		} else {
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				return r.ValueOf(fun(env)).Convert(rt)
			}
		}
	default:
		switch last.NumOut() {
		case 0:
			fun := last.AsX()
			return func(env *Env) {
				for _, f := range funs {
					f(env)
				}
				fun(env)
			}
		case 1:
			var zero r.Value
			if rt != nil {
				zero = r.Zero(rt)
			}
			fun := last.AsX1()
			return func(env *Env) r.Value {
				for _, f := range funs {
					f(env)
				}
				ret := fun(env)
				if ret == Nil {
					ret = zero
				} else if rt != nil && rt != ret.Type() {
					ret = ret.Convert(rt)
				}
				return ret
			}
		default:
			var zero []r.Value
			var rt []r.Type
			for i, t := range last.Types {
				if t != nil {
					rt[i] = t.ReflectType()
					zero[i] = r.Zero(rt[i])
				}
			}
			fun := last.AsXV(opts)
			return func(env *Env) (r.Value, []r.Value) {
				for _, f := range funs {
					f(env)
				}
				_, rets := fun(env)
				for i, ret := range rets {
					if ret == Nil {
						rets[i] = zero[i]
					} else if rt != nil && rt[i] != ret.Type() {
						rets[i] = ret.Convert(rt[i])
					}
				}
				return rets[0], rets
			}
		}
	}
}

// unwrapBind compiles a conversion from a "mis-typed" bind stored in env.Binds[] as reflect.Value
// into a correctly-typed expression
func unwrapBind(bind *Bind, t xr.Type) *Expr {
	idx := bind.Desc.Index()
	var ret I
	switch t.Kind() {
	case r.Bool:
		ret = func(env *Env) bool {
			return env.Binds[idx].Bool()
		}
	case r.Int:
		ret = func(env *Env) int {
			return int(env.Binds[idx].Int())
		}
	case r.Int8:
		ret = func(env *Env) int8 {
			return int8(env.Binds[idx].Int())
		}
	case r.Int16:
		ret = func(env *Env) int16 {
			return int16(env.Binds[idx].Int())
		}
	case r.Int32:
		ret = func(env *Env) int32 {
			return int32(env.Binds[idx].Int())
		}
	case r.Int64:
		ret = func(env *Env) int64 {
			return env.Binds[idx].Int()
		}
	case r.Uint:
		ret = func(env *Env) uint {
			return uint(env.Binds[idx].Uint())
		}
	case r.Uint8:
		ret = func(env *Env) uint8 {
			return uint8(env.Binds[idx].Uint())
		}
	case r.Uint16:
		ret = func(env *Env) uint16 {
			return uint16(env.Binds[idx].Uint())
		}
	case r.Uint32:
		ret = func(env *Env) uint32 {
			return uint32(env.Binds[idx].Uint())
		}
	case r.Uint64:
		ret = func(env *Env) uint64 {
			return env.Binds[idx].Uint()
		}
	case r.Uintptr:
		ret = func(env *Env) uintptr {
			return uintptr(env.Binds[idx].Uint())
		}
	case r.Float32:
		ret = func(env *Env) float32 {
			return float32(env.Binds[idx].Float())
		}
	case r.Float64:
		ret = func(env *Env) float64 {
			return env.Binds[idx].Float()
		}
	case r.Complex64:
		ret = func(env *Env) complex64 {
			return complex64(env.Binds[idx].Complex())
		}
	case r.Complex128:
		ret = func(env *Env) complex128 {
			return env.Binds[idx].Complex()
		}
	case r.String:
		ret = func(env *Env) string {
			return env.Binds[idx].String()
		}
	default:
		rtype := t.ReflectType()
		zero := r.Zero(rtype)
		ret = func(env *Env) r.Value {
			v := env.Binds[idx]
			if !v.IsValid() {
				v = zero
			} else if v.Type() != rtype {
				v = v.Convert(rtype)
			}
			return v
		}
	}
	return exprFun(t, ret)
}

// unwrapBindUp1 compiles a conversion from a "mis-typed" bind stored in env.Outer.Binds[] as reflect.Value
// into a correctly-typed expression
func unwrapBindUp1(bind *Bind, t xr.Type) *Expr {
	idx := bind.Desc.Index()
	var ret I
	switch t.Kind() {
	case r.Bool:
		ret = func(env *Env) bool {
			return env.Outer.Binds[idx].Bool()
		}
	case r.Int:
		ret = func(env *Env) int {
			return int(env.Outer.Binds[idx].Int())
		}
	case r.Int8:
		ret = func(env *Env) int8 {
			return int8(env.Outer.Binds[idx].Int())
		}
	case r.Int16:
		ret = func(env *Env) int16 {
			return int16(env.Outer.Binds[idx].Int())
		}
	case r.Int32:
		ret = func(env *Env) int32 {
			return int32(env.Outer.Binds[idx].Int())
		}
	case r.Int64:
		ret = func(env *Env) int64 {
			return env.Outer.Binds[idx].Int()
		}
	case r.Uint:
		ret = func(env *Env) uint {
			return uint(env.Outer.Binds[idx].Uint())
		}
	case r.Uint8:
		ret = func(env *Env) uint8 {
			return uint8(env.Outer.Binds[idx].Uint())
		}
	case r.Uint16:
		ret = func(env *Env) uint16 {
			return uint16(env.Outer.Binds[idx].Uint())
		}
	case r.Uint32:
		ret = func(env *Env) uint32 {
			return uint32(env.Outer.Binds[idx].Uint())
		}
	case r.Uint64:
		ret = func(env *Env) uint64 {
			return env.Outer.Binds[idx].Uint()
		}
	case r.Uintptr:
		ret = func(env *Env) uintptr {
			return uintptr(env.Outer.Binds[idx].Uint())
		}
	case r.Float32:
		ret = func(env *Env) float32 {
			return float32(env.Outer.Binds[idx].Float())
		}
	case r.Float64:
		ret = func(env *Env) float64 {
			return env.Outer.Binds[idx].Float()
		}
	case r.Complex64:
		ret = func(env *Env) complex64 {
			return complex64(env.Outer.Binds[idx].Complex())
		}
	case r.Complex128:
		ret = func(env *Env) complex128 {
			return env.Outer.Binds[idx].Complex()
		}
	case r.String:
		ret = func(env *Env) string {
			return env.Outer.Binds[idx].String()
		}
	default:
		rtype := t.ReflectType()
		zero := r.Zero(rtype)
		ret = func(env *Env) r.Value {
			v := env.Outer.Binds[idx]
			if !v.IsValid() {
				v = zero
			} else if v.Type() != rtype {
				v = v.Convert(rtype)
			}
			return v
		}
	}
	return exprFun(t, ret)
}
