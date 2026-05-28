#!/usr/bin/env bash

# Backward-compatible test_helper that loads from the new helpers structure
# This file is kept for compatibility with existing test files.
# New tests should load helpers directly from helpers/common.bash

# Get the directory where this helper lives
TEST_HELPER_DIR="$(cd "${BASH_SOURCE%/*}" && pwd)"

# Source the common helpers
source "${TEST_HELPER_DIR}/helpers/common.bash"