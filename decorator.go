package main

import (
	"fmt"
	"reflect"
)

func Decorator(decoPtr, fn interface{}) (err error) {
	var decoratedFun, targetFunc reflect.Value

	decoratedFun = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)

	v := reflect.MakeFunc(targetFunc.Type(),
		func(in []reflect.Value) (out []reflect.Value) {
			fmt.Println("before")
			out = targetFunc.Call(in)
			fmt.Println("after")
			return
		})
	decoratedFun.Set(v)
	return
}

func foo(a, b, c int) int {
	fmt.Printf("%d, %d, %d \n", a, b, c)
	return a + b + c
}

func bar(a, b string) string {
	fmt.Printf("%s, %s \n", a, b)
	return a + b
}

func main() {
	//通过函数签名
	type MyFoo func(int, int, int) int
	var myfoo MyFoo
	err := Decorator(&myfoo, foo)
	if err != nil {
		panic(err)
	}
	i := myfoo(1, 2, 3)
	fmt.Println("result: ", i)
	//不用函数签名
	mybar := bar
	_ = Decorator(&mybar, bar)
	mybar("Hello", "world!")
}
