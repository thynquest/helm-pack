package manager

import (
	"github.com/mitchellh/copystructure"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"sigs.k8s.io/yaml"
)

//MergeChartValues : merge values contained in the chart with the
//ones provided by the user
func MergeChartValues(c *chart.Chart, values map[string]interface{}) (*chart.Chart, error) {
	chartCopy, err := copystructure.Copy(c)
	if err != nil {
		return c, err
	}
	chartResult, ok := chartCopy.(*chart.Chart)
	if !ok {
		return c, errors.New("Cannot merge values")
	}
	for k, file := range c.Raw {
		if file.Name == "values.yaml" {
			cvalues, errValues := chartutil.CoalesceValues(c, values)
			if errValues != nil {
				return nil, errors.Wrap(errValues, "failed to retrieve values")
			}
			m, _ := yaml.Marshal(cvalues)
			chartResult.Raw[k].Data = m
			chartResult.Values = cvalues
			break
		}
	}
	return chartResult, nil
}
