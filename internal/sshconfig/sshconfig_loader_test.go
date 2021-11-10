package sshconfig

import (
	"errors"
	"github.com/spencercjh/sshctx/internal/env"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_getSSHCtxDataDir(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{name: "default", want: filepath.Join(os.Getenv("HOME"), ".sshctx"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSSHCtxDataDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("getSSHCtxDataDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getSSHCtxDataDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSSHCtxDataPath(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{name: "default", want: filepath.Join(os.Getenv("HOME"), ".sshctx", "config.yaml"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSSHCtxDataPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSSHCtxDataPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSSHCtxDataPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSSHConfigPath(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{name: "default", want: filepath.Join(os.Getenv("HOME"), ".ssh", "config"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSSHConfigPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("getSSHConfigPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getSSHConfigPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_openFile(t *testing.T) {
	_ = ioutil.WriteFile("test.txt", []byte("Hello"), 0600)
	test, _ := os.OpenFile("test.txt", os.O_RDONLY, 0755)
	defer func(test *os.File) {
		_ = test.Close()
	}(test)
	defer func() {
		_ = os.Remove("test.txt")
	}()
	type args struct {
		path string
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		{name: "default", args: args{path: "test.txt", name: "test"}, want: test, wantErr: false},
		{name: "not-exist", args: args{path: "not-exist", name: "not-exist"}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openFile(tt.args.path, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("openFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want && got.Name() != tt.want.Name() {
				t.Errorf("openFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var cwd, _ = os.Getwd()

func tearUpSSHCTXData() {
	path, _ := GetSSHCtxDataPath()
	_ = os.Remove(path)
	dir, _ := getSSHCtxDataDir()
	_ = os.Remove(dir)
}

func TestStandardLoader_LoadSSHCTXData(t *testing.T) {

	t.Setenv(env.Debug, "true")

	t.Run("default", func(t *testing.T) {
		st := &StandardLoader{}
		sshctx, err := st.LoadSSHCTXData()
		if err != nil {
			t.Errorf("LoadSSHCTXData() error = %v", err)
		}
		if sshctx == nil {
			t.Errorf("sshctx should not be nil")
		}
	})

	t.Run("exited-specific-file", func(t *testing.T) {
		path := filepath.Join(cwd, "..", "..", "test", "config_example.yaml")
		t.Setenv("SSHCTX", path)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			t.Errorf("test file: %s should exist but not", path)
		}
		st := &StandardLoader{}
		sshctx, err := st.LoadSSHCTXData()
		if err != nil {
			t.Errorf("LoadSSHCTXData() error = %v", err)
		}
		if sshctx == nil {
			t.Errorf("sshctx should not be nil")
		}
	})

	t.Run("unsupported-specific-file", func(t *testing.T) {
		path := filepath.Join(cwd, "..", "..", "test", "config_example.yaml")
		t.Setenv("SSHCTX", path+":"+path)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			t.Errorf("test file: %s should exist but not", path)
		}
		st := &StandardLoader{}
		sshctx, err := st.LoadSSHCTXData()
		if err == nil {
			t.Errorf("LoadSSHCTXData() error should not be nil")
		}
		if sshctx != nil {
			t.Errorf("sshctx should be nil")
		}
	})

	t.Run("non-exited-specific-file", func(t *testing.T) {
		defer tearUpSSHCTXData()
		path := filepath.Join(cwd, "..", "..", "test", "non-existed.yaml")
		t.Setenv("SSHCTX", path)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			if err == nil {
				t.Errorf("test file: %s shouldn't exist but not", path)
			}
		}
		st := &StandardLoader{}
		sshctx, err := st.LoadSSHCTXData()
		if err != nil {
			t.Errorf("LoadSSHCTXData() error = %v", err)
		}
		if sshctx == nil {
			t.Errorf("sshctx should not be nil")
		}
	})

	t.Run("non-exited-specific-file-and-default-file", func(t *testing.T) {
		tearUpSSHCTXData()
		defer tearUpSSHCTXData()
		path := filepath.Join(cwd, "..", "..", "test", "non-existed.yaml")
		t.Setenv("SSHCTX", path)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			if err == nil {
				t.Errorf("test file: %s shouldn't exist but not", path)
			}
		}
		dir, _ := getSSHCtxDataDir()
		_ = os.Remove(filepath.Join(dir, "config.yaml"))
		st := &StandardLoader{}
		sshctx, err := st.LoadSSHCTXData()
		if err != nil {
			t.Errorf("LoadSSHCTXData() error = %v", err)
		}
		if sshctx == nil {
			t.Errorf("sshctx should not be nil")
		}
	})

	t.Run("non-exited-specific-file-and-default-file-but-dir-exist", func(t *testing.T) {
		tearUpSSHCTXData()
		defer tearUpSSHCTXData()
		path := filepath.Join(cwd, "..", "..", "test", "non-existed.yaml")
		t.Setenv("SSHCTX", path)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			if err == nil {
				t.Errorf("test file: %s shouldn't exist but not", path)
			}
		}
		dir, _ := getSSHCtxDataDir()
		_ = os.Mkdir(dir, 0777)
		st := &StandardLoader{}
		sshctx, err := st.LoadSSHCTXData()
		if err != nil {
			t.Errorf("LoadSSHCTXData() error = %v", err)
		}
		if sshctx == nil {
			t.Errorf("sshctx should not be nil")
		}
	})
}

func TestStandardLoader_LoadSSHConfig(t *testing.T) {
	t.Setenv(env.Debug, "true")

	defaultSSHConfigPath, _ := getSSHConfigPath()
	if _, err := os.Stat(defaultSSHConfigPath); err == nil {
		// ~/.ssh/config exist
		t.Run("default", func(t *testing.T) {
			st := &StandardLoader{}
			sshconfig, err := st.LoadSSHConfig()
			if err != nil {
				t.Errorf("LoadSSHConfig() error = %v", err)
			}
			if sshconfig == nil {
				t.Errorf("sshconfig should not be nil")
			}
		})
	}

	t.Run("specific-sshconfig", func(t *testing.T) {
		path := filepath.Join(cwd, "..", "..", "test", "ssh-config-example")
		t.Setenv("SSHCONFIG", path)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", path)
			}
		}
		st := &StandardLoader{}
		sshconfig, err := st.LoadSSHConfig()
		if err != nil {
			t.Errorf("LoadSSHConfig() error = %v", err)
		}
		if sshconfig == nil {
			t.Errorf("sshconfig should not be nil")
		}
	})

	t.Run("non-supported-specific-sshconfig", func(t *testing.T) {
		path := filepath.Join(cwd, "..", "..", "test", "ssh-config-example")
		t.Setenv("SSHCONFIG", path+":"+path)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				t.Errorf("test file: %s should exist but not", path)
			}
		}
		st := &StandardLoader{}
		sshconfig, err := st.LoadSSHConfig()
		if err == nil {
			t.Errorf("LoadSSHConfig() error shouldn't be nil")
		}
		if sshconfig != nil {
			t.Errorf("sshconfig should be nil")
		}
	})

	t.Run("non-existed-specific-sshconfig", func(t *testing.T) {
		path := filepath.Join(cwd, "..", "..", "test", "not-existed")
		t.Setenv("SSHCONFIG", path)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			if err == nil {
				t.Errorf("test file: %s shouldn't exist but not", path)
			}
		}
		st := &StandardLoader{}
		sshconfig, err := st.LoadSSHConfig()
		if err == nil {
			t.Errorf("LoadSSHConfig() error shouldn't be nil")
		}
		if sshconfig != nil {
			t.Errorf("sshconfig should be nil")
		}
	})
}
