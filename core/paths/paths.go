package paths

import (
	"github.com/hashload/boss/consts"
	"github.com/hashload/boss/env"
	"github.com/hashload/boss/models"
	"github.com/hashload/boss/utils"
	"github.com/masterminds/glide/msg"
	"io/ioutil"
	"os"
	"path/filepath"
)

func EnsureCacheDir(dep models.Dependency) {
	cacheDir := filepath.Join(env.GetCacheDir(), dep.GetHashName())

	fi, err := os.Stat(cacheDir)
	if err != nil {
		msg.Debug("Creating %s", cacheDir)
		if err := os.MkdirAll(cacheDir, os.ModeDir|0755); err != nil {
			msg.Die("Could not create %s: %s", cacheDir, err)
		}
	} else if !fi.IsDir() {
		msg.Die(".cache is not a directory")
	}
}

func cleanArtifacts(dir string, lock models.PackageLock) {
	fileInfos, err := ioutil.ReadDir(dir)
	utils.HandleError(err)
	artifactList := lock.GetArtifactList()
	for _, infoArtifact := range fileInfos {
		if infoArtifact.IsDir() {
			continue
		}
		if !utils.Contains(artifactList, infoArtifact.Name()) {
			for {
				err := os.Remove(filepath.Join(dir, infoArtifact.Name()))
				utils.HandleError(err)
				if err == nil {
					break
				}
			}
		}

	}
}

func EnsureCleanModulesDir(dependencies []models.Dependency, lock models.PackageLock) {
	cacheDir := env.GetModulesDir()
	cacheDirInfo, err := os.Stat(cacheDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, os.ModeDir|0755); err != nil {
			msg.Die("Could not create %s: %s", cacheDir, err)
		}
	} else if cacheDirInfo != nil && !cacheDirInfo.IsDir() {
		msg.Die("modules is not a directory")
	} else {
		fileInfos, err := ioutil.ReadDir(cacheDir)
		utils.HandleError(err)
		dependenciesNames := models.GetDependenciesNames(dependencies)
		for _, info := range fileInfos {
			if !info.IsDir() {
				err := os.Remove(info.Name())
				utils.HandleError(err)
			}
			if utils.Contains([]string{consts.BplFolder, consts.DcuFolder, consts.DcpFolder, consts.BinFolder}, info.Name()) {
				cleanArtifacts(filepath.Join(cacheDir, info.Name()), lock)
				continue
			}

			if !utils.Contains(dependenciesNames, info.Name()) {
			remove:
				if err = os.RemoveAll(filepath.Join(cacheDir, info.Name())); err != nil {
					msg.Warn("Failed to remove old cache: %s", err.Error())
					goto remove
				}
			}

		}
	}
	utils.HandleError(os.MkdirAll(filepath.Join(cacheDir, consts.BplFolder), os.ModeDir|0755))
	utils.HandleError(os.MkdirAll(filepath.Join(cacheDir, consts.DcuFolder), os.ModeDir|0755))
	utils.HandleError(os.MkdirAll(filepath.Join(cacheDir, consts.DcpFolder), os.ModeDir|0755))
	utils.HandleError(os.MkdirAll(filepath.Join(cacheDir, consts.BinFolder), os.ModeDir|0755))
}
