package main

import (
	"github.com/kataras/iris"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	"github.com/iris-contrib/middleware/csrf"

	"path"
	"regexp"
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

/*
	XORM:
	go get -u github.com/mattn/go-sqlite3
	go get -u github.com/go-xorm/xorm
*/

type Post struct {
	ID        int64 // auto-increment by-default by xorm
	Ip        string
	Name      string
	Email     string
	Subject   string
	Message   string    `json:"content"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

var orm *xorm.Engine
var app *iris.Application

var isFile, _ = regexp.Compile("^[a-zA-Z0-9_.-]*$")

func main() {
	var err error

	app = newApp()
	// XORM
	orm, err = xorm.NewEngine("sqlite3", "./sqlite3.db")
	if err != nil {
		app.Logger().Fatalf("orm failed to initialized: %v", err)
	}
	defer orm.Close()

	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})

	err = orm.Sync2(new(Post))

	if err != nil {
		app.Logger().Fatalf("orm failed to initialized Post table: %v", err)
	}

	// Because of nginx, need X-Forwarded-For, or IP will be localhost
	// Delete the line: iris.WithRemoteAddrHeader if you don't need to change it.
	app.Run(
		iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithRemoteAddrHeader("X-Forwarded-For"),
	)
}
func newApp() *iris.Application {
	// Setup IRIS
	app := iris.New()

	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	// CSRF
	// RequestHeader and CookieName need to be as same as angular's settings.
	protect := csrf.Protect([]byte("9AB0F421E53A477C084477AEA06096F5"),
		csrf.Secure(false),
		csrf.HTTPOnly(false),
		csrf.RequestHeader("X-Csrf-Token"),
		//csrf.CookieName("csrftoken"),
	)

	// Setup SPA
	app.RegisterView(iris.HTML("./dist", ".html"))

	// token to header on every request except files,
	// before the SPA ofc.
	app.Use(protect, func(ctx iris.Context) {
		if ctx.Method() != iris.MethodGet {
			ctx.Next()
			return
		}

		p := ctx.Path()

		isNotFile := len(p) <= 1 || !isFile.MatchString(path.Base(p))
		if isNotFile {
			// get the generate token by "csrf - protect middleware
			// and saved to the request-scoped context's values.
			tok := csrf.Token(ctx)
			// You don't need those, but good to know:
			// ctx.Header("X-CSRF-Token", tok)
			// ctx.Header("Access-Control-Allow-Origin", "*")
			// ctx.Header("Access-Control-Allow-Headers", "*")
			// ctx.Header("Access-Control-Expose-Headers", "*")
			// ctx.Header("Access-Control-Allow-Credentials", "true")
			// ctx.ViewData("csrf", tok)

			// if don't need other template fields
			// use ctx.View("index.html", tok)
			// and change {{.csrf}} to {{.}} in index.html
			ctx.View("index.html", tok)
			return
		}

		ctx.Next()
	})

	assetHandler := app.StaticHandler("./dist", false, false)

	// Use SPA method to handle most of things like routers, files.
	app.SPA(assetHandler)

	app.PartyFunc("/api", func(apiRouter iris.Party) {
		apiRouter.Post("/message", postMessageForm)
	})

	return app
}

func postMessageForm(ctx iris.Context) {
	p := &Post{}
	if err := ctx.ReadJSON(p); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	p.Ip = ctx.RemoteAddr()

	// app.Logger().Debugf("%#v", p)

	orm.Insert(p)

	ctx.JSON(iris.Map{"success": true})
}
