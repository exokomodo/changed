package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/exokomodo/changed/pkg/changelog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CreateChangeRequest struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Actor     string    `json:"actor"`
	Service   string    `json:"service"`
	Details   string    `json:"details"`
}

type Handler struct {
	repo changelog.ChangeRepository
}

func NewHandler(repo changelog.ChangeRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ListChanges(c *gin.Context) {
	// TODO: Fix to use timestamp
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	if limit > 1000 {
		limit = 1000
	}

	changes, err := h.repo.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, changes)
}

func (h *Handler) GetChange(c *gin.Context) {
	id := c.Param("id")
	change, err := h.repo.Get(id)
	if err != nil {
		if err == changelog.ErrChangeNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "change not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, change)
}

func (h *Handler) CreateChange(c *gin.Context) {
	var req CreateChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	change := changelog.NewChangeWithTimestamp(req.Actor, req.Service, req.Details, req.Timestamp)
	createdChange, err := h.repo.Create(change)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdChange)
}

// Server to setup the router

func (h *Handler) SetupRouter() *gin.Engine {
	// Make cors work for all origins
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	router.GET("/changes", h.ListChanges)
	router.GET("/changes/:id", h.GetChange)
	router.POST("/changes", h.CreateChange)

	return router
}
func Run(server *gin.Engine, addr string) error {
	fmt.Printf("Starting server on %s\n", addr)
	return server.Run(addr)
}
