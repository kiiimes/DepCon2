/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"os"
	"time"
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"

	"github.com/kubernetes-sigs/kube-batch/cmd/kube-batch/app"
	"github.com/kubernetes-sigs/kube-batch/cmd/kube-batch/app/options"

	// Import default actions/plugins.
	_ "github.com/kubernetes-sigs/kube-batch/pkg/scheduler/actions"
	_ "github.com/kubernetes-sigs/kube-batch/pkg/scheduler/plugins"
)

var logFlushFreq = pflag.Duration("log-flush-frequency", 5*time.Second, "Maximum number of seconds between log flushes")

func main() {
	s := options.NewServerOption()
	s.AddFlags(pflag.CommandLine)
	s.RegisterOptions()

	flag.Set("alsologtostderr", "true")
	flag.Parse()
	flag.Set("v","3")
	flag.Parse()

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)

	// Sync the glog and klog flags.
        flag.CommandLine.VisitAll(func(f1 *flag.Flag) {
                f2 := klogFlags.Lookup(f1.Name)
                if f2 != nil {
                        value := f1.Value.String()
                        f2.Value.Set(value)
                }
        })

        glog.Info("hello from glog!")
	klog.Info("main start")
	if err := s.CheckOptionOrDie(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	// The default glog flush interval is 30 seconds, which is frighteningly long.
	go wait.Until(glog.Flush, *logFlushFreq, wait.NeverStop)
	defer glog.Flush()
	klog.Flush()

	if err := app.Run(s); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
