package kafka

import (
	"net"
	"strconv"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// EnsureTopic создаёт топик в Kafka, если он ещё не существует.
// numPartitions — количество партиций, replicationFactor — фактор репликации.
func EnsureTopic(brokers []string, topic string, numPartitions int, replicationFactor int) error {
	// Подключаемся к первому доступному брокеру и узнаём контроллер кластера
	conn, err := kafkago.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))

	controllerConn, err := kafkago.Dial("tcp", controllerAddr)
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafkago.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil {
		return err
	}

	logger.Log.Info("kafka topic ensured",
		zap.String("topic", topic),
		zap.Int("partitions", numPartitions),
		zap.Int("replication", replicationFactor),
	)

	return nil
}

