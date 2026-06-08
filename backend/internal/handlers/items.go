package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/planeta-azul/backend/internal/middleware"
	"github.com/planeta-azul/backend/internal/models"
	"github.com/planeta-azul/backend/internal/repository"
)

type ItemHandler struct {
	store *repository.MemStore
}

func NewItemHandler(store *repository.MemStore) *ItemHandler {
	return &ItemHandler{store: store}
}

// GET /items
func (h *ItemHandler) List(c *gin.Context) {
	userIDStr, _ := c.Get(middleware.ContextUserID)
	roleStr, _ := c.Get(middleware.ContextRole)
	areaIDStr, _ := c.Get(middleware.ContextAreaID)

	userID, _ := uuid.Parse(userIDStr.(string))
	role := models.UserRole(roleStr.(string))

	var areaID *uuid.UUID
	if areaIDStr.(string) != "" {
		id, err := uuid.Parse(areaIDStr.(string))
		if err == nil {
			areaID = &id
		}
	}

	items := h.store.ListItems(areaID, userID, role)
	c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
}

// GET /items/:id
func (h *ItemHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	item, err := h.store.GetItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ítem no encontrado"})
		return
	}

	// Enrich with comments
	item.Comments = h.store.ListCommentsByItem(id)

	c.JSON(http.StatusOK, item)
}

// POST /items
func (h *ItemHandler) Create(c *gin.Context) {
	userIDStr, _ := c.Get(middleware.ContextUserID)
	areaIDStr, _ := c.Get(middleware.ContextAreaID)

	userID, _ := uuid.Parse(userIDStr.(string))
	areaID, _ := uuid.Parse(areaIDStr.(string))

	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Priority == "" {
		req.Priority = models.PriorityMedium
	}

	item := &models.Item{
		ParentID:         req.ParentID,
		Title:            req.Title,
		Description:      req.Description,
		CreatedBy:        userID,
		AssignedTo:       req.AssignedTo,
		AreaID:           areaID,
		Status:           models.StatusDraft,
		Visibility:       req.Visibility,
		RequiresApproval: req.RequiresApproval,
		Priority:         req.Priority,
		Deadline:         req.Deadline,
	}

	if err := h.store.CreateItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear el ítem"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// PATCH /items/:id
func (h *ItemHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.AssignedTo != nil {
		updates["assigned_to"] = *req.AssignedTo
	}

	item, err := h.store.UpdateItem(id, updates)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ítem no encontrado"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DELETE /items/:id  (soft delete)
func (h *ItemHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	userIDStr, _ := c.Get(middleware.ContextUserID)
	userID, _ := uuid.Parse(userIDStr.(string))

	if err := h.store.SoftDeleteItem(id, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ítem no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ítem eliminado"})
}

// POST /items/:id/comments
func (h *ItemHandler) AddComment(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	userIDStr, _ := c.Get(middleware.ContextUserID)
	userID, _ := uuid.Parse(userIDStr.(string))

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &models.Comment{
		ItemID:  itemID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := h.store.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear el comentario"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GET /items/:id/comments
func (h *ItemHandler) ListComments(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	comments := h.store.ListCommentsByItem(itemID)
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

// POST /items/:id/interactions
func (h *ItemHandler) CreateInteraction(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	userIDStr, _ := c.Get(middleware.ContextUserID)
	userID, _ := uuid.Parse(userIDStr.(string))

	var req models.CreateInteractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mark item as waiting
	h.store.UpdateItem(itemID, map[string]interface{}{"status": models.StatusWaiting})

	interaction := &models.Interaction{
		ItemID:      itemID,
		OpenedBy:    userID,
		AddressedTo: req.AddressedTo,
		Message:     req.Message,
	}

	if err := h.store.CreateInteraction(interaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear la interacción"})
		return
	}

	c.JSON(http.StatusCreated, interaction)
}

// GET /items/:id/interactions
func (h *ItemHandler) ListInteractions(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	interactions := h.store.GetInteractionsByItem(itemID)
	c.JSON(http.StatusOK, gin.H{"interactions": interactions})
}
