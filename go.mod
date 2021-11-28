module github.com/kamva/hexa-rpc

go 1.13

require (
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/kamva/gutil v0.0.0-20210827084201-35b6a3421580
	github.com/kamva/hexa v0.0.0-20211128175703-59125a2fe5ec
	github.com/kamva/tracer v0.0.0-20201115122932-ea39052d56cd
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.27.0
	go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/exporters/prometheus v0.25.0 // indirect
	go.opentelemetry.io/otel/metric v0.25.0
	go.opentelemetry.io/otel/sdk v1.2.0 // indirect
	go.opentelemetry.io/otel/sdk/export/metric v0.25.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.25.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.26.0
)
