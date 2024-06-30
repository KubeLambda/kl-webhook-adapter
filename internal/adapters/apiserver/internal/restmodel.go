package internal

type VersionRest struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Build   string `json:"build"`
}
