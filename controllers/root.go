package controllers

import "github.com/gin-gonic/gin"

// RootController provides an interface for "static" pages.
type RootController struct{}

// HandleGet returns a basic string to let us know that the service is working if we don't specify a key.
func (_ *RootController) HandleGet(c *gin.Context) {
	c.String(200, "Welcome to lnk.now. This is a private URL shortener for perpetualtag.com.")
}
