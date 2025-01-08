package api

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"golang.org/x/mod/semver"
	"io"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Stats struct {
	core.Config

	Date          string `json:"date"`
	DevBuild      string `json:"devBuild"`
	Downloads     int    `json:"downloads"`
	OnlinePlayers int    `json:"onlinePlayers"`
	OnlineUUIDs   int    `json:"onlineUUIDs"`

	Builds map[string]int `json:"builds"`
}

var builds map[string]int

func InitStats() {
	t := time.NewTicker(10 * time.Minute)

	go func() {
		for {
			builds = getBuildNumbers()

			<-t.C
		}
	}()
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	if date == "" {
		g := db.GetGlobal()

		core.Json(w, Stats{
			Config:        core.GetConfig(),
			Date:          core.GetDate(),
			DevBuild:      g.DevBuild,
			Downloads:     g.Downloads,
			OnlinePlayers: GetPlayingCount(),
			Builds:        builds,
		})
	} else {
		stats, err := db.GetJoinStats(date)

		if err != nil {
			core.JsonError(w, "Invalid date.")
			return
		}

		core.Json(w, stats)
	}
}

func RecheckMavenHandler(w http.ResponseWriter, _ *http.Request) {
	builds = getBuildNumbers()

	core.Json(w, struct{}{})
}

func GetLatestVersion() (string, int) {
	latest := "0.0.0"
	build := 0

	for version, number := range builds {
		if semver.Compare("v"+version, "v"+latest) == 1 {
			latest = version
			build = number
		}
	}

	return latest, build
}

func GetVersionBuild(version string) int {
	if build, ok := builds[version]; ok {
		return build
	}

	return 0
}

// Maven

type MavenSnapshotVersion struct {
	Extension  string `xml:"extension"`
	Classifier string `xml:"classifier"`
	Value      string `xml:"value"`
}

type MavenVersioning struct {
	Versions         []string               `xml:"versions>version"`
	SnapshotVersions []MavenSnapshotVersion `xml:"snapshotVersions>snapshotVersion"`
}

type MavenMetadata struct {
	Versioning MavenVersioning `xml:"versioning"`
}

type FabricMod struct {
	Version string `json:"version"`
}

func getBuildNumbers() map[string]int {
	builds := make(map[string]int)

	res, err := http.Get("https://maven.meteordev.org/snapshots/meteordevelopment/meteor-client/maven-metadata.xml")
	if err != nil {
		return builds
	}

	var metadata MavenMetadata
	err = xml.NewDecoder(res.Body).Decode(&metadata)
	if err != nil {
		return builds
	}

	mutex := sync.Mutex{}
	group := sync.WaitGroup{}

	for _, version := range metadata.Versioning.Versions {
		version := version
		group.Add(1)

		go func() {
			i := strings.IndexRune(version, '-')
			mcVersion := version[:i]

			if !strings.HasPrefix(mcVersion, "0") {
				build, err := getBuildNumber(version)

				if err == nil {
					mutex.Lock()
					builds[mcVersion] = build
					mutex.Unlock()
				}
			}

			group.Done()
		}()
	}

	group.Wait()

	return builds
}

func getBuildNumber(version string) (int, error) {
	res, err := http.Get("https://maven.meteordev.org/snapshots/meteordevelopment/meteor-client/" + version + "/maven-metadata.xml")
	if err != nil {
		return -1, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	var metadata MavenMetadata
	err = xml.NewDecoder(res.Body).Decode(&metadata)
	if err != nil {
		return -1, err
	}

	var filename = ""

	for _, version := range metadata.Versioning.SnapshotVersions {
		if version.Classifier == "" && version.Extension == "jar" {
			filename = "meteor-client-" + version.Value + ".jar"
			break
		}
	}

	res, err = http.Get("https://maven.meteordev.org/snapshots/meteordevelopment/meteor-client/" + version + "/" + filename)
	if err != nil {
		return -1, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return -1, err
	}

	jar, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return -1, err
	}

	var mod FabricMod

	for _, file := range jar.File {
		if file.Name == "fabric.mod.json" {
			reader, err := file.Open()
			if err != nil {
				continue
			}

			_ = json.NewDecoder(reader).Decode(&mod)

			_ = reader.Close()

			break
		}
	}

	if mod.Version != "" {
		i := strings.IndexRune(mod.Version, '-')

		build, err := strconv.ParseInt(mod.Version[i+1:], 10, 32)
		if err != nil {
			return -1, err
		}

		return int(build), nil
	}

	return -1, errors.New("unknown build")
}
