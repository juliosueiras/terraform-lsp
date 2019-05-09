// Mostly adapted from tfschema - https://github.com/minamijoyo/tfschema
package tfstructs

import (
	"fmt"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/hashicorp/terraform/providers"
	"github.com/mitchellh/go-homedir"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// Client represents a tfschema Client.
type Client struct {
	provider     *plugin.GRPCProvider
	pluginClient interface{}
}

// NewClient creates a new Client instance.
func NewClient(providerName string, targetDir string) (*Client, error) {
	// find a provider plugin
	pluginMeta, err := findPlugin("provider", providerName, targetDir)
	if err != nil {
		return nil, err
	}

	// initialize a plugin Client.
	pluginClient := plugin.Client(*pluginMeta)
	rpcClient, err := pluginClient.Client()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize plugin: %s", err)
	}

	// create a new resource provider.
	raw, err := rpcClient.Dispense(plugin.ProviderPluginName)
	if err != nil {
		return nil, fmt.Errorf("Failed to dispense plugin: %s", err)
	}

	provider := raw.(*plugin.GRPCProvider)

	return &Client{
		provider:     provider,
		pluginClient: pluginClient,
	}, nil
}

// findPlugin finds a plugin with the name specified in the arguments.
func findPlugin(pluginType string, pluginName string, targetDir string) (*discovery.PluginMeta, error) {
	dirs, err := pluginDirs(targetDir)
	if err != nil {
		return nil, err
	}

	pluginMetaSet := discovery.FindPlugins(pluginType, dirs).WithName(pluginName)

	// if pluginMetaSet doesn't have any pluginMeta, pluginMetaSet.Newest() will call panic.
	// so check it here.
	if pluginMetaSet.Count() > 0 {
		ret := pluginMetaSet.Newest()
		return &ret, nil
	}

	return nil, fmt.Errorf("Failed to find plugin: %s. Plugin binary was not found in any of the following directories: [%s]", pluginName, strings.Join(dirs, ", "))
}

// pluginDirs returns a list of directories to find plugin.
// This is almost the same as Terraform, but not exactly the same.
func pluginDirs(targetDir string) ([]string, error) {
	dirs := []string{}

	// current directory
	dirs = append(dirs, ".")

	// same directory as this executable
	exePath, err := os.Executable()
	if err != nil {
		return []string{}, fmt.Errorf("Failed to get executable path: %s", err)
	}
	dirs = append(dirs, filepath.Dir(exePath))

	// user vendor directory
	arch := runtime.GOOS + "_" + runtime.GOARCH
	vendorDir := filepath.Join("terraform.d", "plugins", arch)
	dirs = append(dirs, vendorDir)

	// auto installed directory
	// This does not take into account overriding the data directory.
	autoInstalledDir := filepath.Join(targetDir, ".terraform", "plugins", arch)
	dirs = append(dirs, autoInstalledDir)

	// global plugin directory
	homeDir, err := homedir.Dir()
	if err != nil {
		return []string{}, fmt.Errorf("Failed to get home dir: %s", err)
	}
	configDir := filepath.Join(homeDir, ".terraform.d", "plugins")
	dirs = append(dirs, configDir)
	dirs = append(dirs, filepath.Join(configDir, arch))

	// GOPATH
	// This is not included in the Terraform, but for convenience.
	gopath := build.Default.GOPATH
	dirs = append(dirs, filepath.Join(gopath, "bin"))

	log.Printf("[DEBUG] plugin dirs: %#v", dirs)
	return dirs, nil
}

// GetRawProviderSchema returns a raw type definiton of provider schema.
func (c *Client) GetRawProviderSchema() (*providers.Schema, error) {

	res := c.provider.GetSchema()
	return &res.Provider, nil
}

// GetRawResourceTypeSchema returns a type definiton of resource type.
func (c *Client) GetRawResourceTypeSchema(resourceType string) (*providers.Schema, error) {

	res := c.provider.GetSchema()
	if res.ResourceTypes[resourceType].Block == nil {
		return nil, fmt.Errorf("Failed to find resource type: %s", resourceType)
	}

	b := res.ResourceTypes[resourceType]
	return &b, nil
}

// GetResourceTypes returns a type definiton of resource type.
func (c *Client) GetResourceTypes() ([]string, error) {
	res := c.provider.GetSchema()
	var result []string

	for k, _ := range res.ResourceTypes {
		result = append(result, k)
	}
	return result, nil
}

// GetDataSourceTypes returns a type definiton of resource type.
func (c *Client) GetDataSourceTypes() ([]string, error) {
	res := c.provider.GetSchema()
	var result []string

	for k, _ := range res.DataSources {
		result = append(result, k)
	}
	return result, nil
}

// GetRawDataSourceTypeSchema returns a type definiton of resource type.
func (c *Client) GetRawDataSourceTypeSchema(dataSourceType string) (*providers.Schema, error) {

	res := c.provider.GetSchema()

	if res.DataSources[dataSourceType].Block == nil {
		return nil, fmt.Errorf("Failed to find data source type: %s", dataSourceType)
	}

	b := res.DataSources[dataSourceType]
	return &b, nil
}

// Kill kills a process of the plugin.
func (c *Client) Kill() {
	// We cannot import the vendor version of go-plugin using terraform.
	// So, we call (*go-plugin.Client).Kill() by reflection here.
	v := reflect.ValueOf(c.pluginClient).MethodByName("Kill")
	v.Call([]reflect.Value{})
}
