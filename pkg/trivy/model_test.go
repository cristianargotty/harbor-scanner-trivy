package trivy

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
	"strings"
	"testing"
)

const (
	emptyReport  = `[]`
	sampleReport = `[{
	"Target": "knqyf263/vuln-image (alpine 3.7.1)",
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
},
{
	"Target": "node-app/package-lock.json",
	"Vulnerabilities": [{
		"VulnerabilityID": "CVE-2019-11358",
		"PkgName": "jquery",
		"InstalledVersion": "3.3.9",
		"FixedVersion": "3.4.0",
		"Severity": "MEDIUM",
		"References": [
			"http://packetstormsecurity.com/files/152787/dotCMS-5.1.1-Vulnerable-Dependencies.html",
			"http://packetstormsecurity.com/files/153237/RetireJS-CORS-Issue-Script-Execution.html"
		],
		"LayerID": "sha256:5216338b40a7b96416b8b9858974bbe4acc3096ee60acbc4dfb1ee02aecceb10"
	}]
},
{
	"Target": "php-app/composer.lock",
	"Vulnerabilities": [{
		"VulnerabilityID": "CVE-2016-5385",
		"PkgName": "guzzlehttp/guzzle",
		"InstalledVersion": "6.2.0",
		"FixedVersion": "6.2.1",
		"Severity": "MEDIUM",
		"References": [
			"http://linux.oracle.com/cve/CVE-2016-5385.html"
		],
		"LayerID": "sha256:5216338b40a7b96416b8b9858974bbe4acc3096ee60acbc4dfb1ee02aecceb11"
	}]
},
{
	"Target": "python-app/Pipfile.lock",
	"Vulnerabilities": null
}
]`
)

func TestScanReportFrom(t *testing.T) {
	testCases := []struct {
		name           string
		jsonOutput     string
		expectedReport ScanReport
		expectedError  error
	}{
		{
			name:       "Should parse scan report with application dependencies",
			jsonOutput: sampleReport,
			expectedReport: ScanReport{
				Target: "knqyf263/vuln-image (alpine 3.7.1)",
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
					{
						VulnerabilityID:  "CVE-2019-11358",
						PkgName:          "jquery",
						InstalledVersion: "3.3.9",
						FixedVersion:     "3.4.0",
						Severity:         "MEDIUM",
						References: []string{
							"http://packetstormsecurity.com/files/152787/dotCMS-5.1.1-Vulnerable-Dependencies.html",
							"http://packetstormsecurity.com/files/153237/RetireJS-CORS-Issue-Script-Execution.html",
						},
						LayerID: "sha256:5216338b40a7b96416b8b9858974bbe4acc3096ee60acbc4dfb1ee02aecceb10",
					},
					{
						VulnerabilityID:  "CVE-2016-5385",
						PkgName:          "guzzlehttp/guzzle",
						InstalledVersion: "6.2.0",
						FixedVersion:     "6.2.1",
						Severity:         "MEDIUM",
						References: []string{
							"http://linux.oracle.com/cve/CVE-2016-5385.html",
						},
						LayerID: "sha256:5216338b40a7b96416b8b9858974bbe4acc3096ee60acbc4dfb1ee02aecceb11",
					},
				},
			},
		},
		{
			name:          "Should return error when scan report is empty",
			jsonOutput:    emptyReport,
			expectedError: xerrors.New("expected at least one report"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			report, err := ScanReportFrom(strings.NewReader(tc.jsonOutput))
			assert.Equal(t, tc.expectedReport, report)
			switch {
			case tc.expectedError != nil:
				assert.EqualError(t, err, tc.expectedError.Error())
			default:
				assert.NoError(t, err)
			}
		})
	}
}
