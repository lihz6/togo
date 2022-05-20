# The Way to Go

## Lang

**1. nil 零值都有类型，但接口零值不能直接反射**

```go
sumpfic := []any{
  unsafe.Pointer(nil),
  map[any]any(nil),
  (func())(nil),
  chan any(nil),
  (*any)(nil),
  []any(nil),
  // any(nil),
}
for _, v := range sumpfic {
  v := reflect.ValueOf(v)
  fmt.Printf("%s, %v\n", v.Kind(), v.Type().Size())
}
// Output:
// unsafe.Pointer, 8
// map, 8
// func, 8
// chan, 8
// ptr, 8
// slice, 24
```

```go
// 接口值自动拆箱装箱，不能直接反射
reflect.TypeOf(nil) == nil
reflect.ValueOf(any(nil)).Type() // panic
reflect.ValueOf(nil).Kind() == reflect.Invalid
```

```go
v := reflect.ValueOf((*any)(nil)).Elem()
fmt.Printf("%s, %v\n", v.Kind(), v.Type().Size())
// Output: interface, 16
```

**2. if switch for 会产生一个夹层作用域**

```go
// naked return nil的陷阱
func bad() (err error) {
  if v, err := todo(); err != nil {
    return // bad, return nil
  }
}
func good() (err error) {
  v, err := todo()
  if err != nil {
    return // good, return err
  }
}
```

```go
// for循环指针方法和闭包的陷阱
type Val int
func (v *Val) want(i int) {
  fmt.Printf("want %v, got %v\n", i, *v)
}
for i, v := range [...]Val{0, 1, 2} {
  defer v.want(i)
}
// Output:
// want 2, got 2
// want 1, got 2
// want 0, got 2
```

```go
// for循环指针和数组切片的陷阱
var (
  s1 [][]int
  s2 []*int
  i  int
)
for i, a := range [...][1]int{{0}, {1}, {2}} {
  s1 = append(s1, a[:])
  s2 = append(s2, &i)
}
for i := 0; i < 3; i++ {
  fmt.Printf("want %v, got %v, %v\n", i, s1[i][0], *s2[i])
}
fmt.Printf("want 2, got %v\n", i)
// Output:
// want 0, got 2, 2
// want 1, got 2, 2
// want 2, got 2, 2
// want 2, got 0
```

**3. make 只用于 Map, Slice 和 Chan, new 适用于一切类型**

```go
var m map[comparable]T = make(M, size=0)
var s []T = make(S, len, cap=len)
var c chan T = make(C, buf=0)
```

```go
new(T) != (*T)(nil)
// get zero value of any type
func Zero[T any]() T {
	return *new(T)
}
func Zero[T any]() (t T) {
	return
}
```

## Chan

**1. 通道的两种类型、两类值及其操作**

|    ch     | chan<- T | <-chan T |  nil  | closed |
| :-------: | :------: | :------: | :---: | :----: |
|   <- ch   |  error   |    OK    | never | T, ok  |
|  ch <- T  |    OK    |  error   | never | panic  |
| close(ch) |    OK    |  error   | panic | panic  |

```go
ch := make(chan int, 2)
ch <- 1
close(ch)
fmt.Println(len(ch), cap(ch)) // 1 2
fmt.Println(<- ch, <-ch) // 1 0
ch <- 2 // panic
```

**2. 等待发送中的通道不能关闭**

```go
ch := make(chan int)
go func() {
  ch <- 0 // panic: send on closed channel
}()
// wait a second
close(ch)
// wait a second
```

**3. 在发送端不能无法检测通道是否关闭，接收端可以**

```go
ch <- 0 // not sure ch closed or not

v, ok := ch // ok indicates not closed
```

**4. for range 感知通道关闭，for select 不感知通道关闭**

```go
ch := make(chan int)
close(ch)
for range ch {
  // never
}
for {
  // forever
  select {
  case <- ch:
  }
}
close(ch)
```

**5. 零值通道不能关闭，发送和接收永远阻塞**

```go
ch := chan int(nil)
go func() {
  // goroutine leaks!
  ch <- 0
}()
go func() {
  // goroutine leaks!
  <- ch
}()
close(ch) // panic: close of nil channel
```

**6. 非零通道只能通过 make 创建**

```go
// ch1 := make(chan int, 0)
ch1 := make(chan int)
ch2 := make(chan int, 1)

// wanted:
// ch := chan int{0, 1}
// ch := make(chan int, 2)
// ch <- 0
// ch <- 1
```

**7. 为避免协程泄露通道关闭的最佳实践**

```go
// 哪个函数创建的由哪个函数负责关闭
func closer(done <-chan struct{}) <-chan int {
  ch := make(chan int)
  go func() {
    defer close(ch)
    select {
      case ch <- 0:
      case <- done:
    }
  }()
  return ch
}
```

```go
// 函数的通道参数或返回值必须要有方向
func bad(ch chan int) chan int {}
func good(ch <-chan int) <-chan int {}

// 与其返回发送端，不如接受接收端
func bad() (ch chan<- int) {
  // 如把发送端返回，则由caller负责关闭
}
func good(ch <-chan int) {}
// 避免接受发送端参数，如有需要也不负责关闭
func avoid(ch chan<- int) {
  // 不负责ch的关闭，也不开启协程处理ch
}
```

```go
// 当不负责关闭管道时避免开启协程，由caller开启
func noGoNoClose(ch chan<- int) {
  // send on ch, no go, no close
}
func noGo(ch <-chan int) {
  // receive on ch, no go
}
```

```go
// 确保caller不消费时不泄露的几种方案
func plan1() <-chan int {
  ch := make(ch int, 1) // buffer
  defer close(ch) // 可以不关闭
  ch <- 0
  return ch
}

func plan2() (chan<- struct{}, <-chan int) {
  done := make(chan struct{})
  data := make(chan int)
  go func() {
    defer close(data)
    select {
      case <- done:
      case ch <- 0:
    }
  }()
  return done, data
}

func plan3(done chan<- struct{}) <-chan int {
  ch := make(ch int)
  go func() {
    defer close(ch)
    select {
      case <- done:
      case ch <- 0:
    }
  }()
  return ch
}

func plan4(ctx context.Context) <-chan int {
  return plan3(ctx.Done())
}
```

## Func

**1. func Func 既不像常量也不像变量**

```go
func Func(){}

// error, not constant
// const f = Func

// error, not variable
// var f = &Func

// error, not constant
// const f = func() {}
```

**2. return v... = 返回值地址赋值 v + 返回**

```go
func fn1() (int, bool) {
  // 返回值在开始执行前已自动初始化为零值
  // 以下通过panic再recovery强制返回
  defer func() {
    recover()
  }()
  panic(nil)
}

func fn2(naked bool) (i int, ok bool) {
  if naked {
    // 具名返回值可赋值再naked return
    i, err = 1, true
    return
  }
  // 也可以赋值和return一步到位
  return 1, true
}
```

**3. defer 只能修改具名返回值的内存空间**

```go
func fn1() (i int) {
  defer func() {
    // 具名返回值把返回值内存空间暴露了出来
    // 所以可在return赋值后继续修改
    i++
  }()
  return 1
}

func fn2() *int {
  i := 1
  defer func() {
    // 虽不能修改返回值内存空间的指针
    // 但可修改指针指向的内容
    i++
  }()
  // 匿名返回值的内存空间
  // 只能return赋值更新
  return &i
}
```

**4. 切片传入可变参数时不会深度复制**

```go
slice := []int{0}
update(slice...)

slice[0] == 1 // true

func update(slice ...int) {
  slice[0] = 1
}
```

## Slice

**1. []T 等价于(ptr \*T, len int, cap int)**

```go
s := make([]int, 1) // len(0) is required
s = append(s, 0) // must re-assign to s
```

**2. string 是特型切片(ptr \*byte, len int)**

```go
s := "Hello world"
b := []byte(s)
s = string(b)
s = s[:5]
```

**3. 小小切片可能导致内存大大占用**

```go
func word(file string) string {
  long := readText(file)
  // return strings.Clone(long[:1])
  return long[:1]
}
```

**4. 为什么\[]T(nil)不等价于[]T{}**

```go
empty := new([0]int)
zero := (*[0]int)(nil)
fmt.Printf("%p, %p\n", empty, zero)
// Output: 0x1165fe0, 0x0
```

**5. 比较两个切片是否相等**

```go
func SliceEqual[T comparable](s1, s2 []T) bool {
  if s1 == nil {
    return s2 == nil
  }
  if len(s1) != len(s2) {
    return false
  }
  for i, v := range s1 {
    if v != s2[i] {
      return false
    }
  }
  return true
}
```

## Interface

**1. 接口值包含真实值及其类型的指针**

```go
fmt.Println(
  unsafe.Sizeof(any(nil)),
  unsafe.Sizeof(uintptr(0)),
)
// Output: 16 8
```

**2. 接口值可以是(T, v), (T, nil)或(nil, nil)**

```go
T_v := any(new(any))
T_nil := any((*any)(nil)) // 非零值，T_nil != nil
nil_nil := any(nil) // 零值，nil_nil == nil
```

**3. 似零非零的 error 返回值(T, nil)很危险**

```go
func returnsError() error {
  var p *MyError = nil
  if bad() {
    p = ErrBad
  }
  return p // always return a non-nil error
}
```

**4. 接口值赋值自动拆箱装箱，T 永不为接口类型**

```go
any1 := any(nil) // box(nil, nil)
any1 = false // box(bool, false)
any2 := any(0) // box(int, 0)
any1 = any2 // unbox(int, 0) -> box(int, 0)
any1 = any(any(0)) // box(int, 0) -> unbox(int, 0) -> box(int, 0)
```

**5. 接口值自动拆箱装箱之迷惑**

```go
raw = 0
box = any(raw) // any(int, 0)
fmt.Printf("%T, %v\n", raw, raw) // box as any(int, 0)
fmt.Printf("%T, %v\n", raw, raw) // no need to box
fmt.Println(any(any(nil)) == error(error(nil)))
// Output:
// int, 0
// int, 0
// true
fmt.Println(reflect.TypeOf(any(nil)) == nil)
fmt.Println(reflect.TypeOf(new(any)).Elem())
// Output:
// true
// interface {}
err := errors.New("")
v1 := any(error(err))
v2 := any(error(nil))
_, v1 = v1.(error)
_, v2 = v2.(error)
fmt.Println(v1, v2)
// Output: true false
```

**6. 接口值与非接口值的转换**

```go
if _, ok := any(nil).(any); ok {
  panic("never, 零值不可以转换成功")
}

// 转换失败时没有 ok 守护直接 panic
_ = any(nil).(any) // panic: ...

if _, ok := error(nil).(interface {
  anyMethod()
}); !ok {
  // 可以尝试转换至任何其他接口类型
}

// error, impossible type assertion
// _ = error(nil).(struct{})

for _, v := range []any{new(error), (*int)(nil), 0, nil} {
  switch v.(type) {
  case *error:
    fmt.Println("*error")
  case *int:
    fmt.Println("*int")
  case int:
    fmt.Println("int")
  case nil:
    fmt.Println("nil")
  default:
    fmt.Println("unknown")
  }
}
// Output:
// *error
// *int
// int
// nil
```

**7. 接口零值及其零值类型打印**

```go
fmt.Printf("%T, %T, %v\n",
  any(nil),
  new(error),
  reflect.TypeOf(new(error)).Elem(),
)
// Output: <nil>, *error, error
```
