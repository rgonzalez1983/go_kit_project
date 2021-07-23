package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go_kit_project/internal/static"
)

func KafkaConsumer() *kafka.Consumer {
	fmt.Println(static.MsgReceivingDataFromKafka)
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		static.KeyKafkaBootstrapServer: static.KafkaBootstrapServerValue,
		static.KeyKafkaGroupID:         static.KafkaGroupIDValue,
		static.KeyKafkaAutoOffsetReset: static.KafkaAutoOffsetResetValue,
	})
	if err != nil {
		panic(err)
	}
	return c
}

func KafkaProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{static.KeyKafkaBootstrapServer: static.KafkaBootstrapServerValue})
}

func SaveDataToKafka(data interface{}) {

	fmt.Println(static.MsgSavingDataKafka)

	jsonString, err := json.Marshal(data)

	personString := string(jsonString)
	fmt.Print(personString)

	p, err := KafkaProducer()

	if err != nil {
		panic(err)
	}

	// Produce messages to topic (asynchronously)
	topic := static.KafkaTopicInUse
	for _, word := range []string{personString} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}
	//ReceiveDataFromKafka()
}

func ReceiveDataFromKafka() {
	fmt.Println(static.MsgReceivingDataFromKafka)
	c := KafkaConsumer()
	c.SubscribeTopics([]string{static.KafkaTopicInUse}, nil)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received from Kafka %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}
	defer c.Close()
}
