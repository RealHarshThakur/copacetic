package vex

import (
	"testing"
	"time"

	"github.com/project-copacetic/copacetic/pkg/buildkit"
	"github.com/project-copacetic/copacetic/pkg/pkgmgr"
	"github.com/project-copacetic/copacetic/pkg/types"
)

func TestOpenVex_CreateVEXDocument(t *testing.T) {
	config := &buildkit.Config{}
	workingFolder := "/tmp"
	alpineManager, _ := pkgmgr.GetPackageManager("alpine", config, workingFolder)
	debianManager, _ := pkgmgr.GetPackageManager("debian", config, workingFolder)
	t.Setenv("COPA_VEX_AUTHOR", "test author")

	// mock time
	expectedTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	now = func() time.Time { return expectedTime }

	// mock id
	id = func() (string, error) { return "https://openvex.dev/test", nil }

	type args struct {
		updates *types.UpdateManifest
		pkgmgr  pkgmgr.PackageManager
	}
	tests := []struct {
		name    string
		o       *OpenVex
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid openvex document",
			o:    &OpenVex{},
			args: args{
				updates: &types.UpdateManifest{
					Updates: []types.UpdatePackage{
						{
							Name:             "test1",
							InstalledVersion: "1.0",
							FixedVersion:     "1.1",
							VulnerabilityID:  "CVE-2020-1234",
						},
					},
					OSType: "alpine",
					Arch:   "x86_64",
				},
				pkgmgr: alpineManager,
			},
			want: `{
  "@context": "https://openvex.dev/ns",
  "@id": "https://openvex.dev/test",
  "author": "test author",
  "timestamp": "2009-11-17T20:34:58.651387237Z",
  "version": 1,
  "tooling": "Project Copacetic",
  "statements": [
    {
      "vulnerability": {
        "@id": "CVE-2020-1234"
      },
      "products": [
        {
          "@id": "pkg:apk/alpine/test1@1.1?arch=x86_64"
        }
      ],
      "status": "fixed"
    }
  ]
}
`,
			wantErr: false,
		},
		{
			name: "valid openvex document with multiple statements and multiple vulnerabilities",
			o:    &OpenVex{},
			args: args{
				updates: &types.UpdateManifest{
					Updates: []types.UpdatePackage{
						{
							Name:             "test2",
							InstalledVersion: "1.0",
							FixedVersion:     "1.2",
							VulnerabilityID:  "CVE-2020-1234",
						},
						{
							Name:             "test3",
							InstalledVersion: "1.0",
							FixedVersion:     "1.3",
							VulnerabilityID:  "CVE-2020-1235",
						},
					},
					OSType: "debian",
					Arch:   "x86_64",
				},
				pkgmgr: debianManager,
			},
			want: `{
  "@context": "https://openvex.dev/ns",
  "@id": "https://openvex.dev/test",
  "author": "test author",
  "timestamp": "2009-11-17T20:34:58.651387237Z",
  "version": 1,
  "tooling": "Project Copacetic",
  "statements": [
    {
      "vulnerability": {
        "@id": "CVE-2020-1234"
      },
      "products": [
        {
          "@id": "pkg:apk/alpine/test1@1.1?arch=x86_64"
        },
        {
          "@id": "pkg:deb/debian/test2@1.2?arch=x86_64"
        }
      ],
      "status": "fixed"
    },
    {
      "vulnerability": {
        "@id": "CVE-2020-1235"
      },
      "products": [
        {
          "@id": "pkg:deb/debian/test3@1.3?arch=x86_64"
        }
      ],
      "status": "fixed"
    }
  ]
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OpenVex{}
			got, err := o.CreateVEXDocument(tt.args.updates, tt.args.pkgmgr)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenVex.CreateVEXDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OpenVex.CreateVEXDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}
