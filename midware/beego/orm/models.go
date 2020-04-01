/**
 * @Author: liyalei
 * @Description:
 * @Version:
 * @Date: 2020/3/23 11:32 上午
 */
package main

type User struct {
	Id      int      `orm:"column(id);auto" description:"id"`
	Name    string   `orm:"column(name)" description:"name"`
	Profile *Profile `orm:"rel(one)"` // OneToOne relation
	//Posts   []*Post  `orm:"reverse(many)"`  //one2many情况下，不要写reverse
}

type Profile struct {
	Id   int   `orm:"column(id);auto" description:"id"`
	Age  int16 `orm:"column(age)" description:"age"`
	User *User `orm:"reverse(one)"` // Reverse relationship (optional)
}

type Post struct {
	Id    int    `orm:"column(id);auto" description:"id"`
	Title string `orm:"column(title)" description:"title"`
	User  *User  `orm:"rel(fk)"`  // OneToMany relation
	Tags  []*Tag `orm:"rel(m2m)"` // m2m relation
}

type Tag struct {
	Id    int     `orm:"column(id);auto" description:"id"`
	Name  string  `orm:"column(name)" description:"name"`
	Posts []*Post `orm:"reverse(many)"`
}
