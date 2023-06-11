generate: interface_mocks
	@echo "generated"

interface_mocks:
	go get github.com/golang/mock/mockgen/model
	go install github.com/golang/mock/mockgen@v1.6.0
	mockgen -destination=./mocks/mock_kpprocessor.go -package=mocks github.com/honestbank/kp KPProducer
	mockgen -destination=./mocks/mock_kafkaconsumer.go -package=mocks github.com/honestbank/kp KPConsumer
	mockgen -destination=./mocks/mock_kafkaprocessor.go -package=mocks github.com/honestbank/kp KafkaProcessor
