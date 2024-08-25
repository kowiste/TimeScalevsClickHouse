package timescale

import (
	"fmt"
	"testdb/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TS struct {
	db *gorm.DB
}

const name = "timescale"

func (ts *TS) Connect() (err error) {
	dsn := "host=127.0.0.1 user=postgres password=test dbname=measure port=5432 sslmode=disable connect_timeout=10"
	fmt.Println("Connecting to database with DSN:", dsn)
	ts.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return err
	}
	fmt.Println("Successfully connected to database")
	return
}

func (t TS) Name() string {
	return name
}

func (t TS) IsPopulate() (count int64, err error) {
	err = t.db.Model(&model.Measure{}).Count(&count).Error
	return
}

func (t TS) Save(measures []model.Measure) (err error) {
	const batchSize = 10000 // Adjust the batch size as needed

	for i := 0; i < len(measures); i += batchSize {
		end := i + batchSize
		if end > len(measures) {
			end = len(measures)
		}

		batch := measures[i:end]

		err = t.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&batch).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (t TS) Delete() (err error) {
	err = t.db.Exec("DELETE FROM measures").Error
	return
}
func (t TS) GetAssets() (data []string, err error) {
	err = t.db.Model(&model.Measure{}).Distinct().Pluck("asset", &data).Error
	return
}
func (t TS) GetOldNewTime() (old, new time.Time, err error) {
	type Result struct {
		OldestTime time.Time
		NewestTime time.Time
	}

	var result Result
	err = t.db.Raw(`
        SELECT 
            MIN(time) AS oldest_time, 
            MAX(time) AS newest_time 
        FROM 
            measures
    `).Scan(&result).Error

	return
}

func (t TS) GetByAssets(assetIDs []string) (data []model.Measure, err error) {
	err = t.db.Where("asset IN ?", assetIDs).Find(&data).Error
	return
}

func (t TS) GetByIntervalAndAssets(start, end time.Time, assets []string) (data []model.Measure, err error) {
    err = t.db.Where("time BETWEEN ? AND ?", start, end).
        Where("asset IN ?", assets).
        Find(&data).Error
    return
}

func (t TS) GetByIntervalAndAssets2(start, end time.Time, assets []string) (data []model.Measure, err error) {
    err = t.db.Raw(`
        SELECT * 
        FROM measures 
        WHERE time BETWEEN ? AND ? 
        AND asset IN ?
    `, start, end, assets).Scan(&data).Error
    return
}
