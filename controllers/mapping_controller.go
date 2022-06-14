package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *MappingReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MappingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&debotops.Mapping{}).
		Complete(r)
}
