package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/skuid/helm-value-store/store"
	"github.com/skuid/spec"
	"go.uber.org/zap"
	"k8s.io/helm/pkg/helm"
	hapi_release5 "k8s.io/helm/pkg/proto/hapi/release"
)

var client *helm.Client

func init() {
	client = helm.NewClient(
		helm.Host(
			os.Getenv("TILLER_HOST"),
		),
	)

	err := os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func SyncReleases(rs store.ReleaseStore, labels map[string]string, interval time.Duration, blacklist []string) {
	for {
		spec.Logger.Info("Syncing releases")
		go syncReleases(rs, labels, blacklist)
		time.Sleep(interval)
	}
}

func helmReleasesToMap(releases []*hapi_release5.Release) map[string]*hapi_release5.Release {
	response := map[string]*hapi_release5.Release{}
	for _, r := range releases {
		response[r.Name] = r
	}
	return response
}

func syncReleases(rs store.ReleaseStore, labels map[string]string, blacklist []string) {
	storeReleases, err := rs.List(labels)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	releases, err := client.ListReleases()
	if err != nil {
		spec.Logger.Error("Error listing releases", zap.Error(err))
		return
	}
	releaseMap := helmReleasesToMap(releases.Releases)

	for _, storeRelease := range storeReleases {
		blacklisted := false
		for _, b := range blacklist {
			if storeRelease.Name == b {
				blacklisted = true
				break
			}

		}
		if blacklisted {
			continue
		}

		release, ok := releaseMap[storeRelease.Name]
		if !ok {
			spec.Logger.Info(
				"Release not found, installing ",
				zap.String("release", storeRelease.Name),
				zap.String("namespace", storeRelease.Namespace),
				zap.String("version", storeRelease.Version),
				zap.String("chart", storeRelease.Chart),
			)
			// TODO install release
			continue
		}
		// TODO evaluate if version or values are different

		if release.Namespace != storeRelease.Namespace {
			spec.Logger.Error(
				"Namespaces differ for release!",
				zap.String("release", release.Name),
				zap.String("installed namespace", release.Namespace),
				zap.String("desired namespace", storeRelease.Namespace),
			)
			continue
		}

		if release.Chart.Metadata.Version != storeRelease.Version {
			spec.Logger.Info(
				fmt.Sprintf("Upgrading chart from %s to %s", release.Chart.Metadata.Version, storeRelease.Version),
				zap.String("release", release.Name),
			)
			// TODO upgrade chart version
			continue
		}
		if !(release.Config.Raw == storeRelease.Values || (release.Config.Raw == "{}\n" && len(storeRelease.Values) == 0)) {
			spec.Logger.Info(
				"Values differ for release!",
				zap.String("release", release.Name),
				zap.String("store_values", storeRelease.Values),
				zap.String("release_values", release.Config.Raw),
			)
			continue
		}

	}

	spec.Logger.Info("Completed sync")
}
