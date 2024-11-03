package seeds

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/Amierza/e-wallet/entity"
	"gorm.io/gorm"
)

func ListUserSeeder(db *gorm.DB) error {
	// membuka file json beranam users.json
	jsonFile, err := os.Open("./migrations/json/users.json")
	if err != nil {
		return err
	}

	// membaca file json sebagai byte
	jsonData, _ := io.ReadAll(jsonFile)

	// mengubah file json menjadi slice dari entity user menggunakan json.unmarshal
	var listUser []entity.User
	if err := json.Unmarshal(jsonData, &listUser); err != nil {
		return nil
	}

	// memastikan tabel user berada di database dan jika tidak ditemukan maka akan membuat tabel baru sesuai dengan entity user
	hashTable := db.Migrator().HasTable(&entity.User{})
	if !hashTable {
		if err := db.Migrator().CreateTable(&entity.User{}); err != nil {
			return err
		}
	}

	// memasukkan data user ke tabel user dan melakukan pengecekan terhadap phone number yang tidak boleh sama
	for _, data := range listUser {
		var user entity.User
		err := db.Where(&entity.User{PhoneNumber: data.PhoneNumber}).First(&user).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		isData := db.Find(&user, "phone_number = ?", data.PhoneNumber).RowsAffected
		if isData == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
