package cliTable

import (
	"path/filepath"
	"runtime"
	"testing"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type testCase struct {
	Attrs            map[string]string
	ExpectedTemplate string
}

func TestCliTable_FindTemplate(t *testing.T) {

	tests := []struct {
		IndexFileName string
		TestCases     []testCase
	}{

		{
			IndexFileName: "index",
			TestCases: []testCase{
				{
					Attrs:            map[string]string{"Command": "sh ip bgp sum", "Vendor": "Cisco"},
					ExpectedTemplate: "cisco_bgp_summary_template",
				},
				{
					Attrs:            map[string]string{"Command": "show version", "Vendor": "Force10"},
					ExpectedTemplate: "f10_version_template",
				},
				{
					Attrs:            map[string]string{"Command": "show ip bgp summary", "Vendor": "Force10"},
					ExpectedTemplate: "f10_ip_bgp_summary_template",
				},
				{
					Attrs:            map[string]string{"Command": "show bgp summary", "Vendor": "Juniper"},
					ExpectedTemplate: "juniper_bgp_summary_template",
				},
				{
					Attrs:            map[string]string{"Command": "show bgp summary", "Vendor": "Juniper"},
					ExpectedTemplate: "juniper_bgp_summary_template",
				},
				{
					Attrs:            map[string]string{"Command": "sh ve", "Vendor": "Juniper"},
					ExpectedTemplate: "juniper_version_template",
				},
			},
		},
		{
			IndexFileName: "index2",
			TestCases: []testCase{
				{
					Attrs:            map[string]string{"Command": "sh count global", "Platform": "paloalto_panos"},
					ExpectedTemplate: "paloalto_panos_show_counter_global.template",
				},
				{
					Attrs:            map[string]string{"Command": "show interface", "Platform": "cisco_asa"},
					ExpectedTemplate: "cisco_asa_show_interface.template",
				},
				{
					Attrs:            map[string]string{"Command": "sho reload cause", "Platform": "arista_eos"},
					ExpectedTemplate: "arista_eos_show_reload_cause.template",
				},
				{
					Attrs:            map[string]string{"Command": "show ip bgp neig", "Platform": "vmware_nsxv"},
					ExpectedTemplate: "vmware_nsxv_show_ip_bgp_neighbors.template",
				},
				{
					Attrs:            map[string]string{"Command": "show version", "Platform": "cisco_asa"},
					ExpectedTemplate: "cisco_asa_show_version.template",
				},
			},
		},
	}

	for testIndx, test := range tests {
		templateDir := basepath + "/../../testfiles/"
		T := NewCliTable(templateDir, test.IndexFileName)
		for caseIndx, testCase := range test.TestCases {
			tmp, _ := T.FindTemplate(testCase.Attrs)

			if tmp != testCase.ExpectedTemplate {
				t.Errorf("%d:%d fail: expected template is '%s', got '%s'", testIndx, caseIndx, testCase.ExpectedTemplate, tmp)
			}
		}

	}

}
