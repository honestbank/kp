generate: mocks

mocks:
	go get github.com/golang/mock/mockgen/model
	go install github.com/golang/mock/mockgen@v1.6.0
	mockgen -destination=./mocks/mock_flow.go -package=mocks github.com/honestbank/kp/kp KPProducer
	mockgen -destination=./mocks/mock_challenge.go -package=mocks github.com/honestbank/kp/kp KPConsumer
	mockgen -destination=./mocks/mock_jwt.go -package=mocks github.com/honestbank/kp/kp KafkaProcessor
