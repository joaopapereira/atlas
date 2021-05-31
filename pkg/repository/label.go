/*
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
package repository

import "strings"

type ImageLabel struct {
	Images ImagesJSON
}

type ImagesJSON = []ImageJSON

type ImageJSON struct {
	Product      `json:",inline"`
	Description  string       `json:"description,omitempty"`
	ImageRepo    string       `json:"imageRepo,omitempty"`
	ImageSHA     string       `json:"imageSha,omitempty"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
	UpgradesFrom []Version    `json:"upgradesFrom,omitempty"`
}

type Version struct {
	MajorV string `json:"major,omitempty"`
	MinorV string `json:"minor,omitempty"`
	PatchV string `json:"patch,omitempty"`
}

type Product struct {
	Name    string  `json:"name,omitempty"`
	Slug    string  `json:"slug,omitempty"`
	Version Version `json:"version,omitempty"`
}

type Dependency struct {
	Product `json:",inline"`
}

func NewVersion(version string) Version {
	v := strings.Split(version, ".")

	return Version{
		MajorV: v[0],
		MinorV: v[1],
		PatchV: v[2],
	}
}

func (v Version) Major() string {
	return v.MajorV
}
func (v Version) Minor() string {
	return v.MinorV
}
func (v Version) Patch() string {
	return v.PatchV
}
