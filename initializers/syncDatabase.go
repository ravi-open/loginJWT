package initializers

import "JWTAUTH/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.UserOpen{})
}
