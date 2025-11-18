package telemetry

import (
	"go.mongodb.org/mongo-driver/event"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

// MongoMonitor returns an OpenTelemetry command monitor for the MongoDB driver.
func MongoMonitor() *event.CommandMonitor {
	return otelmongo.NewMonitor()
}
