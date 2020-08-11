package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/operator-framework/operator-sdk/pkg/ready"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/controller"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var log = logf.Log.WithName("cmd")

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("operator-sdk Version: %v", sdkVersion.Version))
}

// Main program of the operator. This instantiates a new manager which
// registers all custom resource definitions under pkg/apis/... and starts
// all controllers under pkg/controllers/...
func main() {
	flag.Parse()

	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	logf.SetLogger(logf.ZapLogger(false))

	printVersion()

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Error(err, "failed to get watch namespace")
		os.Exit(1)
	}
	if namespace == "" {
		log.Info("Running as cluster-scoped operator.")
	} else {
		log.Error(fmt.Errorf("Running as namespace-scoped operator is not allowed."), "Stopping...", "namespace", namespace)
		os.Exit(1)
	}

	// make sure that OPERATOR_NAME environment variable is set. we need that one for separating dev and prod stages
	operatorName, err := k8sutil.GetOperatorName()
	if err != nil {
		log.Error(err, "failed to get operator name")
		os.Exit(1)
	}
	log.Info("Scoped to namespace.", "OperatorNamespace", operatorName)

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// We deactivate the leader election at the moment as we are running only in one replica.
	// The only time when there are multiple operators in the cluster is at the time of update, but here this
	// leader election mechanism is hindering the update: it permits the second (updated) operator to finish start,
	// since the first (older) operator still holds the lock. but the first operator could only be stopped when
	// second operator has started successfully -> deadlock.
	// as soon as we talk about HA-setup of operator, we have to revise this again.
	//	// Become the leader before proceeding
	//	leader.Become(context.TODO(), "iot-operator-lock")

	r := ready.NewFileReady()
	err = r.Set()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	defer r.Unset()

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{Namespace: namespace})
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "manager exited non-zero")
		os.Exit(1)
	}
}
