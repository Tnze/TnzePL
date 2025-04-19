# Announcements

作为一门嵌入式语言，要解决的两大问题是，宿主应该给我提供什么样的接口？我又要给宿主提供什么样的接口？

本语言在设计之初就针对第一个问题给出了自己的解决方案：宿主提供的接口需要显式声明。

而第二个问题，目前的答案是：脚本通过宿主提供的接口向宿主提供功能。即宿主不能主动影响脚本的执行流。
这种设计也许能带来一些好处，但是我还没有细想。

## 声明宿主接口的语法

脑袋一拍，关键字就用 `extern` 吧！

```tn
extern {
    print:  fn(v: string)
    assert: fn(v: bool)
    range:  fn(min: int, max: int) -> Iterator<int>
    type Logger {
        debug: fn(self, v: string)
        info: fn(self, v: string)
        error: fn(self, v: string)
    }
}
```
