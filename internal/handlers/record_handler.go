package handlers

import(
	"finance-backend/internal/models"
	"finance-backend/internal/repository"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RecordHandler struct{
	Repo *repository.RecordRepository
}


func (h *RecordHandler) CreateRecord(c *gin.Context){
	var input models.Record
    
	// parse request
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// basic validation
	if input.Amount <= 0 || (input.Type != "income" && input.Type != "expense") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount or type"})
		return
	}

	// get user_id from middleware
	userIDRaw, _ := c.Get("user_id")
userID := int(userIDRaw.(float64)) 

//record.UserID = userID

	record := models.Record{
		UserID:   userID, // JWT gives float64
		Amount:   input.Amount,
		Type:     input.Type,
		Category: input.Category,
		Date:     input.Date,
		Notes:    input.Notes,
	}

err := h.Repo.Create(&record)
if err != nil {
	println("DB ERROR:", err.Error()) // 🔥 IMPORTANT
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
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
    idParam := c.Param("id")

    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid ID"})
        return
    }

    record, err := h.Repo.GetRecordByID(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "Record not found"})
        return
    }

    // 🔥 Ownership check
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	// Fetch existing record
	record, err := h.Repo.GetRecordByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Record not found"})
		return
	}

	// Ownership check
	if record.UserID != userID {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	// Bind input
	var input models.Record
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Partial update logic
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

func (h *RecordHandler) GetFilteredRecords(c *gin.Context) {
	userID := c.GetInt("user_id")

	// Query params
	recordType := c.Query("type")
	category := c.Query("category")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	records, err := h.Repo.GetFilteredRecords(
		userID,
		recordType,
		category,
		startDate,
		endDate,
		limit,
		offset,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, records)
}