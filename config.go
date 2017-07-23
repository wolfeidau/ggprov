package ggprov

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"

	yaml "gopkg.in/yaml.v2"
)

// ThingConfig used to store thing details
type ThingConfig struct {
	ThingRole   *IamRole
	ThingPolicy *IotPolicy
	Thing       *Thing
	ThingCreds  *ThingCreds
	IotEndpoint *IotEndpoint
}

// NewThingConfig create a new thing config
func NewThingConfig(thingRole *IamRole, thingPolicy *IotPolicy, thing *Thing, thingCreds *ThingCreds, endpoint *IotEndpoint) *ThingConfig {
	return &ThingConfig{
		ThingRole:   thingRole,
		ThingPolicy: thingPolicy,
		Thing:       thing,
		ThingCreds:  thingCreds,
		IotEndpoint: endpoint,
	}
}

// Save the configuration to a YAML file with the name used as the filename
func (tc *ThingConfig) Save(name string) error {

	d, err := yaml.Marshal(&tc)
	if err != nil {
		return err
	}

	return errors.Wrap(ioutil.WriteFile(fmt.Sprintf("%s.yaml", name), d, 0600), "Failed to save thing config")
}

// Load the configuration from the supplied path
func Load(filePath string) (*ThingConfig, error) {

	var config ThingConfig

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load config file")
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Unmarshal config file")
	}

	return &config, nil
}
