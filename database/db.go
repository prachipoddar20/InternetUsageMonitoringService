package database

import (
	"IUMS/models"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	var err error
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting to database : error=%v", err)
		return nil
	}
	return db
}

// structure for db table
type iums_records struct {
	gorm.Model
	Record_id   int       `gorm:"type:int;not null;unique_index"`
	Username    string    `gorm:"type:varchar(40)"`
	Mac_address string    `gorm:"type:varchar(25)"`
	Start_time  time.Time `gorm:"type:datetime"`
	Usage_time  int       `gorm:"type:int"`
	Upload      float64   `gorm:"type:double"`
	Download    float64   `gorm:"type:double"`
}

func InsertDatasetIntoDB(db *gorm.DB) {
	records := []models.UsageRecord{}

	file, err := os.Open("dataset.csv")
	if err != nil {
		log.Fatal(err)
	}

	df := csv.NewReader(file)
	data, _ := df.ReadAll()

	for _, value := range data {
		records = append(records, models.UsageRecord{Username: value[0], MacAddress: value[1], StartTime: value[2], UsageTime: value[3], UploadSize: value[4], DownloadSize: value[5]})
	}

	db.AutoMigrate(&iums_records{})
	fmt.Println("Data migration started...")

	for i := 1; i < len(records); i++ {
		parsedStartTime, _ := time.Parse("2006-01-02 15:04:05", records[i].StartTime)
		upload, _ := strconv.ParseFloat(records[i].UploadSize, 32)
		download, _ := strconv.ParseFloat(records[i].DownloadSize, 32)
		parsedUsageTime, _ := strconv.Atoi(strings.Replace(records[i].UsageTime, ":", "", -1))

		err = db.Exec("insert into iums_records(record_id, username, mac_address, start_time, usage_time, upload, download) values(?,?,?,?,?,?,?)", i, records[i].Username, records[i].MacAddress, parsedStartTime, parsedUsageTime, upload, download).Error

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Data inserted successfully.\n\n")
}
