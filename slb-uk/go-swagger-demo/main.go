//go:generate swag init
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/slb-uk/go-swagger-demo/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Message struct {
    ID      int    `json:"id" example:"1"`
    Message string `json:"message" example:"hello world"`
}

var store = map[int]Message{
    1: {ID: 1, Message: "hello"},
    2: {ID: 2, Message: "namaste"},
}

// @title           Messages API
// @version         1.0
// @description     A simple demo API documented with Swagger 2.0 annotations.
// @termsOfService  https://example.com/terms
// @contact.name   API Support
// @contact.url    https://example.com/support
// @contact.email  support@example.com
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
// @host      localhost:8080
// @BasePath  /v1
func main() {
    r := gin.Default()
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    v1 := r.Group("/v1")
    {
        v1.GET("/hello", helloHandler)
        v1.GET("/messages", listMessages)
        v1.GET("/message/:id", getMessageByID)
        v1.POST("/message", createMessage)
        v1.PUT("/message/:id", updateMessage)
        v1.DELETE("/message/:id", deleteMessage)
    }

    r.Run(":8080")
}

// @Summary      Welcome
// @Description  Returns a friendly welcome message.
// @Tags         misc
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /hello [get]
func helloHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Welcome to Messages API"})
}

// @securityDefinitions.apikey BearerAuth
// @Summary      List messages
// @Description  Get all messages from store.
// @Tags         messages
// @Produce      json
// @Success      200 {array} Message
// @Router       /messages [get]
func listMessages(c *gin.Context) {
    out := make([]Message, 0, len(store))
    for _, m := range store {
        out = append(out, m)
    }
    c.JSON(http.StatusOK, out)
}

// @securityDefinitions.apikey BearerAuth
// @Summary      Get message by ID
// @Description  Get a single message by its ID.
// @Tags         messages
// @Param        id   path      int  true  "Message ID"
// @Produce      json
// @Success      200 {object} Message
// @Failure      404 {object} map[string]string
// @Router       /message/{id} [get]
func getMessageByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    m, ok := store[id]
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, m)
}

// @securityDefinitions.apikey BearerAuth
// @Summary      Create message
// @Description  Create and store a new message.
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        payload body Message true "Message payload (ID optional)"
// @Success      201 {object} Message
// @Failure      400 {object} map[string]string
// @Router       /message [post]
func createMessage(c *gin.Context) {
    var in Message
    if err := c.ShouldBindJSON(&in); err != nil || in.Message == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    next := len(store) + 1
    in.ID = next
    store[next] = in
    c.JSON(http.StatusCreated, in)
}

// @securityDefinitions.apikey BearerAuth
// @Summary      Update message
// @Description  Update message by ID.
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id      path int      true "Message ID"
// @Param        payload body Message  true "Message payload"
// @Success      200 {object} Message
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /message/{id} [put]
func updateMessage(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    _, ok := store[id]
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    var in Message
    if err := c.ShouldBindJSON(&in); err != nil || in.Message == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    in.ID = id
    store[id] = in
    c.JSON(http.StatusOK, in)
}

// @securityDefinitions.apikey BearerAuth
// @Summary      Delete message
// @Description  Delete message by ID.
// @Tags         messages
// @Produce      json
// @Param        id   path int true "Message ID"
// @Success      204 "No Content"
// @Failure      404 {object} map[string]string
// @Router       /message/{id} [delete]
func deleteMessage(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if _, ok := store[id]; !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    delete(store, id)
    c.Status(http.StatusNoContent)
}
