package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserID retrieves the user ID from the gin context
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, gin.Error{}
	}
	return uuid.Parse(userID.(string))
}
