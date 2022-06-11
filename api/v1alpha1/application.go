package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashfunc/debotops/pkg/k8s"
)

func (in *Application) NewDeployment() (*appsv1.Deployment, error) {
	revision, err := Revision(&in.Spec)

	if err != nil {
		return nil, err
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            in.Name,
			Namespace:       in.Namespace,
			OwnerReferences: k8s.NewOwnerReferences(in),
			Annotations: map[string]string{
				RevisionKey(): revision,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &in.Spec.Container.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"debotops/application": in.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"debotops/application": in.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*in.toContainer(),
					},
				},
			},
		},
	}

	if in.Spec.Option.Proxy.Enable {
		labels := deployment.Spec.Template.ObjectMeta.Labels
		labels["sidecar.istio.io/inject"] = "true"
	}

	return deployment, nil
}

func (in *Application) toContainer() *corev1.Container {
	spec := in.Spec.Container

	container := &corev1.Container{
		Name:    spec.Name,
		Image:   spec.Image,
		Ports:   spec.Ports,
		Env:     spec.Environments,
		Command: spec.Command,
		Args:    spec.Args,
	}

	if spec.Resource != nil {
		container.Resources = *spec.Resource
	}

	if spec.Health != nil {
		container.StartupProbe = spec.Health.Startup
		container.ReadinessProbe = spec.Health.Ready
		container.LivenessProbe = spec.Health.Live
	}

	return container
}
