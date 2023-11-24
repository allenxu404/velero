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

const (
	HookSourceAnnotation = "annotation"
	HookSourceSpec       = "spec"
)

// hookTrackerKey identifies a backup/restore hook
type hookTrackerKey struct {
	// PodNamespace indicates the namespace of pod where hooks are executed.
	// For hooks specified in the backup/restore spec, this field is the namespace of an applicable pod.
	// For hooks specified in pod annotation, this field is the namespace of pod where hooks are annotated.
	PodNamespace string
	// PodName indicates the pod where hooks are executed.
	// For hooks specified in the backup/restore spec, this field is an applicable pod name.
	// For hooks specified in pod annotation, this field is the pod where hooks are annotated.
	PodName string
	// HookPhase is only for backup hooks, for restore hooks, this field is empty.
	HookPhase hookPhase
	// HookName is only for hooks specified in the backup/restore spec.
	// For hooks specified in pod annotation, this field is empty or "<from-annotation>".
	HookName string
	// HookSource indicates where hooks come from.
	HookSource string
	// Container indicates the container hooks use.
	// For hooks specified in the backup/restore spec, the container might be the same under different hookName.
	Container string
}

// hookTrackerVal records the execution status of a specific hook.
// hookTrackerVal is extensible to accommodate additional fields as needs develop.
type hookTrackerVal struct {
	// HookFailed indicates if hook failed to execute.
	HookFailed bool
}

// HookTracker tracks all hooks' execution status
type HookTracker struct {
	Tracker map[hookTrackerKey]hookTrackerVal
}

// NewHookTracker creates a hookTracker.
func NewHookTracker() *HookTracker {
	return &HookTracker{
		Tracker: make(map[hookTrackerKey]hookTrackerVal),
	}
}

// Add adds a hook to the tracker
func (this *HookTracker) Add(podNamespace, podName, container, source, hookName string, hookPhase hookPhase) {
	key := hookTrackerKey{
		PodNamespace: podNamespace,
		PodName:      podName,
		HookSource:   source,
		Container:    container,
		HookPhase:    hookPhase,
		HookName:     hookName,
	}

	if _, ok := this.Tracker[key]; !ok {
		this.Tracker[key] = hookTrackerVal{}
	}
}

// Update updates the hook's execution status
func (this *HookTracker) Update(podNamespace, podName, container, source, hookName string, hookPhase hookPhase, hookFailed bool) {
	key := hookTrackerKey{
		PodNamespace: podNamespace,
		PodName:      podName,
		HookSource:   source,
		Container:    container,
		HookPhase:    hookPhase,
		HookName:     hookName,
	}

	if _, ok := this.Tracker[key]; ok {
		this.Tracker[key] = hookTrackerVal{HookFailed: hookFailed}
	}
}

// Stat calculates the number of attempted hooks and failed hooks
func (this *HookTracker) Stat() (hookAttemptedCnt int, hookFailed int) {
	for _, hookInfo := range this.Tracker {
		hookAttemptedCnt++
		if hookInfo.HookFailed {
			hookFailed++
		}
	}
	return
}

func convertHookSource(hookSourceFromHandler string) string {
	if hookSourceFromHandler == "backupSpec" {
		return HookSourceSpec
	}
	return HookSourceAnnotation
}
