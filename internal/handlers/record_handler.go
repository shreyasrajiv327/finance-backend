package handlers

import(
	"finance-backend/internal/models"
	"finance-backend/internal/repository"
	"net/http"
	"github.com/gin-gonic/gin"
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