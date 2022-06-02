package controllers

import (
	"context"

	istioclient "istio.io/client-go/pkg/apis/networking/v1alpha3"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	debotops "github.com/hashfunc/debotops/api/v1alpha1"
	"github.com/hashfunc/debotops/pkg/core"
	"github.com/hashfunc/debotops/pkg/k8s"
)

// ListenerReconciler reconciles a Listener object
type ListenerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=listeners,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=listeners/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=debotops.hashfunc.io,resources=listeners/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create

//+kubebuilder:rbac:groups=networking.istio.io,resources=gateways,verbs=get;list;watch;create;update

func (r *ListenerReconciler) GetListener(ctx context.Context, request ctrl.Request) (*debotops.Listener, error) {
	listener := &debotops.Listener{}

	if err := r.Get(ctx, request.NamespacedName, listener); err != nil {
		return nil, err
	}

	return listener, nil
}

func (r *ListenerReconciler) GetGateway(ctx context.Context, request ctrl.Request) (*istioclient.Gateway, error) {
	gateway := &istioclient.Gateway{}

	if err := r.Get(ctx, request.NamespacedName, gateway); err != nil {
		return nil, err
	}

	return gateway, nil
}

func IgnoreIsNotFound(err error) error {
	if k8serrors.IsNotFound(err) {
		return nil
	}
	return err
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ListenerReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	namespaceForSecret := core.GetNamespaceForSecret()

	_, err := k8s.GetDeBotOpsSecret(r.Client, ctx, namespaceForSecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			rootSecrets, err := core.NewDeBotOpsSecret(namespaceForSecret)
			if err != nil {
				return ctrl.Result{}, err
			}

			for _, secret := range rootSecrets {
				if err := r.Create(ctx, secret); err != nil {
					return ctrl.Result{}, err
				}
			}

			return ctrl.Result{Requeue: true}, nil
		}

		return ctrl.Result{}, err
	}

	listener, err := r.GetListener(ctx, request)

	if err != nil {
		return ctrl.Result{}, IgnoreIsNotFound(err)
	}

	gateway, err := r.GetGateway(ctx, request)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			newGateway, err := listener.NewGateway()

			if err != nil {
				return ctrl.Result{}, err
			}

			if err := r.Create(ctx, newGateway); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	revision, err := listener.Hash()

	if err != nil {
		return ctrl.Result{}, err
	}

	if gateway.Annotations[debotops.GetHashKey()] != revision {
		newGateway, err := listener.NewGateway()

		if err != nil {
			return ctrl.Result{}, err
		}

		newGateway.SetResourceVersion(gateway.GetResourceVersion())

		if err := r.Update(ctx, newGateway); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ListenerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&debotops.Listener{}).
		Owns(&istioclient.Gateway{}).
		Complete(r)
}
