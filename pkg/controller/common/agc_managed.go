/*
Copyright 2026 The opendatahub.io Authors.

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

package common

import (
	inferencev1alpha1 "github.com/opendatahub-io/ai-gateway-payload-processing/api/inference/v1alpha1"
)

// IsAGCManaged reports whether an ExternalModel uses fields that only
// ai-gateway-controller (AGC) interprets (RHOAIENG-90800): explicit gateway
// attachments (spec.gatewayRefs) or a provider reference set that leaves IPP
// nothing to serve in the model namespace.
//
// A gatewayRefs attachment is always AGC-managed: it explicitly selects AGC
// tenant gateways, which AGC serves regardless of the model namespace. A
// cross-namespace ExternalProvider reference is AGC-managed only when no
// same-namespace reference remains: AGC does not resolve cross-namespace
// providers yet (RHOAIENG-90804), so a model that also references a provider
// in its own namespace keeps serving through IPP on that local reference —
// the foreign reference is skipped, not handed off, so the edit cannot
// silently remove service.
//
// IPP must not reconcile, serve, or write status for AGC-managed objects
// (RHOAIENG-99761): AGC owns their routing and status exclusively. Without
// this guard the legacy IPP controller would ignore the attachment, create a
// route on the default gateway and overwrite AGC status.
func IsAGCManaged(model *inferencev1alpha1.ExternalModel) bool {
	if len(model.Spec.GatewayRefs) > 0 {
		return true
	}
	localRef := false
	foreignRef := false
	for _, ref := range model.Spec.ExternalProviderRefs {
		if ref.Ref.Namespace != "" && ref.Ref.Namespace != model.Namespace {
			foreignRef = true
		} else {
			localRef = true
		}
	}
	return foreignRef && !localRef
}
