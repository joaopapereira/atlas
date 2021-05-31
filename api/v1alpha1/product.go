package v1alpha1

func (in *ProductLocation) Equals(other ProductLocation) bool {
	return in.VersionMatch(other) &&
		in.ImageRepo == other.ImageRepo &&
		in.ImageSHA == other.ImageSHA
}

func (in *ProductLocation) VersionMatch(other ProductLocation) bool {
	return in.Version.Equals(other.Version)
}

func (in *Version) Equals(version Version) bool {
	return in.Major == version.Major &&
		in.Minor == version.Minor &&
		in.Patch == version.Patch
}
