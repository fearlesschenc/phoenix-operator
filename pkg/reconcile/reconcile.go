package reconcile

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

type TaskFunc func() (Result, error)

func Run(taskFuncs []TaskFunc) (ctrl.Result, error) {
	for _, taskFunc := range taskFuncs {
		taskResult, err := taskFunc()

		if err != nil || taskResult.RequeueRequest {
			return RequeueRequestAfter(taskResult.RequeueDelay, err)
		}

		if taskResult.CancelReconciliation {
			return DoNotRequeueRequest()
		}
	}

	return DoNotRequeueRequest()
}

func DoNotRequeueRequest() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func RequeueRequestOnErr(err error) (ctrl.Result, error) {
	// note: reconcile will auto requeue failed request
	return ctrl.Result{}, err
}

func RequeueRequestAfter(duration time.Duration, err error) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: duration}, err
}
