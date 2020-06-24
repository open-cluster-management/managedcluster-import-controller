// (c) Copyright IBM Corporation 2019, 2020. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// U.S. Government Users Restricted Rights - Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.
// Licensed Materials - Property of IBM
//
// Copyright (c) 2020 Red Hat, Inc.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	clusterv1 "github.com/open-cluster-management/api/cluster/v1"
	workv1 "github.com/open-cluster-management/api/work/v1"
	certificatesv1beta1 "k8s.io/api/certificates/v1beta1"

	ocinfrav1 "github.com/openshift/api/config/v1"
	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	kubemetrics "github.com/operator-framework/operator-sdk/pkg/kube-metrics"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.
	"k8s.io/client-go/rest"
	rbacv1 "k8s.io/kubernetes/pkg/apis/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	// "github.com/open-cluster-management/rcm-controller/pkg/apis"
	"github.com/open-cluster-management/rcm-controller/pkg/controller"
	// "github.com/open-cluster-management/rcm-controller/pkg/utils"
)

// Change below variables to serve metrics on different host or port.
var (
	metricsHost               = "0.0.0.0"
	metricsPort         int32 = 8383
	operatorMetricsPort int32 = 8686
)
var log = logf.Log.WithName("cmd")

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("Version of operator-sdk: %v", sdkVersion.Version))
}

func main() {
	// Add the zap logger flag set to the CLI. The flag set must
	// be added before calling pflag.Parse().
	pflag.CommandLine.AddFlagSet(zap.FlagSet())

	// Add flags registered by imported packages (e.g. glog and
	// controller-runtime)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Parse()

	// Use a zap logr.Logger implementation. If none of the zap
	// flags are configured (or if the zap flag set is not being
	// used), this defaults to a production zap logger.
	//
	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	logf.SetLogger(zap.Logger())

	printVersion()

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Error(err, "Failed to get watch namespace")
		os.Exit(1)
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	ctx := context.TODO()
	// Become the leader before proceeding
	err = leader.Become(ctx, "rcm-controller-lock")
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          namespace,
		MetricsBindAddress: fmt.Sprintf("%s:%d", metricsHost, metricsPort),
	})
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources

	if err := ocinfrav1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err := hivev1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err := corev1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err := clusterv1.Install(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err := workv1.Install(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err := rbacv1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err := certificatesv1beta1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Setup manager with controllers")
	missingGVS, err := controller.GetMissingGVS(cfg)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	//Channel to stop the manager
	stopMgrCh := make(chan struct{})

	if err := controller.AddToManager(mgr, missingGVS); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	nbOfMissingGVS := len(missingGVS)

	//If some CRD are not yet installled then we will monitor them
	if nbOfMissingGVS != 0 {
		go func() {
			change := false
			for !change {
				time.Sleep(time.Second * 10)
				currentMissingGVS, err := controller.GetMissingGVS(cfg)
				if err != nil {
					log.Error(err, "")
					os.Exit(1)
				}
				change = len(currentMissingGVS) != nbOfMissingGVS
				if change {
					log.Info(fmt.Sprintf("Old missing GVS: %v", missingGVS))
					log.Info(fmt.Sprintf("New missing GVS: %v", currentMissingGVS))
				}
			}
			//Close the manager
			log.Error(fmt.Errorf("new CRD discovered %s", ""),
				"This is an expected behavior, the operator stopped because a new CRD managed by this operator get discovered")
			close(stopMgrCh)
		}()
	}

	// Start the Cmd
	if err := mgr.Start(stopMgrCh); err != nil {
		log.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}
}

// serveCRMetrics gets the Operator/CustomResource GVKs and generates metrics based on those types.
// It serves those metrics on "http://metricsHost:operatorMetricsPort".
func serveCRMetrics(cfg *rest.Config) error {
	// Below function returns filtered operator/CustomResource specific GVKs.
	// For more control override the below GVK list with your own custom logic.
	filteredGVK, err := k8sutil.GetGVKsFromAddToScheme(clusterv1.AddToScheme)
	if err != nil {
		return err
	}
	// Get the namespace the operator is currently deployed in.
	operatorNs, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		return err
	}
	// To generate metrics in other namespaces, add the values below.
	ns := []string{operatorNs}
	// Generate and serve custom resource specific metrics.
	err = kubemetrics.GenerateAndServeCRMetrics(cfg, ns, filteredGVK, metricsHost, operatorMetricsPort)
	if err != nil {
		return err
	}

	return nil
}
