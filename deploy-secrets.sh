#!/usr/bin/env bash
set -euo pipefail

# 1. Define a mapping of Docker Secret Names to Bitwarden Secret IDs
# Replace the UUIDs below with the actual Secret IDs from your Bitwarden UI
declare -A SECRETS=(
  ["token"]="640f4ba2-98b4-4455-b575-b45b00362274"
  ["client-id"]="4b2e8d52-5fb7-41be-ae2d-b45b002faa2c"
  ["guild-id"]="e74aa0bd-9668-4335-947b-b45b002fb545"
  ["main-channel-id"]="d19ad9ab-263d-40fd-b955-b45b002fd3b5"
  ["db-pass"]="8040da69-ecf3-4153-a6bd-b45b0033206a"
  ["redis-url"]="a792e517-a324-4921-9db9-b46a0186d545"
)

echo "Locking down and Syncing Bitwarden Secrets to Docker Swarm..."

for SECRET_NAME in "${!SECRETS[@]}"; do
  SECRET_ID="${SECRETS[$SECRET_NAME]}"
  
  echo "Processing secret: ${SECRET_NAME}..."

  # Remove the secret if it already exists to prevent "already exists" conflicts
  if docker secret inspect "$SECRET_NAME" >/dev/null 2>&1; then
    echo " -> Removing existing Swarm secret: ${SECRET_NAME}"
    docker secret rm "$SECRET_NAME"
  fi

  # Fetch the value via Bitwarden CLI and pipe it securely into Docker Swarm
  # --value returns just the raw string, avoiding json parsing
  bws secret get "$SECRET_ID" | jq -r '.value' | docker secret create "$SECRET_NAME" -

  echo " -> Successfully created Swarm secret: ${SECRET_NAME}"
done

echo "All secrets loaded successfully into Docker Swarm!"
