package manager

import (
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"sigs.k8s.io/yaml"
)

//MergeChartValues : merge values contained in the chart with the
//ones provided by the user
func MergeChartValues(c *chart.Chart, values map[string]interface{}) (*chart.Chart, error) {
	var i int
	var content []byte
	for k, file := range c.Raw {
		if file.Name == "values.yaml" {
			cvalues, errValues := chartutil.CoalesceValues(c, values)
			if errValues != nil {
				return nil, errors.Wrap(errValues, "failed to retrieve values")
			}
			m, err := yaml.Marshal(cvalues)
			if err != nil {
				return nil, err
			}
			content = m
			i = k
			break
		}
	}
	c.Raw[i].Data = content
	return c, nil
}
