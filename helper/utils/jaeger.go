package utils

import (
	"context"
	"encoding/json"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type JaegerTracer struct {
}

func NewJaegerTracer() *JaegerTracer {
	return &JaegerTracer{}
}

func (t *JaegerTracer) StartSpan(ctx context.Context, tracerName, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	// PERBAIKAN: Tambahkan error handling
	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌ Recovered from panic in StartSpan: %v", r)
		}
	}()

	// // Cek apakah tracer provider tersedia
	// if otel.GetTracerProvider() == nil {
	// 	// Return no-op span jika tracer tidak tersedia
	// 	return ctx, trace.spa
	// }

	ctx, span := otel.Tracer(tracerName).Start(ctx, spanName)

	if len(attrs) > 0 && span != nil {
		span.SetAttributes(attrs...)
	}
	return ctx, span
}

func (t *JaegerTracer) RecordSpanError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}

	// PERBAIKAN: Tambahkan error handling
	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌ Recovered from panic in RecordSpanError: %v", r)
		}
	}()

	// Cek apakah span valid
	if !span.SpanContext().IsValid() {
		return
	}

	span.RecordError(err)
	span.SetAttributes(attribute.String("error.message", err.Error()))
	span.SetAttributes(attribute.String("error", "true"))
	span.SetStatus(codes.Error, err.Error())
}

func (t *JaegerTracer) RecordSpanSuccess(span trace.Span, message string) {
	if span == nil {
		return
	}

	// PERBAIKAN: Tambahkan error handling
	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌ Recovered from panic in RecordSpanSuccess: %v", r)
		}
	}()

	// Cek apakah span valid
	if !span.SpanContext().IsValid() {
		return
	}

	span.SetStatus(codes.Ok, message)
	span.SetAttributes(attribute.String("success.message", message))
}

func (t *JaegerTracer) AddObjectAsAttribute(span trace.Span, key string, obj any) {
	if span == nil {
		return
	}

	// PERBAIKAN: Tambahkan error handling
	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌ Recovered from panic in AddObjectAsAttribute: %v", r)
		}
	}()

	// Cek apakah span valid
	if !span.SpanContext().IsValid() {
		return
	}

	jsonData, err := json.Marshal(obj)
	if err != nil {
		span.SetAttributes(attribute.String("json.marshal.error", err.Error()))
		return
	}
	span.SetAttributes(attribute.String(key, string(jsonData)))
}

func (t *JaegerTracer) EndSpanSafe(span trace.Span) {
	if span == nil {
		return
	}

	// PERBAIKAN: Tambahkan lebih banyak validasi
	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌ Recovered from panic in EndSpanSafe: %v", r)
		}
	}()

	// Cek apakah span context valid
	if !span.SpanContext().IsValid() {
		return
	}

	// Cek apakah span sudah ended (jika ada method untuk check)
	// Beberapa implementasi memiliki cara untuk check ini

	span.End()
}
