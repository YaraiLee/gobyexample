/**
 * @Author: liyalei
 * @Description: 整个程序，只会执行onces()方法一次,onced()方法是不会被执行的
 * @Version:
 * @Date: 2020/3/30 11:00 上午
 */
package main

import (
	"fmt"
	"sync"
	"time"
)

var once sync.Once

func main() {

	for i, v := range make([]string, 10) {
		once.Do(onces)
		fmt.Println("count:", v, "---", i)
	}
	for i := 0; i < 10; i++ {

		go func() {
			once.Do(onced)
			fmt.Println("213")
		}()
	}
	time.Sleep(4000)
}
func onces() {
	fmt.Println("onces")
}
func onced() {
	fmt.Println("onced")
}
