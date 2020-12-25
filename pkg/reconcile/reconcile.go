package reconcile

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

type ObjectState bool

var (
	ObjectChanged   ObjectState = true
	ObjectUnchanged ObjectState = false
)

type SubroutineFunc func() (Result, error)

func RunReconcileRoutine(subroutineFuncs []SubroutineFunc) (ctrl.Result, error) {
	for _, subroutineFunc := range subroutineFuncs {
		result, err := subroutineFunc()

		if err != nil || result.RequeueRequest {
			return RequeueRequestAfter(result.RequeueDelay, err)
		}

		if result.CancelReconciliation {
			return DoNotRequeueRequest()
		}
	}

	return DoNotRequeueRequest()
}

func RunSubRoutine(subroutineFuncs []SubroutineFunc) (Result, error) {
	for _, subroutineFunc := range subroutineFuncs {
		result, err := subroutineFunc()
		if err != nil {
			return result, err
		}

		if result.RequeueRequest || result.CancelReconciliation {
			return result, err
		}
	}

	return Continue()
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
