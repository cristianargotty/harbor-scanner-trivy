package trivy

import (
	"github.com/aquasecurity/harbor-scanner-trivy/pkg/etc"
	"github.com/aquasecurity/harbor-scanner-trivy/pkg/ext"
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
)

var (
	expectedReportJSON = `[{
	"Target": "alpine:3.10.2",
	"Vulnerabilities": [{
		"VulnerabilityID": "CVE-2018-6543",
		"PkgName": "binutils",
		"InstalledVersion": "2.30-r1",
		"FixedVersion": "2.30-r2",
		"Severity": "MEDIUM",
		"References": [
			"https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2018-6543"
		],
		"LayerID": "sha256:5216338b40a7b96416b8b9858974bbe4acc3096ee60acbc4dfb1ee02aecceb10"
	}]
}]`
	expectedReport = ScanReport{
		Target: "alpine:3.10.2",
		Vulnerabilities: []Vulnerability{
			{
				VulnerabilityID:  "CVE-2018-6543",
				PkgName:          "binutils",
				InstalledVersion: "2.30-r1",
				FixedVersion:     "2.30-r2",
				Severity:         "MEDIUM",
				References: []string{
					"https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2018-6543",
				},
				LayerID: "sha256:5216338b40a7b96416b8b9858974bbe4acc3096ee60acbc4dfb1ee02aecceb10",
			},
		},
	}
)

func TestWrapper_Scan(t *testing.T) {
	ambassador := ext.NewMockAmbassador()
	ambassador.On("Environ").Return([]string{"HTTP_PROXY=http://someproxy:7777"})
	ambassador.On("LookPath", "trivy").Return("/usr/local/bin/trivy", nil)

	config := etc.Trivy{
		CacheDir:      "/home/scanner/.cache/trivy",
		ReportsDir:    "/home/scanner/.cache/reports",
		DebugMode:     true,
		VulnType:      "os,library",
		Severity:      "CRITICAL,MEDIUM",
		IgnoreUnfixed: true,
		SkipUpdate:    true,
		GitHubToken:   "<github_token>",
	}

	imageRef := ImageRef{
		Name:     "alpine:3.10.2",
		Auth:     RegistryAuth{Username: "dave.loper", Password: "s3cret"},
		Insecure: true,
	}

	expectedCmdArgs := []string{
		"/usr/local/bin/trivy",
		"--skip-update",
		"--debug",
		"--ignore-unfixed",
		"--no-progress",
		"--cache-dir",
		"/home/scanner/.cache/trivy",
		"--severity",
		"CRITICAL,MEDIUM",
		"--vuln-type",
		"os,library",
		"--format",
		"json",
		"--output",
		"/home/scanner/.cache/reports/scan_report_1234567890.json",
		"alpine:3.10.2",
	}

	expectedCmdEnvs := []string{
		"HTTP_PROXY=http://someproxy:7777",
		"TRIVY_USERNAME=dave.loper",
		"TRIVY_PASSWORD=s3cret",
		"TRIVY_NON_SSL=true",
		"GITHUB_TOKEN=<github_token>",
	}

	ambassador.On("TempFile", "/home/scanner/.cache/reports", "scan_report_*.json").
		Return(ext.NewFakeFile("/home/scanner/.cache/reports/scan_report_1234567890.json", expectedReportJSON), nil)
	ambassador.On("Remove", "/home/scanner/.cache/reports/scan_report_1234567890.json").
		Return(nil)
	ambassador.On("RunCmd", &exec.Cmd{
		Path: "/usr/local/bin/trivy",
		Env:  expectedCmdEnvs,
		Args: expectedCmdArgs},
	).Return([]byte{}, nil)

	report, err := NewWrapper(config, ambassador).Scan(imageRef)

	require.NoError(t, err)
	require.Equal(t, expectedReport, report)

	ambassador.AssertExpectations(t)
}
