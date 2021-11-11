package sshconfig

import (
	"errors"
	"github.com/spencercjh/sshctx/internal/testutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestHost_ToSSHParameter(t *testing.T) {
	tests := []struct {
		name string
		host Host
		want string
	}{
		{name: "ipv4", host: Host{Host: "192.168.1.1", Username: "test", Port: 22}, want: "test@192.168.1.1 -p 22"},
		{name: "domain", host: Host{Host: "test.com", Username: "test", Port: 22}, want: "test@test.com -p 22"},
		{name: "localhost", host: Host{Host: "localhost", Username: "test", Port: 22}, want: "test@localhost -p 22"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Host{
				Host:     tt.host.Host,
				Username: tt.host.Username,
				Port:     tt.host.Port,
			}
			if got := h.ToSSHParameter(); got != tt.want {
				t.Errorf("ToSSHParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSSHConfig_Parse(t *testing.T) {
	testutil.SetupSSHConfig(t)
	defer testutil.TearDownSSHConfig()

	sc := new(SSHConfig).WithLoader(DefaultLoader)

	tests := []struct {
		name    string
		fields  SSHConfig
		wantErr bool
	}{
		{name: "default", fields: *sc, wantErr: false},
		{name: "uninitialized", fields: SSHConfig{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SSHConfig{
				loader:        tt.fields.loader,
				sshctxDataRWC: tt.fields.sshctxDataRWC,
				sshconfigRWC:  tt.fields.sshconfigRWC,
				Hosts:         tt.fields.Hosts,
				PreviousHost:  tt.fields.PreviousHost,
				rootNode:      tt.fields.rootNode,
			}
			if err := s.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSSHConfig_Parse_withSpecific(t *testing.T) {
	testutil.SetupSSHConfig(t)
	defer testutil.TearDownSSHConfig()

	t.Run("specific-sshconfig-and-sshctx-data", func(t *testing.T) {
		sshconfigPath := filepath.Join(cwd, "..", "..", "test", "ssh-config-example")
		t.Setenv("SSHCONFIG", sshconfigPath)
		if _, err := os.Stat(sshconfigPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshconfigPath)
			}
		}
		sshctxDataPath := filepath.Join(cwd, "..", "..", "test", "config_example.yaml")
		t.Setenv("SSHCTX", sshctxDataPath)
		if _, err := os.Stat(sshctxDataPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshctxDataPath)
			}
		}
		sc := new(SSHConfig).WithLoader(DefaultLoader)
		defer sc.Close()
		if err := sc.Parse(); err != nil {
			t.Errorf("Parse() error should nil")
		}
		hostCount := 12
		if len(sc.Hosts) != hostCount {
			t.Errorf("Parse() result: Hosts should be: %d but %d", hostCount, len(sc.Hosts))
		}
	})

	t.Run("specific-sshconfig-and-empty-sshctx-data", func(t *testing.T) {
		sshconfigPath := filepath.Join(cwd, "..", "..", "test", "ssh-config-example")
		t.Setenv("SSHCONFIG", sshconfigPath)
		if _, err := os.Stat(sshconfigPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshconfigPath)
			}
		}
		sshctxDataPath := filepath.Join(cwd, "..", "..", "test", "empty_config.yaml")
		t.Setenv("SSHCTX", sshctxDataPath)
		if _, err := os.Stat(sshctxDataPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshctxDataPath)
			}
		}
		sc := new(SSHConfig).WithLoader(DefaultLoader)
		defer sc.Close()
		if err := sc.Parse(); err != nil {
			t.Errorf("Parse() error should nil")
		}
		hostCount := 12
		if len(sc.Hosts) != hostCount {
			t.Errorf("Parse() result: Hosts should be: %d but %d", hostCount, len(sc.Hosts))
		}
		if sc.PreviousHost != EmptyHost {
			t.Errorf("Parse() result: PreviousHost should be Empty but %v", sc.PreviousHost)
		}
	})

	t.Run("specific-sshconfig-and-blank-sshctx-data", func(t *testing.T) {
		sshconfigPath := filepath.Join(cwd, "..", "..", "test", "ssh-config-example")
		t.Setenv("SSHCONFIG", sshconfigPath)
		if _, err := os.Stat(sshconfigPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshconfigPath)
			}
		}
		sshctxDataPath := filepath.Join(cwd, "..", "..", "test", "blank_config.yaml")
		t.Setenv("SSHCTX", sshctxDataPath)
		if _, err := os.Stat(sshctxDataPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshctxDataPath)
			}
		}
		sc := new(SSHConfig).WithLoader(DefaultLoader)
		defer sc.Close()
		if err := sc.Parse(); err != nil {
			t.Errorf("Parse() error should nil")
		}
		hostCount := 12
		if len(sc.Hosts) != hostCount {
			t.Errorf("Parse() result: Hosts should be: %d but %d", hostCount, len(sc.Hosts))
		}
		if sc.PreviousHost != EmptyHost {
			t.Errorf("Parse() result: PreviousHost should be Empty but %v", sc.PreviousHost)
		}
	})

	t.Run("specific-sshconfig-and-missing-sshctx-data", func(t *testing.T) {
		sshconfigPath := filepath.Join(cwd, "..", "..", "test", "ssh-config-example")
		t.Setenv("SSHCONFIG", sshconfigPath)
		if _, err := os.Stat(sshconfigPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshconfigPath)
			}
		}
		sshctxDataPath := filepath.Join(cwd, "..", "..", "test", "missing_config.yaml")
		t.Setenv("SSHCTX", sshctxDataPath)
		if _, err := os.Stat(sshctxDataPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshctxDataPath)
			}
		}
		sc := new(SSHConfig).WithLoader(DefaultLoader)
		defer sc.Close()
		if err := sc.Parse(); err != nil {
			t.Errorf("Parse() error should nil")
		}
		hostCount := 12
		if len(sc.Hosts) != hostCount {
			t.Errorf("Parse() result: Hosts should be: %d but %d", hostCount, len(sc.Hosts))
		}
		if sc.PreviousHost != EmptyHost {
			t.Errorf("Parse() result: PreviousHost should be Empty but %v", sc.PreviousHost)
		}
	})

	t.Run("specific-sshconfig-and-wrong-sshctx-data", func(t *testing.T) {
		sshconfigPath := filepath.Join(cwd, "..", "..", "test", "ssh-config-example")
		t.Setenv("SSHCONFIG", sshconfigPath)
		if _, err := os.Stat(sshconfigPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshconfigPath)
			}
		}
		sshctxDataPath := filepath.Join(cwd, "..", "..", "test", "wrong_config.yaml")
		t.Setenv("SSHCTX", sshctxDataPath)
		if _, err := os.Stat(sshctxDataPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", sshctxDataPath)
			}
		}
		sc := new(SSHConfig).WithLoader(DefaultLoader)
		defer sc.Close()
		if err := sc.Parse(); err != nil {
			t.Errorf("Parse() error should nil")
		}
		hostCount := 12
		if len(sc.Hosts) != hostCount {
			t.Errorf("Parse() result: Hosts should be: %d but %d", hostCount, len(sc.Hosts))
		}
		if sc.PreviousHost != EmptyHost {
			t.Errorf("Parse() result: PreviousHost should be Empty but %v", sc.PreviousHost)
		}
	})

	t.Run("empty-specific-sshconfig", func(t *testing.T) {
		path := filepath.Join(cwd, "..", "..", "test", "empty-ssh-config")
		t.Setenv("SSHCONFIG", path)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", path)
			}
		}
		sc := new(SSHConfig).WithLoader(DefaultLoader)
		defer sc.Close()
		if err := sc.Parse(); err == nil {
			t.Errorf("Parse() error shouldn't nil")
		}
	})
}

func TestSSHConfig_Close(t *testing.T) {
	testutil.SetupSSHConfig(t)
	defer testutil.TearDownSSHConfig()

	sc := new(SSHConfig).WithLoader(DefaultLoader)
	_ = sc.Parse()

	tests := []struct {
		name   string
		fields SSHConfig
		want   []error
	}{
		{name: "default", fields: *sc, want: []error{nil, nil}},
		{name: "uninitialized", fields: SSHConfig{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SSHConfig{
				loader:        tt.fields.loader,
				sshctxDataRWC: tt.fields.sshctxDataRWC,
				sshconfigRWC:  tt.fields.sshconfigRWC,
				Hosts:         tt.fields.Hosts,
				PreviousHost:  tt.fields.PreviousHost,
				rootNode:      tt.fields.rootNode,
			}
			if got := s.Close(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Close() = %v, want %v", got, tt.want)
			}
		})
	}
}
