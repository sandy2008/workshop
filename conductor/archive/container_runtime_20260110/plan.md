# Plan: Container Runtime Workshop (Namespaces & Cgroups)

## Phase 1: Environment & Prerequisites Setup [checkpoint: c2a57c4]

- [x] Task: Create workshop directory structure and initial files
  - [x] Subtask: Create `infra/container_runtime_ja.md` with the "Overview" and "Prerequisites" sections.
  - [x] Subtask: Create a `Makefile` or setup script instructions in the doc to verify the environment (check `podman` version, `go` version, `sudo` access).
- [x] Task: Document Rootfs Preparation
  - [x] Subtask: Write instructions to pull an Alpine Linux image using Podman.
  - [x] Subtask: Write instructions to export the image to a rootfs directory (e.g., `rootfs/`).
- [x] Task: Conductor - User Manual Verification 'Environment & Prerequisites Setup' (Protocol in workflow.md)

## Phase 2: Manual Namespace Isolation (CLI) [checkpoint: f2a619b]

- [x] Task: Document PID Namespace Isolation
  - [x] Subtask: Add explanation of PID namespaces.
  - [x] Subtask: Add step-by-step commands using `unshare --pid --fork` and verification with `ps`.
- [x] Task: Document Mount Namespace Isolation
  - [x] Subtask: Add explanation of Mount namespaces.
  - [x] Subtask: Add commands for `unshare --mount`, `chroot` (or `pivot_root`), and mounting `/proc`.
- [x] Task: Document Network Namespace Isolation
  - [x] Subtask: Add explanation of Network namespaces.
  - [x] Subtask: Add commands using `unshare --net` and verification with `ip link`.
- [x] Task: Conductor - User Manual Verification 'Manual Namespace Isolation (CLI)' (Protocol in workflow.md)

## Phase 3: Resource Control with Cgroups (CLI) [checkpoint: 1b91fbe]

- [x] Task: Document Cgroups Setup
  - [x] Subtask: Explain Cgroups v2 concepts.
  - [x] Subtask: Write instructions to manually create a cgroup directory in `/sys/fs/cgroup`.
- [x] Task: Document Resource Limiting
  - [x] Subtask: Add steps to limit memory or CPU (e.g., writing to `memory.max`).
  - [x] Subtask: Provide a method to verify the limit (e.g., running a stress test process and adding its PID to `cgroup.procs`).
- [x] Task: Conductor - User Manual Verification 'Resource Control with Cgroups (CLI)' (Protocol in workflow.md)

## Phase 4: Networking with Veth [checkpoint: e2be89f]

- [x] Task: Document Veth Pair Creation
  - [x] Subtask: Explain veth pairs and their role in container networking.
  - [x] Subtask: Add commands to create a veth pair and move one end to the container's network namespace.
- [x] Task: Document IP Assignment & Connectivity
  - [x] Subtask: Add commands to assign IPs to both host and container veth ends and bring interfaces up.
  - [x] Subtask: Add instructions to test connectivity (ping/curl) from Host to Container.
- [x] Task: Conductor - User Manual Verification 'Networking with Veth' (Protocol in workflow.md)

## Phase 5: Recreating it in Go [checkpoint: 7a8b9c2]

- [x] Task: Document Go Implementation Setup
  - [x] Subtask: Explain the mapping between CLI commands and Go `syscall` package.
  - [x] Subtask: Create a skeleton Go program (`main.go`) for the container runtime.
- [x] Task: Implement Namespace & Command Execution in Go
  - [x] Subtask: Add code to set `Cloneflags` (CLONE_NEWUTS, CLONE_NEWPID, CLONE_NEWNS, CLONE_NEWNET).
  - [x] Subtask: Add code to execute `/bin/sh` inside the new environment.
  - [x] Subtask: Add code for `chroot` and mounting `/proc` within the Go program.
- [x] Task: Conductor - User Manual Verification 'Recreating it in Go' (Protocol in workflow.md)

## Phase 6: Final Review & Polish [checkpoint: d435ceb]

- [x] Task: Review and Refine Documentation
  - [x] Subtask: Read through the entire `infra/container_runtime_ja.md` for clarity and flow.
  - [x] Subtask: Ensure all commands are copy-pasteable and work on the target environment.
  - [x] Subtask: Add "Summary" and "Next Steps" sections.
- [x] Task: Conductor - User Manual Verification 'Final Review & Polish' (Protocol in workflow.md)
