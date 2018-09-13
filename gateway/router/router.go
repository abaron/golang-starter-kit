package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/qeelyn/golang-starter-kit/gateway/app"
	"github.com/qeelyn/golang-starter-kit/gateway/handle"
	"net/http"
	"path"
	"path/filepath"
)

func DefaultRouter() *gin.Engine {
	var (
		router *gin.Engine
		err    error
	)
	if app.IsDebug {
		gin.SetMode(gin.DebugMode)
		router = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
		router.Use(gin.Recovery())
	}

	staticdir := app.Config.GetString("web.staticdir")
	if staticdir, err = filepath.Abs(staticdir); err != nil {
		panic(err)
	}
	router.LoadHTMLGlob(path.Join(staticdir, "html/*"))
	router.Static("/public", staticdir)

	router.GET("/", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	return router
}

func SetupRouterGroup(router *gin.Engine) *gin.RouterGroup {
	v1 := router.Group("/v1")
	if app.AuthHanlerFunc != nil {
		v1.Use(app.AuthHanlerFunc)
	}
	if app.CheckAccessMiddleware != nil {
		v1.Use(app.CheckAccessMiddleware.CheckAccessHandle())
	}
	//TODO add router
	return v1
}

func SetGraphQlRouterGroup(router *gin.Engine) *gin.RouterGroup {
	if app.IsDebug {
		router.GET("/graphiql", func(c *gin.Context) {
			c.HTML(http.StatusOK, "graphiql.html", nil)
		})
	} else {
		handles := make([]gin.HandlerFunc, 0)
		if app.AuthHanlerFunc != nil {
			handles = append(handles, app.AuthHanlerFunc)
		}
		handles = append(handles, func(c *gin.Context) {
			c.HTML(http.StatusOK, "graphiql.html", nil)
		})
		router.GET("/graphiql",
			handles...,
		)
	}

	v2 := router.Group("/v2")
	if app.AuthHanlerFunc != nil {
		v2.Use(app.AuthHanlerFunc)
	}
	handle.ServeGraphqlResource(v2)
	return v2
}
