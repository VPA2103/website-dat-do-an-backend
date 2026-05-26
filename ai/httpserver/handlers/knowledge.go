package handlers

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/ai/ingest"
)

type KnowledgeHandler struct {
	fs     FileStore
	vector VectorStore
	llm    Gemini
}

type menuItemIn struct {
	ID          *string  `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Tags        any      `json:"tags"`
	Allergens   any      `json:"allergens"`
	Ingredients any      `json:"ingredients"`
}

type restaurantInfoIn struct {
	Name      *string `json:"name"`
	Address   *string `json:"address"`
	OpenHours *string `json:"open_hours"`
	Phone     *string `json:"phone"`
	Style     *string `json:"style"`
	Policies  *string `json:"policies"`
}

func NewKnowledgeHandler(fs FileStore, vector VectorStore, llm Gemini) *KnowledgeHandler {
	return &KnowledgeHandler{fs: fs, vector: vector, llm: llm}
}

func (h *KnowledgeHandler) ImportMenu(c *gin.Context) {
	rid := RestaurantID(c)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Missing file"})
		return
	}
	if file.Filename == "" || !ingest.HasXLSXExtension(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Only .xlsx is supported"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to read file"})
		return
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to read file"})
		return
	}

	res, err := ingest.IngestMenuXLSX(c.Request.Context(), rid, b, h.fs, h.vector, h.llm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Import failed"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *KnowledgeHandler) AddMenuItem(c *gin.Context) {
	rid := RestaurantID(c)
	var in menuItemIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request"})
		return
	}

	item := ingest.NormalizeMenuItemFromJSONIn(ingest.MenuItemJSONIn{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Tags:        in.Tags,
		Allergens:   in.Allergens,
		Ingredients: in.Ingredients,
	})
	if item.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request"})
		return
	}

	created, err := h.fs.UpsertMenuItem(rid, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to save"})
		return
	}

	if err := ingest.EmbedMenuItem(c.Request.Context(), rid, created, h.vector, h.llm); err != nil {
		log.Printf("embed menu item failed (id=%s name=%s): %v", created.ID, created.Name, err)
	}
	c.JSON(http.StatusOK, created)
}

func (h *KnowledgeHandler) ListMenuItems(c *gin.Context) {
	rid := RestaurantID(c)
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.fs.ListMenuItems(rid, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "offset": offset, "limit": limit, "total": total})
}

func (h *KnowledgeHandler) DeleteMenuItem(c *gin.Context) {
	rid := RestaurantID(c)
	id := c.Param("item_id")
	deleted, err := h.fs.DeleteMenuItem(rid, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed"})
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Not found"})
		return
	}
	_ = h.vector.DeleteMenuItem(c.Request.Context(), rid, id)
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

func (h *KnowledgeHandler) UpdateRestaurant(c *gin.Context) {
	rid := RestaurantID(c)
	var in restaurantInfoIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request"})
		return
	}

	info := RestaurantInfo{
		Name:      in.Name,
		Address:   in.Address,
		OpenHours: in.OpenHours,
		Phone:     in.Phone,
		Style:     in.Style,
		Policies:  in.Policies,
	}

	saved, err := h.fs.SetRestaurant(rid, info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed"})
		return
	}

	if err := ingest.EmbedRestaurant(c.Request.Context(), rid, saved, h.vector, h.llm); err != nil {
		name := ""
		if saved.Name != nil {
			name = *saved.Name
		}
		log.Printf("embed restaurant failed (name=%s): %v", name, err)
	}
	c.JSON(http.StatusOK, saved)
}

func (h *KnowledgeHandler) GetRestaurant(c *gin.Context) {
	rid := RestaurantID(c)
	info, err := h.fs.GetRestaurant(rid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed"})
		return
	}
	c.JSON(http.StatusOK, info)
}
