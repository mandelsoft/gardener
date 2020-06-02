// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"errors"
	"sync"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	memcache "k8s.io/client-go/discovery/cached/memory"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FakeDiscovery is used to tests, more specifically - chart applier.
type FakeDiscovery struct {
	*fakediscovery.FakeDiscovery
	lock          sync.Mutex
	GroupListFn   func() *metav1.APIGroupList
	ResourceMapFn func() map[string]*metav1.APIResourceList
}

// ServerResourcesForGroupVersion returns resource lists from ResourceMapFn
func (c *FakeDiscovery) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if rl, ok := c.ResourceMapFn()[groupVersion]; ok {
		return rl, nil
	}

	return nil, errors.New("doesn't exist")
}

// ServerGroups return group lists from ResourceMapFn
func (c *FakeDiscovery) ServerGroups() (*metav1.APIGroupList, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	groupList := c.GroupListFn()
	if groupList == nil {
		return nil, errors.New("doesn't exist")
	}

	return groupList, nil
}

// ServerVersion return empty version.
func (c *FakeDiscovery) ServerVersion() (*version.Info, error) {
	return &version.Info{}, nil
}

// NewTestApplier create a new fake applier.
func NewTestApplier(c client.Client, discovery discovery.DiscoveryInterface) (*kubernetes.Applier, error) {
	tmp := kubernetes.NewControllerClient

	defer func() {
		kubernetes.NewControllerClient = tmp
	}()

	cachedDiscoveryClient := memcache.NewMemCacheClient(discovery)
	kubernetes.NewControllerClient = func(config *rest.Config, options client.Options) (client.Client, error) {
		return c, nil
	}

	return kubernetes.NewApplierInternal(nil, cachedDiscoveryClient)
}
