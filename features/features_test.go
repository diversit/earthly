package features

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFeaturesStringEnabled(t *testing.T) {
	fts := &Features{
		Major:              0,
		Minor:              5,
		ReferencedSaveOnly: true,
	}
	s := fts.String()
	Equal(t, "VERSION --referenced-save-only 0.5", s)
}

func TestFeaturesStringDisabled(t *testing.T) {
	fts := &Features{
		Major:              1,
		Minor:              1,
		ReferencedSaveOnly: false,
	}
	s := fts.String()
	Equal(t, "VERSION 1.1", s)
}

func TestApplyFlagOverrides(t *testing.T) {
	fts := &Features{}
	err := ApplyFlagOverrides(fts, "referenced-save-only")
	Nil(t, err)
	Equal(t, true, fts.ReferencedSaveOnly)
	Equal(t, false, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, false, fts.NoImplicitIgnore)
}

func TestApplyFlagOverridesWithDashDashPrefix(t *testing.T) {
	fts := &Features{}
	err := ApplyFlagOverrides(fts, "--referenced-save-only")
	Nil(t, err)
	Equal(t, true, fts.ReferencedSaveOnly)
	Equal(t, false, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, false, fts.NoImplicitIgnore)
}

func TestApplyFlagOverridesMultipleFlags(t *testing.T) {
	fts := &Features{}
	err := ApplyFlagOverrides(fts, "referenced-save-only,use-copy-include-patterns,no-implicit-ignore")
	Nil(t, err)
	Equal(t, true, fts.ReferencedSaveOnly)
	Equal(t, true, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, true, fts.NoImplicitIgnore)
}

func TestApplyFlagOverridesEmptyString(t *testing.T) {
	fts := &Features{}
	err := ApplyFlagOverrides(fts, "")
	Nil(t, err)
	Equal(t, false, fts.ReferencedSaveOnly)
	Equal(t, false, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, false, fts.NoImplicitIgnore)
}

func TestAvailableFlags(t *testing.T) {
	// This test feels like it may be overkill, but it's nice to know that if we
	// introduce a typo we have to introduce it twice for our tests to still
	// pass.
	for _, tt := range []struct {
		flag  string
		field string
	}{
		// 0.5
		{"exec-after-parallel", "ExecAfterParallel"},
		{"parallel-load", "ParallelLoad"},
		{"use-registry-for-with-docker", "UseRegistryForWithDocker"},

		// 0.6
		{"for-in", "ForIn"},
		{"no-implicit-ignore", "NoImplicitIgnore"},
		{"referenced-save-only", "ReferencedSaveOnly"},
		{"require-force-for-unsafe-saves", "RequireForceForUnsafeSaves"},
		{"use-copy-include-patterns", "UseCopyIncludePatterns"},

		// 0.7
		{"check-duplicate-images", "CheckDuplicateImages"},
		{"ci-arg", "EarthlyCIArg"},
		{"earthly-git-author-args", "EarthlyGitAuthorArgs"},
		{"earthly-locally-arg", "EarthlyLocallyArg"},
		{"earthly-version-arg", "EarthlyVersionArg"},
		{"explicit-global", "ExplicitGlobal"},
		{"git-commit-author-timestamp", "GitCommitAuthorTimestamp"},
		{"new-platform", "NewPlatform"},
		{"no-tar-build-output", "NoTarBuildOutput"},
		{"save-artifact-keep-own", "SaveArtifactKeepOwn"},
		{"shell-out-anywhere", "ShellOutAnywhere"},
		{"use-cache-command", "UseCacheCommand"},
		{"use-chmod", "UseChmod"},
		{"use-copy-link", "UseCopyLink"},
		{"use-host-command", "UseHostCommand"},
		{"use-no-manifest-list", "UseNoManifestList"},
		{"use-pipelines", "UsePipelines"},
		{"use-project-secrets", "UseProjectSecrets"},
		{"wait-block", "WaitBlock"},

		// unreleased
		{"no-use-registry-for-with-docker", "NoUseRegistryForWithDocker"},
		{"try", "TryFinally"},
		{"arg-scope-set", "ArgScopeSet"},
	} {
		tt := tt
		t.Run(tt.flag, func(t *testing.T) {
			t.Parallel()

			var fts Features
			err := ApplyFlagOverrides(&fts, tt.flag)
			Nil(t, err)
			field := reflect.ValueOf(fts).FieldByName(tt.field)
			True(t, field.IsValid(), "field %v does not exist on %T", tt.field, fts)
			val, ok := field.Interface().(bool)
			True(t, ok, "field %v was not a boolean", tt.field)
			True(t, val, "expected field %v to be set to true by flag %v", tt.field, tt.flag)
		})
	}
}

func TestVersionAtLeast(t *testing.T) {
	tests := []struct {
		earthlyVer Features
		major      int
		minor      int
		expected   bool
	}{
		{
			earthlyVer: Features{Major: 0, Minor: 6},
			major:      0,
			minor:      5,
			expected:   true,
		},
		{
			earthlyVer: Features{Major: 0, Minor: 6},
			major:      0,
			minor:      7,
			expected:   false,
		},
		{
			earthlyVer: Features{Major: 0, Minor: 6},
			major:      1,
			minor:      2,
			expected:   false,
		},
		{
			earthlyVer: Features{Major: 1, Minor: 2},
			major:      1,
			minor:      2,
			expected:   true,
		},
	}
	for _, test := range tests {
		title := fmt.Sprintf("earthly version %d.%d is at least %d.%d",
			test.earthlyVer.Major, test.earthlyVer.Minor, test.major, test.minor)
		t.Run(title, func(t *testing.T) {
			actual := versionAtLeast(test.earthlyVer, test.major, test.minor)
			Equal(t, test.expected, actual)
		})
	}
}
