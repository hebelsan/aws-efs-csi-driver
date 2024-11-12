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

package driver

import (
	"flag"
	"fmt"
)

// Options contains options and configuration settings for the driver.
type Options struct {
	Mode Mode

	// #### Server options ####

	// Endpoint is the endpoint for the CSI driver server
	Endpoint string
	// EfsUtilsCfgDirPath is the preferred path for the efs-utils config directory.
	// Efs-utils-config-legacy-dir-path will be used if it is not empty, otherwise efs-utils-config-dir-path will be used.
	EfsUtilsCfgDirPath string
	// EfsUtilsCfgLegacyDirPath ist the path to the legacy efs-utils config directory mounted from the host path /etc/amazon/efs
	EfsUtilsCfgLegacyDirPath string
	// EfsUtilsStaticFilesPath is the path to efs-utils static files directory
	EfsUtilsStaticFilesPath string

	// #### Controller options ####

	// DeleteAccessPointRootDir is the option to delete an access point root directory by DeleteVolume.
	// By default, DeleteVolume will delete the access point behind Persistent Volume and deleting
	// access point will not delete the access point root directory or its contents
	DeleteAccessPointRootDir bool
	// Tags are the space separated key:value pairs which will be added for EFS resources.
	//  For example, 'environment:prod region:us-east-1'
	Tags string

	// #### Node options #####

	// VolMetricsOptIn is the option to emit volume metrics
	VolMetricsOptIn bool
	// VolMetricsRefreshPeriod is the refresh period for volume metrics in minutes
	VolMetricsRefreshPeriod float64
	// VolMetricsFsRateLimit is the volume metrics routines rate limiter per file system
	VolMetricsFsRateLimit int
}

func (o *Options) AddFlags(f *flag.FlagSet) {
	// Server options
	f.StringVar(&o.Endpoint, "endpoint", "unix://tmp/csi.sock", "CSI Endpoint")
	f.StringVar(&o.EfsUtilsCfgDirPath, "efs-utils-config-dir-path", "/var/amazon/efs", "The preferred path for the efs-utils config directory. efs-utils-config-legacy-dir-path will be used if it is not empty, otherwise efs-utils-config-dir-path will be used.")
	f.StringVar(&o.EfsUtilsCfgLegacyDirPath, "efs-utils-config-legacy-dir-path", "/etc/amazon/efs-legacy", "The path to the legacy efs-utils config directory mounted from the host path /etc/amazon/efs")
	f.StringVar(&o.EfsUtilsStaticFilesPath, "efs-utils-static-files-path", "/etc/amazon/efs-static-files/", "The path to efs-utils static files directory")

	// Controller options
	if o.Mode == AllMode || o.Mode == ControllerMode {
		f.BoolVar(&o.DeleteAccessPointRootDir, "delete-access-point-root-dir", false,
			"Opt in to delete access point root directory by DeleteVolume. By default, DeleteVolume will delete the access point behind Persistent Volume and deleting access point will not delete the access point root directory or its contents.")
		f.StringVar(&o.Tags, "tags", "", "Space separated key:value pairs which will be added as tags for EFS resources. For example, 'environment:prod region:us-east-1'")
	}
	// Node options
	if o.Mode == AllMode || o.Mode == NodeMode {
		f.BoolVar(&o.VolMetricsOptIn, "vol-metrics-opt-in", false, "Opt in to emit volume metrics")
		f.Float64Var(&o.VolMetricsRefreshPeriod, "vol-metrics-refresh-period", 240, "Refresh period for volume metrics in minutes")
		f.IntVar(&o.VolMetricsFsRateLimit, "vol-metrics-fs-rate-limit", 5, "Volume metrics routines rate limiter per file system")
	}
}

func (o *Options) AddCmd(cmd string) error {
	switch cmd {
	case string(ControllerMode), string(NodeMode), string(AllMode):
		o.Mode = Mode(cmd)
	default:
		return fmt.Errorf("unknown driver mode %s: expected %s, %s, %s", cmd, ControllerMode, NodeMode, AllMode)
	}
	return nil
}
