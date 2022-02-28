---
title: "Go 语言闭包及其应用"
slug: "goclosure"
date: 2021-10-28T23:59:44+08:00
description: "在学习 Go 语言闭包之前，需要先对匿名函数有所了解。匿名函数与常规函数的不同点在于"
draft: false
tags: ["go"]
---
## 匿名函数
在学习 Go 语言闭包之前，需要先对匿名函数有所了解。匿名函数与常规函数的不同点在于：**匿名函数没有指定函数名称**，一般用法为在原地执行一次或者赋值给一个变量供以后使用。

```go
func main() {
	// execute locally
	func() {
		fmt.Println("Anonymous function execute locally!")
	}()

	// assign to a variable
	var anonyfunc func() = func() {
		fmt.Println("Anoymous function assign to a variable anonyfunc")
	}
	anonyfunc()

	// change value of anonyfunc
	anonyfunc = somefunc
	anonyfunc()
}

func somefunc() {
	fmt.Println("some standard func")
}
```

这段代码简单展示了匿名函数原地执行和赋值给变量执行的例子，并且在运行期间改变了 `anonyfunc` 的值。

## Go 语言闭包
闭包可以看作是一种特殊的匿名函数，该匿名函数**引用了其包裹函数的变量**。
[closure.go:](https://github.com/kfngp/kfngp.github.io/blob/main/examples/closure/closure.go)
```go
func intSeq() func() int {
    i := 0
	
	// the returned anonymous function references variable i
    return func() int {
        i++
        return i
    }
}

func main() {
    nextInt := intSeq()

    fmt.Println(nextInt())
    fmt.Println(nextInt())
    fmt.Println(nextInt())
	
	// call intSeq() generate a new enviroment, i in newINts set to 0
    newInts := intSeq()
    fmt.Println(newInts())
}
```

`intSeq` 返回的匿名函数引用其包裹函数 `intSeq` 的变量 `i`，程序的最终执行结果为：
```
1
2
3
1
```

注意在没有给匿名函数传递参数的情况下，匿名函数却引用到了包裹函数的变量 `i`，这使其称之为闭包。

### 闭包的实现
现在请读者思考一个问题：`intSeq()` 中的局部变量 `i` 应该在内存中的什么位置分配？`intSeq()` 函数的栈上？

首先肯定不能在 `intSeq` 函数的栈上分配，因为函数 `intSeq` 返回后，对应的函数栈就失效了，此时匿名函数中的变量 `i` 就引用了一个失效的位置。所以闭包环境中引用的变量不能在栈上分配。

实际上 Go 语言会进行逃逸分析，将闭包环境中的变量**分配至堆上**。使用 `go build --gcflags=-m examples/closure/closure.go` 查看输出，可以看到如下一行输出信息：
```bash
examples/closure/closure.go:6:2: moved to heap: i
```
使用 `go tool compile -N -l -S examples/closure/closure.go` 查看汇编代码可以看到更详细的实现：
```armasm
...
MOVD    $type.noalg.struct { F uintptr; "".i *int }(SB), R0
MOVD    R0, 8(RSP)
PCDATA  $1, $1
CALL    runtime.newobject(SB)
...
```
可以看到闭包实际上将其函数和其环境中的变量打包成了一个结构体使用 `runtime.newobject` 将其分配到了堆内存上。因此我们也可以得出一个结论**只要闭包函数对象存在，其引用（捕获）的那些变量就会存在，就不会被回收。**
### 闭包的应用
闭包在实际实际开发中有很多用处，这里介绍一下在网络开发中使用闭包的例子。
在使用 Go 进行 Web 开发时，可以在实际路由处理函数前后集成中间件进行额外的逻辑处理，比如使用 Gin 框架进行服务器开发时，默认会使用 [Logger](https://github.com/gin-gonic/gin/blob/1c2aa59b20c8cff5b3c601708afe22100eac25e6/logger.go#L182) 中间件进行日志处理。其实现为：
```go
func Logger() HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

func LoggerWithFormatter(f LogFormatter) HandlerFunc {
	// closure function
	return LoggerWithConfig(LoggerConfig{
		Formatter: f,
	})
}

type HandlerFunc func(*Context)
```
---
[5 Useful Ways to Use Closures in Go](https://www.calhoun.io/5-useful-ways-to-use-closures-in-go/)

[What is a Closure?](https://www.calhoun.io/what-is-a-closure/)

[Go中被闭包捕获的变量何时会被回收](https://tonybai.com/2021/08/09/when-variables-captured-by-closures-are-recycled-in-go/)