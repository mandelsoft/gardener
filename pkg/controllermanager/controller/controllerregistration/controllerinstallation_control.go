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

package controllerregistration

import (
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardencorev1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
)

func (c *Controller) controllerInstallationAdd(obj interface{}) {
	controllerInstallation, ok := obj.(*gardencorev1beta1.ControllerInstallation)
	if !ok {
		return
	}

	c.controllerRegistrationSeedQueue.Add(controllerInstallation.Spec.SeedRef.Name)
}

func (c *Controller) controllerInstallationUpdate(oldObj, newObj interface{}) {
	oldObject, ok := oldObj.(*gardencorev1beta1.ControllerInstallation)
	if !ok {
		return
	}

	newObject, ok := newObj.(*gardencorev1beta1.ControllerInstallation)
	if !ok {
		return
	}

	if gardencorev1beta1helper.IsControllerInstallationRequired(*oldObject) == gardencorev1beta1helper.IsControllerInstallationRequired(*newObject) {
		return
	}

	c.controllerInstallationAdd(newObj)
}
