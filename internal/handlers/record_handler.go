package handlers

import (
	"finance-backend/internal/models"
	"finance-backend/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RecordHandler struct {
	Repo *repository.RecordRepository
}


func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var input models.Record

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Amount <= 0 || (input.Type != "income" && input.Type != "expense") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount or type"})
		return
	}

	userID := c.GetInt("user_id")

	if input.Date.IsZero() {
		input.Date = time.Now()
	}

	record := models.Record{
		UserID:   userID,
		Amount:   input.Amount,
		Type:     input.Type,
		Category: input.Category,
		Date:     input.Date,
		Notes:    input.Notes,
	}

	err := h.Repo.Create(&record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}


func (h *RecordHandler) GetRecords(c *gin.Context) {
	userID := c.GetInt("user_id")

	records, err := h.Repo.GetRecordsByUserID(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch records"})
		return
	}

	c.JSON(200, records)
}


func (h *RecordHandler) GetRecordByID(c *gin.Context) {
	userID := c.GetInt("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	record, err := h.Repo.GetRecordByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Record not found"})
		return
	}

	if record.UserID != userID {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, record)
}


func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	userID := c.GetInt("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	record, err := h.Repo.GetRecordByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Not found"})
		return
	}

	if record.UserID != userID {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	err = h.Repo.DeleteRecord(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Delete failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Deleted"})
}


func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	userID := c.GetInt("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	record, err := h.Repo.GetRecordByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Not found"})
		return
	}

	if record.UserID != userID {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	var input models.Record
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if input.Amount != 0 {
		record.Amount = input.Amount
	}
	if input.Type != "" {
		record.Type = input.Type
	}
	if input.Category != "" {
		record.Category = input.Category
	}
	if input.Notes != "" {
		record.Notes = input.Notes
	}
	if !input.Date.IsZero() {
		record.Date = input.Date
	}

	err = h.Repo.UpdateRecord(record)
	if err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(200, record)
}


func (h *RecordHandler) GetSummary(c *gin.Context) {
	userID := c.GetInt("user_id")

	income, expense, err := h.Repo.GetSummary(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch summary"})
		return
	}

	c.JSON(200, gin.H{
		"total_income":  income,
		"total_expense": expense,
		"balance":       income - expense,
	})
}


func (h *RecordHandler) GetCategorySummary(c *gin.Context) {
	userID := c.GetInt("user_id")

	data, err := h.Repo.GetCategorySummary(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch category summary"})
		return
	}

	c.JSON(200, data)
}


func (h *RecordHandler) GetRecentRecords(c *gin.Context) {
	userID := c.GetInt("user_id")

	data, err := h.Repo.GetRecentRecords(userID, 5)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch recent records"})
		return
	}

	c.JSON(200, data)
}


func (h *RecordHandler) GetMonthlySummary(c *gin.Context) {
	userID := c.GetInt("user_id")

	data, err := h.Repo.GetMonthlySummary(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch monthly summary"})
		return
	}

	c.JSON(200, data)
}