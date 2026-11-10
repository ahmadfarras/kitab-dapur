package platform_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kitab-dapur/api/internal/openapi"
	"kitab-dapur/api/internal/platform"
)

var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    platform.Config
		wantErr bool
	}{
		{name: "defaults the address", env: map[string]string{"DATABASE_URL": "postgres://x"}, want: platform.Config{HTTPAddr: ":8080", DatabaseURL: "postgres://x"}},
		{name: "reads the address", env: map[string]string{"DATABASE_URL": "postgres://x", "HTTP_ADDR": ":9000"}, want: platform.Config{HTTPAddr: ":9000", DatabaseURL: "postgres://x"}},
		{name: "requires the database", env: map[string]string{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := platform.LoadConfig(func(k string) string { return tt.env[k] })
			if (err != nil) != tt.wantErr || got != tt.want {
				t.Fatalf("LoadConfig() = %+v, %v", got, err)
			}
		})
	}
}

func TestOpenDBFailsFastOnBadURL(t *testing.T) {
	if _, err := platform.OpenDB(context.Background(), "postgres://nobody@127.0.0.1:1/x?connect_timeout=1"); err == nil {
		t.Fatal("OpenDB() on a closed port succeeded")
	}
}

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

func TestHealth(t *testing.T) {
	ok, _ := platform.NewHealth(fakePinger{}, discard).GetHealth(context.Background(), openapi.GetHealthRequestObject{})
	if _, isOK := ok.(openapi.GetHealth200JSONResponse); !isOK {
		t.Fatalf("healthy db: got %T", ok)
	}
	down, _ := platform.NewHealth(fakePinger{err: errors.New("down")}, discard).GetHealth(context.Background(), openapi.GetHealthRequestObject{})
	if _, isDown := down.(openapi.GetHealth503JSONResponse); !isDown {
		t.Fatalf("unreachable db: got %T", down)
	}
}

// stubAPI fails every operation the way a broken database would.
type stubAPI struct{ openapi.StrictServerInterface }

func (stubAPI) GetHealth(context.Context, openapi.GetHealthRequestObject) (openapi.GetHealthResponseObject, error) {
	return nil, errors.New("pq: secret internal detail")
}

func (stubAPI) CreateKitab(context.Context, openapi.CreateKitabRequestObject) (openapi.CreateKitabResponseObject, error) {
	return openapi.CreateKitab201JSONResponse{}, nil
}

func TestAPIHandlerErrors(t *testing.T) {
	h := platform.NewAPIHandler(stubAPI{}, discard)
	tests := []struct {
		name, method, path, body string
		wantStatus               int
	}{
		{name: "internal errors are hidden", method: "GET", path: "/health", wantStatus: 500},
		{name: "malformed JSON", method: "POST", path: "/kitabs", body: "{", wantStatus: 400},
		{name: "oversized body", method: "POST", path: "/kitabs", body: `{"name":"` + strings.Repeat("a", platform.MaxBodyBytes) + `"}`, wantStatus: 400},
		{name: "bad uuid in path", method: "GET", path: "/kitabs/not-a-uuid/recipes", wantStatus: 400},
		{name: "valid request passes through", method: "POST", path: "/kitabs", body: `{"name":"x"}`, wantStatus: 201},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body)))
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body %s", rec.Code, tt.wantStatus, rec.Body)
			}
			if strings.Contains(rec.Body.String(), "secret") {
				t.Fatalf("internal error leaked: %s", rec.Body)
			}
			if tt.wantStatus >= 400 && rec.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("error is not JSON: %q", rec.Header().Get("Content-Type"))
			}
		})
	}
}

func TestServeShutsDownOnCancel(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- platform.Serve(ctx, addr, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(204) }), discard)
	}()

	var resp *http.Response
	for range 50 {
		if resp, err = http.Get("http://" + addr); err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond) // polling until the listener is up; no ready signal is exposed
	}
	if err != nil || resp.StatusCode != 204 {
		t.Fatalf("server not serving: %v", err)
	}
	resp.Body.Close()

	cancel()
	if err := <-done; err != nil {
		t.Fatalf("Serve() = %v, want clean shutdown", err)
	}
}

func TestServeReportsListenError(t *testing.T) {
	if err := platform.Serve(context.Background(), "bad-address", http.NotFoundHandler(), discard); err == nil {
		t.Fatal("Serve() on a bad address returned nil")
	}
}
