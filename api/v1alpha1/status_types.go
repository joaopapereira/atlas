package v1alpha1

type AtlasStatus struct {
	Conditions         Conditions `json:"conditions,omitempty"`
	ObservedGeneration int64      `json:"observedGeneration,omitempty"`
}
