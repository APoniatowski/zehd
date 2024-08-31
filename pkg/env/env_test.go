package env

import (
	"os"
	"testing"
)

func TestEnvPort(t *testing.T) {
	// Test with environment variable set
	os.Setenv("PORT", "8080")
	want := ":8080"
	got := EnvPort()
	if got != want {
		t.Errorf("EnvPort() = %q, want %q", got, want)
	}

	// Test with environment variable unset
	os.Unsetenv("PORT")
	want = ":80"
	got = EnvPort()
	if got != want {
		t.Errorf("EnvPort() = %q, want %q", got, want)
	}
}

func TestEnvHostname(t *testing.T) {
	// Save the original hostname environment variable
	originalHostname := os.Getenv("NEWHOSTNAME")
	defer os.Setenv("NEWHOSTNAME", originalHostname)

	testCases := []struct {
		name           string
		setEnvHostname string
		want           string
	}{
		{
			name:           "With NEWHOSTNAME env var",
			setEnvHostname: "custom-hostname",
			want:           "custom-hostname",
		},
		{
			name:           "Without NEWHOSTNAME env var",
			setEnvHostname: "",
			want:           "", // The actual hostname will be set by the OS
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnvHostname != "" {
				os.Setenv("NEWHOSTNAME", tc.setEnvHostname)
			} else {
				os.Unsetenv("NEWHOSTNAME")
			}

			got := EnvHostname()

			if tc.setEnvHostname != "" {
				if got != tc.want {
					t.Errorf("EnvHostname() = %q, want %q", got, tc.want)
				}
			} else {
				// When HOSTNAME is not set, we expect the OS hostname
				osHostname, _ := os.Hostname()
				if got != osHostname {
					t.Errorf("EnvHostname() = %q, want OS hostname %q", got, osHostname)
				}
			}
		})
	}
}

func TestEnvCacheRefresh(t *testing.T) {
	// Test with environment variable set
	os.Setenv("REFRESHCACHE", "120")
	want := "120"
	wantScore := 2
	got, gotScore := EnvCacheRefresh()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvCacheRefresh() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}

	// Test with environment variable unset
	os.Unsetenv("REFRESHCACHE")
	want = "60"
	wantScore = 1
	got, gotScore = EnvCacheRefresh()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvCacheRefresh() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}
}

func TestEnvTemplateDir(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEMPLATEDIRECTORY", "/custom/templates/")
	want := "/custom/templates/"
	wantScore := 2
	got, gotScore := EnvTemplateDir()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvTemplateDir() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}

	// Test with environment variable unset
	os.Unsetenv("TEMPLATEDIRECTORY")
	want = "/var/frontend/templates/"
	wantScore = 1
	got, gotScore = EnvTemplateDir()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvTemplateDir() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}
}

func TestEnvTemplateType(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEMPLATETYPE", "html")
	want := "html"
	wantScore := 2
	got, gotScore := EnvTemplateType()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvTemplateType() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}

	// Test with environment variable unset
	os.Unsetenv("TEMPLATETYPE")
	want = "gohtml"
	wantScore = 1
	got, gotScore = EnvTemplateType()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvTemplateType() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}
}

func TestEnvGitLink(t *testing.T) {
	// Test with environment variable set
	os.Setenv("GITLINK", "https://github.com/example/repo.git")
	want := "https://github.com/example/repo.git"
	wantScore := 2
	got, gotScore := EnvGitLink()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvGitLink() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}

	// Test with environment variable unset
	os.Unsetenv("GITLINK")
	want = ""
	wantScore = 1
	got, gotScore = EnvGitLink()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvGitLink() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}
}

func TestEnvGitToken(t *testing.T) {
	// Test with environment variable set
	os.Setenv("GITTOKEN", "abc123")
	want := "abc123"
	wantScore := 2
	got, gotScore := EnvGitToken()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvGitToken() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}

	// Test with environment variable unset
	os.Unsetenv("GITTOKEN")
	want = ""
	wantScore = 1
	got, gotScore = EnvGitToken()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvGitToken() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}
}

func TestEnvGitUsername(t *testing.T) {
	// Test with environment variable set
	os.Setenv("GITUSERNAME", "example-user")
	want := "example-user"
	wantScore := 2
	got, gotScore := EnvGitUsername()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvGitUsername() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}

	// Test with environment variable unset
	os.Unsetenv("GITUSERNAME")
	want = ""
	wantScore = 1
	got, gotScore = EnvGitUsername()
	if got != want || gotScore != wantScore {
		t.Errorf("EnvGitUsername() = (%q, %d), want (%q, %d)", got, gotScore, want, wantScore)
	}
}

func TestEnvPathDefiner(t *testing.T) {
	testCases := []struct {
		name        string
		envVar      string
		envValue    string
		input       string
		wantPath    string
		wantScore   int
		unsetEnvVar bool
	}{
		{"CSS Path Set", "CSSPATH", "/custom/css", "css", "/custom/css", 2, false},
		{"CSS Path Unset", "CSSPATH", "", "css", "css", 1, true},
		{"JS Path Set", "JSPATH", "/custom/js", "js", "/custom/js", 2, false},
		{"JS Path Unset", "JSPATH", "", "js", "js", 1, true},
		{"Images Path Set", "IMAGESPATH", "/custom/images", "images", "/custom/images", 2, false},
		{"Images Path Unset", "IMAGESPATH", "", "images", "images", 1, true},
		{"Downloads Path Set", "DOWNLOADSPATH", "/custom/downloads", "downloads", "/custom/downloads", 2, false},
		{"Downloads Path Unset", "DOWNLOADSPATH", "", "downloads", "downloads", 1, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.unsetEnvVar {
				os.Unsetenv(tc.envVar)
			} else {
				os.Setenv(tc.envVar, tc.envValue)
			}

			gotPath, gotScore := EnvPathDefiner(tc.input)
			if gotPath != tc.wantPath || gotScore != tc.wantScore {
				t.Errorf("EnvPathDefiner(%q) = (%q, %d), want (%q, %d)",
					tc.input, gotPath, gotScore, tc.wantPath, tc.wantScore)
			}
		})
	}
}
