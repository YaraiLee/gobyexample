/**
 * @Author: liyalei
 * @Description:
 * @Version:
 * @Date: 2020/3/23 11:32 上午
 */
package main

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:12345678@/beego_orm?charset=utf8")
}

func main() {
	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库
	//insert
	//profile := new(Profile)
	//profile.Age = 30
	//
	//user := new(User)
	//user.Profile = profile
	//user.Name = "slene"

	//fmt.Println(o.Insert(profile))
	//fmt.Println(o.Insert(user))

	////InsertMulti 并行插入多条
	//users := []User{
	//	{Name: "lyl", Profile: &Profile{Id: 3}},
	//	{Name: "astaxie", Profile: &Profile{Id: 4}},
	//	{Name: "unknown", Profile: &Profile{Id: 5}},
	//}
	//
	//successNums, err := o.InsertMulti(100, users)
	//if nil != err {
	//	panic(err)
	//}
	//fmt.Println(successNums)
	//read
	//user := User{Id: 1}
	//err := o.Read(&user)
	//if err == orm.ErrNoRows {
	//	fmt.Println("查询不到")
	//} else if err == orm.ErrMissPK {
	//	fmt.Println("找不到主键")
	//} else {
	//	fmt.Println(user.Id, user.Name)
	//}

	////ReadOrCreate
	//user := User{Name: "slene"}
	//profile := new(Profile)
	//profile.Age = 10
	//user.Profile = profile
	//
	//// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	//if created, id, err := o.ReadOrCreate(&user, "Name"); err == nil {
	//	if created {
	//		fmt.Println("New Insert an object. Id:", id)
	//	} else {
	//		fmt.Println("Get an object. Id:", id)
	//	}
	//} else {
	//	panic(err)
	//}

	//update
	//user := User{Name: "slene"}
	//if err := o.Read(&user, "Name"); nil == err {
	//	user.Profile.Id = 2
	//	fmt.Println(o.Update(&user))
	//}

	//del
	//if num, err := o.Delete(&User{Id: 1}); err == nil {
	//	fmt.Println(num)
	//}
}
