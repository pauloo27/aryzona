#!/bin/bash

PROJECT_ID=$1
CI_URL=$2
TOKEN=$3

agola project secret delete --project $PROJECT_ID -u $CI_URL --name DEPLOY --token $TOKEN

agola project secret create --project $PROJECT_ID -u $CI_URL --name DEPLOY --token $TOKEN -f ./deploy-secrets.yml

function recreate_secret {
  agola project variable delete --project $PROJECT_ID -u $CI_URL --name $1 --token $TOKEN
  echo "- secret_name: DEPLOY
  secret_var: $2" | agola project variable create --project $PROJECT_ID -u $CI_URL --name $1 --token $TOKEN -f "-"
}

recreate_secret discordwebhook DISCORD_WEBHOOK
recreate_secret deployuser USER
recreate_secret deployhost HOST
recreate_secret deploypath PATH
