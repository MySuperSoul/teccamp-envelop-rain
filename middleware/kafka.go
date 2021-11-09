package middleware

import (
	"encoding/json"
	"envelop-rain/configs"
	"envelop-rain/constant"

	db "envelop-rain/repository"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type KafkaProducer struct {
	kafka sarama.SyncProducer
	topic string
}

type KafkaConsumer struct {
	kafka sarama.Consumer
}

func GetKafkaProducer(topic string) *KafkaProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	kafka, err := sarama.NewSyncProducer([]string{configs.GlobalConfig.GetString("Kafka.BrokerAddr")}, config)
	if err != nil {
		logrus.Fatal("Kafka producer failed")
	}
	return &KafkaProducer{kafka: kafka, topic: topic}
}

func (producer *KafkaProducer) SendDBMessage(message []byte) error {
	_, _, err := producer.kafka.SendMessage(&sarama.ProducerMessage{
		Topic: producer.topic,
		Value: sarama.ByteEncoder(message),
	})
	return err
}

func GetKafkaConsumer() *KafkaConsumer {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{configs.GlobalConfig.GetString("Kafka.BrokerAddr")}, config)
	if err != nil {
		logrus.Fatal("Kafka consumer failed")
	}
	return &KafkaConsumer{kafka: consumer}
}

func (consumer *KafkaConsumer) StartConsume(topic_name string, mysql *gorm.DB) {
	partitions, err := consumer.kafka.Partitions(topic_name)
	if err != nil {
		logrus.Fatal("Partitions err: ", err)
	}
	for _, partition_id := range partitions {
		go consumer.ConsumeByPartitionID(partition_id, topic_name, mysql)
	}
}

func (consumer *KafkaConsumer) ConsumeByPartitionID(id int32, topic_name string, mysql *gorm.DB) {
	partitionConsumer, err := consumer.kafka.ConsumePartition(topic_name, id, sarama.OffsetOldest)
	if err != nil {
		logrus.Fatal("ConsumePartition err: ", err)
	}
	for message := range partitionConsumer.Messages() {
		var data map[string]interface{}
		err := json.Unmarshal(message.Value, &data)
		if err != nil {
			logrus.Error("Decode message failed")
		}

		switch data["type"] {
		case constant.CREATE_PACKET_TYPE:
			db.CreatePacket(data, mysql)
		case constant.UPDATE_PACKET_TYPE:
			db.UpdatePacket(data, mysql)
		case constant.UPDATE_REMAIN_TYPE:
			db.UpdateRemainToDB(data, mysql)
		default:
			continue
		}
	}
}
