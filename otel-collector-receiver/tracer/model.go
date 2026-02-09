package tracer

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
)

type Atm struct {
	ID           int64
	Version      string
	Name         string
	StateID      string
	SerialNumber string
	ISPNetwork   string
}

type BackendSystem struct {
	Version       string
	ProcessName   string
	OSType        string
	OSVersion     string
	CloudProvider string
	CloudRegion   string
	Endpoint      string
}

func generateAtm() Atm {
	i := getRandomNumber(1, 2)
	var newAtm Atm

	switch i {
	case 1:
		newAtm = Atm{
			ID:           111,
			Name:         "ATM-111-IL",
			SerialNumber: "atmxph-2022-111",
			Version:      "v1.0",
			ISPNetwork:   "comcast-chicago",
			StateID:      "IL",
		}

	case 2:
		newAtm = Atm{
			ID:           222,
			Name:         "ATM-222-CA",
			SerialNumber: "atmxph-2022-222",
			Version:      "v1.0",
			ISPNetwork:   "comcast-sanfrancisco",
			StateID:      "CA",
		}
	}

	return newAtm
}

func generateBackendSystem() BackendSystem {
	i := getRandomNumber(1, 3)

	newBackend := BackendSystem{
		ProcessName:   "accounts",
		Version:       "v2.5",
		OSType:        "lnx",
		OSVersion:     "4.16.10-300.fc28.x86_64",
		CloudProvider: "amzn",
		CloudRegion:   "us-east-2",
	}

	switch i {
	case 1:
		newBackend.Endpoint = "api/v2.5/balance"
	case 2:
		newBackend.Endpoint = "api/v2.5/deposit"
	case 3:
		newBackend.Endpoint = "api/v2.5/withdrawn"
	}

	return newBackend
}

func getRandomNumber(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

// OTel part
func generateTraces(numberOfTraces int) ptrace.Traces {
	traces := ptrace.NewTraces()

	for range numberOfTraces {
		newAtm := generateAtm()
		newBackendSystem := generateBackendSystem()

		resourceSpan := traces.ResourceSpans().AppendEmpty()
		atmResource := resourceSpan.Resource()
		fillResourceWithAtm(&atmResource, newAtm)
		atmInstScope := appendAtmSystemInstrScopeSpans(&resourceSpan)

		resourceSpan = traces.ResourceSpans().AppendEmpty()
		backendResource := resourceSpan.Resource()
		fillResourceWithBackendSystem(&backendResource, newBackendSystem)
		backendInstScope := appendAtmSystemInstrScopeSpans(&resourceSpan)

		appendTraceSpans(&newBackendSystem, &backendInstScope, &atmInstScope)
	}

	return traces
}

// Notice that we didn’t add an attribute named “atm.name” or “backendsystem.name”
// to the pcommon.Resource representing the Atm and BackendSystem entity names.
// This is because most (not to say all) distributed tracing backend systems compatible
// with the OTel trace specification interpret the pcommon.Resource described in a
// trace as a Service. Therefore, they expect the pcommon.Resource to have a required
// attribute named service.name as prescribed by the resource semantic convention.
func fillResourceWithAtm(resource *pcommon.Resource, atm Atm) {
	atmAttrs := resource.Attributes()
	atmAttrs.PutInt("atm.id", atm.ID)
	atmAttrs.PutStr("atm.stateid", atm.StateID)
	atmAttrs.PutStr("atm.ispnetwork", atm.ISPNetwork)
	atmAttrs.PutStr("atm.serialnumber", atm.SerialNumber)
	atmAttrs.PutStr(semconv.AttributeServiceName, atm.Name)
	atmAttrs.PutStr(semconv.AttributeServiceVersion, atm.Version)
}

func fillResourceWithBackendSystem(resource *pcommon.Resource, backend BackendSystem) {
	backendAttrs := resource.Attributes()
	var osType, cloudProvider string

	switch {
	case backend.CloudProvider == "amzn":
		cloudProvider = semconv.AttributeCloudProviderAWS
	case backend.OSType == "mcrsft":
		cloudProvider = semconv.AttributeCloudProviderAzure
	case backend.OSType == "gogl":
		cloudProvider = semconv.AttributeCloudProviderGCP
	}

	backendAttrs.PutStr(semconv.AttributeCloudProvider, cloudProvider)
	backendAttrs.PutStr(semconv.AttributeCloudRegion, backend.CloudRegion)

	switch {
	case backend.OSType == "lnx":
		osType = semconv.AttributeOSTypeLinux
	case backend.OSType == "wndws":
		osType = semconv.AttributeOSTypeWindows
	case backend.OSType == "slrs":
		osType = semconv.AttributeOSTypeSolaris
	}

	backendAttrs.PutStr(semconv.AttributeOSType, osType)
	backendAttrs.PutStr(semconv.AttributeOSVersion, backend.OSVersion)

	backendAttrs.PutStr(semconv.AttributeServiceName, backend.ProcessName)
	backendAttrs.PutStr(semconv.AttributeServiceVersion, backend.Version)
}

func appendAtmSystemInstrScopeSpans(resourceSpans *ptrace.ResourceSpans) ptrace.ScopeSpans {
	scopeSpans := resourceSpans.ScopeSpans().AppendEmpty()
	scopeSpans.Scope().SetName("atm-system")
	scopeSpans.Scope().SetVersion("v1.0")
	return scopeSpans
}

func NewTraceID() pcommon.TraceID {
	return pcommon.TraceID(uuid.New())
}

func NewSpanID() pcommon.SpanID {
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	randSource := rand.New(rand.NewSource(rngSeed))

	var sid [8]byte
	randSource.Read(sid[:])
	spanID := pcommon.SpanID(sid)

	return spanID
}

func appendTraceSpans(backend *BackendSystem, backendScopeSpans *ptrace.ScopeSpans, atmScopeSpans *ptrace.ScopeSpans) {
	traceId := NewTraceID()

	var atmOperationName string

	switch {
	case strings.Contains(backend.Endpoint, "balance"):
		atmOperationName = "Check Balance"
	case strings.Contains(backend.Endpoint, "deposit"):
		atmOperationName = "Make Deposit"
	case strings.Contains(backend.Endpoint, "withdraw"):
		atmOperationName = "Fast Cash"
	}

	atmSpanId := NewSpanID()
	atmSpanStartTime := time.Now()
	atmDuration, _ := time.ParseDuration("4s")
	atmSpanFinishTime := atmSpanStartTime.Add(atmDuration)

	atmSpan := atmScopeSpans.Spans().AppendEmpty()
	atmSpan.SetTraceID(traceId)
	atmSpan.SetSpanID(atmSpanId)
	atmSpan.SetName(atmOperationName)
	atmSpan.SetKind(ptrace.SpanKindClient)
	atmSpan.Status().SetCode(ptrace.StatusCodeOk)
	atmSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(atmSpanStartTime))
	atmSpan.SetEndTimestamp(pcommon.NewTimestampFromTime(atmSpanFinishTime))

	backendSpanId := NewSpanID()

	backendDuration, _ := time.ParseDuration("2s")
	backendSpanStartTime := atmSpanStartTime.Add(backendDuration)

	backendSpan := backendScopeSpans.Spans().AppendEmpty()
	backendSpan.SetTraceID(atmSpan.TraceID())
	backendSpan.SetSpanID(backendSpanId)
	backendSpan.SetParentSpanID(atmSpan.SpanID())
	backendSpan.SetName(backend.Endpoint)
	backendSpan.SetKind(ptrace.SpanKindServer)
	backendSpan.Status().SetCode(ptrace.StatusCodeOk)
	backendSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(backendSpanStartTime))
	backendSpan.SetEndTimestamp(atmSpan.EndTimestamp())
}
