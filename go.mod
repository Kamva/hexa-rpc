module github.com/kamva/hexa-rpc

go 1.13

require (
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/kamva/gutil v0.0.0-20210827084201-35b6a3421580
	github.com/kamva/hexa v0.0.0-20210911195556-8149f7a3fc8c
	github.com/kamva/tracer v0.0.0-20201115122932-ea39052d56cd
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.23.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC3 // indirect
	go.opentelemetry.io/otel/sdk v1.0.0-RC3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.26.0
)
