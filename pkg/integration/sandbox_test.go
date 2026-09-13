package integration

// Tests for NewIntegration's handling of the environment the remote's custom-installed-driver
// sandbox sets - see doc/integration-driver/driver-installation.md's "Runtime environment"
// section. These env vars are the *only* way the remote can tell a driver.json-installed driver
// what to bind to (no CLI arguments are passed), so getting the names and precedence right here
// matters in a way that's easy to get wrong silently.

import (
	"testing"
)

func TestNewIntegration_UCIntegrationEnvVarsOverrideConfig(t *testing.T) {
	t.Setenv(UCIntegrationHTTPPortEnv, "9999")
	t.Setenv(UCIntegrationInterfaceEnv, "127.0.0.1")

	i, err := NewIntegration(Config{ListenPort: 8080, WebsocketPath: "/ws", ConfigHome: t.TempDir()})
	if err != nil {
		t.Fatalf("NewIntegration: %v", err)
	}

	if i.Config.ListenPort != 9999 {
		t.Errorf("Config.ListenPort = %d, want 9999 (from %s)", i.Config.ListenPort, UCIntegrationHTTPPortEnv)
	}
	if i.Config.BindInterface != "127.0.0.1" {
		t.Errorf("Config.BindInterface = %q, want \"127.0.0.1\" (from %s)", i.Config.BindInterface, UCIntegrationInterfaceEnv)
	}
	if i.listenAddress != "127.0.0.1:9999" {
		t.Errorf("listenAddress = %q, want \"127.0.0.1:9999\"", i.listenAddress)
	}
}

func TestNewIntegration_MalformedPortEnvIsIgnored(t *testing.T) {
	t.Setenv(UCIntegrationHTTPPortEnv, "not-a-number")

	i, err := NewIntegration(Config{ListenPort: 8080, WebsocketPath: "/ws", ConfigHome: t.TempDir()})
	if err != nil {
		t.Fatalf("NewIntegration: %v", err)
	}

	if i.Config.ListenPort != 8080 {
		t.Errorf("Config.ListenPort = %d, want the configured 8080 to survive an unparsable %s", i.Config.ListenPort, UCIntegrationHTTPPortEnv)
	}
}

func TestNewIntegration_DefaultBindAddressIsAllInterfaces(t *testing.T) {
	i, err := NewIntegration(Config{ListenPort: 8080, WebsocketPath: "/ws", ConfigHome: t.TempDir()})
	if err != nil {
		t.Fatalf("NewIntegration: %v", err)
	}

	if i.listenAddress != "0.0.0.0:8080" {
		t.Errorf("listenAddress = %q, want \"0.0.0.0:8080\" when BindInterface isn't set", i.listenAddress)
	}
}

func TestNewIntegration_SandboxDetectionDisablesMDNSAndRegistration(t *testing.T) {
	t.Setenv(UCIntegrationHTTPPortEnv, "8080")

	i, err := NewIntegration(Config{
		ListenPort:         8080,
		WebsocketPath:      "/ws",
		ConfigHome:         t.TempDir(),
		DisableMDNS:        false, // explicitly requesting mDNS...
		EnableRegistration: true,  // ...and registration - both should be overridden.
	})
	if err != nil {
		t.Fatalf("NewIntegration: %v", err)
	}

	if !i.Config.DisableMDNS {
		t.Error("Config.DisableMDNS = false, want true when running sandboxed (mDNS discovery is meaningless there)")
	}
	if i.Config.EnableRegistration {
		t.Error("Config.EnableRegistration = true, want false when running sandboxed (self-registration is meaningless there)")
	}
}

func TestNewIntegration_NotSandboxedLeavesConfigAlone(t *testing.T) {
	// Explicitly not setting UCIntegrationHTTPPortEnv - the default state for a container or
	// standalone deployment.
	i, err := NewIntegration(Config{
		ListenPort:         8080,
		WebsocketPath:      "/ws",
		ConfigHome:         t.TempDir(),
		DisableMDNS:        false,
		EnableRegistration: true,
	})
	if err != nil {
		t.Fatalf("NewIntegration: %v", err)
	}

	if i.Config.DisableMDNS {
		t.Error("Config.DisableMDNS = true, want false to be left alone outside the sandbox")
	}
	if !i.Config.EnableRegistration {
		t.Error("Config.EnableRegistration = false, want true to be left alone outside the sandbox")
	}
}

func TestNewIntegration_DataHomeDefaultsToConfigHome(t *testing.T) {
	dir := t.TempDir()
	i, err := NewIntegration(Config{ListenPort: 8080, WebsocketPath: "/ws", ConfigHome: dir})
	if err != nil {
		t.Fatalf("NewIntegration: %v", err)
	}

	if i.Config.DataHome != dir {
		t.Errorf("Config.DataHome = %q, want it to default to ConfigHome (%q)", i.Config.DataHome, dir)
	}
}

func TestNewIntegration_DataHomeExplicitValueIsKept(t *testing.T) {
	dataDir := t.TempDir()
	i, err := NewIntegration(Config{ListenPort: 8080, WebsocketPath: "/ws", ConfigHome: t.TempDir(), DataHome: dataDir})
	if err != nil {
		t.Fatalf("NewIntegration: %v", err)
	}

	if i.Config.DataHome != dataDir {
		t.Errorf("Config.DataHome = %q, want the explicitly configured %q to be kept", i.Config.DataHome, dataDir)
	}
}

func TestPersistAndLoadSetupData_CreatesConfigDirectory(t *testing.T) {
	// A subdirectory that doesn't exist yet - PersistSetupData must create it, not just fail
	// silently (the sandbox auto-creates UC_CONFIG_HOME for a custom-installed driver, but nothing
	// guarantees a directory exists for other deployments, e.g. a bare binary run outside the
	// Dockerfile's pre-created /app/ucconfig).
	configHome := t.TempDir() + "/does/not/exist/yet"

	i, err := NewIntegration(Config{ListenPort: 8080, WebsocketPath: "/ws", ConfigHome: configHome})
	if err != nil {
		t.Fatalf("NewIntegration: %v", err)
	}
	i.Metadata = &DriverMetadata{DriverId: "test-driver", Name: LanguageText{En: "Test"}, Version: "0.0.0"}
	i.SetupData = SetupData{"key": "value"}

	if err := i.PersistSetupData(); err != nil {
		t.Fatalf("PersistSetupData: %v", err)
	}

	i.SetupData = nil
	i.LoadSetupData()

	if i.SetupData["key"] != "value" {
		t.Errorf("SetupData[\"key\"] = %q after round-tripping through Persist/LoadSetupData, want \"value\"", i.SetupData["key"])
	}
}
