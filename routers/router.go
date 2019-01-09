package routers

import (
	"devops-handson/controllers"

	"github.com/astaxie/beego"

	"go.opencensus.io/trace"

	"contrib.go.opencensus.io/exporter/stackdriver"

	"cloud.google.com/go/profiler"

	log "github.com/sirupsen/logrus"

	"fmt"
	"time"
	"os"
)

func init() {
	gcpProject := beego.AppConfig.String("gcpproject")

	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
        log.FieldKeyTime:  "time",
        log.FieldKeyLevel: "severity",
        log.FieldKeyMsg:   "message",
    },
    TimestampFormat: time.RFC3339Nano,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	if err := profiler.Start(profiler.Config{
  	Service:        "devops-demo",
    ServiceVersion: "1.0.0",
    // ProjectID must be set if not running on GCP.
    ProjectID: gcpProject,
  }); err != nil {
		fmt.Println(err)
  }

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: gcpProject,
		MetricPrefix: "devops-",
	})
	if err != nil {
		fmt.Println(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	beego.Router("/bench1", &controllers.BenchController{})
	beego.Router("/", &controllers.MainController{})
}
