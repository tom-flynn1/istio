// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/types"

	adptTmpl "istio.io/api/mixer/adapter/model/v1beta1"
	"istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/pkg/lang/ast"
	"istio.io/istio/mixer/pkg/protobuf/yaml/dynamic"
	"istio.io/istio/mixer/pkg/template"
)

type (

	// Snapshot view of configuration.
	Snapshot struct {
		ID int64

		// Static information
		Templates map[string]*template.Info
		Adapters  map[string]*adapter.Info

		// Config store based information
		Attributes ast.AttributeDescriptorFinder

		HandlersStatic  map[string]*HandlerStatic
		InstancesStatic map[string]*InstanceStatic

		//  TemplateMetadatas contains template descriptors loaded from the store
		TemplateMetadatas map[string]*TemplateMetadata
		//  AdapterMetadatas contains adapter metadata loaded from the store
		AdapterMetadatas map[string]*AdapterMetadata

		HandlersDynamic  map[string]*HandlerDynamic
		InstancesDynamic map[string]*InstanceDynamic
		Rules            []*Rule

		// Perf Counters relevant to configuration.
		Counters Counters
	}

	// HandlerDynamic configuration for dynamically loaded, grpc adapters. Fully resolved.
	HandlerDynamic struct {
		Name string

		Adapter *Adapter

		// AdapterConfig used to construct the Handler. This is passed in verbatim to the remote adapter.
		AdapterConfig *types.Any

		// Connection information for the handler.
		Connection *v1beta1.Connection
	}

	// HandlerStatic configuration for compiled in adapters. Fully resolved.
	HandlerStatic struct {

		// Name of the Handler. Fully qualified.
		Name string

		// Associated adapter. Always resolved.
		Adapter *adapter.Info

		// parameters used to construct the Handler.
		Params proto.Message
	}

	// InstanceDynamic configuration for dynamically loaded templates. Fully resolved.
	InstanceDynamic struct {
		Name string

		Template *Template

		// Encoder to create request instance bytes from attributes
		Encoder dynamic.Encoder
	}

	// InstanceStatic configuration for compiled templates. Fully resolved.
	InstanceStatic struct {
		// Name of the instance. Fully qualified.
		Name string

		// Associated template. Always resolved.
		Template *template.Info

		// parameters used to construct the instance.
		Params proto.Message

		// inferred type for the instance.
		InferredType proto.Message
	}

	// Rule configuration. Fully resolved.
	Rule struct {
		// Name of the rule
		Name string

		// Namespace of the rule
		Namespace string

		// Match condition
		Match string

		ActionsDynamic []*ActionDynamic

		ActionsStatic []*ActionStatic

		ResourceType ResourceType
	}

	// ActionDynamic configuration. Fully resolved.
	ActionDynamic struct {
		// Handler that this action is resolved to.
		Handler *HandlerDynamic
		// Instances that should be generated as part of invoking action.
		Instances []*InstanceDynamic
	}

	// ActionStatic configuration. Fully resolved.
	ActionStatic struct {
		// Handler that this action is resolved to.
		Handler *HandlerStatic

		// Instances that should be generated as part of invoking action.
		Instances []*InstanceStatic
	}

	// Template contains info about a template
	Template struct {
		// Name of the template.
		//
		// Note this is the template's resource name and not the template's internal name that adapter developer
		// uses to implement adapter service.
		Name string

		// Variety of this template
		Variety adptTmpl.TemplateVariety

		// InternalPackageDerivedName is the name of the template from adapter developer point of view.
		// The service and functions implemented by the adapter is based on this name
		// NOTE: This name derived from template proto package and not the resource name.
		InternalPackageDerivedName string

		// Template's file descriptor set.
		FileDescSet *descriptor.FileDescriptorSet

		// package name of the `Template` message
		PackageName string
	}

	// Adapter contains info about an adapter
	Adapter struct {
		Name string

		// Adapter's file descriptor set.
		ConfigDescSet *descriptor.FileDescriptorSet

		// package name of the `Params` message
		PackageName string

		SupportedTemplates []*Template

		SessionBased bool

		Description string
	}
)

// Empty returns a new, empty configuration snapshot.
func Empty() *Snapshot {
	return &Snapshot{
		ID:       -1,
		Rules:    []*Rule{},
		Counters: newCounters(-1),
	}
}
