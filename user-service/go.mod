module github.com/mahendraintelops/home-automation-solution-v2/user-service

go 1.20

require (

    github.com/go-sql-driver/mysql v1.7.1
    gorm.io/driver/mysql v1.5.1
	github.com/mattn/go-sqlite3 v1.14.16
	gorm.io/driver/sqlite v1.5.2
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.2.2
    github.com/sirupsen/logrus v1.9.3
    go.mongodb.org/mongo-driver v1.12.1
    github.com/uptrace/opentelemetry-go-extra/otelgorm v0.2.2
    go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.42.0
    go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.42.0
    go.opentelemetry.io/otel v1.16.0
    go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.16.0
    go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.16.0
    go.opentelemetry.io/otel/sdk v1.16.0
    google.golang.org/grpc v1.57.0
    google.golang.org/protobuf v1.31.0

)

require (

	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect

)
