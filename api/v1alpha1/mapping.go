package v1alpha1

import (
	"fmt"

	istio "istio.io/api/networking/v1alpha3"
	istioclient "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashfunc/debotops/pkg/k8s"
)

func (in *Mapping) NewVirtualService(application *Application) (*istioclient.VirtualService, error) {
	revision, err := Revision(&in.Spec)

	if err != nil {
		return nil, err
	}

	host := fmt.Sprintf(
		"%s.%s.svc.cluster.local.",
		application.Name,
		application.Namespace,
	)

	virtualService := &istioclient.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:            in.Name,
			Namespace:       in.Namespace,
			OwnerReferences: k8s.NewOwnerReferences(in),
			Annotations: map[string]string{
				RevisionKey(): revision,
			},
		},
		Spec: istio.VirtualService{
			Hosts: in.Spec.Hosts,
			Http: []*istio.HTTPRoute{
				{
					Name: application.Name,
					Route: []*istio.HTTPRouteDestination{
						{
							Destination: &istio.Destination{
								Host: host,
								Port: &istio.PortSelector{
									Number: in.Spec.Application.Port,
								},
							},
						},
					},
				},
			},
		},
	}

	return virtualService, nil
}
