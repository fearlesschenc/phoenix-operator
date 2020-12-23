package reconcile

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

type TaskFunc func() (Result, error)

func Run(taskFuncs []TaskFunc) (ctrl.Result, error) {
	for _, taskFunc := range taskFuncs {
		taskResult, err := taskFunc()

		if err != nil || taskResult.RequeueRequest {
			return ctrl.Result{RequeueAfter: taskResult.RequeueDelay}, err
		}

		if taskResult.CancelReconciliation {
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

//func DoNotRequeue() (ctrl.Result, error) {
//	return ctrl.Result{}, nil
//}

//
//func RequeueOnErr(err error) (ctrl.Result, error) {
//	// note: reconcile will auto requeue failed request
//	return ctrl.Result{}, err
//}

//func RequeueAfter(duration time.Duration, err error) (ctrl.Result, error) {
//	return ctrl.Result{RequeueAfter: duration}, err
//}
