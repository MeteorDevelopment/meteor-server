package api

import (
	"fmt"
	"io/ioutil"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
	"os"
	"strings"
	"time"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	version := core.GetConfig().Version
	devBuild := r.URL.Query().Get("devBuild")
	url := ""

	if devBuild != "" {
		version = core.GetConfig().DevBuildVersion

		if devBuild == "latest" {
			devBuild = core.GetConfig().DevBuildVersion
		}

		http.ServeFile(w, r, fmt.Sprintf("meteor-client-%s-%s.jar", version, devBuild))
	} else {
		url = fmt.Sprintf("https://maven.meteordev.org/releases/meteordevelopment/meteor-client/%s/meteor-client-%s.jar", version, version)
	}

	http.Redirect(w, r, url, http.StatusPermanentRedirect)
	db.IncrementDownloads()
}

func UploadDevBuildHandler(w http.ResponseWriter, r *http.Request) {
	// Validate file
	formFile, header, err := r.FormFile("file")
	if err != nil {
		core.JsonError(w, "Invalid file.")
		return
	}

	if !strings.HasSuffix(header.Filename, ".jar") {
		core.JsonError(w, "File needs to be a JAR.")
		return
	}

	// Save file
	_ = os.Mkdir("dev_builds", 0755)

	file, err := os.Create("dev_builds/" + header.Filename)
	if err != nil {
		core.JsonError(w, "Server error. Failed to create file.")
		return
	}

	core.DownloadFile(formFile, file, w)

	// Delete old file if needed
	files, _ := ioutil.ReadDir("dev_builds")

	if len(files) > core.GetConfig().MaxDevBuilds {
		oldestTime := time.Now()
		oldest := ""

		for _, file := range files {
			time_ := file.ModTime()

			if time_.Before(oldestTime) {
				oldestTime = time_
				oldest = file.Name()
			}
		}

		if oldest != "" {
			_ = os.Remove("dev_builds/" + oldest)
		}
	}
}
