/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"

	commonv1pb "github.com/dapr/dapr/pkg/proto/common/v1"
	internalv1pb "github.com/dapr/dapr/pkg/proto/internals/v1"
)

// InvokeMethodResponse holds InternalInvokeResponse protobuf message
// and provides the helpers to manage it.
type InvokeMethodResponse struct {
	r *internalv1pb.InternalInvokeResponse
}

// NewInvokeMethodResponse returns new InvokeMethodResponse object with status.
func NewInvokeMethodResponse(statusCode int32, statusMessage string, statusDetails []*anypb.Any) *InvokeMethodResponse {
	return &InvokeMethodResponse{
		r: &internalv1pb.InternalInvokeResponse{
			Status:  &internalv1pb.Status{Code: statusCode, Message: statusMessage, Details: statusDetails},
			Message: &commonv1pb.InvokeResponse{},
		},
	}
}

// InternalInvokeResponse returns InvokeMethodResponse for InternalInvokeResponse pb to use the helpers.
func InternalInvokeResponse(resp *internalv1pb.InternalInvokeResponse) (*InvokeMethodResponse, error) {
	rsp := &InvokeMethodResponse{r: resp}
	if resp.Message == nil {
		resp.Message = &commonv1pb.InvokeResponse{Data: &anypb.Any{Value: []byte{}}}
	}

	return rsp, nil
}

// WithMessage sets InvokeResponse pb object to Message field.
func (imr *InvokeMethodResponse) WithMessage(pb *commonv1pb.InvokeResponse) *InvokeMethodResponse {
	imr.r.Message = pb
	return imr
}

// WithRawData sets Message using byte data and content type.
func (imr *InvokeMethodResponse) WithRawData(data []byte, contentType string) *InvokeMethodResponse {
	imr.r.Message.ContentType = contentType

	imr.r.Message.Data = &anypb.Any{
		Value: data,
	}

	return imr
}

// WithHeaders sets gRPC response header metadata.
func (imr *InvokeMethodResponse) WithHeaders(headers metadata.MD) *InvokeMethodResponse {
	imr.r.Headers = MetadataToInternalMetadata(headers)
	return imr
}

// WithFastHTTPHeaders populates HTTP response header to gRPC header metadata.
func (imr *InvokeMethodResponse) WithHTTPHeaders(headers map[string][]string) *InvokeMethodResponse {
	imr.r.Headers = MetadataToInternalMetadata(headers)
	return imr
}

// WithFastHTTPHeaders populates fasthttp response header to gRPC header metadata.
func (imr *InvokeMethodResponse) WithFastHTTPHeaders(header *fasthttp.ResponseHeader) *InvokeMethodResponse {
	md := DaprInternalMetadata{}
	header.VisitAll(func(key []byte, value []byte) {
		md[string(key)] = &internalv1pb.ListStringValue{
			Values: []string{string(value)},
		}
	})
	if len(md) > 0 {
		imr.r.Headers = md
	}
	return imr
}

// WithTrailers sets Trailer in internal InvokeMethodResponse.
func (imr *InvokeMethodResponse) WithTrailers(trailer metadata.MD) *InvokeMethodResponse {
	imr.r.Trailers = MetadataToInternalMetadata(trailer)
	return imr
}

// Status gets Response status.
func (imr *InvokeMethodResponse) Status() *internalv1pb.Status {
	if imr.r == nil {
		return nil
	}
	return imr.r.Status
}

// IsHTTPResponse returns true if response status code is http response status.
func (imr *InvokeMethodResponse) IsHTTPResponse() bool {
	if imr.r == nil {
		return false
	}
	// gRPC status code <= 15 - https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
	// HTTP status code >= 100 - https://tools.ietf.org/html/rfc2616#section-10
	return imr.r.Status.Code >= 100
}

// Proto returns the internal InvokeMethodResponse Proto object.
func (imr *InvokeMethodResponse) Proto() *internalv1pb.InternalInvokeResponse {
	return imr.r
}

// Headers gets Headers metadata.
func (imr *InvokeMethodResponse) Headers() DaprInternalMetadata {
	if imr.r == nil {
		return nil
	}
	return imr.r.Headers
}

// Trailers gets Trailers metadata.
func (imr *InvokeMethodResponse) Trailers() DaprInternalMetadata {
	if imr.r == nil {
		return nil
	}
	return imr.r.Trailers
}

// Message returns message field in InvokeMethodResponse.
func (imr *InvokeMethodResponse) Message() *commonv1pb.InvokeResponse {
	if imr.r == nil {
		return nil
	}
	return imr.r.Message
}

// RawData returns content_type and byte array body.
func (imr *InvokeMethodResponse) RawData() (string, []byte) {
	m := imr.Message()
	if m == nil || m.Data == nil {
		return "", nil
	}

	contentType := m.ContentType
	dataTypeURL := m.Data.TypeUrl

	if dataTypeURL != "" {
		contentType = ProtobufContentType
	}

	return contentType, m.Data.Value
}
