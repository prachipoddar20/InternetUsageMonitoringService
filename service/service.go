package service

import (
	"IUMS/models"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func GetUsageDetails(db *gorm.DB, date time.Time, limit int, page int) (*[]models.UsageDetails, error) {

	var records []record

	dateBefore1Day := date.AddDate(0, 0, -1)
	dateBefore7Days := date.AddDate(0, 0, -7)
	dateBefore30Days := date.AddDate(0, 0, -30)

	//fetching data from database
	rows, err := db.Table("iums_records").Select("username, SUM(CASE WHEN start_time BETWEEN ? AND ? THEN usage_time ELSE 0 END) AS last_day_usage, SUM(CASE WHEN start_time BETWEEN ? AND ? THEN usage_time ELSE 0 END) AS last_7_days_usage, SUM(CASE WHEN start_time BETWEEN ? AND ? THEN usage_time ELSE 0 END) AS last_30_days_usage", dateBefore1Day, date, dateBefore7Days, date, dateBefore30Days, date).Group("username").Order("MAX(start_time)").Limit(limit).Offset((page - 1) * limit).Rows()

	if err != nil {
		return nil, err
	}

	//storing sql rows in a Go object
	if rows != nil {
		for rows.Next() {
			var record record
			if err := rows.Scan(&record.Username, &record.LastDayUsage, &record.Last7DaysUsage, &record.Last30DaysUsage); err != nil {
				return nil, err
			}
			records = append(records, record)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	Data := make([]models.UsageDetails, len(records))

	//storing formatted data
	if len(records) != 0 {
		for i := 0; i < len(records); i++ {
			Data[i].Username = records[i].Username
			Data[i].LastDayUsage = formatDuration(records[i].LastDayUsage)
			Data[i].Last7DaysUsage = formatDuration(records[i].Last7DaysUsage)
			Data[i].Last30DaysUsage = formatDuration(records[i].Last30DaysUsage)
		}
	}

	return &Data, nil
}

func GetUserDetails(db *gorm.DB, username string) (*models.UserDetails, error) {
	var response models.UserDetails

	if !checkUserExists(db, username) {
		return &response, errors.New("user not found")
	}

	//currentTime := time.Date(2022, time.December, 10, 23, 10, 5, 0, time.Local)
	currentTime := time.Now()
	timeBefore1Hour := currentTime.Add(time.Hour * -1)
	timeBefore6Hours := currentTime.Add(time.Hour * -6)
	timeBefore24Hours := currentTime.Add(time.Hour * -24)

	//db call
	rows, err := db.Raw(`SELECT COALESCE(SUM(usage_time), 0) AS total_usage_time, COALESCE(SUM(upload), 0) AS total_upload_size, COALESCE(SUM(download), 0) AS total_download_size FROM iums_records WHERE username = ? AND start_time >= ? UNION ALL SELECT COALESCE(SUM(usage_time), 0) AS total_usage_time, COALESCE(SUM(upload), 0) AS total_upload_size, COALESCE(SUM(download), 0) AS total_download_size FROM iums_records WHERE username = ? AND start_time >= ? UNION ALL SELECT COALESCE(SUM(usage_time), 0) AS total_usage_time, COALESCE(SUM(upload), 0) AS total_upload_size, COALESCE(SUM(download), 0) AS total_download_size FROM iums_records WHERE username = ? AND start_time >= ? `, username, timeBefore1Hour, username, timeBefore6Hours, username, timeBefore24Hours).Rows()

	if err != nil {
		return &response, err
	}

	var userRecords []userRecord

	//storing sql rows in a Go object
	for rows.Next() {
		var userRecord userRecord
		if err := rows.Scan(&userRecord.usageTime, &userRecord.uploadSize, &userRecord.downloadSize); err != nil {
			return &response, err
		}
		userRecords = append(userRecords, userRecord)
	}

	//handling edge case
	if userRecords[2].usageTime == 0 {
		return &response, errors.New("user has not used internet in last 24 hours.")
	}

	var userRecentInternetUsage = make([]models.Usage, 3)
	for i := 0; i < 3; i++ {

		userRecentInternetUsage[i].Time = formatDuration(userRecords[i].usageTime)
		userRecentInternetUsage[i].Upload = formatSize(userRecords[i].uploadSize)
		userRecentInternetUsage[i].Download = formatSize(userRecords[i].downloadSize)
	}
	response.Username = username
	response.LastHourUsage = userRecentInternetUsage[0]
	response.Last6HourUsage = userRecentInternetUsage[1]
	response.Last24HourUsage = userRecentInternetUsage[2]

	return &response, nil
}

// private structures
type record struct {
	Username        string
	LastDayUsage    int
	Last7DaysUsage  int
	Last30DaysUsage int
}

type userRecord struct {
	usageTime    int
	uploadSize   float64
	downloadSize float64
}

// private functions
func formatDuration(totalUsage int) string {
	duration := strconv.Itoa(totalUsage)
	var length int = len(duration)
	if length >= 2 {
		seconds, _ := strconv.Atoi(duration[(length - 2):])

		if length >= 4 {
			minutes, _ := strconv.Atoi(duration[(length - 4):(length - 2)])
			minutes += seconds / 60
			hours, _ := strconv.Atoi(duration[:(length - 4)])
			hours += minutes / 60
			return strconv.Itoa(hours) + "h" + strconv.Itoa(minutes) + "m"
		}
		return strconv.Itoa(seconds) + "s"
	}
	return "0h0m"
}

func formatSize(size float64) string {
	mb := 1e+6
	gb := 1e+9
	tb := 1e+12
	if size >= gb {
		if size >= tb {
			return strconv.FormatFloat((size/tb), 'f', 1, 64) + "TB"
		}
		return strconv.FormatFloat((size/gb), 'f', 1, 64) + "GB"
	}
	return strconv.FormatFloat((size/mb), 'f', 1, 64) + "MB"
}

func checkUserExists(db *gorm.DB, username string) bool {
	var userExists bool
	db.Table("iums_records").
		Select("count(1) > 0").
		Where("iums_records.username = ?", username).
		Find(&userExists)
	return userExists
}
