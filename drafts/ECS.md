# Integrated Entity-Component-System

一门天生适配 ECS 的编程理念的语言是什么样的？

## Component

组件是一种类型，通过关键字`component`定义：

```tn
component Pos(x: f32, y: f32);   // 有多个值的组件，每个值都要有自己的名字
component Health(u8);            // 只有一个值的组件，这个值可以是匿名的
component Player(): Pos, Health; // 组件可以依赖其他组件，
component Enemy(): Pos, Health;  // 当add_comp组件时如果它的依赖组件不存在会报错
```

## Entity

实体可以通过`entity()`创建，新创建的实体不包含任何组件。
随后可以通过`.add_comp<T: component>(...values)`方法添加组件。

```tn
player := entity();
player.add_comp<Pos>(x=0.0, y=0.0);
player.add_comp<Health>(20);
player.add_comp<Player>();
```

## Query

可以通过`query<...T>()`查询所有“包含泛型参数中全部组件的”实体，返回的是一个迭代器。
随后可以通过`.get_comp<T>()`方法获得组件的值。

```tn
for e : query<Pos, Health>() {
    print(e.get_comp<Pos>().x);
    print(e.get_comp<Pos>().y);
    print(e.get_comp<Health>());
}
```

或者把语法设计成这样可以提高性能？

```tn
for Pos(x, y), Health(health) : query() {
    print(x);
    print(y);
    print(health);
}
```
