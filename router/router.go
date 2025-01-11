package router

import (
    "github.com/gin-gonic/gin"
    "github.com/qich3n/crypto-sentiment-analyzer/controllers"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Serve static files (CSS, JS) from the "./static" folder
    r.Static("/static", "./static")

    // Load HTML templates from the "./templates" folder
    r.LoadHTMLGlob("templates/*")

    // Routes
    r.GET("/", controllers.RenderIndex)
    r.GET("/api/analyze", controllers.GetSentimentAnalysis)

    return r
}
