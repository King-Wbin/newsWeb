package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"encoding/base64"
)

//创建控制器
type UserController struct {
	beego.Controller
}

/*type UserController struct {
	beego.Controller
} */

//展示注册页面
func (this *UserController) ShowReg() {
	this.TplName = "register.html"
}

/*func (this*UserController) ShowReg()  {
	this.TplName = "register.html"
}*/

//处理注册业务
func (this *UserController) HandleReg() {

	//获取注册信息
	userName := this.GetString("userName")
	password := this.GetString("password")
	/*userName := this.GetString("userName")
	password := this.GetString("password")*/

	//校检数据
	if userName == "" || password == "" {
		this.Data["errmsg"] = "用户名或密码不能为空"
		this.TplName = "register.html"
		return
	}
	/*if userName == "" || password == "" {
	    this.Data["errmsg"] = "用户名或密码不能为空"
	    this.TplName = "register.html"
	    return
	  }*/

	//操作数据
	o := orm.NewOrm()
	var cl models.User
	cl.UserName = userName
	cl.Pwd = password
	_, err := o.Insert(&cl)
	if err != nil {
		this.Data["errmsg"] = "注册失败"
		this.TplName = "register.html"
		return
	}
	/*o := orm.NewOrm()
	  var cl models.User
	  cl.UserName = userName
	  cl.Pwd = password
	  _,err := o.Insert(&cl)
	  if err != nil{
	     this.Data["errmsg"] = "注册失败"
	     this.TplName = "register.html"
	     return
	  }*/

	//返回数据
	this.Redirect("/login", 302)
	//this.Redirect("/login",302)
}

//展示登录页面
func (this *UserController) ShowLog() {

	dec := this.Ctx.GetCookie("userName")
	userName, _ := base64.StdEncoding.DecodeString(dec)
	/*dec := this.Ctx.GetCookie("userName")
	  userName,_ ;= base64.StdEncoding.DecodeString(dec)*/

	if string(userName) != "" {
		this.Data["userName"] = string(userName)
		this.Data["checked"] = "checked"
	} else {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}

	/*if string(userName) != "" {
	    this.Data["userName"] = string(userName)
	    this.Data["checked"] = "checked"
	  } else {
	    this.Data["userName"] = ""
	    this.Data["checked"] = ""
	  }*/

	this.TplName = "login.html"
	//this.TplName = "login.html"

}

//处理登录业务
func (this *UserController) HandleLog() {
    //获取数据
	userName := this.GetString("userName")
	password := this.GetString("password")
	/*userName := this.GetString("userName")
	  password := this.GetString("password")*/

	//检查数据
	if userName == "" || password == "" {
		this.Data["errmsg"] = "用户名或者密码不能为空"
		this.TplName = "login.html"
		return
	}
	/*if userName == "" || password == "" {
	    this.Data["errmsg"] = "用户名或者密码不能为空"
	    this.TplNmae = "login.html"
	    return
	  }*/

	//操作数据
	o := orm.NewOrm()
	var cl models.User
	cl.UserName = userName
	/*o := orm.NewOrm()
	  var cl models.User
	  cl.UserName = userName*/

	err := o.Read(&cl, "userName")
	if err != nil {
		this.Data["errmsg"] = "用户名不存在"
		this.TplName = "login.html"
		return
	}
	/*err := o.Read(&cl,"userName")
	  if err != nil {
	    this.Data["errmsg"] = "用户名不存在"
	    this.TplName = "login.html"
	    return
	  }*/

	if password != cl.Pwd {
		this.Data["errmsg"] = "密码错误，请重新输入"
		this.TplName = "login.html"
		return
	}
	/*if password != cl.Pwd {
	    this.Data["errmsg"] = "密码错误，请重新输入"
	    this.TplName = "login.html"
	    return
	  }*/

	remember := this.GetString("remember")
	enc := base64.StdEncoding.EncodeToString([]byte(userName))
	/*remember := this.GetName("remember")
	  enc := base64.StEncodeing.EncodeToString([]byte(userNmae))*/

	if remember == "on" {
		this.Ctx.SetCookie("userName", enc, 3600*1)
	} else {
		this.Ctx.SetCookie("userName", enc, -1)
	}
	/*if remember == "on" {
	    this.Cxt.SetCookie("userName,enc,3600*1")
	  } else {
	    this.Cxt.SetCookie("userName",enc,-1)
	  }*/

	this.SetSession("userName", userName)
	this.Redirect("/article/articleList", 302)
	//this.Ctx.WriteString("登录成功！")
	/*this.SetSession("userName",userNmae)
	  this.Redirect("/article/articleList",302)*/

}

//退出登录
func (this *UserController) LogOut() {
	//删除session
	this.DelSession("userName")
	//this.DelSession("userName")
	//返回页面
	this.Redirect("/login", 302)
	//this.Redirect("/login",302)
}
