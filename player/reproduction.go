package player

import (
	metrics "github.com/massimo-gollo/DASHpher/models"
	"github.com/massimo-gollo/DASHpher/reproduction"
	log "github.com/sirupsen/logrus"
	"load-generator/models"
	"load-generator/utils"
	"sync"
	"time"
)

func Play(counter *utils.Counter, metric *metrics.ReproductionMetrics, nreq uint64, u uint64, list []models.VideoMetadata, wg *sync.WaitGroup, dryMode bool, concurrentGoroutines chan struct{}) {
	defer utils.HandleError()
	defer wg.Done()
	//log.Printf("[Req#%d] Reproducing video n. %d => %s - goroutine number: %d - startTime %s", nreq, u, list[u].Id, runtime.NumGoroutine(), st.Format("2006-01-02 15:04:05"))
	st := time.Now()
	_ = st
	counter.Inc("total")
	counter.Inc("active")
	duration := utils.GetRandomDurationBetween(4, 240)
	if dryMode {
		time.Sleep(time.Second * 10)
	} else {
		_, directUrl := utils.GetVideoUrl(list[u])
		metric.ContentUrl = directUrl
		metric.ReproductionID = nreq
		//log.Infof("randomized duration %ds", duration)
		err := reproduction.Stream(metric, "h264", "conventional", 1080, duration, 2, 5, nreq)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
		}
		switch metric.Status {
		case metrics.Aborted:
			counter.Inc("error")
		case metrics.Error:
			counter.Inc("witherror")
		default:
			counter.Inc("success")
		}
	}
	counter.Dec("active")
	<-concurrentGoroutines
	defer log.Printf("[Req#%d] End video n. %d => %s - endTime %s duration: %s - Video legth: %d", nreq, u, list[u].Id, time.Now().Format("2006-01-02 15:04:05"),
		time.Since(st).String(), time.Duration(duration))
}
