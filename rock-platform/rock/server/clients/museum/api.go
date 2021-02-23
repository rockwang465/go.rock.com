package museum

import "fmt"

const (
	chartList             = "%s/api/charts"
	namedChartVersionList = "%s/api/charts/%s"
	nameChartVersion      = "%s/api/charts/%s/%s"
)

// Returns all charts list.
func (c *client) Charts() (*ChartMapper, error) {
	var out = make(ChartMapper)
	uri := fmt.Sprintf(chartList, c.addr)
	err := c.get(uri, &out)
	return &out, err
}

// Returns a named chart's version list.
func (c *client) Versions(name string) ([]*ChartVersion, error) {
	var out = make([]*ChartVersion, 0)
	uri := fmt.Sprintf(namedChartVersionList, c.addr, name)
	err := c.get(uri, &out)
	return out, err
}

// Returns a named chart's version list.
func (c *client) Version(name, version string) (*ChartVersion, error) {
	var out = ChartVersion{}
	uri := fmt.Sprintf(nameChartVersion, c.addr, name, version)
	err := c.get(uri, &out)
	return &out, err
}

// Delete a named chart's version.
func (c *client) DeleteVersion(name, version string) error {
	uri := fmt.Sprintf(nameChartVersion, c.addr, name, version)
	err := c.delete(uri)
	return err
}
