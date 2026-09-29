package integration

// Checks the setup-flow wire types against the vendored Core-API spec (see internal/spec). Unlike
// entity features, the setup error enum is small and shared with every driver, so goucrt declares
// all of it: a code missing on either side fails here.

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/splattner/goucrt/internal/spec"
)

var allDriverSetupErrors = []DriverSetupError{
	NoneError,
	NotFoundError,
	ConnectionRefusedError,
	AuthErrorError,
	TimeoutError,
	DriverUnavailableError,
	InvalidInputError,
	AbortedError,
	AlreadyConfiguredError,
	NotSupportedError,
	OtherError,
}

func TestDriverSetupErrorsMatchSpec(t *testing.T) {
	want, err := spec.SetupErrorCodes()
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for _, e := range allDriverSetupErrors {
		got = append(got, string(e))
	}

	for _, code := range got {
		if !slices.Contains(want, code) {
			t.Errorf("goucrt declares setup error %q, which is not in the spec's integrationSetupError enum", code)
		}
	}
	for _, code := range want {
		if !slices.Contains(got, code) {
			t.Errorf("spec defines setup error %q, which goucrt does not declare", code)
		}
	}
}

func TestDriverSetupError_RequiresCoreAPI019(t *testing.T) {
	legacy := []DriverSetupError{NoneError, NotFoundError, ConnectionRefusedError, AuthErrorError, TimeoutError, OtherError}
	for _, e := range allDriverSetupErrors {
		if got, want := e.requiresCoreAPI019(), !slices.Contains(legacy, e); got != want {
			t.Errorf("%s.requiresCoreAPI019() = %v, want %v", e, got, want)
		}
	}
}

func TestVersionLess(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"", "0.19.0", true},
		{"garbage", "0.19.0", true},
		{"0.18.9", "0.19.0", true},
		{"0.19.0", "0.19.0", false},
		{"0.19", "0.19.0", false},
		{"v0.19.1", "0.19.0", false},
		{"0.19.0-beta", "0.19.0", false},
		{"1.0.0", "0.19.0", false},
	}
	for _, c := range cases {
		if got := versionLess(c.a, c.b); got != c.want {
			t.Errorf("versionLess(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestDriverSetupChangeData_ErrorMessage(t *testing.T) {
	without, err := json.Marshal(DriverSetupChangeData{EventType: StopEvent, State: ErrorState, Error: OtherError})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(without, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["error_message"]; ok {
		t.Errorf("error_message present without being set: %s", without)
	}

	with, err := json.Marshal(DriverSetupChangeData{
		EventType:    SetupEvent,
		State:        WaitUserActionState,
		Error:        InvalidInputError,
		ErrorMessage: &LanguageText{En: "PIN rejected"},
	})
	if err != nil {
		t.Fatal(err)
	}
	m = nil
	if err := json.Unmarshal(with, &m); err != nil {
		t.Fatal(err)
	}
	if got, _ := m["error_message"].(map[string]interface{}); got["en"] != "PIN rejected" {
		t.Errorf("error_message = %v, want {en: PIN rejected}", m["error_message"])
	}
	if m["error"] != "INVALID_INPUT" {
		t.Errorf("error = %v, want INVALID_INPUT", m["error"])
	}
}

func TestSetDriverSetupStateWithMessage_NewErrorWithoutMetadata(t *testing.T) {
	i := newTestIntegration(t)

	// Must not panic on nil Metadata while checking min_core_api.
	i.SetDriverSetupStateWithMessage(SetupEvent, WaitUserActionState, InvalidInputError, &LanguageText{En: "PIN rejected"}, &RequireUserAction{})

	if i.SetupState != WaitUserActionState {
		t.Errorf("SetupState = %q, want %q", i.SetupState, WaitUserActionState)
	}
}

func TestHandleSetupDriverRequest_StoresLanguage(t *testing.T) {
	i := newTestIntegration(t)
	i.Metadata = &DriverMetadata{DriverId: "test-driver", Name: LanguageText{En: "Test"}, Version: "0.0.0"}

	raw := []byte(`{"kind":"req","id":3,"msg":"setup_driver","msg_data":{"setup_data":{"ipaddr":"10.0.0.5"},"language":"de"}}`)
	var req SetupDriverMessageReq
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatal(err)
	}

	if res := i.handleSetupDriverRequest(&req); res.Code != 200 {
		t.Fatalf("response code = %d, want 200", res.Code)
	}
	if i.SetupLanguage != "de" {
		t.Errorf("SetupLanguage = %q, want %q", i.SetupLanguage, "de")
	}
	if i.SetupData["ipaddr"] != "10.0.0.5" {
		t.Errorf("SetupData[ipaddr] = %q, want 10.0.0.5", i.SetupData["ipaddr"])
	}
}
