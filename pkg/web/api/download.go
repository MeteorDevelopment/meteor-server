package api

import (
	"fmt"
	"log"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	devBuild := r.URL.Query().Get("devBuild")

	if devBuild != "" {
		version := db.GetGlobal().DevBuildVersion

		if devBuild == "latest" {
			devBuild = db.GetGlobal().DevBuild
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=meteor-client-%s-%s.jar", version, devBuild))
		http.ServeFile(w, r, fmt.Sprintf("data/dev_builds/meteor-client-%s-%s.jar", version, devBuild))
	} else {
		version := core.GetConfig().Version
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
	_ = os.Mkdir("data/dev_builds", 0755)

	devBuild := header.Filename[strings.LastIndex(header.Filename, "-")+1 : len(header.Filename)-4]
	devBuildNum, _ := strconv.Atoi(devBuild)
	devBuildVersion := header.Filename[strings.Index(header.Filename, ".")-1 : strings.LastIndex(header.Filename, "-")]

	currDevBuild, _ := strconv.Atoi(db.GetGlobal().DevBuild)

	if currDevBuild < devBuildNum {
		db.SetDevBuild(devBuild)
		db.SetDevBuildVersion(devBuildVersion)
	}

	file, err := os.Create("data/dev_builds/meteor-client-" + devBuildVersion + "-" + devBuild + ".jar")
	if err != nil {
		core.JsonError(w, "Server error. Failed to create file.")
		return
	}

	core.DownloadFile(formFile, file, w)

	// Delete old file if needed
	files, _ := os.ReadDir("data/dev_builds")

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
			_ = os.Remove("data/dev_builds/" + oldest)
		}
	}

	// Response
	core.Json(w, core.J{
		"version": devBuildVersion,
		"number":  devBuild,
	})

	err = formFile.Close()
	if err != nil {
		log.Println("Error closing input file from updateDevBuild")
	}
}
