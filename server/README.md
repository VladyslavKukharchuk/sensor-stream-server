# Sensor stream server

### 🚀 Local run (without Docker)

```bash
go mod download
go run main.go
```

### 🐳 Local run in Docker Container:

1. Build image
    ```bash
    docker build -f build/Dockerfile -t sensor-stream-server .
    ```

2. Run container
    ```bash
    docker run -p 8080:8080 sensor-stream-server
    ```

### How to push to artifact registry (GCloud)

1. Authenticate Docker for Artifact Registry
   Follow GCP [docs](https://cloud.google.com/artifact-registry/docs/docker/pushing-and-pulling#pushing) or run:
   ```bash
   gcloud auth configure-docker us-central1-docker.pkg.dev
   ```

2. Build for linux/amd64 machine
    ```bash
    docker build \
      --platform linux/amd64 \
      -f build/Dockerfile \
      -t sensor-stream-server . 
    ```
   
3. Create new artifacts repository 
    ```bash
    gcloud artifacts repositories create sensor-repo \
    --repository-format=docker \
    --location=us-central1 \
    --project=opportune-scope-320904
    ```

4. Tag image for Artifact Registry
    ```bash
    docker tag sensor-stream-server \
    us-central1-docker.pkg.dev/opportune-scope-320904/sensor-repo/sensor-stream-server:latest
    ```

5. Push this image to artifact registry
    ```bash
    docker push \
    us-central1-docker.pkg.dev/opportune-scope-320904/sensor-repo/sensor-stream-server:latest
    ```

6. Deploy on Cloud Run
    ```bash
    gcloud run deploy sensor-stream-server \
      --image=us-central1-docker.pkg.dev/opportune-scope-320904/sensor-repo/sensor-stream-server:latest \
      --platform=managed \
      --region=us-central1 \
      --allow-unauthenticated \
      --port=8080
    ```
