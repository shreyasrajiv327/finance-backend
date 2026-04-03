package models

import "time"

type Record struct{
	ID  int          
	UserID int
	Amount    float64   
	Type      string    
	Category  string    
	Date      time.Time 
	Notes     string
	CreatedAt time.Time
}