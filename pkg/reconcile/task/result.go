package task

import "time"

type Result struct {
	RequeueDelay   time.Duration
	RequeueRequest bool
	CancelRequest  bool
}

func continueOperationResult() Result {
	return Result{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  false,
	}
}

func ContinueProcessing() (result Result, err error) {
	result = continueOperationResult()
	return
}

func stopOperationResult() Result {
	return Result{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  true,
	}
}

func StopProcessing() (result Result, err error) {
	result = stopOperationResult()
	return
}

func RequeueWithError(errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:   0,
		RequeueRequest: true,
		CancelRequest:  false,
	}
	err = errIn
	return
}

func RequeueOnErrorOrStop(errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  true,
	}
	err = errIn
	return
}

func RequeueOnErrorOrContinue(errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  false,
	}
	err = errIn
	return
}

func RequeueAfter(delay time.Duration, errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:   delay,
		RequeueRequest: true,
		CancelRequest:  false,
	}
	err = errIn
	return
}
