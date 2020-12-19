package reconcile

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

// common used reconcile result

func DoNotRequeue() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func RequeueOnErr(err error) (ctrl.Result, error) {
	// note: reconcile will auto requeue failed request
	return ctrl.Result{}, err
}

func RequeueAfter(duration time.Duration, err error) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: duration}, err
}
