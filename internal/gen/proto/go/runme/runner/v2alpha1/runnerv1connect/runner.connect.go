// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: runme/runner/v2alpha1/runner.proto

package runnerv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/stateful/runme/internal/gen/proto/go/runme/runner/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// RunnerServiceName is the fully-qualified name of the RunnerService service.
	RunnerServiceName = "runme.runner.v2alpha1.RunnerService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// RunnerServiceCreateSessionProcedure is the fully-qualified name of the RunnerService's
	// CreateSession RPC.
	RunnerServiceCreateSessionProcedure = "/runme.runner.v2alpha1.RunnerService/CreateSession"
	// RunnerServiceGetSessionProcedure is the fully-qualified name of the RunnerService's GetSession
	// RPC.
	RunnerServiceGetSessionProcedure = "/runme.runner.v2alpha1.RunnerService/GetSession"
	// RunnerServiceListSessionsProcedure is the fully-qualified name of the RunnerService's
	// ListSessions RPC.
	RunnerServiceListSessionsProcedure = "/runme.runner.v2alpha1.RunnerService/ListSessions"
	// RunnerServiceDeleteSessionProcedure is the fully-qualified name of the RunnerService's
	// DeleteSession RPC.
	RunnerServiceDeleteSessionProcedure = "/runme.runner.v2alpha1.RunnerService/DeleteSession"
	// RunnerServiceExecuteProcedure is the fully-qualified name of the RunnerService's Execute RPC.
	RunnerServiceExecuteProcedure = "/runme.runner.v2alpha1.RunnerService/Execute"
)

// RunnerServiceClient is a client for the runme.runner.v2alpha1.RunnerService service.
type RunnerServiceClient interface {
	CreateSession(context.Context, *connect_go.Request[v1.CreateSessionRequest]) (*connect_go.Response[v1.CreateSessionResponse], error)
	GetSession(context.Context, *connect_go.Request[v1.GetSessionRequest]) (*connect_go.Response[v1.GetSessionResponse], error)
	ListSessions(context.Context, *connect_go.Request[v1.ListSessionsRequest]) (*connect_go.Response[v1.ListSessionsResponse], error)
	DeleteSession(context.Context, *connect_go.Request[v1.DeleteSessionRequest]) (*connect_go.Response[v1.DeleteSessionResponse], error)
	// Execute executes a program. Examine "ExecuteRequest" to explore
	// configuration options.
	//
	// It's a bidirectional stream RPC method. It expects the first
	// "ExecuteRequest" to contain details of a program to execute.
	// Subsequent "ExecuteRequest" should only contain "input_data" as
	// other fields will be ignored.
	Execute(context.Context) *connect_go.BidiStreamForClient[v1.ExecuteRequest, v1.ExecuteResponse]
}

// NewRunnerServiceClient constructs a client for the runme.runner.v2alpha1.RunnerService service.
// By default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped
// responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewRunnerServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) RunnerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &runnerServiceClient{
		createSession: connect_go.NewClient[v1.CreateSessionRequest, v1.CreateSessionResponse](
			httpClient,
			baseURL+RunnerServiceCreateSessionProcedure,
			opts...,
		),
		getSession: connect_go.NewClient[v1.GetSessionRequest, v1.GetSessionResponse](
			httpClient,
			baseURL+RunnerServiceGetSessionProcedure,
			opts...,
		),
		listSessions: connect_go.NewClient[v1.ListSessionsRequest, v1.ListSessionsResponse](
			httpClient,
			baseURL+RunnerServiceListSessionsProcedure,
			opts...,
		),
		deleteSession: connect_go.NewClient[v1.DeleteSessionRequest, v1.DeleteSessionResponse](
			httpClient,
			baseURL+RunnerServiceDeleteSessionProcedure,
			opts...,
		),
		execute: connect_go.NewClient[v1.ExecuteRequest, v1.ExecuteResponse](
			httpClient,
			baseURL+RunnerServiceExecuteProcedure,
			opts...,
		),
	}
}

// runnerServiceClient implements RunnerServiceClient.
type runnerServiceClient struct {
	createSession *connect_go.Client[v1.CreateSessionRequest, v1.CreateSessionResponse]
	getSession    *connect_go.Client[v1.GetSessionRequest, v1.GetSessionResponse]
	listSessions  *connect_go.Client[v1.ListSessionsRequest, v1.ListSessionsResponse]
	deleteSession *connect_go.Client[v1.DeleteSessionRequest, v1.DeleteSessionResponse]
	execute       *connect_go.Client[v1.ExecuteRequest, v1.ExecuteResponse]
}

// CreateSession calls runme.runner.v2alpha1.RunnerService.CreateSession.
func (c *runnerServiceClient) CreateSession(ctx context.Context, req *connect_go.Request[v1.CreateSessionRequest]) (*connect_go.Response[v1.CreateSessionResponse], error) {
	return c.createSession.CallUnary(ctx, req)
}

// GetSession calls runme.runner.v2alpha1.RunnerService.GetSession.
func (c *runnerServiceClient) GetSession(ctx context.Context, req *connect_go.Request[v1.GetSessionRequest]) (*connect_go.Response[v1.GetSessionResponse], error) {
	return c.getSession.CallUnary(ctx, req)
}

// ListSessions calls runme.runner.v2alpha1.RunnerService.ListSessions.
func (c *runnerServiceClient) ListSessions(ctx context.Context, req *connect_go.Request[v1.ListSessionsRequest]) (*connect_go.Response[v1.ListSessionsResponse], error) {
	return c.listSessions.CallUnary(ctx, req)
}

// DeleteSession calls runme.runner.v2alpha1.RunnerService.DeleteSession.
func (c *runnerServiceClient) DeleteSession(ctx context.Context, req *connect_go.Request[v1.DeleteSessionRequest]) (*connect_go.Response[v1.DeleteSessionResponse], error) {
	return c.deleteSession.CallUnary(ctx, req)
}

// Execute calls runme.runner.v2alpha1.RunnerService.Execute.
func (c *runnerServiceClient) Execute(ctx context.Context) *connect_go.BidiStreamForClient[v1.ExecuteRequest, v1.ExecuteResponse] {
	return c.execute.CallBidiStream(ctx)
}

// RunnerServiceHandler is an implementation of the runme.runner.v2alpha1.RunnerService service.
type RunnerServiceHandler interface {
	CreateSession(context.Context, *connect_go.Request[v1.CreateSessionRequest]) (*connect_go.Response[v1.CreateSessionResponse], error)
	GetSession(context.Context, *connect_go.Request[v1.GetSessionRequest]) (*connect_go.Response[v1.GetSessionResponse], error)
	ListSessions(context.Context, *connect_go.Request[v1.ListSessionsRequest]) (*connect_go.Response[v1.ListSessionsResponse], error)
	DeleteSession(context.Context, *connect_go.Request[v1.DeleteSessionRequest]) (*connect_go.Response[v1.DeleteSessionResponse], error)
	// Execute executes a program. Examine "ExecuteRequest" to explore
	// configuration options.
	//
	// It's a bidirectional stream RPC method. It expects the first
	// "ExecuteRequest" to contain details of a program to execute.
	// Subsequent "ExecuteRequest" should only contain "input_data" as
	// other fields will be ignored.
	Execute(context.Context, *connect_go.BidiStream[v1.ExecuteRequest, v1.ExecuteResponse]) error
}

// NewRunnerServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewRunnerServiceHandler(svc RunnerServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	runnerServiceCreateSessionHandler := connect_go.NewUnaryHandler(
		RunnerServiceCreateSessionProcedure,
		svc.CreateSession,
		opts...,
	)
	runnerServiceGetSessionHandler := connect_go.NewUnaryHandler(
		RunnerServiceGetSessionProcedure,
		svc.GetSession,
		opts...,
	)
	runnerServiceListSessionsHandler := connect_go.NewUnaryHandler(
		RunnerServiceListSessionsProcedure,
		svc.ListSessions,
		opts...,
	)
	runnerServiceDeleteSessionHandler := connect_go.NewUnaryHandler(
		RunnerServiceDeleteSessionProcedure,
		svc.DeleteSession,
		opts...,
	)
	runnerServiceExecuteHandler := connect_go.NewBidiStreamHandler(
		RunnerServiceExecuteProcedure,
		svc.Execute,
		opts...,
	)
	return "/runme.runner.v2alpha1.RunnerService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case RunnerServiceCreateSessionProcedure:
			runnerServiceCreateSessionHandler.ServeHTTP(w, r)
		case RunnerServiceGetSessionProcedure:
			runnerServiceGetSessionHandler.ServeHTTP(w, r)
		case RunnerServiceListSessionsProcedure:
			runnerServiceListSessionsHandler.ServeHTTP(w, r)
		case RunnerServiceDeleteSessionProcedure:
			runnerServiceDeleteSessionHandler.ServeHTTP(w, r)
		case RunnerServiceExecuteProcedure:
			runnerServiceExecuteHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedRunnerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedRunnerServiceHandler struct{}

func (UnimplementedRunnerServiceHandler) CreateSession(context.Context, *connect_go.Request[v1.CreateSessionRequest]) (*connect_go.Response[v1.CreateSessionResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runme.runner.v2alpha1.RunnerService.CreateSession is not implemented"))
}

func (UnimplementedRunnerServiceHandler) GetSession(context.Context, *connect_go.Request[v1.GetSessionRequest]) (*connect_go.Response[v1.GetSessionResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runme.runner.v2alpha1.RunnerService.GetSession is not implemented"))
}

func (UnimplementedRunnerServiceHandler) ListSessions(context.Context, *connect_go.Request[v1.ListSessionsRequest]) (*connect_go.Response[v1.ListSessionsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runme.runner.v2alpha1.RunnerService.ListSessions is not implemented"))
}

func (UnimplementedRunnerServiceHandler) DeleteSession(context.Context, *connect_go.Request[v1.DeleteSessionRequest]) (*connect_go.Response[v1.DeleteSessionResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runme.runner.v2alpha1.RunnerService.DeleteSession is not implemented"))
}

func (UnimplementedRunnerServiceHandler) Execute(context.Context, *connect_go.BidiStream[v1.ExecuteRequest, v1.ExecuteResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runme.runner.v2alpha1.RunnerService.Execute is not implemented"))
}
