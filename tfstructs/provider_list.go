package tfstructs

type TerraformProvider struct {
	Name string
	Type string
}

var OfficialProviders = map[string]TerraformProvider{
	"aws": TerraformProvider{
		Name: "AWS",
		Type: "Major Cloud",
	},
	"azurerm": TerraformProvider{
		Name: "AzureRM",
		Type: "Major Cloud",
	},
	"arukas": TerraformProvider{
		Name: "Arukas",
		Type: "Cloud",
	},
	"clc": TerraformProvider{
		Name: "CenturyLinkCloud",
		Type: "Cloud",
	},
}

// Finish the rest
//	brightbox](/docs/providers/brightbox/index.html)
//	centurylinkcloud](/docs/providers/clc/index.html)
//	CloudScale.ch](/docs/providers/cloudscale/index.html)
//	CloudStack](/docs/providers/cloudstack/index.html)
//	DigitalOcean](/docs/providers/do/index.html)
//	Fastly](/docs/providers/fastly/index.html)
//	FlexibleEngine](/docs/providers/flexibleengine/index.html)
//	Gridscale](/docs/providers/gridscale/index.html)
//	Hedvig](/docs/providers/hedvig/index.html)
//	Heroku](/docs/providers/heroku/index.html)
//	Hetzner Cloud](/docs/providers/hcloud/index.html)
//	HuaweiCloud](/docs/providers/huaweicloud/index.html)
//	JDCloud](/docs/providers/jdcloud/index.html)
//	Linode](/docs/providers/linode/index.html)
//	Naver Cloud](/docs/providers/ncloud/index.html)
//	Nutanix](/docs/providers/nutanix/index.html)
//	OpenStack](/docs/providers/openstack/index.html)
//	OpenTelekomCloud](/docs/providers/opentelekomcloud/index.html)
//	OVH](/docs/providers/ovh/index.html)
//	Packet](/docs/providers/packet/index.html)
//	ProfitBricks](/docs/providers/profitbricks/index.html)
//	Scaleway](/docs/providers/scaleway/index.html)
//	Skytap](/docs/providers/skytap/index.html)
//	Selectel](/docs/providers/selectel/index.html)
//	SoftLayer](/docs/providers/softlayer/index.html)
//	TelefonicaOpenCloud](/docs/providers/telefonicaopencloud/index.html)
//	TencentCloud](/docs/providers/tencentcloud/index.html)
//	Triton](/docs/providers/triton/index.html)
//	UCloud](/docs/providers/ucloud/index.html)
//	Yandex](/docs/providers/yandex/index.html)
//	1&1](/docs/providers/oneandone/index.html)
