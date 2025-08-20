// Copyright 2025 Hedgehog
// SPDX-License-Identifier: Apache-2.0

package alloy_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.githedgehog.com/libmeta/pkg/alloy"
)

func TestConfigRender(t *testing.T) {
	tests := []struct {
		name   string
		config *alloy.Config
	}{
		{
			name:   "empty",
			config: &alloy.Config{},
		},
		{
			name: "targets-hostname-proxy",
			config: &alloy.Config{
				Hostname: "node-1",
				Targets: alloy.Targets{
					Prometheus: map[string]alloy.PrometheusTarget{
						"grafana_cloud": {
							Target: alloy.Target{
								URL: "https://prometheus-prod-36-prod-us-west-0.grafana.net/api/prom/push",
								BasicAuth: &alloy.TargetBasicAuth{
									Username: "username",
									Password: "password",
								},
								Labels: map[string]string{
									"a": "b",
								},
							},
						},
						"grafana_cloud_interval": {
							Target: alloy.Target{
								URL: "https://prometheus-prod-36-prod-us-west-0.grafana.net/api/prom/push",
								BasicAuth: &alloy.TargetBasicAuth{
									Username: "username",
									Password: "password",
								},
								Labels: map[string]string{
									"a": "b",
								},
							},
							SendIntervalSeconds: 42,
						},
						"another": {
							Target: alloy.Target{
								URL:                "https://another.com/prom/push",
								BearerToken:        "secret",
								InsecureSkipVerify: true,
							},
						},
					},
					Loki: map[string]alloy.LokiTarget{
						"grafana_cloud": {
							Target: alloy.Target{
								URL: "https://logs-prod-021.grafana.net/loki/api/v1/push",
								BasicAuth: &alloy.TargetBasicAuth{
									Username: "username",
									Password: "password",
								},
								Labels: map[string]string{
									"a": "b",
								},
							},
						},
						"another": {
							Target: alloy.Target{
								URL:                "https://another.com/loki/push",
								BearerToken:        "secret",
								InsecureSkipVerify: true,
							},
						},
					},
				},
				ProxyURL: "the-proxy",
			},
		},
		{
			name: "targets-no-hostname-no-proxy",
			config: &alloy.Config{
				Targets: alloy.Targets{
					Prometheus: map[string]alloy.PrometheusTarget{
						"grafana_cloud": {
							Target: alloy.Target{
								URL: "https://prometheus-prod-36-prod-us-west-0.grafana.net/api/prom/push",
								BasicAuth: &alloy.TargetBasicAuth{
									Username: "username",
									Password: "password",
								},
								Labels: map[string]string{
									"a": "b",
								},
							},
						},
						"grafana_cloud_interval": {
							Target: alloy.Target{
								URL: "https://prometheus-prod-36-prod-us-west-0.grafana.net/api/prom/push",
								BasicAuth: &alloy.TargetBasicAuth{
									Username: "username",
									Password: "password",
								},
								Labels: map[string]string{
									"a": "b",
								},
							},
							SendIntervalSeconds: 42,
						},
						"another": {
							Target: alloy.Target{
								URL:                "https://another.com/prom/push",
								BearerToken:        "secret",
								InsecureSkipVerify: true,
							},
						},
					},
					Loki: map[string]alloy.LokiTarget{
						"grafana_cloud": {
							Target: alloy.Target{
								URL: "https://logs-prod-021.grafana.net/loki/api/v1/push",
								BasicAuth: &alloy.TargetBasicAuth{
									Username: "username",
									Password: "password",
								},
								Labels: map[string]string{
									"a": "b",
								},
							},
						},
						"another": {
							Target: alloy.Target{
								URL:                "https://another.com/loki/push",
								BearerToken:        "secret",
								InsecureSkipVerify: true,
							},
						},
					},
				},
			},
		},
		{
			name: "full",
			config: &alloy.Config{
				Targets: alloy.Targets{
					Prometheus: map[string]alloy.PrometheusTarget{
						"grafana_cloud": {
							Target: alloy.Target{
								URL: "https://prometheus-prod-36-prod-us-west-0.grafana.net/api/prom/push",
								BasicAuth: &alloy.TargetBasicAuth{
									Username: "username",
									Password: "password",
								},
								Labels: map[string]string{
									"a": "b",
								},
							},
						},
						"grafana_cloud_interval": {
							Target: alloy.Target{
								URL: "https://prometheus-prod-36-prod-us-west-0.grafana.net/api/prom/push",
								BasicAuth: &alloy.TargetBasicAuth{
									Username: "username",
									Password: "password",
								},
								Labels: map[string]string{
									"a": "b",
								},
							},
							SendIntervalSeconds: 42,
						},
						"another": {
							Target: alloy.Target{
								URL:                "https://another.com/prom/push",
								BearerToken:        "secret",
								InsecureSkipVerify: true,
							},
						},
					},
					Loki: map[string]alloy.LokiTarget{
						"grafana_cloud": {
							Target: alloy.Target{
								URL: "https://logs-prod-021.grafana.net/loki/api/v1/push",
								BasicAuth: &alloy.TargetBasicAuth{
									Username: "username",
									Password: "password",
								},
								Labels: map[string]string{
									"a": "b",
								},
							},
						},
						"another": {
							Target: alloy.Target{
								URL:                "https://another.com/loki/push",
								BearerToken:        "secret",
								InsecureSkipVerify: true,
							},
						},
					},
				},
				Scrapes: map[string]alloy.Scrape{
					"test_address": {
						IntervalSeconds: 42,
						Address:         "localhost:12345",
					},
					"test_relabel": {
						IntervalSeconds: 42,
						Address:         "localhost:12345",
						Relabel: []alloy.ScrapeRelabelRule{
							{
								SourceLabels: []string{"l1", "l2"},
								TargetLabel:  "t1",
								Separator:    ";",
								Replacement:  "$1",
								Regex:        "r2",
								Action:       "drop",
							},
							{
								SourceLabels: []string{"l1"},
								Action:       "drop",
							},
						},
					},
					"test_self": {
						IntervalSeconds: 43,
						Self: alloy.ScrapeSelf{
							Enable: true,
						},
					},
					"test_unix": {
						IntervalSeconds: 44,
						Unix: alloy.ScrapeUnix{
							Enable:     true,
							Collectors: []string{"asd"},
						},
					},
				},
				LogFiles: map[string]alloy.LogFile{
					"syslog": {
						PathTargets: []alloy.LogFilePathTarget{
							{
								Path: "/var/log/syslog",
							},
						},
					},
					"varlog": {
						PathTargets: []alloy.LogFilePathTarget{
							{
								Path:        "/var/log/*.log",
								PathExclude: "/var/log/agent.log",
							},
						},
					},
				},
				Kube: alloy.Kube{
					PodLogs: true,
					Events:  true,
				},
			},
		},
		{
			name: "no-targets",
			config: &alloy.Config{
				Scrapes: map[string]alloy.Scrape{
					"test_address": {
						IntervalSeconds: 42,
						Address:         "localhost:12345",
					},
					"test_self": {
						IntervalSeconds: 43,
						Self: alloy.ScrapeSelf{
							Enable: true,
						},
					},
					"test_unix": {
						IntervalSeconds: 44,
						Unix: alloy.ScrapeUnix{
							Enable:     true,
							Collectors: []string{"asd"},
						},
					},
				},
				LogFiles: map[string]alloy.LogFile{
					"syslog": {
						PathTargets: []alloy.LogFilePathTarget{
							{
								Path: "/var/log/syslog",
							},
						},
					},
				},
				Kube: alloy.Kube{
					PodLogs: true,
					Events:  true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileName := strings.ReplaceAll(t.Name(), "/", "_")
			expectedFileName := filepath.Join("testdata", fileName+".expected")
			actualFileName := filepath.Join("testdata", fileName+".actual")

			actual, err := tt.config.Render()
			require.NoError(t, err)

			if os.Getenv("UPDATE") == "true" {
				err = os.WriteFile(expectedFileName, []byte(actual), 0o644)
				require.NoError(t, err)

				currentDir, err := os.Getwd()
				require.NoError(t, err)

				out, err := exec.CommandContext(t.Context(),
					"docker", "run",
					"--rm", "--pull", "always",
					"-v", currentDir+":/config:ro",
					"grafana/alloy:latest",
					"validate", filepath.Join("/config", expectedFileName)).CombinedOutput()
				if !assert.NoError(t, err) {
					fmt.Println("Error:", err)
					fmt.Println("Output:", string(out))
				}

				return
			}

			expected, err := os.ReadFile(expectedFileName)
			require.NoError(t, err)

			assert.Equal(t, strings.TrimSpace(string(expected)), strings.TrimSpace(string(actual)))

			err = os.WriteFile(actualFileName, []byte(actual), 0o644)
			require.NoError(t, err)
		})
	}
}
