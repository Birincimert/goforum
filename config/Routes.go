package config

import (
	admin "goforum/admin/controllers"
	site "goforum/site/controllers"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

func Routes(store *sessions.CookieStore) *httprouter.Router {
	r := httprouter.New()

	// SITE
	r.GET("/", site.Homepage{Store: store}.Index)
	r.GET("/post/:slug", site.Homepage{Store: store}.Detail)
	r.GET("/yazilar/:slug", site.Homepage{Store: store}.Detail)
	r.GET("/about", site.Homepage{Store: store}.About)
	r.GET("/contact", site.Homepage{Store: store}.Contact)
	r.POST("/contact/submit", site.HandleContactForm)

	// Site Auth
	r.GET("/login", site.Userauth{Store: store}.LoginRegisterPage)
	r.POST("/site/login", site.Userauth{Store: store}.DoLogin)
	r.POST("/site/register", site.Userauth{Store: store}.DoRegister)
	r.GET("/site/logout", site.Userauth{Store: store}.Logout)

	//ADMIN
	r.GET("/admin", admin.Dashboard{Store: store}.Index)
	//Blog Posts
	r.GET("/admin/yeni-ekle", admin.Dashboard{Store: store}.NewItem)
	r.POST("/admin/add", admin.Dashboard{Store: store}.Add)
	r.GET("/admin/delete/:id", admin.Dashboard{Store: store}.Delete)
	r.GET("/admin/edit/:id", admin.Dashboard{Store: store}.Edit)
	r.POST("/admin/update/:id", admin.Dashboard{Store: store}.Update)

	//Categories
	r.GET("/admin/kategoriler", admin.Categories{Store: store}.Index)
	r.POST("/admin/kategoriler/add", admin.Categories{Store: store}.Add)
	r.GET("/admin/kategoriler/delete/:id", admin.Categories{Store: store}.Delete)

	//Userops
	r.GET("/admin/login", admin.Userops{Store: store}.Index)
	r.POST("/admin/do_login", admin.Userops{Store: store}.Login)
	r.GET("/admin/logout", admin.Userops{Store: store}.Logout)

	//Contact
	r.GET("/admin/contacts", admin.Contacts{Store: store}.Index)
	r.GET("/admin/contact/delete/:id", admin.Contacts{Store: store}.Delete)

	//About us
	r.GET("/admin/about", admin.About{Store: store}.Index)
	// Hakkında sayfasındaki formu kaydetmek için POST metodu
	r.POST("/admin/about/save", admin.About{Store: store}.Save)

	//COMMENT
	comments := site.SiteComments{Store: store}
	r.POST("/comment/add", comments.CommentAdd)
	r.POST("/comment/upvote/:id", comments.CommentUpvote)
	r.GET("/comment/likes/:id", comments.CommentLikeCount)
	r.GET("/comment/liked/:id", comments.CommentIsLiked)

	// Admin Comment Management Routes
	r.GET("/admin/comments", admin.Comments{Store: store}.Index)
	r.GET("/admin/comments/post/:id", admin.Comments{Store: store}.Post)
	r.GET("/admin/comments/delete/:id", admin.Comments{Store: store}.Delete)

	//SERVE FILES
	r.ServeFiles("/admin/assets/*filepath", http.Dir("admin/assets"))
	r.ServeFiles("/assets/*filepath", http.Dir("site/assets"))
	r.ServeFiles("/uploads/*filepath", http.Dir("uploads"))
	return r
}
