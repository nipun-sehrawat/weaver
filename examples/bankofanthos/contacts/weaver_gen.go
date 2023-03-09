// go:build !ignoreWeaverGen

package contacts

// Code generated by "weaver generate". DO NOT EDIT.
import (
	"context"
	"errors"
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"time"
)

func init() {
	codegen.Register(codegen.Registration{
		Name:        "github.com/ServiceWeaver/weaver/examples/bankofanthos/contacts/T",
		Iface:       reflect.TypeOf((*T)(nil)).Elem(),
		New:         func() any { return &impl{} },
		ConfigFn:    func(i any) any { return i.(*impl).WithConfig.Config() },
		LocalStubFn: func(impl any, tracer trace.Tracer) any { return t_local_stub{impl: impl.(T), tracer: tracer} },
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return t_client_stub{stub: stub, getContactsMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/ServiceWeaver/weaver/examples/bankofanthos/contacts/T", Method: "GetContacts"}), addContactMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/ServiceWeaver/weaver/examples/bankofanthos/contacts/T", Method: "AddContact"})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return t_server_stub{impl: impl.(T), addLoad: addLoad}
		},
	})
}

// Local stub implementations.

type t_local_stub struct {
	impl   T
	tracer trace.Tracer
}

func (s t_local_stub) GetContacts(ctx context.Context, a0 string) (r0 []Contact, err error) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "contacts.T.GetContacts", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.GetContacts(ctx, a0)
}

func (s t_local_stub) AddContact(ctx context.Context, a0 string, a1 Contact) (err error) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "contacts.T.AddContact", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.AddContact(ctx, a0, a1)
}

// Client stub implementations.

type t_client_stub struct {
	stub               codegen.Stub
	getContactsMetrics *codegen.MethodMetrics
	addContactMetrics  *codegen.MethodMetrics
}

func (s t_client_stub) GetContacts(ctx context.Context, a0 string) (r0 []Contact, err error) {
	// Update metrics.
	start := time.Now()
	s.getContactsMetrics.Count.Add(1)

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "contacts.T.GetContacts", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.getContactsMetrics.ErrorCount.Add(1)
		}
		span.End()

		s.getContactsMetrics.Latency.Put(float64(time.Since(start).Microseconds()))
	}()

	// Preallocate a buffer of the right size.
	size := 0
	size += (4 + len(a0))
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.String(a0)
	var shardKey uint64

	// Call the remote method.
	s.getContactsMetrics.BytesRequest.Put(float64(len(enc.Data())))
	var results []byte
	results, err = s.stub.Run(ctx, 1, enc.Data(), shardKey)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}
	s.getContactsMetrics.BytesReply.Put(float64(len(results)))

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = serviceweaver_dec_slice_Contact_d00a3378(dec)
	err = dec.Error()
	return
}

func (s t_client_stub) AddContact(ctx context.Context, a0 string, a1 Contact) (err error) {
	// Update metrics.
	start := time.Now()
	s.addContactMetrics.Count.Add(1)

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "contacts.T.AddContact", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.addContactMetrics.ErrorCount.Add(1)
		}
		span.End()

		s.addContactMetrics.Latency.Put(float64(time.Since(start).Microseconds()))
	}()

	// Preallocate a buffer of the right size.
	size := 0
	size += (4 + len(a0))
	size += serviceweaver_size_Contact_15811618(&a1)
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.String(a0)
	(a1).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	s.addContactMetrics.BytesRequest.Put(float64(len(enc.Data())))
	var results []byte
	results, err = s.stub.Run(ctx, 0, enc.Data(), shardKey)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}
	s.addContactMetrics.BytesReply.Put(float64(len(results)))

	// Decode the results.
	dec := codegen.NewDecoder(results)
	err = dec.Error()
	return
}

// Server stub implementations.

type t_server_stub struct {
	impl    T
	addLoad func(key uint64, load float64)
}

// GetStubFn implements the stub.Server interface.
func (s t_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "GetContacts":
		return s.getContacts
	case "AddContact":
		return s.addContact
	default:
		return nil
	}
}

func (s t_server_stub) getContacts(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 string
	a0 = dec.String()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.GetContacts(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_slice_Contact_d00a3378(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s t_server_stub) addContact(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 string
	a0 = dec.String()
	var a1 Contact
	(&a1).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.AddContact(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

// AutoMarshal implementations.

var _ codegen.AutoMarshal = &Contact{}

func (x *Contact) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Contact.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Username)
	enc.String(x.Label)
	enc.String(x.AccountNum)
	enc.String(x.RoutingNum)
	enc.Bool(x.IsExternal)
}

func (x *Contact) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Contact.WeaverUnmarshal: nil receiver"))
	}
	x.Username = dec.String()
	x.Label = dec.String()
	x.AccountNum = dec.String()
	x.RoutingNum = dec.String()
	x.IsExternal = dec.Bool()
}

// Encoding/decoding implementations.

func serviceweaver_enc_slice_Contact_d00a3378(enc *codegen.Encoder, arg []Contact) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_Contact_d00a3378(dec *codegen.Decoder) []Contact {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]Contact, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}

// Size implementations.

// serviceweaver_size_Contact_15811618 returns the size (in bytes) of the serialization
// of the provided type.
func serviceweaver_size_Contact_15811618(x *Contact) int {
	size := 0
	size += 0
	size += (4 + len(x.Username))
	size += (4 + len(x.Label))
	size += (4 + len(x.AccountNum))
	size += (4 + len(x.RoutingNum))
	size += 1
	return size
}