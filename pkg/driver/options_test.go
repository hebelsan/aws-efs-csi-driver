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
	"testing"
)

func TestAddFlags(t *testing.T) {
	o := &Options{}
	o.Mode = AllMode

	f := flag.NewFlagSet("test", flag.ExitOnError)
	o.AddFlags(f)

	if err := f.Set("endpoint", "custom-endpoint"); err != nil {
		t.Errorf("error setting endpoint: %v", err)
	}
	if err := f.Set("efs-utils-config-dir-path", "custom-efs-utils-config-dir-path"); err != nil {
		t.Errorf("error setting efs-utils-config-dir-path: %v", err)
	}
	if err := f.Set("efs-utils-config-legacy-dir-path", "custom-efs-utils-config-legacy-dir-path"); err != nil {
		t.Errorf("error setting efs-utils-config-legacy-dir-path: %v", err)
	}
	if err := f.Set("efs-utils-static-files-path", "custom-efs-utils-static-files-path"); err != nil {
		t.Errorf("error setting efs-utils-static-files-path: %v", err)
	}
	if err := f.Set("delete-access-point-root-dir", "true"); err != nil {
		t.Errorf("error setting delete-access-point-root-dir: %v", err)
	}
	if err := f.Set("tags", "environment:prod"); err != nil {
		t.Errorf("error setting tags: %v", err)
	}
	if err := f.Set("vol-metrics-opt-in", "true"); err != nil {
		t.Errorf("error setting vol-metrics-opt-in: %v", err)
	}
	if err := f.Set("vol-metrics-refresh-period", "230"); err != nil {
		t.Errorf("error setting vol-metrics-refresh-period: %v", err)
	}
	if err := f.Set("vol-metrics-fs-rate-limit", "6"); err != nil {
		t.Errorf("error setting vol-metrics-fs-rate-limit: %v", err)
	}

	if o.Endpoint != "custom-endpoint" {
		t.Errorf("unexpected Endpoint: got %s, want custom-endpoint", o.Endpoint)
	}
	if o.EfsUtilsCfgDirPath != "custom-efs-utils-config-dir-path" {
		t.Errorf("unexpected Efs utils cfg path: got %s, want custom-efs-utils-config-dir-path", o.EfsUtilsCfgDirPath)
	}
	if o.EfsUtilsCfgLegacyDirPath != "custom-efs-utils-config-legacy-dir-path" {
		t.Errorf("unexpected Efs utils cfg legacy path: got %s, want custom-efs-utils-config-legacy-dir-path", o.EfsUtilsCfgLegacyDirPath)
	}
	if o.EfsUtilsStaticFilesPath != "custom-efs-utils-static-files-path" {
		t.Errorf("unexpected Efs utils static file path: got %s, want custom-efs-utils-static-files-path", o.EfsUtilsStaticFilesPath)
	}
	if !o.DeleteAccessPointRootDir {
		t.Errorf("unexpected delete access point root dir: got %t, want true", o.DeleteAccessPointRootDir)
	}
	if o.Tags != "environment:prod" {
		t.Errorf("unexpected tags: got %s, want environment:prod", o.Tags)
	}
	if !o.VolMetricsOptIn {
		t.Errorf("unexpected tags: got %t, want true", o.VolMetricsOptIn)
	}
	if o.VolMetricsRefreshPeriod != 230 {
		t.Errorf("unexpected vol-metrics-refresh-period: got %v, want 230", o.VolMetricsRefreshPeriod)
	}
	if o.VolMetricsFsRateLimit != 6 {
		t.Errorf("unexpected vol-metrics-fs-rate-limit: got %d, want 6", o.VolMetricsFsRateLimit)
	}
}

func TestAddCmd(t *testing.T) {
	o := &Options{}
	if err := o.AddCmd(string(AllMode)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := o.AddCmd(string(NodeMode)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := o.AddCmd(string(ControllerMode)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := o.AddCmd("invalid-mode"); err == nil {
		t.Errorf("Expected an error when passing invalid mode")
	}
}
