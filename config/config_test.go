package config

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	A  string
	B  int `json:",string"`
	KV map[string]string

	configFileName string
	isCreated      bool
}

func (t TestConfig) ConfigFileName() string {
	return t.configFileName
}

func (t *TestConfig) SetConfigFileName(c string) {
	t.configFileName = c
}

func (t TestConfig) IsCreated() bool {
	return t.isCreated
}

func (t *TestConfig) SetIsCreated(i bool) {
	t.isCreated = i
}

func (t TestConfig) AppName() string {
	return "test"
}

func (t TestConfig) DefaultFileName() string {
	return "test.cfg.json"
}

var defaultCfg = &TestConfig{"AVal", 1234, map[string]string{"A": "a", "B": "b"}, "", false}

var defaultCfgJSON = []byte(`{
  "A": "AVal",
  "B": "1234",
  "KV": {
    "A": "a",
    "B": "b"
  }
}`)

func TestCreatesDefaultConfigFile(t *testing.T) {
	testCfgFile := filepath.Join(os.TempDir(), "test.cfg.json")
	// remove tmp if exists
	_, err := os.Stat(testCfgFile)
	if err != nil && !os.IsNotExist(err) {
		t.Error(err)
	}
	if err == nil {
		err = os.Remove(testCfgFile)
		if err != nil {
			t.Errorf("could not remove test cfg file %s: %v", testCfgFile, err)
		}
	}
	LoadConfig(testCfgFile, defaultCfg)
	f, err := os.Open(testCfgFile)
	if err != nil {
		t.Errorf("could not open test cfg file %s: %v", testCfgFile, err)
	}
	defer os.Remove(testCfgFile)
	defer f.Close()
	buf := make([]byte, len(defaultCfgJSON))
	_, err = io.ReadFull(f, buf)
	if err != nil {
		t.Errorf("could not read created cfg file %s: %v", testCfgFile, err)
	}
	if !bytes.Equal(defaultCfgJSON, buf) {
		t.Errorf("invalid default cfg\nexpect:\n%s (%x)\n\ngot:\n%s (%x)", defaultCfgJSON, defaultCfgJSON, buf, buf)
	}
}
