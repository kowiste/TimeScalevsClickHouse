package main

import (
	"fmt"
	"math/rand"
	"testdb/db/clickhouse"
	"testdb/db/timescale"
	"testdb/model"
	"testdb/util"
	"time"
)

type Measure struct {
	DataEntry
	Value int
	Time  time.Time
}

type DataEntry struct {
	ID    string
	Asset string
}

const (
	assetCont     = 500
	measureCont   = 1000
	dataCont      = 1000000
	getAssetsCont = 10
)

var dbs []model.IDatabase = []model.IDatabase{new(timescale.TS), new(clickhouse.CH)}

func main() {

	var db model.IDatabase

	for i := range dbs {
		db = dbs[i]
		name := db.Name()
		start := time.Now()
		err := db.Connect()
		if err != nil {
			fmt.Printf("%s - Error connecting %s\n", name, err)
			continue
		}
		duration := time.Since(start)
		count, _ := db.IsPopulate()
		selectedAsset := make([]string, 0)
		if count != dataCont {
			start = time.Now()
			data, assets := util.CreateData(assetCont, measureCont, dataCont)
			duration = time.Since(start)
			fmt.Printf("Generated %d measures in %s\n", len(data), duration)
			selectedAsset = getRandomAsset(assets, getAssetsCont)
			fmt.Printf("Regenerate %d measures\n", count)
			err = db.Delete()
			if err != nil {
				fmt.Printf("%s - Error deleting %s\n", name, err)
				continue
			}
			start = time.Now()
			err = db.Save(data)
			if err != nil {
				fmt.Printf("%s - Error saving %s\n", name, err)
				continue
			}
			duration := time.Since(start)
			fmt.Printf("%s - Saving %d measures in %s\n", name, len(data), duration)
		}
		if len(selectedAsset) == 0 {
			assets, _ := db.GetAssets()
			selectedAsset = getRandomAsset(assets, getAssetsCont)
		}
		start = time.Now()
		m, err := db.GetByAssets(selectedAsset)
		duration = time.Since(start)
		if err != nil {
			fmt.Printf("%s - Error get by asset %s\n", name, err)
			continue
		}
		fmt.Printf("%s - get %d Measure by Assets in %s \n", name, len(m), duration)
		//time get
		// old, new, err := db.GetOldNewTime()
		// if err != nil {
		// 	fmt.Printf("%s - Error get old new %s\n", name, err)
		// 	continue
		// }
		// startDate, endDate := randomTimeInterval(old, new)
		// start = time.Now()
		// m, err = db.GetByIntervalAndAssets(startDate, endDate, selectedAsset)
		// if err != nil {
		// 	fmt.Printf("%s - Error get by asset %s\n", name, err)
		// 	continue
		// }
		// duration = time.Since(start)
		// fmt.Printf("%s - get %d Measure by Day in %s \n", name, len(m), duration)

	}
}

func getRandomAsset(assets []string, quantity int) (selectedAsset []string) {
	selectedAsset = make([]string, 0)

	for range quantity {
		selectedAsset = append(selectedAsset, assets[rand.Intn(len(assets))])
	}
	return
}

func randomTimeInterval(start, end time.Time) (time.Time, time.Time) {
	duration := end.Sub(start)

	randomDuration := time.Duration(rand.Int63n(int64(duration)))

	randomStart := start.Add(randomDuration)
	remainingDuration := end.Sub(randomStart)
	randomEndDuration := time.Duration(rand.Int63n(int64(remainingDuration)))
	randomEnd := randomStart.Add(randomEndDuration)

	return randomStart, randomEnd
}
