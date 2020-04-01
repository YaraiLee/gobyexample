/**
 * @Author: liyalei
 * @Description:
 * @Version:
 * @Date: 2020/4/1 2:53 下午
 */
package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

//初始化数据
func datainit(O orm.Ormer) {
	//rel  : 自动生成外键为 表名_id
	sql1 := "insert into user (name,profile_id) values ('ming',1),('hua',2),('qiang',3);"
	sql2 := "insert into profile (age) values (16),(14),(15);"
	sql3 := "insert into tag (name) values ('offical'),('beta'),('dev');"
	sql4 := "insert into post (title,user_id) values ('paper1',1),('paper2',1),('paper3',2),('paper4',3),('paper5',3);"
	// m2m 生成的 表名：子表_主表s  主键自增
	sql5 := "insert into post_tags (tag_id, post_id) values (1,1),(1,3),(2,2),(3,3),(2,4),(3,4),(3,5); "

	//使用Raw（）.Exec（）执行sql
	O.Raw(sql1).Exec()
	O.Raw(sql2).Exec()
	O.Raw(sql3).Exec()
	O.Raw(sql4).Exec()
	O.Raw(sql5).Exec()
}

//one2one查询
func one2one(O orm.Ormer) {
	//one to one :主表Profile 子表User [常用方式：使用级联查询全部数据]
	fmt.Println("one2one-------------------------------")
	//1.通过已知的 子表User数据，查询主表Profile数据
	user := &User{Id: 1}
	O.Read(user) //查询子表
	if user.Profile != nil {
		O.Read(user.Profile)
	}
	fmt.Println("1.二次查询 user now:", user)
	fmt.Println("1.二次查询 profile now:", user.Profile)
	fmt.Println("-------------------------------")
	//2.级联查询
	user = &User{}
	O.QueryTable("user").Filter("Id", 1).RelatedSel().One(user)
	fmt.Println("2.级联查询 user:", user)
	fmt.Println("2.级联查询 profile:", user.Profile)
	fmt.Println("-------------------------------")
	//3.reverse查询 通过子表条件 查询主表 ,此时并没有获取另一个表的数据
	profile := Profile{}
	O.QueryTable("profile").Filter("User__Id", 1).One(&profile)
	fmt.Println("3.reserve 查询 profile:", profile, "条件 user id:1")

	var profiles []*Profile
	O.QueryTable("profile").Filter("User__Name", "ming").One(&profiles)
	for _, a := range profiles {
		fmt.Println("3.reserve 查询 profile:", a, "条件 user name:ming")
	}
	fmt.Println("-------------------------------")
}

//one2many查询
func one2many(O orm.Ormer) {
	//one to many : 主表User 子表Post [常用方式：使用级联查询全部数据]
	fmt.Println("one2many-------------------------------")
	//1.级联查询
	var posts []*Post
	O.QueryTable("post").Filter("User__Id", 1).RelatedSel().All(&posts)
	for _, v := range posts {
		fmt.Println("1.级联查询 post:", v)
		fmt.Println("1.级联查询 post.user.name:", v.User.Name)
	}
	fmt.Println("-------------------------------")
	//2.reverse 查询
	var user User
	err := O.QueryTable("user").Filter("Post__Title", "paper1").Limit(1).One(&user)
	if err == nil {
		fmt.Println("2.reverse 查询 user:", user)
	} else {
		fmt.Println("err:", err)
	}
	fmt.Println("-------------------------------")

}

//many2many查询
func many2many(O orm.Ormer) {
	//many to many : 主表 Tag 子表Post
	//1.reverse 查询
	fmt.Println("many2many-------------------------------")
	var posts []*Post
	O.QueryTable("post").Filter("Tags__Tag__Name", "offical").All(&posts)
	for _, v := range posts {
		fmt.Println("1.reverse 查询 post:", v)
	}
	fmt.Println("-------------------------------")

	//reverse 查询
	var tags []*Tag
	O.QueryTable("tag").Filter("Posts__Post__Title", "paper1").All(&tags)
	for _, x := range tags {
		fmt.Println("2.reverse 查询 tag:", x)
	}
	fmt.Println("-------------------------------")
	//3.级联查询
	//可以创建post_tags表的结构体MapPostTag，包含Post、Tag
	//使用O.QueryTable("post_tags").RelatedSel().All(&MapPostTag)进行查询
}

//one2one one2many插入
func o2o_o2mInsert(O orm.Ormer) {
	//one2one插入
	err := O.Begin() //开启事物
	//插入主表
	profile := Profile{Age: 19}
	id, err := O.Insert(&profile)
	if err != nil {
		O.Rollback()
	} else {
		fmt.Println("success insert profile")
	}
	//插入子表
	user := User{Name: "kakaxi", Profile: &Profile{Id: int(id)}}
	_, err = O.Insert(&user)
	if err != nil {
		O.Rollback()
	} else {
		fmt.Println("success insert user")
	}
	err = O.Commit()
}

//m2m插入
func m2mInsert(O orm.Ormer) error {
	//m2m
	//比如知道一个邮件对象post，一个标记对象tag，插入他们的对应关系
	O.Begin()
	//step1:插入post，获取m2m对象
	user := User{Id: 1}
	O.Read(&user)                                    //读取user
	post := Post{Title: "kakaxi开启写轮眼了", User: &user} //拼接post
	id, err := O.Insert(&post)                       //插入post
	if err != nil {
		O.Rollback()
		return err
	} else {
		fmt.Println("insert post success")
	}
	post = Post{Id: int(id)}
	m2m := O.QueryM2M(&post, "Tags") //创建m2m对象，Tags为关联属性
	//step2:插入tag
	tag := Tag{Name: "娱乐"}
	_, err = O.Insert(&tag)
	if err != nil {
		O.Rollback()
		return err
	} else {
		fmt.Println("insert tag success")
	}
	//step3:插入关系表
	num, err := m2m.Add(&tag)
	if err != nil {
		O.Rollback()
		return err
	} else {
		fmt.Println("insert post_tags success,added nums:", num)
	}
	O.Commit()
	return nil
}
