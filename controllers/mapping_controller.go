package controllers

import (
	"context"

	istioclient "istio.io/client-go/pkg/apis/networking/v1alpha3"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	debotops "github.com/hashfunc/debotops/api/v1alpha1"
)

// MappingReconciler reconciles a Mapping object
type MappingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=mappings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=mappings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=mappings/finalizers,verbs=update
//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=applications,verbs=get;list;watch;create;update

//+kubebuilder:rbac:groups=networking.istio.io,resources=virtualservices,verbs=get;list;watch;create;update

func (r *MappingReconciler) GetMapping(ctx context.Context, request ctrl.Request) (*debotops.Mapping, error) {
	mapping := &debotops.Mapping{}

	if err := r.Get(ctx, request.NamespacedName, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

func (r *MappingReconciler) GetApplication(ctx context.Context, mapping *debotops.Mapping) (*debotops.Application, error) {
	application := &debotops.Application{}

	namespacedName := types.NamespacedName{
		Namespace: mapping.Namespace,
		Name:      mapping.Spec.Application.Name,
	}

	if err := r.Get(ctx, namespacedName, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (r *MappingReconciler) GetVirtualService(ctx context.Context, request ctrl.Request) (*istioclient.VirtualService, error) {
	virtualService := &istioclient.VirtualService{}

	if err := r.Get(ctx, request.NamespacedName, virtualService); err != nil {
		return nil, err
	}

	return virtualService, nil
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *MappingReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	mapping, err := r.GetMapping(ctx, request)

	if err != nil {
		return ctrl.Result{}, IgnoreIsNotFound(err)
	}

	application, err := r.GetApplication(ctx, mapping)

	if err != nil {
		return ctrl.Result{}, err
	}

	virtualService, err := r.GetVirtualService(ctx, request)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			newVirtualService, err := mapping.NewVirtualService(application)

			if err != nil {
				return ctrl.Result{}, err
			}

			if err := r.Create(ctx, newVirtualService); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}

		return ctrl.Result{}, err
	}

	revision, err := debotops.Revision(&mapping.Spec)

	if err != nil {
		return ctrl.Result{}, err
	}

	if virtualService.Annotations[debotops.RevisionKey()] != revision {
		newVirtualService, err := mapping.NewVirtualService(application)

		if err != nil {
			return ctrl.Result{}, err
		}

		newVirtualService.SetResourceVersion(virtualService.GetResourceVersion())

		if err := r.Update(ctx, newVirtualService); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MappingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&debotops.Mapping{}).
		Owns(&debotops.Application{}).
		Owns(&istioclient.VirtualService{}).
		Complete(r)
}
