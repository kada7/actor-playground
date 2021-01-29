module actor-playground

go 1.14

require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20210104230532-90c8f2d201d8
	github.com/artyomturkin/protoactor-go-persistence-boltdb v0.0.0-20180517054913-d2b934b1aaab
	github.com/boltdb/bolt v1.3.1
	github.com/couchbase/gocb v1.6.7
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/atomic v1.7.0
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	gopkg.in/couchbaselabs/gojcbmock.v1 v1.0.4 // indirect
)

replace github.com/artyomturkin/protoactor-go-persistence-boltdb => github.com/kada7/protoactor-go-persistence-boltdb v0.2.0
