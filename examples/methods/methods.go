// Go 支持在结构体类型中定义_方法(methods)_ 。

package main

import "fmt"

type rect struct {
    width, height int
    eleArea       int
    elePerimeter  int
}

// 这里的 `area` 方法有一个_接收器(receiver)类型_ `rect`。
func (r *rect) area() int {
    r.eleArea = r.width * r.height
    return r.eleArea
}

// 可以为值类型或者指针类型的接收器定义方法。这里是一个
// 值类型接收器的例子。
func (r rect) perim() int {
    r.elePerimeter = 2 * (r.width + r.height)
    return r.elePerimeter
}

func main() {
    r := rect{width: 10, height: 5}

    // 这里我们调用上面为结构体定义的两个方法。
    fmt.Println("area: ", r.area())
    fmt.Println(r.eleArea)
    fmt.Println("perim:", r.perim())
    fmt.Println(r.elePerimeter)

    // Go 自动处理方法调用时的值和指针之间的转化。你可以使
    // 用指针来调用方法来避免在方法调用时产生一个拷贝，或者
    // 让方法能够改变接受的机构体。
    rp := &r
    fmt.Println("area: ", rp.area())
    fmt.Println(r.eleArea)
    fmt.Println("perim:", rp.perim())
    fmt.Println(r.elePerimeter)
}
