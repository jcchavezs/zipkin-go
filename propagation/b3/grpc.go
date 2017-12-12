package b3

import (
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation"
)

// ExtractGRPC will extract a span.Context from the gRPC Request metadata if
// found in B3 header format.
func ExtractGRPC(md *metadata.MD) propagation.Extractor {
	return func() (*model.SpanContext, error) {
		var (
			traceIDHeader      = GetGRPCHeader(md, TraceID)
			spanIDHeader       = GetGRPCHeader(md, SpanID)
			parentSpanIDHeader = GetGRPCHeader(md, ParentSpanID)
			sampledHeader      = GetGRPCHeader(md, Sampled)
			flagsHeader        = GetGRPCHeader(md, Flags)
		)

		return parseHeaders(
			traceIDHeader, spanIDHeader, parentSpanIDHeader, sampledHeader,
			flagsHeader,
		)
	}
}

// InjectGRPC will inject a span.Context into gRPC metadata.
func InjectGRPC(md *metadata.MD) propagation.Injector {
	return func(sc model.SpanContext) error {
		if (model.SpanContext{}) == sc {
			return ErrEmptyContext
		}

		if sc.Debug {
			setGRPCHeader(md, Flags, "1")
		} else if sc.Sampled != nil {
			// Debug is encoded as X-B3-Flags: 1. Since Debug implies Sampled,
			// we don't send "X-B3-Sampled" if Debug is set.
			if *sc.Sampled {
				setGRPCHeader(md, Sampled, "1")
			} else {
				setGRPCHeader(md, Sampled, "0")
			}
		}

		if !sc.TraceID.Empty() && sc.ID > 0 {
			// set identifiers
			setGRPCHeader(md, TraceID, sc.TraceID.ToHex())
			setGRPCHeader(md, SpanID, fmt.Sprintf("%016x", sc.ID))
			if sc.ParentID != nil {
				setGRPCHeader(md, ParentSpanID, fmt.Sprintf("%016x", *sc.ParentID))
			}
		}

		return nil
	}
}

// GetGRPCHeader retrieves the last value found for a particular key. If key is
// not found it returns an empty string.
func GetGRPCHeader(md *metadata.MD, key string) string {
	v := (*md)[key]
	if len(v) < 1 {
		return ""
	}
	return v[len(v)-1]
}

func setGRPCHeader(md *metadata.MD, key, value string) {
	(*md)[key] = append((*md)[key], value)
}
