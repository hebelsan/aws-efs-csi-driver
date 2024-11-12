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
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/kubernetes-sigs/aws-efs-csi-driver/pkg/cloud"
	"github.com/kubernetes-sigs/aws-efs-csi-driver/pkg/util"
)

// Mode is the operating mode of the CSI driver.
type Mode string

const (
	// ControllerMode is the mode that only starts the controller service.
	ControllerMode Mode = "controller"
	// NodeMode is the mode that only starts the node service.
	NodeMode Mode = "node"
	// AllMode is the mode that only starts both the controller and the node service.
	AllMode Mode = "all"

	driverName = "efs.csi.aws.com"

	// AgentNotReadyNodeTaintKey contains the key of taints to be removed on driver startup
	AgentNotReadyNodeTaintKey = "efs.csi.aws.com/agent-not-ready"
)

type Driver struct {
	options      *Options
	nodeID       string
	srv          *grpc.Server
	mounter      Mounter
	efsWatchdog  Watchdog
	cloud        cloud.Cloud
	nodeCaps     []csi.NodeServiceCapability_RPC_Type
	volStatter   VolStatter
	gidAllocator GidAllocator
	tags         map[string]string
}

func NewDriver(efsUtilsCfgPath string, o *Options) *Driver {
	cloud, err := cloud.NewCloud(o.Mode == ControllerMode)
	if err != nil {
		klog.Fatalln(err)
	}

	nodeCaps := SetNodeCapOptInFeatures(o.VolMetricsOptIn)
	watchdog := newExecWatchdog(efsUtilsCfgPath, o.EfsUtilsStaticFilesPath, "amazon-efs-mount-watchdog")
	return &Driver{
		options:      o,
		nodeID:       cloud.GetMetadata().GetInstanceID(),
		mounter:      newNodeMounter(),
		efsWatchdog:  watchdog,
		cloud:        cloud,
		nodeCaps:     nodeCaps,
		volStatter:   NewVolStatter(),
		gidAllocator: NewGidAllocator(),
		tags:         parseTagsFromStr(strings.TrimSpace(o.Tags)),
	}
}

func SetNodeCapOptInFeatures(volMetricsOptIn bool) []csi.NodeServiceCapability_RPC_Type {
	var nCaps = []csi.NodeServiceCapability_RPC_Type{}
	if volMetricsOptIn {
		klog.V(4).Infof("Enabling Node Service capability for Get Volume Stats")
		nCaps = append(nCaps, csi.NodeServiceCapability_RPC_GET_VOLUME_STATS)
	} else {
		klog.V(4).Infof("Node Service capability for Get Volume Stats Not enabled")
	}
	return nCaps
}

func (d *Driver) Run() error {
	scheme, addr, err := util.ParseEndpoint(d.options.Endpoint)
	if err != nil {
		return err
	}

	listener, err := net.Listen(scheme, addr)
	if err != nil {
		return err
	}

	logErr := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			klog.Errorf("GRPC error: %v", err)
		}
		return resp, err
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logErr),
	}
	d.srv = grpc.NewServer(opts...)

	csi.RegisterIdentityServer(d.srv, d)

	switch d.options.Mode {
	case ControllerMode:
		klog.Info("Registering Controller Server")
		csi.RegisterControllerServer(d.srv, d)
	case NodeMode:
		klog.Info("Registering Node Server")
		csi.RegisterNodeServer(d.srv, d)
	case AllMode:
		klog.Info("Registering Node Server")
		csi.RegisterNodeServer(d.srv, d)
		klog.Info("Registering Controller Server")
		csi.RegisterControllerServer(d.srv, d)
	default:
		return fmt.Errorf("unknown mode: %s", d.options.Mode)
	}

	klog.Info("Starting efs-utils watchdog")
	if err := d.efsWatchdog.start(); err != nil {
		return err
	}

	reaper := newReaper()
	klog.Info("Starting reaper")
	reaper.start()

	// Remove taint from node to indicate driver startup success
	// This is done at the last possible moment to prevent race conditions or false positive removals
	if d.options.Mode != ControllerMode {
		go tryRemoveNotReadyTaintUntilSucceed(time.Second, func() error {
			return removeNotReadyTaint(cloud.DefaultKubernetesAPIClient)
		})
	}

	klog.Infof("Listening for connections on address: %#v", listener.Addr())
	return d.srv.Serve(listener)
}

func parseTagsFromStr(tagStr string) map[string]string {
	defer func() {
		if r := recover(); r != nil {
			klog.Errorf("Failed to parse input tag string: %v", tagStr)
		}
	}()

	m := make(map[string]string)
	if tagStr == "" {
		klog.Infof("Did not find any input tags.")
		return m
	}
	tagsSplit := strings.Split(tagStr, " ")
	for _, pair := range tagsSplit {
		p := strings.Split(pair, ":")
		m[p[0]] = p[1]
	}
	return m
}
