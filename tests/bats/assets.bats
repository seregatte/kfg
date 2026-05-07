#!/usr/bin/env bats

load 'test_helper/bats-support/load'
load 'test_helper/bats-assert/load'

KFG="${BATS_TEST_DIRNAME}/../../bin/kfg"

setup() {
  TMPDIR=$(mktemp -d)
}

teardown() {
  rm -rf "$TMPDIR"
}

# Helper: create a manifest with Asset and Converter
create_manifest() {
  cat > "$TMPDIR/manifest.yaml" << 'EOF'
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-data
spec:
  input:
    format: yaml
  data:
    server: production
    port: 8080
    region: us-east-1
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: to-json
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: json
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: to-yaml
spec:
  input:
    format: yaml
  engine:
    expression: "."
  output:
    format: yaml
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: raw-output
spec:
  input:
    format: yaml
  engine:
    expression: ".server"
  output:
    format: raw
EOF
}

# Test: successful conversion with --convert and --use
@test "kfg apply --convert --use succeeds" {
  create_manifest

  result="$($KFG apply -f "$TMPDIR/manifest.yaml" --convert test-data --use to-json)"
  
  assert_success
  # Output should be JSON
  assert_output --partial '"server"'
  assert_output --partial '"production"'
}

# Test: --convert without --use fails
@test "kfg apply --convert without --use fails" {
  create_manifest

  run $KFG apply -f "$TMPDIR/manifest.yaml" --convert test-data
  
  assert_failure 2
}

# Test: --use without --convert fails
@test "kfg apply --use without --convert fails" {
  create_manifest

  run $KFG apply -f "$TMPDIR/manifest.yaml" --use to-json
  
  assert_failure 2
}

# Test: --convert/--use with -w fails
@test "kfg apply --convert/--use with -w fails" {
  create_manifest

  run $KFG apply -f "$TMPDIR/manifest.yaml" --convert test-data --use to-json -w my-workflow
  
  assert_failure 2
}

# Test: --convert/--use with -c fails
@test "kfg apply --convert/--use with -c fails" {
  create_manifest

  run $KFG apply -f "$TMPDIR/manifest.yaml" --convert test-data --use to-json -c my-cmd
  
  assert_failure 2
}

# Test: Asset not found
@test "kfg apply --convert with unknown asset fails" {
  create_manifest

  run $KFG apply -f "$TMPDIR/manifest.yaml" --convert nonexistent --use to-json
  
  assert_failure 1
  assert_output --partial "Asset not found"
}

# Test: Converter not found
@test "kfg apply --use with unknown converter fails" {
  create_manifest

  run $KFG apply -f "$TMPDIR/manifest.yaml" --convert test-data --use nonexistent
  
  assert_failure 1
  assert_output --partial "Converter not found"
}

# Test: output to file with -o
@test "kfg apply --convert --use -o writes to file" {
  create_manifest

  run $KFG apply -f "$TMPDIR/manifest.yaml" --convert test-data --use to-json -o "$TMPDIR/output.json"
  
  assert_success
  [ -f "$TMPDIR/output.json" ]
  # File should contain JSON
  grep -q '"server"' "$TMPDIR/output.json"
}

# Test: raw output format end-to-end
@test "kfg apply --convert --use raw output" {
  create_manifest

  result="$($KFG apply -f "$TMPDIR/manifest.yaml" --convert test-data --use raw-output)"
  
  assert_success
  assert_output "production"
}
