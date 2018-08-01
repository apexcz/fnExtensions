package models

type GeneratedReport struct {
	LayerCount      int `json:"LayerCount"`
	Vulnerabilities struct {
		High []struct {
			Name          string `json:"Name"`
			NamespaceName string `json:"NamespaceName"`
			Description   string `json:"Description"`
			Link          string `json:"Link"`
			Severity      string `json:"Severity"`
			Metadata      struct {
				NVD struct {
					CVSSv2 struct {
						Score   float64 `json:"Score"`
						Vectors string  `json:"Vectors"`
					} `json:"CVSSv2"`
				} `json:"NVD"`
			} `json:"Metadata"`
			FixedBy        string `json:"FixedBy,omitempty"`
			FeatureName    string `json:"FeatureName"`
			FeatureVersion string `json:"FeatureVersion"`
		} `json:"High"`
		Low []struct {
			Name          string `json:"Name"`
			NamespaceName string `json:"NamespaceName"`
			Description   string `json:"Description"`
			Link          string `json:"Link"`
			Severity      string `json:"Severity"`
			Metadata      struct {
				NVD struct {
					CVSSv2 struct {
						Score   int    `json:"Score"`
						Vectors string `json:"Vectors"`
					} `json:"CVSSv2"`
				} `json:"NVD"`
			} `json:"Metadata"`
			FeatureName    string `json:"FeatureName"`
			FeatureVersion string `json:"FeatureVersion"`
			FixedBy        string `json:"FixedBy,omitempty"`
		} `json:"Low"`
		Medium []struct {
			Name          string `json:"Name"`
			NamespaceName string `json:"NamespaceName"`
			Description   string `json:"Description"`
			Link          string `json:"Link"`
			Severity      string `json:"Severity"`
			Metadata      struct {
				NVD struct {
					CVSSv2 struct {
						Score   int    `json:"Score"`
						Vectors string `json:"Vectors"`
					} `json:"CVSSv2"`
				} `json:"NVD"`
			} `json:"Metadata"`
			FeatureName    string `json:"FeatureName"`
			FeatureVersion string `json:"FeatureVersion"`
			FixedBy        string `json:"FixedBy,omitempty"`
		} `json:"Medium"`
	} `json:"Vulnerabilities"`
}