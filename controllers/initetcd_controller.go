/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"cmit.com/crd/domain-config/etcd"
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	domainconfigv1alpha1 "cmit.com/crd/domain-config/api/v1alpha1"
)

// InitEtcdReconciler reconciles a InitEtcd object
type InitEtcdReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=domain-config.cmit.com,resources=initetcds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=domain-config.cmit.com,resources=initetcds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=domain-config.cmit.com,resources=initetcds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the InitEtcd object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *InitEtcdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("initetcd", req.NamespacedName)

	// your logic here
	reqLogger := r.Log.WithValues("dc", req.NamespacedName)
	reqLogger.Info("----- Reconciling InitConfig -----")
	// Create DnsConfig instance
	atInstance := &domainconfigv1alpha1.InitEtcd{}

	// Try to get cloud native dc instance.
	err := r.Get(ctx, req.NamespacedName, atInstance)
	if err != nil {
		// Request object not found.
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Other error.
		return ctrl.Result{}, err
	}
	//etcd init
	initAdds := ""
	certDir := ""

	if atInstance.Spec.InitConf != "" {
		initAdds = atInstance.Spec.InitConf
		certDir = atInstance.Spec.CertDir
	}
	//decAdds, _ := base64.StdEncoding.DecodeString(initAdds)
	etcdAdds := strings.Split(string(initAdds), ",")
	err = etcd.InitConn(etcdAdds, 5*time.Second, certDir)
	if err != nil {
		// Request object not found.
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Other error.
		return ctrl.Result{}, err
	}
	reqLogger.Info("etcd conn  ", "host: ", etcdAdds)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InitEtcdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&domainconfigv1alpha1.InitEtcd{}).
		Complete(r)
}
