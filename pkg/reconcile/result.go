package reconcile

import "time"

type Result struct {
	RequeueDelay         time.Duration
	RequeueRequest       bool
	CancelReconciliation bool
}

func continueResult() Result {
	return Result{
		RequeueDelay:         0,
		RequeueRequest:       false,
		CancelReconciliation: false,
	}
}

// Continue continue the execution of the whole reconciliation
func Continue() (result Result, err error) {
	result = continueResult()
	return
}

func stopResult() Result {
	return Result{
		RequeueDelay:         0,
		RequeueRequest:       false,
		CancelReconciliation: true,
	}
}

// Stop stop the whole reconciliation
func Stop() (result Result, err error) {
	result = stopResult()
	return
}

// RequeueWithError will always requeue the request whether
// errIn is nil or not.
func RequeueWithError(errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:         0,
		RequeueRequest:       true,
		CancelReconciliation: false,
	}
	err = errIn
	return
}

// RequeueOnErrorOrStop will requeue request if errIn is not nil.
// When errIn is nil, will stop the whole reconciliation.
func RequeueOnErrorOrStop(errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:         0,
		RequeueRequest:       false,
		CancelReconciliation: true,
	}
	err = errIn
	return
}

// RequeueOnErrorOrContinue will requeue request if the errIn is not nil.
// When errIn is nil, will continue the whole reconciliation.
func RequeueOnErrorOrContinue(errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:         0,
		RequeueRequest:       false,
		CancelReconciliation: false,
	}
	err = errIn
	return
}

// RequeueAfter will requeue request after delay.
func RequeueAfter(delay time.Duration, errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:         delay,
		RequeueRequest:       true,
		CancelReconciliation: false,
	}
	err = errIn
	return
}
