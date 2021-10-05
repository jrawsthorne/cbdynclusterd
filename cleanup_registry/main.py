import requests

DOCKER_REGISTRY = "172.23.104.149"

headers = {
    "Accept": "application/vnd.docker.distribution.manifest.v2+json",
    "Content-Type": "application/json"
}

repos = requests.get(f"http://{DOCKER_REGISTRY}/v2/_catalog?n=1000", headers=headers).json()["repositories"]

for repo in repos:
    if repo.startswith("dynclsr-couchbase"):
        repo_url = f"http://{DOCKER_REGISTRY}/v2/{repo}/manifests"
        docker_image_sha = requests.head(f"{repo_url}/latest", headers=headers).headers.get("Docker-Content-Digest")
        if docker_image_sha:
            try:
                requests.delete(f"{repo_url}/{docker_image_sha}", headers=headers).raise_for_status()
                print(f"Deleted {repo}: {docker_image_sha}")
            except Exception:
                print(f"Couldn't delete {repo}: {docker_image_sha}")
