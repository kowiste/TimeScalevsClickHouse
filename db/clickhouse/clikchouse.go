package clickhouse

import (
	"fmt"
	"testdb/model"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type CH struct {
	db *gorm.DB
}

const name = "clickhouse"

func (ch *CH) Connect() (err error) {
dsn := "tcp://default:test@127.0.0.1:9000/measures"
	fmt.Println("Connecting to database with DSN:", dsn)
	ch.db, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return err
	}
	fmt.Println("Successfully connected to database")
	return
}

func (c CH) Name() string {
	return name
}

func (c CH) IsPopulate() (count int64, err error) {
	err = c.db.Model(&model.Measure{}).Count(&count).Error
	return
}

func (c CH) Save(measures []model.Measure) (err error) {
	const batchSize = 10000 // Adjust the batch size as needed

	for i := 0; i < len(measures); i += batchSize {
		end := i + batchSize
		if end > len(measures) {
			end = len(measures)
		}

		batch := measures[i:end]

		err = c.db.Transaction(func(tx *gorm.DB) error {
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

func (c CH) Delete() (err error) {
	err = c.db.Exec("ALTER TABLE measures DELETE WHERE 1").Error
	return
}

func (c CH) GetAssets() (data []string, err error) {
	err = c.db.Model(&model.Measure{}).Distinct().Pluck("asset", &data).Error
	return
}

func (c CH) GetOldNewTime() (old, new time.Time, err error) {
	type Result struct {
		OldestTime time.Time
		NewestTime time.Time
	}

	var result Result
	err = c.db.Raw(`
        SELECT 
            MIN(time) AS oldest_time, 
            MAX(time) AS newest_time 
        FROM 
            measures
    `).Scan(&result).Error

	return
}

func (c CH) GetByAssets(assetIDs []string) (data []model.Measure, err error) {
	err = c.db.Where("asset IN ?", assetIDs).Find(&data).Error
	return
}

func (c CH) GetByIntervalAndAssets(start, end time.Time, assets []string) (data []model.Measure, err error) {
    err = c.db.Where("time BETWEEN ? AND ?", start, end).
        Where("asset IN ?", assets).
        Find(&data).Error
    return
}

func (c CH) GetByIntervalAndAssets2(start, end time.Time, assets []string) (data []model.Measure, err error) {
    err = c.db.Raw(`
        SELECT * 
        FROM measures 
        WHERE time BETWEEN ? AND ? 
        AND asset IN ?
    `, start, end, assets).Scan(&data).Error
    return
}
