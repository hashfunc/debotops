package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	debotops "github.com/hashfunc/debotops/api/v1alpha1"
)

// ApplicationReconciler reconciles an Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=applications/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch

func (r *ApplicationReconciler) GetApplication(ctx context.Context, request ctrl.Request) (*debotops.Application, error) {
	application := &debotops.Application{}

	if err := r.Get(ctx, request.NamespacedName, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (r *ApplicationReconciler) GetDeployment(ctx context.Context, request ctrl.Request) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}

	if err := r.Get(ctx, request.NamespacedName, deployment); err != nil {
		return nil, err
	}

	return deployment, nil
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ApplicationReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	application, err := r.GetApplication(ctx, request)

	if err != nil {
		return ctrl.Result{}, IgnoreIsNotFound(err)
	}

	deployment, err := r.GetDeployment(ctx, request)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			newDeployment, err := application.NewDeployment()

			if err != nil {
				return ctrl.Result{}, err
			}

			if err := r.Create(ctx, newDeployment); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	revision, err := debotops.Revision(&application.Spec)

	if err != nil {
		return ctrl.Result{}, err
	}

	if deployment.Annotations[debotops.RevisionKey()] != revision {
		newDeployment, err := application.NewDeployment()

		if err != nil {
			return ctrl.Result{}, err
		}

		newDeployment.SetResourceVersion(deployment.GetResourceVersion())

		if err := r.Update(ctx, newDeployment); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&debotops.Application{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
