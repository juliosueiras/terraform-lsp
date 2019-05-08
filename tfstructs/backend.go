package tfstructs

import (
	//"github.com/hashicorp/terraform/backend"
	//	backendAtlas "github.com/hashicorp/terraform/backend/atlas"
	//	backendLocal "github.com/hashicorp/terraform/backend/local"
	//	backendRemote "github.com/hashicorp/terraform/backend/remote"
	//	backendArtifactory "github.com/hashicorp/terraform/backend/remote-state/artifactory"
	//	backendAzure "github.com/hashicorp/terraform/backend/remote-state/azure"
	//	backendConsul "github.com/hashicorp/terraform/backend/remote-state/consul"
	"github.com/hashicorp/terraform/configs/configschema"
	//	//	backendEtcdv2 "github.com/hashicorp/terraform/backend/remote-state/etcdv2"
	//	//	backendEtcdv3 "github.com/hashicorp/terraform/backend/remote-state/etcdv3"
	//	backendGCS "github.com/hashicorp/terraform/backend/remote-state/gcs"
	//	backendHTTP "github.com/hashicorp/terraform/backend/remote-state/http"
	//	backendInmem "github.com/hashicorp/terraform/backend/remote-state/inmem"
	//	backendManta "github.com/hashicorp/terraform/backend/remote-state/manta"
	//	backendPg "github.com/hashicorp/terraform/backend/remote-state/pg"
	//	backendS3 "github.com/hashicorp/terraform/backend/remote-state/s3"
	//	backendSwift "github.com/hashicorp/terraform/backend/remote-state/swift"
)

var TerraformBackends = map[string]*configschema.Block{
	// Enhanced backends.
	//	"local":  func() backend.Backend { return backendLocal.New() }().ConfigSchema(),
	//	"remote": func() backend.Backend { return backendRemote.New(nil) }().ConfigSchema(),
	//	// Remote State backends.).ConfigSchema(),
	//	"artifactory": func() backend.Backend { return backendArtifactory.New() }().ConfigSchema(),
	//	"atlas":       func() backend.Backend { return backendAtlas.New() }().ConfigSchema(),
	//	"azurerm":     func() backend.Backend { return backendAzure.New() }().ConfigSchema(),
	//	"consul":      func() backend.Backend { return backendConsul.New() }().ConfigSchema(),
	//	//"etcd":        func() backend.Backend { return backendEtcdv2.New() }().ConfigSchema(),
	//	//"etcdv3":      func() backend.Backend { return backendEtcdv3.New() }().ConfigSchema(),
	//	"gcs":   func() backend.Backend { return backendGCS.New() }().ConfigSchema(),
	//	"http":  func() backend.Backend { return backendHTTP.New() }().ConfigSchema(),
	//	"inmem": func() backend.Backend { return backendInmem.New() }().ConfigSchema(),
	//	"manta": func() backend.Backend { return backendManta.New() }().ConfigSchema(),
	//	"pg":    func() backend.Backend { return backendPg.New() }().ConfigSchema(),
	//	"s3":    func() backend.Backend { return backendS3.New() }().ConfigSchema(),
	//"swift": func() backend.Backend { return backendSwift.New() }().ConfigSchema(),
}
