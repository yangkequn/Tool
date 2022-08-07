package Tool

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

func RedisServerStart(c context.Context, rds *redis.Client, serviceNameKey string, callback func(string) (channel string, content string)) (err error) {
	const BATCH_SIZE = 100
	var (
		jobOfLastMinutes int = 0
		// minutes_now is the minutes of current time
		minutes_now int = int(time.Now().Unix() / 60)
		ins         []string
	)

	for true {
		pipe := rds.Pipeline()
		//redis 中的下标前后可达，不像c++,第一个可以到达，最后一个不可到达
		ins, err = pipe.LRange(c, serviceNameKey, 0, BATCH_SIZE-1).Result()
		pipe.LTrim(c, serviceNameKey, BATCH_SIZE, -1)
		if _, err = pipe.Exec(c); err != nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		//calculate the job
		for _, param := range ins {
			channel, content := callback(param)
			if channel == "" || content == "" {
				continue
			}
			pipe.RPush(c, channel, content)
			jobOfLastMinutes++
		}
		pipe.Exec(c)

		//print Job infomation
		if int(time.Now().Unix()/60) != minutes_now {
			minutes_now = int(time.Now().Unix() / 60)
			jobOfLastMinutes = 0
			fmt.Println(time.Now().String(), " redis service ", serviceNameKey, " last minutes job:", jobOfLastMinutes)
		}

	}
	return
}
func JsonToStruct(value *string, result interface{}) error {
	if err := json.Unmarshal([]byte(*value), result); err != nil {
		return err
	}
	return nil
}

//RedisCall: 1.use RPush to push data to redis. 2.use BLPop to pop data from selected channel
//return: error
func RedisDo(c context.Context, rds *redis.Client, ServiceKey string, paramIn map[string]interface{}, paramOut interface{}) (err error) {
	var (
		b       []byte
		BackTo  string = Int64ToString(rand.Int63())
		results []string
	)
	paramIn["BackTo"] = BackTo

	if b, err = msgpack.Marshal(paramIn); err != nil {
		return err
	}
	ppl := rds.Pipeline()
	ppl.RPush(c, ServiceKey, b)
	//长期不执行的任务，抛弃
	ppl.Expire(c, ServiceKey, time.Second*60)
	if _, err := ppl.Exec(c); err != nil {
		return err
	}
	//BLPop 返回结果 [key1,value1,key2,value2]
	if results, err = rds.BLPop(c, time.Second*20, BackTo).Result(); err != nil {
		return err
	}
	return msgpack.Unmarshal([]byte(results[1]), paramOut)
}

var DefaultRedisClient *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "docker.vm:6379", // use default Addr
	Password: "",               // no password set
	DB:       15,               // use default DB
})

func QueryHeartBeat(c context.Context, Data []int64) (heartbeat float32, err error) {
	if err = RedisDo(c, DefaultRedisClient, "heart_beat", map[string]interface{}{"Data": Data}, &heartbeat); err != nil {
		return 0, err
	}
	return heartbeat, nil
}

// TextToMeaning 将文本转换为语义向量
func TextToMeaning(c context.Context, rds *redis.Client, text string) (result []float32, err error) {
	if rds == nil {
		rds = DefaultRedisClient
	}
	if err = RedisDo(c, rds, "text_to_meaning", map[string]interface{}{"Text": text}, &result); err != nil {
		return nil, err
	}
	return result, nil
}
