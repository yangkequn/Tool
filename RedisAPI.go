package Tool

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
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
func RedisDo(c context.Context, rds *redis.Client, ServiceKey string, param map[string]interface{}) (value string, err error) {
	var (
		b       []byte
		results []string
	)
	param["BackTo"] = Int64ToString(rand.Int63())

	if b, err = json.Marshal(param); err != nil {
		return "", err
	}
	ppl := rds.Pipeline()
	ppl.RPush(c, ServiceKey, b)
	//长期不执行的任务，抛弃
	ppl.Expire(c, ServiceKey, time.Second*60)
	if _, err := ppl.Exec(c); err != nil {
		return "", err
	}
	//BLPop 返回结果 [key1,value1,key2,value2]
	if results, err = rds.BLPop(c, time.Second*5, param["BackTo"].(string)).Result(); err != nil {
		return "", err
	}
	return results[1], nil
}

var heartbeatRedisClient *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "docker.vm:6379", // use default Addr
	Password: "",               // no password set
	DB:       15,               // use default DB
})

func QueryHeartBeat(c context.Context, Data []int64) (heartbeat float32, err error) {
	result, err := RedisDo(c, heartbeatRedisClient, "heart_beat", map[string]interface{}{
		"BackTo": Int64ToString(rand.Int63()),
		"Data":   Data,
	})
	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal([]byte(result), &heartbeat); err != nil {
		return 0, err
	}
	return heartbeat, nil
}

// RedisAPI 通过redis调用API的一般框架
func RedisAPI(c context.Context, rds *redis.Client, PubChannel string, SubChannel string, param interface{}, result interface{}) error {

	//https://redis.uptrace.dev/guide/server.html#connecting-to-redis-server
	// return nil if not allowed  subscribing result failed,
	pubsub := rds.Subscribe(c, SubChannel)
	defer pubsub.Unsubscribe(c, SubChannel)
	defer pubsub.Close()

	if _, errRcv := pubsub.Receive(c); errRcv != nil {
		return errRcv
	}

	b, jError := json.Marshal(param)
	if jError != nil {
		return jError
	}
	err := rds.Publish(c, PubChannel, b).Err()
	if err != nil {
		return err
	}
	msg, errRcv := pubsub.ReceiveTimeout(c, time.Second*30)
	var jsonString string = msg.(*redis.Message).Payload
	parseErr := json.Unmarshal([]byte(jsonString), result)
	if errRcv != nil || parseErr != nil {
		return err
	}

	return nil
}

// TextToMeaning 将文本转换为语义向量
func TextToMeaning(c context.Context, rds *redis.Client, text string) (result []float32, err error) {
	type Input struct {
		Channel string //接收消息的channel
		Text    string
	}
	type Output struct {
		Vectors []float32 //接收消息的channel
	}
	type MeaningVector []float32
	param := Input{Channel: Int64ToString(rand.Int63()), Text: text}
	output := &Output{}
	err = RedisAPI(c, rds, "text_to_meaning", param.Channel, param, output)
	if err != nil || len(output.Vectors) == 0 {
		return nil, err
	}
	return output.Vectors, nil
}

func TextToTopics(c context.Context, rds *redis.Client, text string) (result [][]float32, err error) {
	type Input struct {
		Channel string //接收消息的channel
		Text    string
	}
	//array of MeaningVector
	type MeaningVectors struct {
		Vectors [][]float32
	}
	param := Input{Channel: Int64ToString(rand.Int63()), Text: text}
	output := &MeaningVectors{}
	err = RedisAPI(c, rds, "text_to_topics_vectors", param.Channel, param, output)
	return output.Vectors, err
}

func EvalCoverage(c context.Context, rds *redis.Client, diveVectors *[][]float32, result interface{}) error {
	type Input struct {
		Channel string
		Dive    [][]float32
	}
	param := Input{Channel: Int64ToString(rand.Int63()), Dive: *diveVectors}

	return RedisAPI(c, rds, "eval_coverage_quality", param.Channel, param, result)
}
