heroku container:push web --arg SPOTIFY_SLACK_APP_SLACK_CLIENT_ID=${SPOTIFY_SLACK_APP_SLACK_CLIENT_ID},SPOTIFY_SLACK_APP_SLACK_CLIENT_SECRET=${SPOTIFY_SLACK_APP_SLACK_CLIENT_SECRET},SPOTIFY_SLACK_APP_SPOTIFY_REDIRECT_URL=${SPOTIFY_SLACK_APP_SPOTIFY_REDIRECT_URL},SPOTIFY_SLACK_APP_NEW_RELIC_APP_NAME=${SPOTIFY_SLACK_APP_NEW_RELIC_APP_NAME},SPOTIFY_SLACK_APP_NEW_RELIC_LICENSE=${SPOTIFY_SLACK_APP_NEW_RELIC_LICENSE},SPOTIFY_SLACK_APP_DATABASE_URL=${SPOTIFY_SLACK_APP_DATABASE_URL},SPOTIFY_SLACK_APP_SLACK_AUTH_URL=${SPOTIFY_SLACK_APP_SLACK_AUTH_URL}

heroku container:release web