# 🧱 Modern Docker Architecture

## 🔧 Docker Architecture Diagram

```
+-------------------------------------------------------------+
|                        Docker CLI / API                     |
|                     (docker, docker-compose)                |
+------------------------+------------------------------+-----+
                         |                              |
                         v                              v
               +------------------+            +---------------------+
               |   Docker Daemon  |<---------->| REST API Server     |
               |   (dockerd)      |            | (Exposes Docker API)|
               +--------+---------+            +---------------------+
                        |
                        v
         +-------------------------------+
         |  Containerd (Container Runtime)|
         +---------------+---------------+
                         |
            +------------+------------+
            |                         |
            v                         v
  +--------------------+    +------------------------+
  |  runc (OCI CLI)    |    | Container Network Stack|
  +--------------------+    +------------------------+
            |                         |
            v                         v
  +--------------------+     +------------------------------+
  | Namespaces, cgroups|     |  Bridge, Overlay, Host, etc. |
  | SELinux, AppArmor  |     |  (via libnetwork or CNI)     |
  +--------------------+     +------------------------------+

              +--------------------------------------+
              |       Docker Image Registry          |
              |     (DockerHub / Private Registry)   |
              +--------------------------------------+

         +--------------------------------------+
         |        Orchestration Layer           |
         | (Docker Swarm / Kubernetes via CRI)  |
         +--------------------------------------+
```

## ⚙️ Components Explained

### 1. Docker CLI / API

* **docker**, **docker-compose**: Command-line interface for building, running, and managing containers.
* Communicates with `dockerd` (Docker daemon) via REST API.

### 2. Docker Daemon (`dockerd`)

* Central component of Docker.
* Listens to Docker API requests (from CLI or tools).
* Manages images, containers, volumes, networks.
* Can communicate with **other Docker daemons** in a swarm setup.

### 3. Docker API Server

* REST API exposed by Docker daemon.
* Enables programmatic control and integration.
* Tools like Portainer, Jenkins, and Kubernetes Docker shim used this to control Docker.

### 4. Containerd

* Industry-standard container runtime (now graduated CNCF project).
* Manages container lifecycle: start, stop, pause, image pull/push, etc.
* **dockerd** delegates to **containerd** for container operations.

### 5. runc

* Low-level OCI-compliant CLI tool to spawn and run containers.
* Created by Docker, now maintained as an open-source project under **Open Container Initiative (OCI)**.
* Executes containers using Linux namespaces, cgroups, and seccomp.

### 6. Kernel Features

* `Namespaces`: Process isolation (PID, NET, IPC, MNT, etc.)
* `cgroups`: Resource control (CPU, memory, I/O, etc.)
* `AppArmor / SELinux`: Mandatory access controls (security)

### 7. Networking Layer

* Managed by **libnetwork** or **CNI plugins**.
* Network drivers: `bridge`, `overlay`, `macvlan`, `host`, `none`.
* Allows multi-host networking (e.g., overlay in Swarm or CNI in Kubernetes)

### 8. Storage Layer

* Docker volumes or bind mounts.
* Used for persistent or shared data.
* Volume drivers supported (e.g., NFS, EBS, Azure Disk)

### 9. Image Registry

* Stores and distributes Docker images.
* Public: DockerHub, GitHub Container Registry
* Private: Harbor, AWS ECR, JFrog Artifactory

### 10. Orchestration Layer

* **Docker Swarm** (native) – clustering, load balancing, service deployment.
* **Kubernetes** (modern standard) – through **Container Runtime Interface (CRI)** and **containerd shim**.
* Modern Docker supports Kubernetes integration via Docker Desktop.

## 🛠 Modern Enhancements (Post-2020)

| Feature               | Description                                  |
| --------------------- | -------------------------------------------- |
| **Docker Desktop**    | GUI + Kubernetes + Dev Tools on Mac/Windows. |
| **BuildKit**          | Faster, more powerful image builds.          |
| **Compose v2**        | Docker Compose as a plugin, written in Go.   |
| **containerd & runc** | Decoupled runtimes for OCI compatibility.    |
| **Docker Extensions** | Marketplace-like plugin system for Desktop.  |

## 🚀 Workflow Summary (from Docker CLI to Running Container)

1. Developer runs a command like `docker run -d nginx`
2. `Docker CLI` sends request to `dockerd` via REST API
3. `dockerd` delegates to `containerd`
4. `containerd` uses `runc` to create and run the container
5. Kernel isolation is applied (namespaces, cgroups)
6. Networking is configured via libnetwork/CNI
7. Logs and metrics are collected
8. Container is visible via `docker ps`

## 🎯 Top 25 Docker Interview Questions with Answers

### 1. **What is Docker and why is it used?**

Docker is a platform to develop, ship, and run applications inside containers. It allows consistent environments, quick deployment, and resource isolation.

### 2. **What is the difference between a container and a virtual machine?**

Containers share the host OS kernel and are lightweight. VMs run full guest OS and are heavier.

### 3. **What are Docker images and containers?**

An image is a snapshot or template used to create containers. A container is a running instance of an image.

### 4. **How does Docker achieve isolation?**

Using Linux kernel features like namespaces and cgroups for process, network, and resource isolation.

### 5. **What is the role of Dockerfile?**

Dockerfile is a script that contains instructions to build a Docker image.

### 6. **What is the difference between CMD and ENTRYPOINT?**

Both define the container's executable. `CMD` can be overridden with CLI args. `ENTRYPOINT` is not easily overridden.

### 7. **How does Docker manage networking?**

Through network drivers like bridge, host, overlay. It can use `libnetwork` or CNI plugins.

### 8. **What is the difference between Docker volumes and bind mounts?**

Volumes are managed by Docker, portable and secure. Bind mounts are host-specific.

### 9. **What is Docker Compose?**

Tool for defining and running multi-container Docker applications using a YAML file.

### 10. **How can you persist data in Docker containers?**

By using volumes or bind mounts.

### 11. **What is Docker Swarm?**

Native clustering and orchestration tool in Docker for managing services across a cluster.

### 12. **How do you scale services using Docker?**

Using `docker-compose scale` or Docker Swarm’s `docker service scale`.

### 13. **What is a multi-stage build in Docker?**

Technique to reduce image size by separating build and runtime dependencies.

### 14. **What is Docker Hub?**

A public registry for sharing Docker images.

### 15. **How do you optimize Docker images?**

* Use slim base images
* Combine RUN instructions
* Leverage `.dockerignore`
* Use multi-stage builds

### 16. **What is the difference between `COPY` and `ADD` in Dockerfile?**

`COPY` is for basic copying. `ADD` supports remote URLs and automatic archive extraction.

### 17. **What is container orchestration?**

Managing containers at scale—scheduling, deployment, scaling, health checks. Examples: Kubernetes, Docker Swarm.

### 18. **What are Docker layers?**

Each instruction in Dockerfile creates a layer. Layers help in caching and reducing redundancy.

### 19. **What are health checks in Docker?**

Docker can periodically run a command to determine if a container is healthy.

### 20. **What is the Docker context?**

Context is the set of files and directories sent to Docker daemon when building an image.

### 21. **How do you remove dangling images and containers?**

```bash
# Dangling images
$ docker image prune

# Stopped containers
$ docker container prune
```

### 22. **What is the difference between `EXPOSE` and `-p`?**

`EXPOSE` documents port in Dockerfile. `-p` maps container port to host during runtime.

### 23. **How does Docker handle logs?**

Docker captures stdout/stderr. Logs can be accessed via `docker logs`. Supports drivers like `json-file`, `fluentd`, `syslog`.

### 24. **Can containers communicate with each other?**

Yes, if they’re on the same network (e.g., Docker bridge or custom overlay).

### 25. **What is the difference between `RUN`, `CMD`, and `ENTRYPOINT`?**

* `RUN`: Executes during image build
* `CMD`: Default command during container start
* `ENTRYPOINT`: Defines executable, often used with CMD as arguments
