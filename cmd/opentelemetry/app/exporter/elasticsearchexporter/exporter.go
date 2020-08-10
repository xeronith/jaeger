// Copyright (c) 2020 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elasticsearchexporter

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/jaegertracing/jaeger/plugin/storage/es"
)

// new creates Elasticsearch exporter/storage.
func new(config *Config, params component.ExporterCreateParams) (component.TraceExporter, error) {
	esCfg := config.GetPrimary()
	w, err := newEsSpanWriter(*esCfg, params.Logger)
	if err != nil {
		return nil, err
	}
	if config.Primary.IsCreateIndexTemplates() {
		spanMapping, serviceMapping := es.GetSpanServiceMappings(esCfg.GetNumShards(), esCfg.GetNumReplicas(), esCfg.GetVersion())
		if err = w.CreateTemplates(spanMapping, serviceMapping); err != nil {
			return nil, err
		}
	}
	return exporterhelper.NewTraceExporter(
		config,
		w.WriteTraces)
}