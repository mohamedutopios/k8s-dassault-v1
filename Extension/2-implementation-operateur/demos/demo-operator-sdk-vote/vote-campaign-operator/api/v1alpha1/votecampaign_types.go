package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VoteCampaignSpec defines the desired state of VoteCampaign
type VoteCampaignSpec struct {
    StartTime string `json:"start_time,omitempty"`
    EndTime   string `json:"end_time,omitempty"`
    Options   []Option `json:"options,omitempty"`
}

// Option defines a voting option
type Option struct {
    Name        string `json:"name,omitempty"`
    Description string `json:"description,omitempty"`
}

// VoteCampaignStatus defines the observed state of VoteCampaign
type VoteCampaignStatus struct {
    Active bool `json:"active,omitempty"`
    Votes  map[string]int `json:"votes,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VoteCampaign is the Schema for the votecampaigns API
type VoteCampaign struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   VoteCampaignSpec   `json:"spec,omitempty"`
    Status VoteCampaignStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VoteCampaignList contains a list of VoteCampaign
type VoteCampaignList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []VoteCampaign `json:"items"`
}

func init() {
    SchemeBuilder.Register(&VoteCampaign{}, &VoteCampaignList{})
}
