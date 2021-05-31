package v1alpha1

type version interface {
	Major() string
	Minor() string
	Patch() string
}

func (in *ProductRelease) CompareVersion(version version) bool {
	return in.Spec.Version.Major == version.Major() &&
		in.Spec.Version.Minor == version.Minor() &&
		in.Spec.Version.Patch == version.Patch()
}
