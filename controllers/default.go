package controllers

import (
	"github.com/astaxie/beego"
	bctx "github.com/astaxie/beego/context"

	"go.opencensus.io/trace"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"

	log "github.com/sirupsen/logrus"

	"math/rand"
	"time"
	"strconv"
	"context"
)

type MainController struct {
	beego.Controller
}

type BenchController struct {
	beego.Controller
}

func MethodLogging(traceId string, spanId string, msg string){
	gcpProject := beego.AppConfig.String("gcpproject")

	log.WithFields(log.Fields{
		"logging.googleapis.com/trace": "projects/" + gcpProject + "/traces/" + traceId,
		"logging.googleapis.com/spanId": spanId,
	}).Info(msg)
}

func HttpLogging(input *bctx.BeegoInput, traceId string, spanId string, msg string) {
	gcpProject := beego.AppConfig.String("gcpproject")

	log.WithFields(log.Fields{
		"UserAgent": input.UserAgent(),
    "RequestURL": input.URL(),
		"RequestURI": input.URI(),
		"RequestMethod": input.Method(),
		"X-Cloud-Trace-Context": input.Header("X-Cloud-Trace-Context"),
		"Proxy": input.Proxy(),
		"Protocol": input.Protocol(),
		"logging.googleapis.com/trace": "projects/" + gcpProject + "/traces/" + traceId,
		"logging.googleapis.com/spanId": spanId,
	}).Info(msg)
}

func fibonacci(ctx context.Context) {
	_, span := trace.StartSpan(ctx, "fibonacci")
	s := span.SpanContext()

	MethodLogging(s.TraceID.String(), s.SpanID.String(), "hello from fibonacci()")

  prev, next := 0, 1
	for i := 0; i < 3000000000; i++ {
       prev, next = next, prev+next
			 if i % 1000000000 == 0 {
				 MethodLogging(s.TraceID.String(), s.SpanID.String(), "fibonacci calcuration processing " + strconv.Itoa(i) + " times")
			 }
  }
	MethodLogging(s.TraceID.String(), s.SpanID.String(), "fibonacci calcuration finished")
	defer span.End()
}

func wait500ms(ctx context.Context){
	_, span := trace.StartSpan(ctx, "wait500ms")
	s := span.SpanContext()

	MethodLogging(s.TraceID.String(), s.SpanID.String(), "hello from wait500ms()")

	time.Sleep(500 * time.Millisecond)
	defer span.End()
}

func (c *BenchController) Get() {
	input := c.Ctx.Input

	if input.Header("X-Cloud-Trace-Context") != "" {
		httpFormat := &propagation.HTTPFormat{}
    parent, _ := httpFormat.SpanContextFromRequest(c.Ctx.Request)
		ctx, span := trace.StartSpanWithRemoteParent(context.Background(), "/bench1", parent)
		s := span.SpanContext()
		HttpLogging(input, s.TraceID.String(), s.SpanID.String(), "Hello from BenchController")
		defer span.End()

		fibonacci(ctx)
		wait500ms(ctx)
		wait500ms(ctx)
	} else{
		ctx, span := trace.StartSpan(context.Background(), "/bench1")
		s := span.SpanContext()
		HttpLogging(input, s.TraceID.String(), s.SpanID.String(), "Hello from BenchController")
		defer span.End()

		fibonacci(ctx)
		wait500ms(ctx)
		wait500ms(ctx)
	}

	rand.Seed(time.Now().UnixNano())
	r := strconv.Itoa(rand.Intn(100))
	c.Data["WebsiteTitle"] = "BenchController"
	c.Data["RandomNumber"] = r
	c.TplName = "index.tpl"
}

func (c *MainController) Get() {
	input := c.Ctx.Input

	if input.Header("X-Cloud-Trace-Context") != "" {
		httpFormat := &propagation.HTTPFormat{}
    parent, _ := httpFormat.SpanContextFromRequest(c.Ctx.Request)
		_, span := trace.StartSpanWithRemoteParent(context.Background(), "/", parent)
		s := span.SpanContext()
		HttpLogging(input, s.TraceID.String(), s.SpanID.String(), "Hello from MainController")
		defer span.End()
	} else{
		_, span := trace.StartSpan(context.Background(), "/")
		s := span.SpanContext()
		HttpLogging(input, s.TraceID.String(), s.SpanID.String(), "Hello from MainController")
		defer span.End()
	}

	rand.Seed(time.Now().UnixNano())
	r := strconv.Itoa(rand.Intn(100))
	c.Data["WebsiteTitle"] = "MainController"
	c.Data["RandomNumber"] = r
	c.TplName = "index.tpl"
}
