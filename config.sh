# SPDX-FileCopyrightText: AISEC Pentesting Team
#
# SPDX-License-Identifier: Apache-2.0

################################################################################
# Penrun config
# NOTE: These variables are sourced by the shell.
# They are not required to be exported, but could be.
# Please note that bash arrays cannot be exported!

# If this variable is set then penrun creates the artifacts folder hierarchy
# at this location instead of $PWD.
#
# shellcheck disable=SC2034
PENRUN_ARTIFACTS_BASE="$HOME/penrun-artifacts"

# This variable specifies the compression tool where the output is piped
# to. Since zstd provides multithreading capabilities, it is the fastest
# in all out-of-the-box scenarios. Argument must be a bash array.
#
# PENRUN_COMPRESSION_COMMAND=("gzip" "--stdout")

# shellcheck disable=SC2034
PENRUN_COMPRESSION_COMMAND=("zstd" "-T0")

# Specify the extension that is added to the OUTPUT file.
# shellcheck disable=SC2034
PENRUN_OUTPUT_EXTENSION=".zst"

# Pipe penrun output to HR
# Argument must be a bash array.
# PENRUN_PIPE_COMMAND=("hr" "-p" "info")

# Add Default arguments which be appended to each penrun run
# Argument must be a bash array.
# PENRUN_DEFAULT_ARGS=(--verbose)

# Lock on a specifc file
# Useful to make sure, that only one instance of penrun can access a particular resource
# PENRUN_LOCK="/path/to/lock"

################################################################################
# Miscellaneous config

# Perform this function prior to each penrun run.
# Can be used e.g. for power cycles or sending triggers
pre_run() {
	echo "I am a pre_run hook!"
}

# Perform this function after to each penrun run.
post_run() {
	echo "I am a post_run hook!"
}
