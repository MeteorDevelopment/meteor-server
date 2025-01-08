package api

import (
	"encoding/xml"
	"fmt"
	"io"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
)

type mavenMetadata struct {
	ArtifactId string `xml:"artifactId"`

	Versioning struct {
		SnapshotVersions struct {
			List []mavenMetadataSnapshotVersion `xml:"snapshotVersion"`
		} `xml:"snapshotVersions"`
	} `xml:"versioning"`
}

type mavenMetadataSnapshotVersion struct {
	Extension string `xml:"extension"`
	Value     string `xml:"value"`
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	version, build := GetLatestVersion()

	if v := r.URL.Query().Get("version"); v != "" {
		version = v
		build = GetVersionBuild(v)
	}

	downloadFromMaven(
		w, r,
		"https://maven.meteordev.org/snapshots/meteordevelopment/meteor-client/"+version+"-SNAPSHOT",
		fmt.Sprintf("meteor-client-%s-%d.jar", version, build),
	)

	db.IncrementDownloads()
}

func DownloadBaritoneHandler(w http.ResponseWriter, r *http.Request) {
	version := r.URL.Query().Get("version")

	if version == "" {
		version = core.GetConfig().BaritoneMcVersion
	}

	downloadFromMaven(
		w, r,
		"https://maven.meteordev.org/snapshots/meteordevelopment/baritone/"+version+"-SNAPSHOT",
		fmt.Sprintf("baritone-meteor-%s.jar", version),
	)
}

func downloadFromMaven(w http.ResponseWriter, r *http.Request, url string, filename string) {
	// Get maven version
	res, err := http.Get(url + "/maven-metadata.xml")
	if err != nil {
		core.JsonError(w, "Failed to get maven version.")
		return
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		core.JsonError(w, "Failed to get maven version.")
		return
	}

	var metadata mavenMetadata
	err = xml.NewDecoder(res.Body).Decode(&metadata)
	if err != nil {
		core.JsonError(w, "Failed to decode maven metadata.")
		return
	}

	// Get file url
	fileUrl := ""

	for _, snapshotVersion := range metadata.Versioning.SnapshotVersions.List {
		if snapshotVersion.Extension == "jar" {
			fileUrl = fmt.Sprintf("%s/%s-%s.jar", url, metadata.ArtifactId, snapshotVersion.Value)
			break
		}
	}

	if fileUrl == "" {
		core.JsonError(w, "Failed to find jar file.")
	}

	// Server file
	res, err = http.Get(fileUrl)
	if err != nil {
		core.JsonError(w, "Failed to get file from maven.")
		return
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", res.Header.Get("Content-Length"))

	_, _ = io.Copy(w, res.Body)
}
