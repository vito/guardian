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
	bundle := base.WithNamespace(specs.Namespace{Type: specs.NetworkNamespace, Path: spec.NetworkPath})

	if spec.Privileged {
		bundle = bundle.WithCapabilities("CAP_SYS_ADMIN")
	} else {
		bundle = bundle.WithCapabilities("CAP_CHOWN", "CAP_DAC_OVERRIDE", "CAP_FSETID", "CAP_FOWNER", "CAP_MKNOD", "CAP_NET_RAW", "CAP_SETGID", "CAP_SETUID", "CAP_SETFCAP", "CAP_SETPCAP", "CAP_NET_BIND_SERVICE", "CAP_SYS_CHROOT", "CAP_KILL", "CAP_AUDIT_WRITE")
	}

	return bundle
}
