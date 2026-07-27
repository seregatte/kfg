package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveAIOverlay_localExists(t *testing.T) {
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origWd)

	localPath := filepath.FromSlash(localAIOverlay)
	err := os.MkdirAll(localPath, 0755)
	assert.NoError(t, err)

	overlay, err := resolveAIOverlay()
	assert.NoError(t, err)
	assert.Equal(t, localAIOverlay, overlay)
}

func TestResolveAIOverlay_remoteFallback(t *testing.T) {
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origWd)

	overlay, err := resolveAIOverlay()
	assert.NoError(t, err)
	assert.Equal(t, remoteAIOverlay, overlay)
}

func TestResolveAIOverlay_inspectionError(t *testing.T) {
	dir := t.TempDir()
	origWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origWd)

	firstDir := filepath.FromSlash("packages")
	err := os.MkdirAll(firstDir, 0755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.FromSlash("packages/domains"), []byte("not-a-directory"), 0644)
	assert.NoError(t, err)

	overlay, err := resolveAIOverlay()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot inspect")
	assert.Empty(t, overlay)
}

func TestBuildRunArgs_localOverlay(t *testing.T) {
	args := buildRunArgs(localAIOverlay, "pi", nil)
	assert.Equal(t, []string{"run", "-k", localAIOverlay, "pi"}, args)
}

func TestBuildRunArgs_remoteOverlay(t *testing.T) {
	args := buildRunArgs(remoteAIOverlay, "opencode", nil)
	assert.Equal(t, []string{"run", "-k", remoteAIOverlay, "opencode"}, args)
}

func TestBuildRunArgs_withExtraArgs(t *testing.T) {
	args := buildRunArgs(localAIOverlay, "pi", []string{"--model", "sonnet", "create project"})
	assert.Equal(t, []string{"run", "-k", localAIOverlay, "pi", "--", "--model", "sonnet", "create project"}, args)
}

func TestBuildRunArgs_localOverlayWithOpencodeAndArgs(t *testing.T) {
	args := buildRunArgs(localAIOverlay, "opencode", []string{"create a new project"})
	assert.Equal(t, []string{"run", "-k", localAIOverlay, "opencode", "--", "create a new project"}, args)
}

func TestResolveAgent_default(t *testing.T) {
	orig := os.Getenv("KFG_AI_AGENT")
	os.Unsetenv("KFG_AI_AGENT")
	defer os.Setenv("KFG_AI_AGENT", orig)

	agent := resolveAgent()
	assert.Equal(t, "pi", agent)
}

func TestResolveAgent_custom(t *testing.T) {
	orig := os.Getenv("KFG_AI_AGENT")
	os.Setenv("KFG_AI_AGENT", "opencode")
	defer os.Setenv("KFG_AI_AGENT", orig)

	agent := resolveAgent()
	assert.Equal(t, "opencode", agent)
}

func TestResolveAgent_unsupportedFallback(t *testing.T) {
	orig := os.Getenv("KFG_AI_AGENT")
	os.Setenv("KFG_AI_AGENT", "invalid-agent")
	defer os.Setenv("KFG_AI_AGENT", orig)

	agent := resolveAgent()
	assert.Equal(t, "pi", agent)
}
