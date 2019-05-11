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
	"alicloud": TerraformProvider{
		Name: "AliCloud",
		Type: "Major Cloud",
	},
	"aws": TerraformProvider{
		Name: "AWS",
		Type: "Major Cloud",
	},
	"azurerm": TerraformProvider{
		Name: "AzureRM",
		Type: "Major Cloud",
	},
	"azurestack": TerraformProvider{
		Name: "Azure Stack",
		Type: "Major Cloud",
	},
	"google": TerraformProvider{
		Name: "Google Cloud Platform",
		Type: "Major Cloud",
	},
	"oci": TerraformProvider{
		Name: "Oracle Cloud Infrastructure",
		Type: "Major Cloud",
	},
	"oraclepaas": TerraformProvider{
		Name: "Oracle Cloud Platform",
		Type: "Major Cloud",
	},
	"opc": TerraformProvider{
		Name: "Oracle Public Cloud",
		Type: "Major Cloud",
	},
	"nsxt": TerraformProvider{
		Name: "VMware NSX-T",
		Type: "Major Cloud",
	},
	"vcd": TerraformProvider{
		Name: "vCloud Director",
		Type: "Major Cloud",
	},
	"vra7": TerraformProvider{
		Name: "VMware vRA7",
		Type: "Major Cloud",
	},
	"vsphere": TerraformProvider{
		Name: "VMware vSphere",
		Type: "Major Cloud",
	},
}

var database = map[string]TerraformProvider{
	"influxdb": TerraformProvider{
		Name: "InfluxDB",
		Type: "Database",
	},
	"mysql": TerraformProvider{
		Name: "MySQL",
		Type: "Database",
	},
	"postgresql": TerraformProvider{
		Name: "PostgreSQL",
		Type: "Database",
	},
}

var infra = map[string]TerraformProvider{
	"chef": TerraformProvider{
		Name: "Chef",
		Type: "Infra",
	},
	"consul": TerraformProvider{
		Name: "Consul",
		Type: "Infra",
	},
	"docker": TerraformProvider{
		Name: "Docker",
		Type: "Infra",
	},
	"helm": TerraformProvider{
		Name: "Helm",
		Type: "Infra",
	},
	"kubernetes": TerraformProvider{
		Name: "Kubernetes",
		Type: "Infra",
	},
	"mailgun": TerraformProvider{
		Name: "Mailgun",
		Type: "Infra",
	},
	"nomad": TerraformProvider{
		Name: "Nomad",
		Type: "Infra",
	},
	"rabbitmq": TerraformProvider{
		Name: "RabbitMQ",
		Type: "Infra",
	},
	"rancher": TerraformProvider{
		Name: "Rancher",
		Type: "Infra",
	},
	"rightscale": TerraformProvider{
		Name: "RightScale",
		Type: "Infra",
	},
	"rundeck": TerraformProvider{
		Name: "Rundeck",
		Type: "Infra",
	},
	"spotinst": TerraformProvider{
		Name: "Spotinst",
		Type: "Infra",
	},
	"terraform": TerraformProvider{
		Name: "Terraform",
		Type: "Infra",
	},
	"tfe": TerraformProvider{
		Name: "Terraform Enterprise",
		Type: "Infra",
	},
	"vault": TerraformProvider{
		Name: "Vault",
		Type: "Infra",
	},
}

var cloud = map[string]TerraformProvider{
	"arukas": TerraformProvider{
		Name: "Arukas",
		Type: "Cloud",
	},
	"brightbox": TerraformProvider{
		Name: "Brightbox",
		Type: "Cloud",
	},
	"clc": TerraformProvider{
		Name: "CenturyLinkCloud",
		Type: "Cloud",
	},
	"cloudscale": TerraformProvider{
		Name: "CloudScale.ch",
		Type: "Cloud",
	},
	"cloudstack": TerraformProvider{
		Name: "CloudStack",
		Type: "Cloud",
	},
	"digitalocean": TerraformProvider{
		Name: "DigitalOcean",
		Type: "Cloud",
	},
	"fastly": TerraformProvider{
		Name: "Fastly",
		Type: "Cloud",
	},
	"flexibleengine": TerraformProvider{
		Name: "FlexibleEngine",
		Type: "Cloud",
	},
	"gridscale": TerraformProvider{
		Name: "Gridscale",
		Type: "Cloud",
	},
	"hedvig": TerraformProvider{
		Name: "Hedvig",
		Type: "Cloud",
	},
	"heroku": TerraformProvider{
		Name: "Heroku",
		Type: "Cloud",
	},
	"hcloud": TerraformProvider{
		Name: "Hetzner Cloud",
		Type: "Cloud",
	},
	"huaweicloud": TerraformProvider{
		Name: "HuaweiCloud",
		Type: "Cloud",
	},
	"jdcloud": TerraformProvider{
		Name: "JDCloud",
		Type: "Cloud",
	},
	"linode": TerraformProvider{
		Name: "Linode",
		Type: "Cloud",
	},
	"ncloud": TerraformProvider{
		Name: "Naver Cloud",
		Type: "Cloud",
	},
	"nutanix": TerraformProvider{
		Name: "Nutanix",
		Type: "Cloud",
	},
	"openstack": TerraformProvider{
		Name: "OpenStack",
		Type: "Cloud",
	},
	"opentelekomcloud": TerraformProvider{
		Name: "OpenTelekomCloud",
		Type: "Cloud",
	},
	"ovh": TerraformProvider{
		Name: "OVH",
		Type: "Cloud",
	},
	"packet": TerraformProvider{
		Name: "Packet",
		Type: "Cloud",
	},
	"profitbricks": TerraformProvider{
		Name: "ProfitBricks",
		Type: "Cloud",
	},
	"scaleway": TerraformProvider{
		Name: "Scaleway",
		Type: "Cloud",
	},
	"skytap": TerraformProvider{
		Name: "Skytap",
		Type: "Cloud",
	},
	"selectel": TerraformProvider{
		Name: "Selectel",
		Type: "Cloud",
	},
	"softlayer": TerraformProvider{
		Name: "SoftLayer",
		Type: "Cloud",
	},
	"telefonicaopencloud": TerraformProvider{
		Name: "TelefonicaOpenCloud",
		Type: "Cloud",
	},
	"tencentcloud": TerraformProvider{
		Name: "TencentCloud",
		Type: "Cloud",
	},
	"triton": TerraformProvider{
		Name: "Triton",
		Type: "Cloud",
	},
	"ucloud": TerraformProvider{
		Name: "UCloud",
		Type: "Cloud",
	},
	"yandex": TerraformProvider{
		Name: "Yandex",
		Type: "Cloud",
	},
	"oneandone": TerraformProvider{
		Name: "1&1",
		Type: "Cloud",
	},
}

var misc = map[string]TerraformProvider{
	"acme":     TerraformProvider{Name: "ACME", Type: "Misc"},
	"archive":  TerraformProvider{Name: "Archive", Type: "Misc"},
	"cobbler":  TerraformProvider{Name: "Cobbler", Type: "Misc"},
	"external": TerraformProvider{Name: "External", Type: "Misc"},
	"ignition": TerraformProvider{Name: "Ignition", Type: "Misc"},
	"local":    TerraformProvider{Name: "Local", Type: "Misc"},
	"netlify":  TerraformProvider{Name: "Netlify", Type: "Misc"},
	"null":     TerraformProvider{Name: "Null", Type: "Misc"},
	"random":   TerraformProvider{Name: "Random", Type: "Misc"},
	"template": TerraformProvider{Name: "Template", Type: "Misc"},
	"tls":      TerraformProvider{Name: "TLS", Type: "Misc"},
}

var monitoring = map[string]TerraformProvider{
	"circonus":     TerraformProvider{Name: "Circonus", Type: "Monitoring"},
	"datadog":      TerraformProvider{Name: "Datadog", Type: "Monitoring"},
	"dyn":          TerraformProvider{Name: "Dyn", Type: "Monitoring"},
	"grafana":      TerraformProvider{Name: "Grafana", Type: "Monitoring"},
	"icinga2":      TerraformProvider{Name: "Icinga2", Type: "Monitoring"},
	"librato":      TerraformProvider{Name: "Librato", Type: "Monitoring"},
	"logentries":   TerraformProvider{Name: "Logentries", Type: "Monitoring"},
	"logicmonitor": TerraformProvider{Name: "LogicMonitor", Type: "Monitoring"},
	"newrelic":     TerraformProvider{Name: "New Relic", Type: "Monitoring"},
	"opsgenie":     TerraformProvider{Name: "OpsGenie", Type: "Monitoring"},
	"pagerduty":    TerraformProvider{Name: "PagerDuty", Type: "Monitoring"},
	"runscope":     TerraformProvider{Name: "Runscope", Type: "Monitoring"},
	"statuscake":   TerraformProvider{Name: "StatusCake", Type: "Monitoring"},
}

var network = map[string]TerraformProvider{
	"cloudflare": TerraformProvider{Name: "Cloudflare", Type: "Network"},
	"ciscoasa":   TerraformProvider{Name: "Cisco ASA", Type: "Network"},
	"dns":        TerraformProvider{Name: "DNS", Type: "Network"},
	"dnsimple":   TerraformProvider{Name: "DNSimple", Type: "Network"},
	"dme":        TerraformProvider{Name: "DNSMadeEasy", Type: "Network"},
	"bigip":      TerraformProvider{Name: "F5 BIG-IP", Type: "Network"},
	"fortios":    TerraformProvider{Name: "FortiOS", Type: "Network"},
	"http":       TerraformProvider{Name: "HTTP", Type: "Network"},
	"ns1":        TerraformProvider{Name: "NS1", Type: "Network"},
	"panos":      TerraformProvider{Name: "Palo Alto Networks", Type: "Network"},
	"powerdns":   TerraformProvider{Name: "PowerDNS", Type: "Network"},
	"ultradns":   TerraformProvider{Name: "UltraDNS", Type: "Network"},
}

var vcs = map[string]TerraformProvider{
	"gitlab":    TerraformProvider{Name: "GitLab", Type: "VCS"},
	"github":    TerraformProvider{Name: "GitHub", Type: "VCS"},
	"bitbucket": TerraformProvider{Name: "Bitbucket", Type: "VCS"},
}
