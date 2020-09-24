package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/provenance"
)

// Package is the action for packaging a chart.
//
// It provides the implementation of 'helm package'.
type Package struct {
	Sign             bool
	Key              string
	Keyring          string
	Version          string
	AppVersion       string
	Destination      string
	DependencyUpdate bool
	RepositoryConfig string
	RepositoryCache  string
}

// NewPackage creates a new Package object with the given configuration.
func NewPackage() *Package {
	return &Package{}
}

// Run executes 'helm package' against the given chart and returns the path to the packaged chart.
func (p *Package) Run(path string, vals map[string]interface{}) (string, error) {
	ch, err := loader.LoadDir(path)
	if err != nil {
		return "", err
	}

	// If version is set, modify the version.
	if p.Version != "" {
		if err := setVersion(ch, p.Version); err != nil {
			return "", err
		}
	}

	if p.AppVersion != "" {
		ch.Metadata.AppVersion = p.AppVersion
	}

	if reqs := ch.Metadata.Dependencies; reqs != nil {
		if err := CheckDependencies(ch, reqs); err != nil {
			return "", err
		}
	}

	var dest string
	if p.Destination == "." {
		// Save to the current working directory.
		dest, err = os.Getwd()
		if err != nil {
			return "", err
		}
	} else {
		// Otherwise save to set destination
		dest = p.Destination
	}
	chartMerged, errMerge := MergeChartValues(ch, vals)
	if errMerge != nil {
		return "", errors.Wrap(errMerge, "failed to merge values")
	}
	name, err := chartutil.Save(chartMerged, dest)
	//name, err := chartutil.Save(ch, dest)

	if err != nil {
		return "", errors.Wrap(err, "failed to save")
	}

	if p.Sign {
		err = p.Clearsign(name)
	}

	return name, err
}

func setVersion(ch *chart.Chart, ver string) error {
	// Verify that version is a Version, and error out if it is not.
	if _, err := semver.NewVersion(ver); err != nil {
		return err
	}

	// Set the version field on the chart.
	ch.Metadata.Version = ver
	return nil
}

// Clearsign signs a chart
func (p *Package) Clearsign(filename string) error {
	// Load keyring
	signer, err := provenance.NewFromKeyring(p.Keyring, p.Key)
	if err != nil {
		return err
	}

	if err := signer.DecryptKey(promptUser); err != nil {
		return err
	}

	sig, err := signer.ClearSign(filename)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename+".prov", []byte(sig), 0644)
}

// promptUser implements provenance.PassphraseFetcher
func promptUser(name string) ([]byte, error) {
	fmt.Printf("Password for key %q >  ", name)
	// syscall.Stdin is not an int in all environments and needs to be coerced
	// into one there (e.g., Windows)
	pw, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return pw, err
}
