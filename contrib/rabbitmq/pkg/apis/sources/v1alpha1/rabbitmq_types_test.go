/*
Copyright 2019 The Knative Authors

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

package v1alpha1

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func TestRabbitmqSourceStatusIsReady(t *testing.T) {
	tests := []struct {
		name string
		s    *RabbitmqSourceStatus
		want bool
	}{{
		name: "uninitialized",
		s:    &RabbitmqSourceStatus{},
		want: false,
	}, {
		name: "initialized",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			return s
		}(),
		want: false,
	}, {
		name: "mark deployed",
		s: func()  *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkDeployed()
			return s
		}(),
		want: false,
	}, {
		name: "mark sink",
		s: func()  *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			return s
		}(),
		want: false,
	}, {
		name: "mark sink and deployed",
		s: func()  *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkDeployed()
			return s
		}(),
		want: true,
	}, {
		name: "mark sink and deployed and then no sink",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkDeployed()
			s.MarkNoSink("Testing", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark sink and deployed then deploying",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkDeployed()
			s.MarkDeploying("Testing", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark sink and deployed then not deployed",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkDeployed()
			s.MarkNotDeployed("Testing", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark sink and not deployed then deploying then deployed",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkNotDeployed("MarkNotDeployed", "")
			s.MarkDeploying("MarkDeploying", "")
			s.MarkDeployed()
			return s
		}(),
		want: true,
	}, {
		name: "mark sink empty and deployed",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("")
			s.MarkDeployed()
			return s
		}(),
		want: false,
	}, {
		name: "mark sink empty and deployed then sink",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("")
			s.MarkDeployed()
			s.MarkSink("uri://example")
			return s
		}(),
		want: true,
	}}

	for
	_, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.s.IsReady()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("%s: unexpected condition (-want, +got) = %v", test.name, diff)
			}
		})
	}
}

func TestRabbitmqSourceStatusGetCondition(t *testing.T)  {
	tests := []struct{
		name      string
		s         *RabbitmqSourceStatus
		condQuery duckv1alpha1.ConditionType
		want      *duckv1alpha1.Condition
	}{{
		name:      "uninitialized",
		s:         &RabbitmqSourceStatus{},
		condQuery: RabbitmqConditionReady,
		want:      nil,
	}, {
		name: "initialized",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   RabbitmqConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark deployed",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkDeployed()
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   RabbitmqConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark sink",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   RabbitmqConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark sink and deployed",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkDeployed()
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   RabbitmqConditionReady,
			Status: corev1.ConditionTrue,
		},
	}, {
		name: "mark sink and deployed then no sink",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkDeployed()
			s.MarkNoSink("Testing", "hi%s", "")
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:    RabbitmqConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "Testing",
			Message: "hi",
		},
	}, {
		name: "mark sink and deployed then deploying",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkDeployed()
			s.MarkDeploying("Testing", "hi%s", "")
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:    RabbitmqConditionReady,
			Status:  corev1.ConditionUnknown,
			Reason:  "Testing",
			Message: "hi",
		},
	}, {
		name: "mark sink and deployed then not deployed",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkDeployed()
			s.MarkNotDeployed("Testing", "hi%s", "")
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:    RabbitmqConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "Testing",
			Message: "hi",
		},
	}, {
		name: "mark sink and not deployed then deploying then deployed",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkNotDeployed("MarkNotDeployed", "%s", "")
			s.MarkDeploying("MarkDeploying", "%s", "")
			s.MarkDeployed()
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   RabbitmqConditionReady,
			Status: corev1.ConditionTrue,
		},
	}, {
		name: "mark sink empty and deployed",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("")
			s.MarkDeployed()
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:    RabbitmqConditionReady,
			Status:  corev1.ConditionUnknown,
			Reason:  "SinkEmpty",
			Message: "Sink has resolved to empty.",
		},
	}, {
		name: "mark sink empty and deployed then sink",
		s: func() *RabbitmqSourceStatus {
			s := &RabbitmqSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("")
			s.MarkDeployed()
			s.MarkSink("uri://example")
			return s
		}(),
		condQuery: RabbitmqConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   RabbitmqConditionReady,
			Status: corev1.ConditionTrue,
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.s.GetCondition(test.condQuery)
			ignoreTime := cmpopts.IgnoreFields(duckv1alpha1.Condition{},
				"LastTransitionTime", "Severity")
			if diff := cmp.Diff(test.want, got, ignoreTime); diff != "" {
				t.Errorf("unexpected condition (-want, +got) = %v", diff)
			}
		})
	}
}