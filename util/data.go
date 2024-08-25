package util

import (
	"math/rand"
	"sync"
	"testdb/model"
	"time"

	"github.com/google/uuid"
)

func CreateData(assetCont, measureCont, dataCont int) (data []model.Measure, assets []string) {
	assets = make([]string, assetCont)
	measures := make(map[int]model.DataEntry)

	for i := range assets {
		assets[i] = uuid.NewString()
	}
	for i := 0; i < measureCont; i++ {

		measures[i] = model.DataEntry{
			ID:    uuid.NewString(),
			Asset: assets[rand.Intn(len(assets))],
		}
	}

	data = make([]model.Measure, dataCont)
	uniqueMeasures := make(map[string]struct{})
	var mu sync.Mutex
	var wg sync.WaitGroup
	numWorkers := 10
	chunkSize := dataCont / numWorkers

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for i := start; i < start+chunkSize; i++ {
				var measure model.Measure
				for {
					randomMeasure := measures[rand.Intn(measureCont)]
					measure = model.Measure{
						ID:    randomMeasure.ID,
						Asset: randomMeasure.Asset,
						Value: rand.Intn(1001),
						Time:  randomTimeInLastMonth(),
					}
					key := measure.ID + measure.Time.String()
					mu.Lock()
					if _, exists := uniqueMeasures[key]; !exists {
						uniqueMeasures[key] = struct{}{}
						mu.Unlock()
						break
					}
					mu.Unlock()
				}
				data[i] = measure
			}
		}(w * chunkSize)
	}

	wg.Wait()
	return
}

func randomTimeInLastMonth() time.Time {
	now := time.Now()
	oneMonthAgo := now.AddDate(0, -1, 0)
	randomTime := rand.Int63n(now.Unix()-oneMonthAgo.Unix()) + oneMonthAgo.Unix()
	return time.Unix(randomTime, 0)
}
