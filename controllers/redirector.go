package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/matoous/go-nanoid/v2"
	"github.com/skip2/go-qrcode"

	"redirector/helpers"
	"redirector/logging"
	"redirector/models"
)

// RedirectorController is responsible for handling redirection-related operations using key-value storage.
type RedirectorController struct {
	KV models.KV
}

// HandleGet handles GET requests to fetch and process a URL key, providing responses in various formats or performing redirects.
func (r *RedirectorController) HandleGet(c *gin.Context) {
	// Get our Path params
	key := c.Param("key")
	mode := c.Param("mode") // Optional, for non-redirection operations.

	// Try to get the requested key from the DB
	value, err := r.KV.Get([]byte(key))
	if err != nil {
		logging.GetLogger().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// A `nil` value means that it doesn't exist in the DB
	if value == nil {
		c.String(http.StatusNotFound, "not found")
		return
	}

	// Check if we want to do something other than redirect
	switch mode {
	case "/": // Default state, send them out!
		logging.GetLogger().Debugf("redirecting %q to %q", key, string(value))
		c.Redirect(http.StatusTemporaryRedirect, string(value))
		return
	case "/json": // Where does this url actually go? JSON edition!
		c.JSON(http.StatusOK, gin.H{"key": key, "url": string(value)})
		return
	case "/text": // Where does this url actually go? Text edition!
		c.String(http.StatusOK, string(value))
		return
	case "/qr": // Generate a QR code for this URL
		// Figure out the URL for the QR Code.
		// Note: This is not compatible with situations where the redirector is proxied behind a sub-path.
		qrURL := url.URL{
			Scheme: "https",
			Host:   c.Request.Host,
			Path:   key,
		}

		// Get the QR Code configuration based on optional parameters
		qrConfig, err := helpers.GetQRParamsFromContext(c)
		if err != nil {
			logging.GetLogger().Error(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Create a QR Code struct
		qrCode, err := qrcode.New(qrURL.String(), qrConfig.Level)
		if err != nil {
			logging.GetLogger().Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Apply possible configurations to the QR Code
		qrCode.DisableBorder = !qrConfig.Border
		qrCode.BackgroundColor = qrConfig.BgColor
		qrCode.ForegroundColor = qrConfig.FgColor

		// Generate a PNG at the requested size.
		qrData, err := qrCode.PNG(qrConfig.Size)
		if err != nil {
			logging.GetLogger().Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Return the QR Code as a PNG
		c.Data(http.StatusOK, "image/png", qrData)
		return
	default:
		// Not a supported mode :(
		logging.GetLogger().Infof("unknown mode %q for key %q", mode, key)
		c.String(http.StatusNotFound, "not found")
		return
	}
}

// HandlePost processes POST requests to create a redirection entry, using a generated or provided key.
func (r *RedirectorController) HandlePost(c *gin.Context) {
	// Grab the key, if available.
	key := c.Param("key")

	// Bind the Redirect request; This also performs validations against the URL
	var value models.Redirect
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// The key in the path takes precedence, use that if it's provided.
	if key != "" {
		value.Key = key
	}

	// If no key was provided, generate a nanoid to use.
	if value.Key == "" {
		// Note: We don't check if the key exists already since nanoid has enough entropy that a collision is unlikely...
		var err error
		value.Key, err = gonanoid.New(12)
		if err != nil {
			logging.GetLogger().Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	// Make sure that the key is a "sane" length... this is a URL "shortener" after all...
	if len(value.Key) > 100 {
		err := fmt.Errorf("key too long (%d)", len(value.Key))
		logging.GetLogger().Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Attempt to write the key to the DB, This will fail if the key already exists.
	err := r.KV.ExclusivePut([]byte(value.Key), []byte(value.URL))
	if err != nil {
		var ae *helpers.AlreadyExistsError
		if errors.As(err, &ae) {
			c.JSON(http.StatusConflict, gin.H{"error": ae.Error()})
			return
		} else {
			logging.GetLogger().Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	// Return a summary of the new redirect.
	c.JSON(http.StatusOK, gin.H{"status": "success", "redirect": value})
}

// HandlePutWithKey handles PUT requests to update a redirection entry identified by a specified key.
// Replaces the existing URL value and returns both the new and replaced redirect details.
// Responds with a 409 status if the key does not exist or a 500 status for internal server errors.
func (r *RedirectorController) HandlePutWithKey(c *gin.Context) {
	// Grab the key
	key := c.Param("key")

	// Bind the Redirect request; This also performs validations against the URL
	var value, replaced models.Redirect
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value.Key = key
	replaced.Key = key

	// Attempt to replace the existing key
	replacedURL, err := r.KV.Replace([]byte(value.Key), []byte(value.URL))
	if err != nil {
		// If the key doesn't already exist, raise a 409, otherwise report a 500
		var dne *helpers.DoesNotExistError
		if errors.As(err, &dne) {
			c.JSON(http.StatusConflict, gin.H{"error": dne.Error()})
			return
		} else {
			logging.GetLogger().Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	// Populate the value of the URL that was replaced.
	replaced.URL = string(replacedURL)

	// Return a summary of the changes made.
	c.JSON(http.StatusOK, gin.H{"status": "success", "redirect": value, "replaced": replaced})
}
