package helmpack

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

const packUsage = `Helm plugin to pack a helm chart. it uses the same options as the package command but it allows to inject values before packaging
Examples:
  $ helm pack .        # like package command, it creates a helm chart archive file
  $ helm pack . --set mykey=myvalue  # inject/update the mykey property with the myvalue value before creating the helm archive file
`

func newPackCmd(out io.Writer) *cobra.Command {
	client := action.NewPackage()
	valueOpts := &values.Options{}
	cmd := &cobra.Command{
		Use:   "pack [CHART_PATH] [...]",
		Short: "pack a chart directory into a chart archive",
		Long:  packUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.Errorf("need at least one argument, the path to the chart")
			}
			if client.Sign {
				if client.Key == "" {
					return errors.New("--key is required for signing a package")
				}
				if client.Keyring == "" {
					return errors.New("--keyring is required for signing a package")
				}
			}
			client.RepositoryConfig = settings.RepositoryConfig
			client.RepositoryCache = settings.RepositoryCache
			p := getter.All(settings)
			vals, err := valueOpts.MergeValues(p)
			if err != nil {
				return err
			}

			for i := 0; i < len(args); i++ {
				path, err := filepath.Abs(args[i])
				if err != nil {
					return err
				}
				if _, err := os.Stat(args[i]); err != nil {
					return err
				}

				if client.DependencyUpdate {
					downloadManager := &downloader.Manager{
						Out:              ioutil.Discard,
						ChartPath:        path,
						Keyring:          client.Keyring,
						Getters:          p,
						Debug:            settings.Debug,
						RepositoryConfig: settings.RepositoryConfig,
						RepositoryCache:  settings.RepositoryCache,
					}

					if err := downloadManager.Update(); err != nil {
						return err
					}
				}
				p, err := client.Run(path, vals)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "Successfully packaged chart and saved it to: %s\n", p)
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.BoolVar(&client.Sign, "sign", false, "use a PGP private key to sign this package")
	f.StringVar(&client.Key, "key", "", "name of the key to use when signing. Used if --sign is true")
	f.StringVar(&client.Keyring, "keyring", defaultKeyring(), "location of a public keyring")
	f.StringVar(&client.Version, "version", "", "set the version on the chart to this semver version")
	f.StringVar(&client.AppVersion, "app-version", "", "set the appVersion on the chart to this version")
	f.StringVarP(&client.Destination, "destination", "d", ".", "location to write the chart.")
	f.BoolVarP(&client.DependencyUpdate, "dependency-update", "u", false, `update dependencies from "Chart.yaml" to dir "charts/" before packaging`)
	return cmd
}
