/*
Adapted from https://github.com/kubernetes/client-go/blob/master/examples/leader-election/main.go
*/

package election

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

// K8sLeaderElectionMgr manages leader election among cloud connectors by talking to the leader election side car.
type K8sLeaderElectionMgr struct {
	// The name of the election to run.
	name string
	// The namespace to run the election in.
	namespace string
	// Kube config to use.
	kubeConfig string
	// Duration in ms that a lease on a lock lasts.
	leaseDuration time.Duration
	// Duration in ms before a campaign renews their lease.
	renewDeadline time.Duration
	// Duration in ms before attempting to retry acquiring the lease.
	retryPeriod time.Duration
}

// NewK8sLeaderElectionMgr creates a K8sLeaderElectionMgr.
func NewK8sLeaderElectionMgr(electionNamespace string, expectedMaxSkewMS, renewDeadlineMS time.Duration, electionName string) (*K8sLeaderElectionMgr, error) {
	if electionNamespace == "" {
		return nil, errors.New("namespace must be specified for leader election")
	}
	// Might add more complex logic that necessitates errors, but for now not included.
	return &K8sLeaderElectionMgr{
		namespace:     electionNamespace,
		name:          electionName,
		leaseDuration: expectedMaxSkewMS + renewDeadlineMS,
		renewDeadline: renewDeadlineMS,
		retryPeriod:   renewDeadlineMS / 4,
		kubeConfig:    "",
	}, nil
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// runElection manages the election.
func (le *K8sLeaderElectionMgr) runElection(ctx context.Context, callback func(string), id string) {

	// leader election uses the Kubernetes API by writing to a
	// lock object, which can be a LeaseLock object (preferred),
	// a ConfigMap, or an Endpoints (deprecated) object.
	// Conflicting writes are detected and each client handles those actions
	// independently.
	config, err := buildConfig(le.kubeConfig)
	if err != nil {
		log.Fatal(err)
	}
	client := clientset.NewForConfigOrDie(config)

	// We use the EndpointsLock to remain compatible with older versions of K8s, rather than the LeaseLock (k8s >=v1.14) as suggested.
	lock := &resourcelock.EndpointsLock{
		EndpointsMeta: metav1.ObjectMeta{
			Name:      le.name,
			Namespace: le.namespace,
		},
		Client: client.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	// start the leader election code loop
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   le.leaseDuration * time.Millisecond,
		RenewDeadline:   le.renewDeadline * time.Millisecond,
		RetryPeriod:     le.retryPeriod * time.Millisecond,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
			},
			OnStoppedLeading: func() {
				// we can do cleanup here
				log.WithField("id", id).Error("leadership lost")
				panic(fmt.Errorf("leadership lost"))
			},
			OnNewLeader: callback,
		},
	})
}

// Campaign blocks until this node is the leader.
func (le *K8sLeaderElectionMgr) Campaign(ctx context.Context) error {
	ch := make(chan string)
	fn := func(str string) {
		ch <- str
	}

	// Get the hostname of this pod.
	podHostname, err := os.Hostname()
	if err != nil {
		return err
	}

	go le.runElection(ctx, fn, podHostname)

	// Loop until we are elected.
	for newName := range ch {
		// If the leader is the pod then we exit the loop.
		if newName == podHostname {
			log.Infof("Pod '%s' is now leader.", podHostname)
			return nil
		}
		log.Infof("Pod '%s' not yet elected leader. Current leader is '%s'.", podHostname, newName)
	}
	return fmt.Errorf("Channel closed before '%s' elected leader", podHostname)
}
