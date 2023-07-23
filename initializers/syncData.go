package initializers

import "golang/jwt/models"

func SyncData() {
	DB.AutoMigrate(&models.User{})

}
