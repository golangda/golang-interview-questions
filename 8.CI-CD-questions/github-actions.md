## 🚀 How GitHub Actions Work — With Detailed Sequence Diagram

GitHub Actions is a **CI/CD platform built into GitHub** that allows you to **automate workflows** triggered by events like `push`, `pull_request`, `schedule`, etc. You define automation as code in YAML files stored in `.github/workflows/`.

---

### 🔁 High-Level Workflow of GitHub Actions

1. **Trigger**: GitHub detects an event (e.g., `push`, `pull_request`, etc.).
2. **Workflow Dispatch**: It checks if any workflow YAML is configured for that event.
3. **Job Scheduling**: GitHub allocates runners (VMs or containers) per job.
4. **Step Execution**: Jobs execute defined steps (checkout code, run tests, etc.).
5. **Artifact Upload** *(Optional)*: Artifacts or logs are stored/uploaded.
6. **Status Update**: GitHub updates the commit/PR status based on the outcome.

---

### 🧱 Components of GitHub Actions

| Component    | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| **Workflow** | Top-level automation defined in `.github/workflows/*.yml`    |
| **Event**    | Trigger that starts the workflow (e.g., push, pull\_request) |
| **Job**      | A set of steps run in the same runner instance               |
| **Runner**   | A GitHub-hosted or self-hosted virtual machine               |
| **Step**     | A command or action run as part of a job                     |
| **Action**   | Reusable custom logic (e.g., `actions/checkout`)             |

---

### 📊 Sequence Diagram: GitHub Actions Workflow

```
          +----------------+                 +------------------+
          |   Developer    |                 |    GitHub.com     |
          +----------------+                 +------------------+
                  |                                   |
                  |      Push / PR / Trigger Event    |
                  |----------------------------------->|
                  |                                   |
                  |     Match with workflow YAML      |
                  |<----------------------------------|
                  |                                   |
                  |      Schedule Workflow Run        |
                  |----------------------------------->|
                  |                                   |
                  |     Assign Runners (VMs)          |
                  |----------------------------------->|
                  |                                   |
                  |     Fetch Repo / Checkout Code    |
                  |<----------------------------------|
                  |                                   |
                  |         Run Build Steps           |
                  |<----------------------------------|
                  |                                   |
                  |         Run Tests / Linting       |
                  |<----------------------------------|
                  |                                   |
                  |     Store Logs / Artifacts        |
                  |<----------------------------------|
                  |                                   |
                  |     Update PR / Commit Status     |
                  |<----------------------------------|
                  |                                   |
                  |     Notify via Email / Slack      |
                  |<----------------------------------|
```

---

### 🧬 Sample `.github/workflows/go-ci.yml` File

```yaml
name: Go CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v5
      with:
        go-version: 1.21

    - name: Install dependencies
      run: go mod tidy

    - name: Run tests
      run: go test ./...

    - name: Lint
      run: golangci-lint run
```

---

### 🔐 Secure Secrets Handling

GitHub Actions supports **secrets**:

```yaml
env:
  DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
```

Stored in **Settings → Secrets and Variables → Actions**.

---

### 📌 Notes on GitHub Actions

* Supports **matrix builds** for testing against multiple OS/versions.
* Offers **reusable workflows** and **composite actions**.
* Has **environment protections**, **manual approvals**, and **branch filtering**.
* GitHub-hosted runners are billed by usage (minutes).
* Integrates with GitHub UI for logs, status badges, artifacts, and annotations.

---

## 🧪 Top 25 CI/CD Interview Questions with Detailed Answers

1. **What is CI/CD and how does it benefit software development?**
   CI/CD stands for Continuous Integration and Continuous Delivery/Deployment. CI ensures frequent integration and automated testing, while CD automates the deployment process. Benefits include faster releases, fewer bugs, and better team collaboration.

2. **What are the key stages of a CI/CD pipeline?**
   Stages include: Source, Build, Test, Artifact creation, Deploy, Monitor.

3. **What tools are commonly used for CI/CD?**
   CI: Jenkins, GitHub Actions, GitLab CI; CD: ArgoCD, Spinnaker; Supporting: Docker, Kubernetes, Helm.

4. **What is a build artifact?**
   A compiled output of your code (e.g., Docker image, .jar file) stored in artifact repositories like Artifactory or GitHub Packages.

5. **How do you ensure zero-downtime deployments?**
   Strategies include Blue-Green, Canary, Rolling updates, Feature flags, and readiness probes.

6. **What is the difference between Continuous Delivery and Continuous Deployment?**
   Delivery requires manual approval for prod deploy; Deployment pushes directly to production without approval.

7. **What is a Jenkinsfile?**
   A declarative or scripted file in Jenkins that defines CI pipeline as code.

8. **How do you trigger a CI/CD pipeline?**
   Via events like push, pull\_request, schedule, manual trigger, or webhook.

9. **How do you handle secrets in CI/CD pipelines?**
   Use encrypted secrets via Secret Managers or CI secrets configuration; never hardcode.

10. **What are environment promotion strategies?**
    Manually or automatically move artifacts from dev → staging → prod using tags, approvals, or pipeline triggers.

11. **How do you implement canary deployment in Kubernetes?**
    Deploy a new version to a small percentage of traffic, gradually increase using Istio or weighted services.

12. **How do you test Infrastructure as Code (IaC)?**
    Use tools like Terraform validate, Terratest, ansible-lint, and Molecule.

13. **What is Blue-Green deployment?**
    Maintain two environments (Blue: live, Green: new), switch traffic after validation.

14. **What is GitOps?**
    A CD strategy where Git is the single source of truth; ArgoCD or FluxCD reconcile from Git to Kubernetes.

15. **How do you roll back a deployment?**
    Use Git tags, Docker images, or Helm rollbacks to deploy previous working version.

16. **How do you handle flaky tests?**
    Retry logic, isolate tests, use quarantine pipelines, root cause analysis.

17. **How do you perform static code analysis in CI?**
    Use tools like SonarQube, golangci-lint, eslint, and run them in pipeline.

18. **How do you secure a CI/CD pipeline?**
    Encrypt credentials, isolate runners, use signed artifacts, scan dependencies, audit logs.

19. **How do you monitor CI/CD pipelines?**
    Use built-in dashboards, Prometheus/Grafana, log aggregators, and alerting tools.

20. **Pipeline-as-Code vs GUI Pipelines?**
    Pipeline-as-code (YAML/Jenkinsfile) is version-controlled, reproducible; GUI pipelines are quick but harder to maintain.

21. **Common CI/CD Anti-patterns?**
    Long builds, hardcoded secrets, manual steps, deploying on Friday, ignoring failed tests.

22. **How does Docker help in CI/CD?**
    Provides isolated, reproducible environments; simplifies building, testing, and deployment.

23. **How do you handle parallel jobs in CI?**
    Use job matrices or define multiple jobs that can run in parallel to speed up pipelines.

24. **What are matrix builds?**
    Run a job across multiple OS or language versions using combinations (e.g., Go 1.20, 1.21 on Linux and Windows).

25. **How do you implement approvals in CI/CD?**
    Use GitHub Environments, Jenkins `input`, or GitLab manual job approvals before prod deploy.
