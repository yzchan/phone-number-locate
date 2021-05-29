淘宝手机归属地
------

优点是免费，并且抓取没有什么限制，开200协程抓取也没遇到问题。最大的缺点就是返回结果中只包含省，不包含归属城市。

### 返回值解析

返回结果为jsonp格式的，而不是更通用的json。所以解析起来稍微麻烦点。但是返回的格式固定，可以尝试直接处理string为我们需要的格式。或者直接使用正则匹配。

最后还可以利用"rogchap.com/v8go"这个库提供的能力直接在golang中执行js语句从而得到有效的json数据。这种方法相对上面两种方法来说效率会低很多。

### 解析器

运行基准测试脚本可以看出3种解析器的性能

```shell
go test -v -run="none" -bench=. -benchmem -benchtime=3s
```

```text
goos: darwin
goarch: amd64
pkg: phone-number-locate/query/tcc.taobao.com
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkStringParser
BenchmarkStringParser-4          3319981              1114 ns/op             688 B/op          5 allocs/op
BenchmarkV8Parser
BenchmarkV8Parser-4               253009             13779 ns/op             784 B/op         18 allocs/op
BenchmarkRegexpParser
BenchmarkRegexpParser-4          1501740              2402 ns/op             256 B/op          2 allocs/op
PASS
ok      phone-number-locate/query/tcc.taobao.com        14.489s

```

> 因为需要处理的字符串很短，所以字符串处理的方法效率最高
> 正则解析次之（但正则解析内存开销更小）。
> V8Go解析效率最差。