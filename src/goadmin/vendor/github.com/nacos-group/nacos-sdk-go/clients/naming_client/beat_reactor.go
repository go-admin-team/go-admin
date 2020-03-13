package naming_client

import (
	"github.com/nacos-group/nacos-sdk-go/clients/cache"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/utils"
	nsema "github.com/toolkits/concurrent/semaphore"
	"log"
	"strconv"
	"time"
)

type BeatReactor struct {
	beatMap             cache.ConcurrentMap
	serviceProxy        NamingProxy
	clientBeatInterval  int64
	beatThreadCount     int
	beatThreadSemaphore *nsema.Semaphore
	beatRecordMap       cache.ConcurrentMap
}

const Default_Beat_Thread_Num = 20

func NewBeatReactor(serviceProxy NamingProxy, clientBeatInterval int64) BeatReactor {
	br := BeatReactor{}
	if clientBeatInterval <= 0 {
		clientBeatInterval = 5 * 1000
	}
	br.beatMap = cache.NewConcurrentMap()
	br.serviceProxy = serviceProxy
	br.clientBeatInterval = clientBeatInterval
	br.beatThreadCount = Default_Beat_Thread_Num
	br.beatRecordMap = cache.NewConcurrentMap()
	br.beatThreadSemaphore = nsema.NewSemaphore(br.beatThreadCount)
	return br
}

func buildKey(serviceName string, ip string, port uint64) string {
	return serviceName + constant.NAMING_INSTANCE_ID_SPLITTER + ip + constant.NAMING_INSTANCE_ID_SPLITTER + strconv.Itoa(int(port))
}

func (br *BeatReactor) AddBeatInfo(serviceName string, beatInfo model.BeatInfo) {
	log.Printf("[INFO] adding beat: <%s> to beat map.\n", utils.ToJsonString(beatInfo))
	k := buildKey(serviceName, beatInfo.Ip, beatInfo.Port)
	br.beatMap.Set(k, &beatInfo)
	go br.sendInstanceBeat(k, &beatInfo)
}

func (br *BeatReactor) RemoveBeatInfo(serviceName string, ip string, port uint64) {
	log.Printf("[INFO] remove beat: %s@%s:%d from beat map.\n", serviceName, ip, port)
	k := buildKey(serviceName, ip, port)
	data, exist := br.beatMap.Get(k)
	if exist {
		beatInfo := data.(*model.BeatInfo)
		beatInfo.Stopped = true
	}
	br.beatMap.Remove(k)
}

func (br *BeatReactor) sendInstanceBeat(k string, beatInfo *model.BeatInfo) {
	for {
		br.beatThreadSemaphore.Acquire()
		//进行心跳通信
		beatInterval, err := br.serviceProxy.SendBeat(*beatInfo)
		if err != nil {
			log.Printf("[ERROR]:beat to server return error:%s \n", err.Error())
			br.beatThreadSemaphore.Release()
			t := time.NewTimer(beatInfo.Period)
			<-t.C
			continue
		}
		if beatInterval > 0 {
			beatInfo.Period = time.Duration(time.Millisecond.Nanoseconds() * beatInterval)
		}

		//如果当前实例注销，则进行停止心跳
		if beatInfo.Stopped {
			log.Printf("[INFO] intance[%s] stop heartBeating\n", k)
			br.beatThreadSemaphore.Release()
			return
		}

		br.beatRecordMap.Set(k, utils.CurrentMillis())
		br.beatThreadSemaphore.Release()

		t := time.NewTimer(beatInfo.Period)
		<-t.C
	}
}
