package gel

import (
	"testing"
)

func Benchmark_expr(b *testing.B) {
	params := make(map[string]any)
	params["Origin"] = "MOW"
	params["Country"] = "RU"
	params["Adults"] = 1
	params["Value"] = 100

	el, err := Compile(`(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`)
	if err != nil {
		b.Fatal(err)
	}

	var out any

	b.ResetTimer()
	for b.Loop() {
		out, err = el.Calculate(params)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if !(out.(bool)) {
		b.Fatal(out)
	}
}

func Benchmark_expr_eval(b *testing.B) {
	params := make(map[string]any)
	params["Origin"] = "MOW"
	params["Country"] = "RU"
	params["Adults"] = 1
	params["Value"] = 100

	var out any
	var err error

	b.ResetTimer()
	for b.Loop() {
		out, err = Calculate(`(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`, params)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if !(out.(bool)) {
		b.Fatal(out)
	}
}

func Benchmark_arrayIndex(b *testing.B) {
	env := map[string]any{
		"arr": make([]int, 100),
	}
	for i := 0; i < 100; i++ {
		env["arr"].([]int)[i] = i
	}

	el, err := Compile(`arr[50]`)
	if err != nil {
		b.Fatal(err)
	}

	var out any
	b.ResetTimer()
	for b.Loop() {
		out, err = el.Calculate(env)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if out != 50 {
		b.Fatal(out)
	}
}

func Benchmark_envStruct(b *testing.B) {
	type Price struct {
		Value int
	}
	type Env struct {
		Price Price
	}

	el, err := Compile(`Price.Value > 0`)
	if err != nil {
		b.Fatal(err)
	}

	env := Env{Price: Price{Value: 1}}

	var out any
	b.ResetTimer()
	for b.Loop() {
		out, err = el.Calculate(env)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if !(out.(bool)) {
		b.Fatal(out)
	}
}

func Benchmark_envMap(b *testing.B) {
	type Price struct {
		Value int
	}
	env := map[string]any{
		"price": Price{Value: 1},
	}

	el, err := Compile(`price.Value > 0`)
	if err != nil {
		b.Fatal(err)
	}

	var out any
	b.ResetTimer()
	for b.Loop() {
		out, err = el.Calculate(env)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if !(out.(bool)) {
		b.Fatal(out)
	}
}

type CallEnv struct {
	A      int
	B      int
	C      int
	Fn     func() bool
	FnFast func(...any) any
	Foo    CallFoo
}

func (CallEnv) Func() string {
	return "func"
}

type CallFoo struct {
	D int
	E int
	F int
}

func (CallFoo) Method() string {
	return "method"
}

func Benchmark_callMethod(b *testing.B) {
	el, err := Compile(`Foo.Method()`)
	if err != nil {
		b.Fatal(err)
	}

	env := CallEnv{}

	var out any
	b.ResetTimer()
	for b.Loop() {
		out, err = el.Calculate(env)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if out != "method" {
		b.Fatal(out)
	}
}

func Benchmark_largeStructAccess(b *testing.B) {
	type Env struct {
		Data  [1024 * 1024 * 10]byte
		Field int
	}

	el, err := Compile(`Field > 0 && Field > 1 && Field < 99`)
	if err != nil {
		b.Fatal(err)
	}

	env := Env{Field: 21}

	var out any
	b.ResetTimer()
	for b.Loop() {
		out, err = el.Calculate(&env)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if !(out.(bool)) {
		b.Fatal(out)
	}
}

func Benchmark_largeNestedStructAccess(b *testing.B) {
	type Env struct {
		Inner struct {
			Data  [1024 * 1024 * 10]byte
			Field int
		}
	}

	el, err := Compile(`Inner.Field > 0 && Inner.Field > 1 && Inner.Field < 99`)
	if err != nil {
		b.Fatal(err)
	}

	env := Env{}
	env.Inner.Field = 21

	var out any
	b.ResetTimer()
	for b.Loop() {
		out, err = el.Calculate(&env)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if !(out.(bool)) {
		b.Fatal(out)
	}
}

func Benchmark_largeNestedArrayAccess(b *testing.B) {
	type Env struct {
		Data [1][1024 * 1024 * 10]byte
	}

	el, err := Compile(`Data[0][0] > 0`)
	if err != nil {
		b.Fatal(err)
	}

	env := Env{}
	env.Data[0][0] = 1

	var out any
	b.ResetTimer()
	for b.Loop() {
		out, err = el.Calculate(&env)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if !(out.(bool)) {
		b.Fatal(out)
	}
}
