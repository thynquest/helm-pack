package helmpack

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestPack(t *testing.T) {
	abs, _ := filepath.Abs("../../testitems/samplechart")
	os.Args = []string{"pack", abs, "--set", "deployment.version=myvalue123", "--destination", abs}
	cmd := NewPackCmd(os.Args[1:], os.Stdout)
	if err := cmd.Execute(); err != nil {
		t.Errorf("test package failed: %v", err)
	}
	tarfile, errOpen := os.Open(abs + "/samplechart-0.1.0.tgz")
	if errOpen != nil {
		t.Errorf("opening tar file failed: %v", errOpen)
	}
	uncompressed, errStream := gzip.NewReader(tarfile)
	if errStream != nil {
		t.Errorf("error when uncompressign file: %v", errStream)
	}
	defer tarfile.Close()
	tr := tar.NewReader(uncompressed)
	testOk := false
	for {
		hdr, errHdr := tr.Next()
		if errHdr == io.EOF {
			break
		}
		if errHdr != nil {
			t.Errorf("error when reading tar content: %v", errHdr)
		}
		if strings.Contains(hdr.Name, "values.yaml") {
			bs, _ := ioutil.ReadAll(tr)
			values := map[string]interface{}{}
			if errUnmarshal := yaml.Unmarshal(bs, &values); errUnmarshal != nil {
				t.Errorf("failed reading values file: %v", errUnmarshal)
			}
			deploymentValues, ok := values["deployment"].(map[string]interface{})
			if !ok {
				t.Error("deployment key not found")
			}
			version, okVersion := deploymentValues["version"].(string)
			if !okVersion {
				t.Error("version key not found")
			}
			if version == "myvalue123" {
				testOk = true
			} else {
				t.Errorf("got %s expected %s", version, "myvalue123")
			}
			break
		}
	}
	if !testOk {
		t.Error("value myvalue123 not found in the created chart")
	}
}

func TestPackWithDependencies(t *testing.T) {
	abs, _ := filepath.Abs("../../testitems/sampledeps")
	if _, err := os.Stat(abs + "/charts"); err == nil {
		os.Remove(abs + "/charts")
	}
	os.Args = []string{"pack", abs, "--set", "deployment.version=myvalue123", "--destination", abs, "--dependency-update"}
	cmd := NewPackCmd(os.Args[1:], os.Stdout)
	if err := cmd.Execute(); err != nil {
		t.Errorf("test package failed: %v", err)
	}
	tarfile, errOpen := os.Open(abs + "/sampledeps-0.1.0.tgz")
	if errOpen != nil {
		t.Errorf("opening tar file failed: %v", errOpen)
	}
	uncompressed, errStream := gzip.NewReader(tarfile)
	if errStream != nil {
		t.Errorf("error when uncompressign file: %v", errStream)
	}
	defer tarfile.Close()
	tr := tar.NewReader(uncompressed)
	testDepsOk := false
	testvalueOk := false
	for {
		hdr, errHdr := tr.Next()
		if errHdr == io.EOF {
			break
		}
		if errHdr != nil {
			t.Errorf("error when reading tar content: %v", errHdr)
		}
		//test if dependencies has been loaded
		if strings.Contains(hdr.Name, "charts") {
			testDepsOk = true
		}
		//if the values has been modified
		if strings.Contains(hdr.Name, "values.yaml") {
			bs, _ := ioutil.ReadAll(tr)
			values := map[string]interface{}{}
			if errUnmarshal := yaml.Unmarshal(bs, &values); errUnmarshal != nil {
				t.Errorf("failed reading values file: %v", errUnmarshal)
			}
			deploymentValues, ok := values["deployment"].(map[string]interface{})
			if !ok {
				t.Error("deployment key not found")
			}
			version, okVersion := deploymentValues["version"].(string)
			if !okVersion {
				t.Error("version key not found")
			}
			if version == "myvalue123" {
				testvalueOk = true
			} else {
				t.Errorf("got %s expected %s", version, "myvalue123")
			}
		}

		if testDepsOk && testvalueOk {
			break
		}
	}
	if !testDepsOk {
		t.Error("charts dependency folder not found")
	}
	if !testvalueOk {
		t.Error("property deployment.version not modified")
	}
}
