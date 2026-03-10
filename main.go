package main

import (
  "net/http"
  "strings"
  "url-shortner/config"
  "url-shortner/service"

  "github.com/gin-gonic/gin"
)

type POSTURLRequest struct {
  URL string `json:"long_url" binding:"required"`
}

func main() {
	cfg := config.LoadEnv()

  // Create a Gin router with default middleware (logger and recovery)
  r := gin.Default()

  rdb := config.NewRedisClient(cfg)
  urlService := service.NewURLService(rdb.Client, cfg.BaseURL)

  r.POST("/shorten", func(c *gin.Context) {
    var req POSTURLRequest
    err := c.ShouldBindJSON(&req)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
    shortURL, err := urlService.ShortenURL(req.URL)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
  })

  // Define a simple GET endpoint
  r.GET("/ping", func(c *gin.Context) {
    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  // Catch-all route must be last
  r.GET("/get/*shortURL", func(c *gin.Context) {
    shortURLPath := c.Param("shortURL")
    
    // Extract code from "short.url/code" format or direct code
    var code string
    if strings.Contains(shortURLPath, "/") {
      parts := strings.Split(shortURLPath, "/")
      code = parts[len(parts)-1] // Get the last part after the slash
    } else {
      code = shortURLPath // If no slash, use the whole thing
    }
    
    longURL, err := urlService.GetLongURL(c.Request.Context(), code)
    if err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
      return
    }
    c.Redirect(http.StatusMovedPermanently, longURL)
  })

  // Start server on port 8080 (default)
  // Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
  r.Run(":" + cfg.ServerPort)
}