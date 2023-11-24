/*
Copyright 2020 the Velero contributors.

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

package hook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHookTracker(t *testing.T) {
	tracker := NewHookTracker()

	assert.NotNil(t, tracker)
	assert.Empty(t, tracker.Tracker)
}

func TestHookTracker_Add(t *testing.T) {
	tracker := NewHookTracker()

	tracker.Add("ns1", "pod1", "container1", HookSourceAnnotation, "h1", PhasePre)

	key := hookTrackerKey{
		PodNamespace: "ns1",
		PodName:      "pod1",
		Container:    "container1",
		HookPhase:    PhasePre,
		HookSource:   HookSourceAnnotation,
		HookName:     "h1",
	}

	_, ok := tracker.Tracker[key]
	assert.True(t, ok)
}

func TestHookTracker_Update(t *testing.T) {
	tracker := NewHookTracker()
	tracker.Add("ns1", "pod1", "container1", HookSourceAnnotation, "h1", PhasePre)
	tracker.Update("ns1", "pod1", "container1", HookSourceAnnotation, "h1", PhasePre, true)

	key := hookTrackerKey{
		PodNamespace: "ns1",
		PodName:      "pod1",
		Container:    "container1",
		HookPhase:    PhasePre,
		HookSource:   HookSourceAnnotation,
		HookName:     "h1",
	}

	info := tracker.Tracker[key]
	assert.True(t, info.HookFailed)
}

func TestHookTracker_Stat(t *testing.T) {
	tracker := NewHookTracker()

	tracker.Add("ns1", "pod1", "container1", HookSourceAnnotation, "h1", PhasePre)
	tracker.Add("ns2", "pod2", "container1", HookSourceAnnotation, "h2", PhasePre)
	tracker.Update("ns1", "pod1", "container1", HookSourceAnnotation, "h1", PhasePre, true)

	attempted, failed := tracker.Stat()
	assert.Equal(t, 2, attempted)
	assert.Equal(t, 1, failed)
}
