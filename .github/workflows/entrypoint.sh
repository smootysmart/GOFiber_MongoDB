#!/bin/bash
set -e
echo "Check running containers..."
if [ "$(docker ps -q)" ]; then
  docker compose down
  docker system prune -af
fi
docker compose pull
docker compose up -d