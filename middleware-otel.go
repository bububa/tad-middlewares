package tadmw

import (
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv10 "go.opentelemetry.io/otel/semconv/v1.10.0"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const instrumName = "github.com/tencentad/marketing-api-go-sdk"

// OtelMiddleware ...
type OtelMiddleware struct {
	traceProvider trace.TracerProvider
	tracer        trace.Tracer //nolint:structcheck
	meterProvider metric.MeterProvider
	meter         metric.Meter
	histogram     metric.Int64Histogram
	counter       metric.Int64Counter
	attrs         []attribute.KeyValue
}

func NewOtelMiddleware(namespace string) *OtelMiddleware {
	if namespace == "" {
		namespace = instrumName
	}
	ret := &OtelMiddleware{
		traceProvider: otel.GetTracerProvider(),
		meterProvider: otel.GetMeterProvider(),
		attrs: []attribute.KeyValue{
			attribute.String("sdk", "tads"),
		},
	}
	ret.tracer = ret.traceProvider.Tracer(namespace)
	ret.meter = ret.meterProvider.Meter(namespace)
	if histogram, err := ret.meter.Int64Histogram(
		semconv.HTTPClientRequestDurationName,
		metric.WithDescription(semconv.HTTPClientRequestDurationDescription),
		metric.WithUnit("milliseconds"),
	); err == nil {
		ret.histogram = histogram
	}
	if counter, err := ret.meter.Int64Counter(
		semconv.HTTPClientActiveRequestsName,
		metric.WithDescription(semconv.HTTPClientActiveRequestsDescription),
		metric.WithUnit(semconv.HTTPClientActiveRequestsUnit),
	); err == nil {
		ret.counter = counter
	}
	return ret
}

// Handle ...
func (o *OtelMiddleware) Handle(
	req *http.Request,
	next func(req *http.Request) (rsp *http.Response, err error),
) (rsp *http.Response, err error) {
	startTime := time.Now()
	attrs := append(o.attrs,
		semconv10.HTTPURLKey.String(req.URL.String()),
		semconv10.HTTPMethodKey.String(req.Method),
		semconv10.HTTPTargetKey.String(req.URL.Path),
		semconv.URLFull(req.URL.String()),
		semconv.HTTPRequestMethodKey.String(req.Method),
		semconv.URLPath(req.URL.Path),
		semconv.URLDomain(req.URL.Host),
		semconv.HTTPRequestSizeKey.Int64(req.ContentLength),
	)
	ctx, span := o.tracer.Start(req.Context(), fmt.Sprintf("http.%s", req.Method),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()
	req = req.WithContext(ctx)
	rsp, err = next(req)
	if o.histogram != nil {
		o.histogram.Record(ctx, time.Since(startTime).Microseconds(), metric.WithAttributes(o.attrs...))
	}
	if o.counter != nil {
		counterAttrs := append(o.attrs, semconv.URLPath(req.URL.Path))
		o.counter.Add(ctx, 1, metric.WithAttributes(counterAttrs...))
	}
	if !span.IsRecording() {
		return rsp, err
	}
	if rsp != nil {
		span.SetAttributes(semconv.HTTPResponseStatusCode(rsp.StatusCode), semconv.HTTPResponseSizeKey.Int64(rsp.ContentLength), semconv.NetworkProtocolName(rsp.Proto))
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return rsp, err
}
