package reconcile

import (
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile/task"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Playbook []task.Func

func Run(playbook Playbook) (ctrl.Result, error) {
	for _, taskFunc := range playbook {
		result, err := taskFunc()

		if err != nil || result.RequeueRequest {
			return RequeueAfter(result.RequeueDelay, err)
		}

		if result.CancelRequest {
			return DoNotRequeue()
		}
	}

	return DoNotRequeue()
}
