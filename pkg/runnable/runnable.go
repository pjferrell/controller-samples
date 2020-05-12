package runnable

import (
	"time"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// controller-runtime allows custom controllers to be registered and started by
// the controller-runtime.Manager if they implement the manager.Runnable interface
// additionally, if they implement the manager.LeaderElectionRunnable interface
// the controller start will be delayed until leader-election lock is attained.

type IntervalLogger struct {
	Log logr.Logger

	RequireLeaderElection bool
	Name                  string
	Interval              time.Duration
}

// Start implments the manager.Runnable interface
func (i *IntervalLogger) Start(stopCh <-chan struct{}) error {
	go func() {
		i.Log.Info("starting", "controllerName", i.Name, "RequireLeaderElection", i.RequireLeaderElection)
		for {
			select {
			case <-time.After(i.Interval):
				i.Log.Info("test message", "controllerName", i.Name, "RequireLeaderElection", i.RequireLeaderElection)
			case <-stopCh:
				i.Log.Info("stopping", "controllerName", i.Name, "RequireLeaderElection", i.RequireLeaderElection)
				return
			}
		}
	}()

	return nil
}

func (i *IntervalLogger) NeedLeaderElection() bool {
	return i.RequireLeaderElection
}

var (
	_ manager.Runnable               = &IntervalLogger{}
	_ manager.LeaderElectionRunnable = &IntervalLogger{}
)

func (i *IntervalLogger) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(i)
}
