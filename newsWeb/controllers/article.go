package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"math"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

//展示文章列表页
func (this *ArticleController) ShowArticleList() {

	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	this.Data["userName"] = userName.(string)
	/*userName := this.GetSession("userName")
	  if userName == nil {
	    this.Redirect("/login",302)
	    return
	  }
	this.Data["userName"] = userName.(string)*/

	//获取文章
	o := orm.NewOrm()
	var articles []models.Article
	qs := o.QueryTable("Article")
	/*o := orm.NewOrm()
	  var articles []models.Article
	  qs := o.QueryTable("Article")*/

	//文章页数设定
	count, _ := qs.Count()
	pageSize := int64(2)
	pageCount := float64(count) / float64(pageSize)
	pageCount = math.Ceil(pageCount)
	/*count,_ := qs.Count()
	  pageSize := int64(2)
	  pageCount := float64(count) / float64(pageSize)
	  pageCount = math.Ceil(pageCount)*/

	//把数据传递给视图
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	/*this.Data["count"] = count
	  this.Data["pageCount"] = pageCount*/

	//获取首页末页数据
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	/*pageIndex,err := this.GetInt("pageIndex")
	  if err != nil {
	   pageIdex = 1
	}*/

	//获取分页的数据
	start := pageSize * (int64(pageIndex) - 1 )
	qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	/*start := pageSize * (int64(pageIndex) - 1)
	  qs.Limit(pageSize,start).RelatedSel("ArtcleType").All(&articles)*/

	//根据传递的类型获取相应的文章
	//获取数据
	typeName := this.GetString("select")
	this.Data["typeName"] = typeName
	qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	/*typeName := this.GetString("select")
	  this.Data["typeNmae"] = typeNmae
	  qs.Limit(pageSie,start).RelatedSel("ArticleType").Filter("Article__TypeName",typeName).All(&articles)*/

	this.Data["pageIndex"] = pageIndex
	this.Data["articles"] = articles
	/*this.Data["pageIndex"] = pageIndex
	  this.Data["articles"] = articles*/

	//获取文章类型
	var cl []models.ArticleType
	o.QueryTable("ArticleType").All(&cl)
	this.Data["cl"] = cl
	/*var cl []models.ArticleType
	  o.QueryTable("ArticleType").All(&cl)
	  this.Data["cl"] = cl*/

	this.Layout = "layout.html"
	this.TplName = "index.html"
	/*this.Layout = "layout.html"
	  this.TplName = "index.html"*/

}

//展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	userName := this.GetSession("userName")
	this.Data["userName"] = userName.(string)
	/*userName := this.GetSession("userName")
	  this.Data["userName"] = userName.(string)*/

	//获取文章类型
	o := orm.NewOrm()
	var cl []models.ArticleType
	o.QueryTable("ArticleType").All(&cl)
	/*o := orm.NewOrm()
	  var cl []models.ArticleType
	  o.QueryTable("ArticleType").All(&cl)*/

	this.Data["cl"] = cl
	this.Layout = "layout.html"
	this.TplName = "add.html"
	/*this.Data["cl"] = cl
	  this.Layout = "layout.html"
	  this.TplName = "add.html"*/
}

//处理添加文章业务
func (this *ArticleController) HandleAddArticle() {
	//接受数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	typename := this.GetString("select")
	/*articleName := this.GetString("articleName")
	  content := this.GetString("content")
	  typename := this.GetString("select")*/

	//校验数据
	if articleName == "" || content == "" {
		this.Data["errmsg"] = "文章标题或内容不能为空"
		this.TplName = "add.html"
		return
	}
	/*if articleName == "" || content == "" {
	    this.Data[errmsg] = "文章标题或内容不能为空"
	    this.TplName = "add.html"
	    return
	  }*/

	//接收图片
	file, head, err := this.GetFile("uploadname")
	if err != nil {
		//beego.Info("2")
		this.Data["errmsg"] = "获取文件失败"
		this.TplName = "add.html"
		return
	}
	defer file.Close()
	/*file,head,err := this.GetFile("uploadname")
	  if err != nil {
	     this.Data["errmsg"] = "获取文件失败"
	     this.TplName = "add.html"
	     return
	  }
	  defer file.Close()*/

	//1.判断文件大小
	if head.Size > 500000 {
		this.Data["errmsg"] = "文件太大，上传失败"
		this.TplName = "add.html"
		return
	}
	/*if head.Size > 500000 {
	    this.Data["errmsg"] = "文件太大，上传失败"
	    this.TplName = "add.html"
		return
	  }*/

	//2.判断图片格式
	//1.jpg
	fileExt := path.Ext(head.Filename)
	if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg" {
		this.Data["errmsg"] = "文件格式不正确，请重新上传"
		this.TplName = "add.html"
		return
	}
	/*fileExt := path.Ext(head.Filename)
	  if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg" {
	    this.Data["errmsg"] = "文件格式不正确，请重新上传"
		this.TplName = "add.html"
		return
	  }*/

	//3.文件名防止重复
	fileName := time.Now().Format("2006-01-02-15-04-05") + fileExt
	this.SaveToFile("uploadname", "./static/"+fileName)
	/*fileName := time.Now().Format("2006-01-02-15:04:05") + fileExt
	  this.SaveToFile("uploadname","./static/" + fileName)*/

	//处理数据
	//数据库的插入操作
	//获取orm对象
	o := orm.NewOrm()
	//获取插入对象
	var article models.Article
	//给插入对象赋值
	article.Title = articleName
	article.Content = content
	article.Image = "/static/" + fileName
	/*o := orm.NewOrm()
	  var article models.Article
	  article.Title = articleName
	  article.Content = content
	  article.Image = "/static/" + fileName"*/

	var typenames models.ArticleType
	typenames.TypeName = typename
	o.Read(&typenames, "TypeName")
	article.ArticleType = &typenames
	/*var typenames models.ArticleType
	  typenames.TypeName = typename
	  o.Read(&typenames,"Typename")
	  article.ArticleType = &typenames*/

	//插入
	_, err = o.Insert(&article)
	if err != nil {
		this.Data["errmsg"] = "添加文章失败，请重新添加"
		this.TplName = "add.html"
		return
	}
	/*_,err = o.Insert(&article)
	  if err != nil {
	    this.Data["errmsg"] = "添加文章失败，请重新添加"
	    this.TplName = "add.html"
	    return
	  }*/

	//返回页面
	this.Redirect("/article/articleList", 302)
	//this.Redirect("/article/articleList",302)
}

//展示文章详情页
func (this *ArticleController) ShowArticleDetail() {
	//获取文章Id
	articleid, err := this.GetInt("id")
	if err != nil {
		this.Data["errmsg"] = "请求路径错误"
		this.TplName = "index.html"
		return
	}
	/*articleid,err := this.GetInt("id")
	  if err != nil {
	    this.Data["errmsg"] =  "请求路径错误"
	    this.TplName = "index.html"
		return
	  }*/

	//出路数据
	o := orm.NewOrm()
	var cl models.Article
	//获取条件
	cl.Id = articleid
	//查询数据
	err = o.Read(&cl)
	if err != nil {
		this.Data["errmsg"] = "请求内容失败"
		this.TplName = "index.html"
		return
	}
	/*o := orm.NewOrm()
	  var cl models.Article
	  cl.Id = articleid
	  err := o.Read(&cl)
	  if err != nil {
		this.Data["errmsg"] = "请求内容失败"
		this.TplName = "index.html"
		return
	  }*/

	//获取多对多操作对象
	m2m := o.QueryM2M(&cl, "Users")
	//m2m := o.QueryM2M(&cl,"Users")

	//获取要插入的数据
	var user models.User
	userName := this.GetSession("userName")
	user.UserName = userName.(string)
	o.Read(&user, "UserName")
	this.Data["userName"] = userName.(string)
	/*var user models.User
	  userName := this.GetSession("userName")
	  user.UserName = userName.(string)
	  o.Read(&user,"UserName")
	  this.Data["userName"] = userName.(string)*/

	//插入多对多关系
	m2m.Add(user)
	//m2m.Add(user)

	//第一种多对多查询
	//o.LoadRelated(&cl, "Users")

	////第二种多对多关系查询
	////filter  过滤器  指定查询条件，进行过滤查找
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id", articleid).Distinct().All(&users)
	this.Data["users"] = users
	/*var users []models.User
	  o.QueryTable("User").Filter("Articles__Article__Id").Distinct().All(&users)
	  this.Data["users"] = users*/

	this.Data["article"] = cl
	this.Layout = "layout.html"
	this.TplName = "content.html"
	/*this.Data["article"] = cl
	  this.Layout = "layout.html"
	  this.TplName = "content.html"*/
}

//展示编辑文章页面
func (this *ArticleController) ShowUpdateArticle() {
	userName := this.GetSession("userName")
	this.Data["userName"] = userName.(string)

	articleid, err := this.GetInt("id")
	errmsg := this.GetString("errmsg")
	if errmsg != "" {
		this.Data["errmsg"] = errmsg
	}
	if err != nil {
		beego.Info("路径请求失败")
		this.Redirect("/article/articleList?errmsg", 302)
		return
	}

	o := orm.NewOrm()
	var cl models.Article
	cl.Id = articleid
	err = o.Read(&cl)

	this.Data["article"] = cl
	this.Layout = "layout.html"
	this.TplName = "update.html"

}

//文件上传函数
func UploadFile(this *ArticleController, filePath string) string {
	//接收图片
	file, head, err := this.GetFile(filePath)
	if err != nil {
		this.Data["errmsg"] = "获取文件失败"
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()
	//1.判断文件大小
	if head.Size > 500000 {
		this.Data["errmsg"] = "文件太大，上传失败"
		this.TplName = "add.html"
		return ""
	}

	//2.判断图片格式
	//1.jpg
	fileExt := path.Ext(head.Filename)
	if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg" {
		this.Data["errmsg"] = "文件格式不正确，请重新上传"
		this.TplName = "add.html"
		return ""
	}

	//3.文件名防止重复
	fileName := time.Now().Format("2006-01-02-15-04-05") + fileExt
	this.SaveToFile(filePath, "./static/"+fileName)
	return "/static/" + fileName
}

//处理编辑文章业务
func (this *ArticleController) HandleUpdateArticle() {
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	fileName := UploadFile(this, "uploadname")
	articleId, err2 := this.GetInt("id")

	if articleName == "" || content == "" || fileName == "" || err2 != nil {
		errmsg := "内容不能为空"
		this.Redirect("/article/UpdateArticle?id="+strconv.Itoa(articleId)+"&errmsg="+errmsg, 302)
		return
	}

	//更新操作
	o := orm.NewOrm()
	var cl models.Article
	cl.Id = articleId
	err := o.Read(&cl)
	if err != nil {
		errmsg := "更新文章不存在"
		this.Redirect("/article/UpdateArticle?id="+strconv.Itoa(articleId)+"&errmsg="+errmsg, 302)
		return

	}

	cl.Title = articleName
	cl.Content = content
	cl.Image = fileName

	o.Update(&cl)

	this.Redirect("/article/articleList", 302)
}

//删除业务处理
func (this *ArticleController) DeleteArticle() {
	articleId, err := this.GetInt("id")
	if err != nil {
		beego.Info("路径读取失败")
		this.Redirect("/article/articleList?errmsg", 302)
		return
	}

	o := orm.NewOrm()
	var cl models.Article
	cl.Id = articleId
	_, err = o.Delete(&cl)
	if err != nil {
		beego.Info("删除失败")
		this.Redirect("/article/articleList?errmsg", 302)
		return
	}

	this.Redirect("/article/articleList", 302)
}

//展示添加类型界面
func (this *ArticleController) ShowAddType() {
	userName := this.GetSession("userName")
	this.Data["userName"] = userName.(string)

	o := orm.NewOrm()
	var cl []models.ArticleType
	o.QueryTable("ArticleType").All(&cl)
	this.Data["cl"] = cl

	//this.Redirect("/addType", 302)
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

//处理类型添加业务
func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		beego.Info("文章类型不能为空")
		this.Redirect("/article/addType", 302)
		return
	}

	o := orm.NewOrm()
	var cl models.ArticleType
	cl.TypeName = typeName
	_, err := o.Insert(&cl)
	if err != nil {
		beego.Info("数据插入失败")
		this.Redirect("/article/addType", 302)
		return
	}

	this.Redirect("/article/addType", 302)
}

//删除类型
func (this *ArticleController) DeleteType() {
	TypeId, err := this.GetInt("id")
	if err != nil {
		beego.Error("删除失败")
		this.Redirect("/article/addType", 302)
		return
	}

	o := orm.NewOrm()
	var cl models.ArticleType
	cl.Id = TypeId
	_, err = o.Delete(&cl)
	if err != nil {
		beego.Error("删除失败")
		this.Redirect("/article/addType", 302)
		return
	}

	this.Redirect("/article/addType", 302)
}
