package rundmc

import (
	"github.com/cloudfoundry-incubator/goci"
	"github.com/cloudfoundry-incubator/goci/specs"
	"github.com/cloudfoundry-incubator/guardian/gardener"
)

type BundleTemplate struct {
	*goci.Bndl
}

func (base BundleTemplate) Bundle(spec gardener.DesiredContainerSpec) *goci.Bndl {
	return base.
		WithRootFS(spec.RootFSPath).
		WithNamespace(specs.Namespace{Type: specs.NetworkNamespace, Path: spec.NetworkPath})
}
