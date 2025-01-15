package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"redirector/logging"
	"redirector/models"
)

// HealthController is responsible for handling health-related API requests.
type HealthController struct {
	KV models.KV
}

// HandleGet processes a GET request to check the health status and responds with status "ok" or an error message.
func (h *HealthController) HandleGet(c *gin.Context) {
	// Get the stored "health" value from the DB, this should always be "ok"
	val, err := h.KV.Get([]byte("health"))
	if err != nil {
		// Something went wrong :( return a sane error.
		logging.GetLogger().Error(err)
		c.JSON(500, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	} else if string(val) != "ok" {
		// Did someone modify the "health" value outside the app?
		err = fmt.Errorf("expected \"ok\", got: %q", val)
		logging.GetLogger().Error(err)
		c.JSON(500, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// All is well :)
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
