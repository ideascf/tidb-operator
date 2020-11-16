// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/openshift/generic-admission-server/pkg/cmd"
	"github.com/pingcap/tidb-operator/pkg/features"
	"github.com/pingcap/tidb-operator/pkg/version"
	"github.com/pingcap/tidb-operator/pkg/webhook"
	"github.com/pingcap/tidb-operator/pkg/webhook/pod"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
)

var (
	printVersion         bool
	extraServiceAccounts string
)

func init() {
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	flag.BoolVar(&printVersion, "V", false, "Show version and quit")
	flag.BoolVar(&printVersion, "version", false, "Show version and quit")
	flag.StringVar(&extraServiceAccounts, "extraServiceAccounts", "", "comma-separated, extra Service Accounts the Webhook should control. The full pattern for each common service account is system:serviceaccount:<namespace>:<serviceaccount-name>")
	features.DefaultFeatureGate.AddFlag(flag.CommandLine)
}

func main() {

	flag.Parse()

	logs.InitLogs()
	defer logs.FlushLogs()

	if printVersion {
		version.PrintVersionInfo()
		os.Exit(0)
	}
	version.LogVersionInfo()

	flag.CommandLine.VisitAll(func(flag *flag.Flag) {
		klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})

	ah := &webhook.AdmissionHook{
		ExtraServiceAccounts: extraServiceAccounts,
	}
	ns := os.Getenv("NAMESPACE")
	if len(ns) < 1 {
		klog.Fatal("ENV NAMESPACE should be set.")
	}
	pod.AstsControllerServiceAccounts = fmt.Sprintf("system:serviceaccount:%s:advanced-statefulset-controller", ns)

	cmd.RunAdmissionServer(ah)
}
