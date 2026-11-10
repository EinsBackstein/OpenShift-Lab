# OpenShift-Custom-Image-Deployment

A Vite single-page site packaged in one image with Traefik. Traefik listens on port 8080 and forwards requests to a small static server bound to loopback inside the container. The static server serves `index.html` for HTML navigation to SPA routes and returns 404 for missing assets. The image runs as a non-root user.

## Run locally

```sh
docker compose up --build -d
```

Open <http://localhost:8088>. Stop it with `docker compose down`. To change the host port, edit the `ports` mapping in `compose.yaml`.

For frontend development, run `npm ci` and `npm run dev`.

## Prepare a GitHub Container Registry release

The local image is `openshift-custom-image-deployment:local`; it has not been pushed or registered. When ready to publish, replace `OWNER` with the lowercase GitHub user or organization and choose a version:

```sh
docker tag openshift-custom-image-deployment:local ghcr.io/OWNER/openshift-custom-image-deployment:1.0.0
# Authenticate to ghcr.io using the GitHub-supported method for your account.
docker push ghcr.io/OWNER/openshift-custom-image-deployment:1.0.0
```

The local build targets the architecture of the Docker host. For a multi-platform release, use a suitable Buildx builder and build explicitly for the desired platforms before pushing. Configure package visibility and repository permissions in GitHub after publication if needed.
