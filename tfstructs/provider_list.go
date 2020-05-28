package tfstructs

type TerraformProvider struct {
	Name string
	Type string
}

var OfficialProviders = make(map[string]TerraformProvider)

func init() {
	for _, category := range []map[string]TerraformProvider{
		majorCloud,
		cloud,
		infra,
		network,
		vcs,
		monitoring,
		database,
		misc,
	} {
		for k, v := range category {
			OfficialProviders[k] = v
		}
	}

}

var majorCloud = map[string]TerraformProvider{
	"alicloud": {
		Name: "AliCloud",
		Type: "Major Cloud",
	},
	"aws": {
		Name: "AWS",
		Type: "Major Cloud",
	},
	"azurerm": {
		Name: "AzureRM",
		Type: "Major Cloud",
	},
	"azurestack": {
		Name: "Azure Stack",
		Type: "Major Cloud",
	},
	"google": {
		Name: "Google Cloud Platform",
		Type: "Major Cloud",
	},
	"oci": {
		Name: "Oracle Cloud Infrastructure",
		Type: "Major Cloud",
	},
	"oraclepaas": {
		Name: "Oracle Cloud Platform",
		Type: "Major Cloud",
	},
	"opc": {
		Name: "Oracle Public Cloud",
		Type: "Major Cloud",
	},
	"nsxt": {
		Name: "VMware NSX-T",
		Type: "Major Cloud",
	},
	"vcd": {
		Name: "vCloud Director",
		Type: "Major Cloud",
	},
	"vra7": {
		Name: "VMware vRA7",
		Type: "Major Cloud",
	},
	"vsphere": {
		Name: "VMware vSphere",
		Type: "Major Cloud",
	},
}

var database = map[string]TerraformProvider{
	"influxdb": {
		Name: "InfluxDB",
		Type: "Database",
	},
	"mysql": {
		Name: "MySQL",
		Type: "Database",
	},
	"postgresql": {
		Name: "PostgreSQL",
		Type: "Database",
	},
}

var infra = map[string]TerraformProvider{
	"chef": {
		Name: "Chef",
		Type: "Infra",
	},
	"consul": {
		Name: "Consul",
		Type: "Infra",
	},
	"docker": {
		Name: "Docker",
		Type: "Infra",
	},
	"helm": {
		Name: "Helm",
		Type: "Infra",
	},
	"kubernetes": {
		Name: "Kubernetes",
		Type: "Infra",
	},
	"mailgun": {
		Name: "Mailgun",
		Type: "Infra",
	},
	"nomad": {
		Name: "Nomad",
		Type: "Infra",
	},
	"rabbitmq": {
		Name: "RabbitMQ",
		Type: "Infra",
	},
	"rancher": {
		Name: "Rancher",
		Type: "Infra",
	},
	"rightscale": {
		Name: "RightScale",
		Type: "Infra",
	},
	"rundeck": {
		Name: "Rundeck",
		Type: "Infra",
	},
	"spotinst": {
		Name: "Spotinst",
		Type: "Infra",
	},
	"terraform": {
		Name: "Terraform",
		Type: "Infra",
	},
	"tfe": {
		Name: "Terraform Enterprise",
		Type: "Infra",
	},
	"vault": {
		Name: "Vault",
		Type: "Infra",
	},
}

var cloud = map[string]TerraformProvider{
	"arukas": {
		Name: "Arukas",
		Type: "Cloud",
	},
	"brightbox": {
		Name: "Brightbox",
		Type: "Cloud",
	},
	"clc": {
		Name: "CenturyLinkCloud",
		Type: "Cloud",
	},
	"cloudscale": {
		Name: "CloudScale.ch",
		Type: "Cloud",
	},
	"cloudstack": {
		Name: "CloudStack",
		Type: "Cloud",
	},
	"digitalocean": {
		Name: "DigitalOcean",
		Type: "Cloud",
	},
	"fastly": {
		Name: "Fastly",
		Type: "Cloud",
	},
	"flexibleengine": {
		Name: "FlexibleEngine",
		Type: "Cloud",
	},
	"gridscale": {
		Name: "Gridscale",
		Type: "Cloud",
	},
	"hedvig": {
		Name: "Hedvig",
		Type: "Cloud",
	},
	"heroku": {
		Name: "Heroku",
		Type: "Cloud",
	},
	"hcloud": {
		Name: "Hetzner Cloud",
		Type: "Cloud",
	},
	"huaweicloud": {
		Name: "HuaweiCloud",
		Type: "Cloud",
	},
	"jdcloud": {
		Name: "JDCloud",
		Type: "Cloud",
	},
	"linode": {
		Name: "Linode",
		Type: "Cloud",
	},
	"ncloud": {
		Name: "Naver Cloud",
		Type: "Cloud",
	},
	"nutanix": {
		Name: "Nutanix",
		Type: "Cloud",
	},
	"openstack": {
		Name: "OpenStack",
		Type: "Cloud",
	},
	"opentelekomcloud": {
		Name: "OpenTelekomCloud",
		Type: "Cloud",
	},
	"ovh": {
		Name: "OVH",
		Type: "Cloud",
	},
	"packet": {
		Name: "Packet",
		Type: "Cloud",
	},
	"profitbricks": {
		Name: "ProfitBricks",
		Type: "Cloud",
	},
	"scaleway": {
		Name: "Scaleway",
		Type: "Cloud",
	},
	"skytap": {
		Name: "Skytap",
		Type: "Cloud",
	},
	"selectel": {
		Name: "Selectel",
		Type: "Cloud",
	},
	"softlayer": {
		Name: "SoftLayer",
		Type: "Cloud",
	},
	"telefonicaopencloud": {
		Name: "TelefonicaOpenCloud",
		Type: "Cloud",
	},
	"tencentcloud": {
		Name: "TencentCloud",
		Type: "Cloud",
	},
	"triton": {
		Name: "Triton",
		Type: "Cloud",
	},
	"ucloud": {
		Name: "UCloud",
		Type: "Cloud",
	},
	"yandex": {
		Name: "Yandex",
		Type: "Cloud",
	},
	"oneandone": {
		Name: "1&1",
		Type: "Cloud",
	},
}

var misc = map[string]TerraformProvider{
	"acme":     {Name: "ACME", Type: "Misc"},
	"archive":  {Name: "Archive", Type: "Misc"},
	"cobbler":  {Name: "Cobbler", Type: "Misc"},
	"external": {Name: "External", Type: "Misc"},
	"ignition": {Name: "Ignition", Type: "Misc"},
	"local":    {Name: "Local", Type: "Misc"},
	"netlify":  {Name: "Netlify", Type: "Misc"},
	"null":     {Name: "Null", Type: "Misc"},
	"random":   {Name: "Random", Type: "Misc"},
	"template": {Name: "Template", Type: "Misc"},
	"tls":      {Name: "TLS", Type: "Misc"},
}

var monitoring = map[string]TerraformProvider{
	"circonus":     {Name: "Circonus", Type: "Monitoring"},
	"datadog":      {Name: "Datadog", Type: "Monitoring"},
	"dyn":          {Name: "Dyn", Type: "Monitoring"},
	"grafana":      {Name: "Grafana", Type: "Monitoring"},
	"icinga2":      {Name: "Icinga2", Type: "Monitoring"},
	"librato":      {Name: "Librato", Type: "Monitoring"},
	"logentries":   {Name: "Logentries", Type: "Monitoring"},
	"logicmonitor": {Name: "LogicMonitor", Type: "Monitoring"},
	"newrelic":     {Name: "New Relic", Type: "Monitoring"},
	"opsgenie":     {Name: "OpsGenie", Type: "Monitoring"},
	"pagerduty":    {Name: "PagerDuty", Type: "Monitoring"},
	"runscope":     {Name: "Runscope", Type: "Monitoring"},
	"statuscake":   {Name: "StatusCake", Type: "Monitoring"},
}

var network = map[string]TerraformProvider{
	"cloudflare": {Name: "Cloudflare", Type: "Network"},
	"ciscoasa":   {Name: "Cisco ASA", Type: "Network"},
	"dns":        {Name: "DNS", Type: "Network"},
	"dnsimple":   {Name: "DNSimple", Type: "Network"},
	"dme":        {Name: "DNSMadeEasy", Type: "Network"},
	"bigip":      {Name: "F5 BIG-IP", Type: "Network"},
	"fortios":    {Name: "FortiOS", Type: "Network"},
	"http":       {Name: "HTTP", Type: "Network"},
	"ns1":        {Name: "NS1", Type: "Network"},
	"panos":      {Name: "Palo Alto Networks", Type: "Network"},
	"powerdns":   {Name: "PowerDNS", Type: "Network"},
	"ultradns":   {Name: "UltraDNS", Type: "Network"},
}

var vcs = map[string]TerraformProvider{
	"gitlab":    {Name: "GitLab", Type: "VCS"},
	"github":    {Name: "GitHub", Type: "VCS"},
	"bitbucket": {Name: "Bitbucket", Type: "VCS"},
}
