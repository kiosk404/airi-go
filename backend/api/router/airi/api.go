package airi

import (
	"github.com/gin-gonic/gin"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *gin.Engine) {
	root := r.Group("/", rootMw()...)
	{
		root.GET("/health", Health)
	}
}

func Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
