package cors

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
	"github.com/kliuchnikovv/engi/internal/types"
	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	header http.Header
}

func (m *mockWriter) Header() http.Header {
	return m.header
}
func (m *mockWriter) Write([]byte) (int, error) { return 0, nil }
func (m *mockWriter) WriteHeader(statusCode int) {
	m.header.Set("Status", http.StatusText(statusCode))
}

func TestCorsAllowedOrigins_Handle(t *testing.T) {
	const (
		corsOriginHeader      = "Origin"
		corsAllowOriginHeader = "Access-Control-Allow-Origin"
		corsOriginMatchAll    = "*"
	)
	var (
		allowedOrigin    = "https://allowed.com"
		notAllowedOrigin = "https://notallowed.com"
	)

	type fields struct {
		origins corsAllowedOrigins
	}
	type args struct {
		origin string
	}
	tests := []struct {
		name          string
		fields        string
		args          string
		wantForbidden bool
		wantHeaderSet bool
	}{
		{
			name:          "origin allowed",
			fields:        allowedOrigin,
			args:          allowedOrigin,
			wantForbidden: false,
			wantHeaderSet: true,
		},
		{
			name:          "origin not allowed",
			fields:        allowedOrigin,
			args:          notAllowedOrigin,
			wantForbidden: true,
			wantHeaderSet: false,
		},
		{
			name:          "origin match all",
			fields:        corsOriginMatchAll,
			args:          notAllowedOrigin,
			wantForbidden: false,
			wantHeaderSet: true,
		},
		{
			name:          "origin not allowed and no match all",
			fields:        "",
			args:          notAllowedOrigin,
			wantForbidden: true,
			wantHeaderSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				req = request.New(&http.Request{
					Header: http.Header{
						corsOriginHeader: []string{tt.args},
					},
				})
				writer = &mockWriter{header: http.Header{}}
				resp   = response.New(
					writer,
					types.NewJSONMarshaler(),
					&types.ResponseAsIs{},
				)
			)

			err := AllowedOrigins(tt.fields).Handle(context.Background(), req, resp)
			assert.NoError(t, err)

			if tt.wantForbidden {
				assert.Equal(t, "Forbidden", resp.ResponseWriter().Header().Get("Status"))
			} else {
				assert.Equal(t, "", resp.ResponseWriter().Header().Get("Status"))
			}
			if tt.wantHeaderSet {
				assert.Equal(t, tt.args, writer.header.Get(corsAllowOriginHeader))
			} else {
				assert.Empty(t, writer.header.Get(corsAllowOriginHeader))
			}
		})
	}
}

func TestCorsAllowedHeaders_Handle(t *testing.T) {
	type fields struct {
		headers corsAllowedHeaders
	}
	type args struct {
		method        string
		requestHeader string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantForbidden     bool
		wantHeaderSet     bool
		expectedHeaderVal string
	}{
		{
			name:   "method not OPTIONS returns nil",
			fields: fields{headers: corsAllowedHeaders{"X-Test-Header"}},
			args: args{
				method:        "GET",
				requestHeader: "X-Test-Header",
			},
			wantForbidden: false,
			wantHeaderSet: false,
		},
		{
			name:   "empty request headers",
			fields: fields{headers: corsAllowedHeaders{"X-Test-Header"}},
			args: args{
				method:        corsOptionMethod,
				requestHeader: "",
			},
			wantForbidden: false,
			wantHeaderSet: false,
		},
		{
			name:   "header in defaultCorsHeaders is skipped",
			fields: fields{headers: corsAllowedHeaders{"X-Test-Header"}},
			args: args{
				method:        corsOptionMethod,
				requestHeader: "Accept",
			},
			wantForbidden: false,
			wantHeaderSet: false,
		},
		{
			name:   "header not allowed triggers forbidden",
			fields: fields{headers: corsAllowedHeaders{"X-Test-Header"}},
			args: args{
				method:        corsOptionMethod,
				requestHeader: "X-Not-Allowed",
			},
			wantForbidden: true,
			wantHeaderSet: false,
		},
		{
			name:   "allowed header is set",
			fields: fields{headers: corsAllowedHeaders{"X-Test-Header"}},
			args: args{
				method:        corsOptionMethod,
				requestHeader: "X-Test-Header",
			},
			wantForbidden:     false,
			wantHeaderSet:     true,
			expectedHeaderVal: "X-Test-Header",
		},
		{
			name:   "multiple allowed headers are set",
			fields: fields{headers: corsAllowedHeaders{"X-Test-Header", "X-Another-Header"}},
			args: args{
				method:        corsOptionMethod,
				requestHeader: "X-Test-Header, X-Another-Header",
			},
			wantForbidden:     false,
			wantHeaderSet:     true,
			expectedHeaderVal: "X-Test-Header,X-Another-Header",
		},
		{
			name:   "empty canonical header is skipped",
			fields: fields{headers: corsAllowedHeaders{"X-Test-Header"}},
			args: args{
				method:        corsOptionMethod,
				requestHeader: "   ,X-Test-Header",
			},
			wantForbidden:     false,
			wantHeaderSet:     true,
			expectedHeaderVal: "X-Test-Header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := request.New(&http.Request{
				Method: tt.args.method,
				Header: http.Header{
					corsRequestHeadersHeader: []string{tt.args.requestHeader},
				},
			})
			writer := &mockWriter{header: http.Header{}}
			resp := response.New(
				writer,
				types.NewJSONMarshaler(),
				&types.ResponseAsIs{},
			)

			err := tt.fields.headers.Handle(context.Background(), req, resp)
			assert.NoError(t, err)

			if tt.wantForbidden {
				assert.Equal(t, "Forbidden", resp.ResponseWriter().Header().Get("Status"))
			} else {
				assert.NotEqual(t, "Forbidden", resp.ResponseWriter().Header().Get("Status"))
			}

			if tt.wantHeaderSet {
				got := writer.header.Get(corsAllowHeadersHeader)
				// Remove spaces for comparison
				got = strings.ReplaceAll(got, " ", "")
				expected := strings.ReplaceAll(tt.expectedHeaderVal, " ", "")
				assert.Equal(t, expected, got)
			} else {
				assert.Empty(t, writer.header.Get(corsAllowHeadersHeader))
			}
		})
	}
}

func TestCorsAllowedMethods_Handle(t *testing.T) {
	tests := []struct {
		name          string
		methods       corsAllowedMethods
		headerValue   string
		headerPresent bool
		wantStatus    string
	}{
		{
			name:          "missing header returns bad request",
			methods:       corsAllowedMethods{"GET", "POST"},
			headerValue:   "",
			headerPresent: false,
			wantStatus:    "Bad Request",
		},
		{
			name:          "method not allowed returns method not allowed",
			methods:       corsAllowedMethods{"GET", "POST"},
			headerValue:   "PUT",
			headerPresent: true,
			wantStatus:    "Method Not Allowed",
		},
		{
			name:          "method allowed returns nil",
			methods:       corsAllowedMethods{"GET", "POST"},
			headerValue:   "POST",
			headerPresent: true,
			wantStatus:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			if tt.headerPresent {
				header.Set(corsRequestMethodHeader, tt.headerValue)
			}
			req := request.New(&http.Request{
				Header: header,
			})
			writer := &mockWriter{header: http.Header{}}
			resp := response.New(
				writer,
				types.NewJSONMarshaler(),
				&types.ResponseAsIs{},
			)

			err := tt.methods.Handle(context.Background(), req, resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, writer.header.Get("Status"))
		})
	}
}
