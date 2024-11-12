/*
Copyright 2019 The Kubernetes Authors.

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
	"flag"
	"fmt"
	"os"
	"strings"

	"k8s.io/klog/v2"

	"github.com/kubernetes-sigs/aws-efs-csi-driver/pkg/driver"
)

// etcAmazonEfs is the non-negotiable directory that the mount.efs will use for config files. We will create a symlink here.
const etcAmazonEfs = "/etc/amazon/efs"

func main() {
	fs := flag.NewFlagSet("aws-ebs-csi-driver", flag.ExitOnError)

	var (
		version = fs.Bool("version", false, "Print the version and exit")
		args    = os.Args[1:]
		cmd     = string(driver.AllMode)
		options = driver.Options{}
	)

	klog.InitFlags(nil)
	flag.Parse()

	if *version {
		info, err := driver.GetVersionJSON()
		if err != nil {
			klog.Fatalln(err)
		}
		fmt.Println(info)
		os.Exit(0)
	}

	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		cmd = os.Args[1]
		args = os.Args[2:]
	}

	if err := options.AddCmd(cmd); err != nil {
		klog.ErrorS(err, "Failed to parse cmd")
		klog.FlushAndExit(klog.ExitFlushTimeout, 0)
	}

	options.AddFlags(fs)

	if err := fs.Parse(args); err != nil {
		klog.ErrorS(err, "Failed to parse options")
		klog.FlushAndExit(klog.ExitFlushTimeout, 0)
	}

	// chose which configuration directory we will use and create a symlink to it
	err := driver.InitConfigDir(options.EfsUtilsCfgLegacyDirPath, options.EfsUtilsCfgDirPath, etcAmazonEfs)
	if err != nil {
		klog.Fatalln(err)
	}
	drv := driver.NewDriver(etcAmazonEfs, &options)
	if err := drv.Run(); err != nil {
		klog.Fatalln(err)
	}
}
