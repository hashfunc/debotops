package v1alpha1

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"

	istio "istio.io/api/networking/v1alpha3"
	istioclient "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/hashfunc/debotops/pkg/k8s"
)

func GetHashKey() string {
	return fmt.Sprintf("%s/revision", GroupVersion.Group)
}

func (in *Listener) Hash() (string, error) {
	data, err := json.Marshal(in.Spec)

	if err != nil {
		return "", err
	}

	hash64a := fnv.New64()
	_, err = hash64a.Write(data)

	if err != nil {
		return "", err
	}

	hashString := strconv.FormatUint(hash64a.Sum64(), 10)

	return rand.SafeEncodeString(hashString), nil
}

func (in *Listener) NewGateway() (*istioclient.Gateway, error) {
	hash, err := in.Hash()

	if err != nil {
		return nil, err
	}

	gateway := &istioclient.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:            in.Name,
			Namespace:       in.Namespace,
			OwnerReferences: k8s.NewOwnerReferences(in),
			Annotations: map[string]string{
				GetHashKey(): hash,
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
		mode := istio.ServerTLSSettings_TLSmode(
			istio.ServerTLSSettings_TLSmode_value[bind.TLS.Mode],
		)
		servers[index] = &istio.Server{
			Hosts: bind.Hosts,
			Port: &istio.Port{
				Name:     bind.Name,
				Number:   bind.Port,
				Protocol: bind.Protocol,
			},
			Tls: &istio.ServerTLSSettings{
				Mode:           mode,
				CredentialName: bind.TLS.Credential,
			},
		}
	}

	return servers
}
