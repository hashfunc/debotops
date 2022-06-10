package v1alpha1

import (
	istio "istio.io/api/networking/v1alpha3"
	istioclient "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashfunc/debotops/pkg/k8s"
)

func (in *Listener) NewGateway() (*istioclient.Gateway, error) {
	hash, err := Revision(&in.Spec)

	if err != nil {
		return nil, err
	}

	gateway := &istioclient.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:            in.Name,
			Namespace:       in.Namespace,
			OwnerReferences: k8s.NewOwnerReferences(in),
			Annotations: map[string]string{
				RevisionKey(): hash,
			},
		},
		Spec: istio.Gateway{
			Selector: in.Spec.Gateway.Selector,
			Servers:  in.NewGatewayServers(),
		},
	}

	return gateway, nil
}

func (in *Listener) NewGatewayServers() []*istio.Server {
	servers := make([]*istio.Server, len(in.Spec.Bind))

	for index, bind := range in.Spec.Bind {
		servers[index] = &istio.Server{
			Hosts: bind.Hosts,
			Port: &istio.Port{
				Name:     bind.Name,
				Number:   bind.Port,
				Protocol: bind.Protocol,
			},
		}

		if bind.TLS != nil {
			mode := istio.ServerTLSSettings_TLSmode(
				istio.ServerTLSSettings_TLSmode_value[bind.TLS.Mode],
			)
			servers[index].Tls = &istio.ServerTLSSettings{
				Mode:           mode,
				CredentialName: bind.TLS.Credential,
			}
		}
	}

	return servers
}
