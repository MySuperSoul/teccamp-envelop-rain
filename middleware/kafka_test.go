package middleware

import (
	"encoding/json"
	"envelop-rain/repository"
	"testing"
	"time"
)

func TestKafkaProducer(t *testing.T) {
	producer := GetKafkaProducer("test_map")
	defer producer.kafka.Close()
	db := repository.GetMySQLCursor()
	defer repository.CloseMySQL(db)

	consumer := GetKafkaConsumer()
	consumer.StartConsume("test_map", db)

	c1 := map[string]interface{}{
		"type":      "CreatePacket",
		"value":     10,
		"uid":       1234,
		"packet_id": 23456,
		"timestamp": 12345,
	}
	c1_data, _ := json.Marshal(c1)
	err := producer.SendDBMessage(c1_data)
	if err != nil {
		t.Fatal(err.Error())
	}

	time.Sleep(10 * time.Second)
}
