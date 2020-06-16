#!/bin/bash

set -o nounset
set -o pipefail

GITLAB_URL=${GITLAB_URL:=http://gitlab.com}
GITLAB_TOKEN=${GITLAB_TOKEN:=}

# Count inactive users
users=$(env | grep LABEL_USERNAME | wc -l | tr -d ' ')

# Blocks ony by one
for (( i = 0; i < users; i++ )); do
  name="$(printenv LABEL_USERNAME_$i)"
  if [[ -z "$name" ]]; then
    continue
  fi

  echo "Found: $name"
  if ! user_list=$(curl -sfL --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
      "$GITLAB_URL/api/v4/users?username=$name"); then
    echo "Failed to find $name user"
    continue
  fi

  user_id=$(echo "$user_list" | jq -erM .[0].id)
  if ! curl -sfL --header "PRIVATE-TOKEN: $GITLAB_TOKEN" -X POST \
    "$GITLAB_URL/api/v4/users/$user_id/block"; then
    ecoh "Failed to block $name user"
  fi
done
