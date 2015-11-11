package rundmc

import (
	"github.com/cloudfoundry-incubator/goci"
	"github.com/cloudfoundry-incubator/goci/specs"
	"github.com/cloudfoundry-incubator/guardian/gardener"
)

type BundleTemplate struct {
	*goci.Bndl
}

var defaultCapabilities = []string{
	"CAP_CHOWN",
	"CAP_DAC_OVERRIDE",
	"CAP_FSETID",
	"CAP_FOWNER",
	"CAP_MKNOD",
	"CAP_NET_RAW",
	"CAP_SETGID",
	"CAP_SETUID",
	"CAP_SETFCAP",
	"CAP_SETPCAP",
	"CAP_NET_BIND_SERVICE",
	"CAP_SYS_CHROOT",
	"CAP_KILL",
	"CAP_AUDIT_WRITE",
}

func (base BundleTemplate) Bundle(spec gardener.DesiredContainerSpec) *goci.Bndl {
	bndl := base.WithNamespace(specs.Namespace{Type: specs.NetworkNamespace, Path: spec.NetworkPath})

	if !spec.Privileged {
		bndl = bndl.WithCapabilities(defaultCapabilities)
	}

	return bndl
}
