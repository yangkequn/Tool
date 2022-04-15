package Tools

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

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
