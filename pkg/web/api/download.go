package api

import (
	"fmt"
	"io/ioutil"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	version := core.GetConfig().Version
	devBuild := r.URL.Query().Get("devBuild")

	if devBuild != "" {
		version = core.GetConfig().DevBuildVersion

		if devBuild == "latest" {
			devBuild = db.GetGlobal().DevBuild
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=meteor-client-%s-%s.jar", version, devBuild))
		http.ServeFile(w, r, fmt.Sprintf("dev_builds/meteor-client-%s-%s.jar", version, devBuild))
	} else {
		url := fmt.Sprintf("https://maven.meteordev.org/releases/meteordevelopment/meteor-client/%s/meteor-client-%s.jar", version, version)
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}

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

	global := db.GetGlobal()
	d, _ := strconv.Atoi(global.DevBuild)
	global.DevBuild = strconv.Itoa(d + 1)
	db.SetDevBuild(global.DevBuild)

	file, err := os.Create("dev_builds/meteor-client-" + core.GetConfig().DevBuildVersion + "-" + global.DevBuild + ".jar")
	if err != nil {
		core.JsonError(w, "Server error. Failed to create file.")
		return
	}

	core.DownloadFile(formFile, file, w)

	// Delete old file if needed
	files, _ := ioutil.ReadDir("dev_builds")

	if len(files) > core.GetConfig().MaxDevBuilds {
		oldestBuild := 6666
		oldest := ""

		for _, file := range files {
			s := strings.TrimSuffix(file.Name(), ".jar")
			build, _ := strconv.Atoi(s[strings.LastIndex(s, "-")+1:])

			if build < oldestBuild {
				oldestBuild = build
				oldest = file.Name()
			}
		}

		if oldest != "" {
			_ = os.Remove("dev_builds/" + oldest)
		}
	}

	// Response
	core.Json(w, core.J{
		"version": core.GetConfig().DevBuildVersion,
		"number":  global.DevBuild,
	})
}
