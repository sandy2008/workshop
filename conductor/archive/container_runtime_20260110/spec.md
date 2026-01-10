# Specification: Container Runtime Workshop (Namespaces & Cgroups)

## 1. Overview

Create a new hands-on workshop document `infra/container_runtime_ja.md` that explains the internal mechanisms of containers using `namespaces` and `cgroups`. The workshop will focus on manual implementation using standard Linux commands (`unshare`, `ip`, etc.) and Go language, avoiding high-level container runtimes like Docker for the core logic, although Podman may be used for environment setup or comparison.

## 2. Target Audience & Prerequisites

* **Target Audience:** Software engineers and system administrators who want to understand how containers work under the hood. Familiarity with Go language is required.
* **Prerequisites:**
  * **OS:** Linux environment (e.g., Ubuntu 24.04 VM).
  * **Permissions:** Root (`sudo`) access is mandatory.
  * **Tools:** `podman` (for rootfs extraction or comparison), `go` (for implementation exercises).
  * **Knowledge:** Basic Linux commands (`ls`, `ps`, `mount`, etc.) and basic Go programming.
  * **Context:** Familiarity with concepts from `./infra/vlan_ja.md` is beneficial for the networking section.

## 3. Functional Requirements (Workshop Content)

The workshop MUST cover the following steps:

### 3.1. Filesystem & Rootfs Preparation

* Guide the user to prepare a minimal root filesystem (rootfs) using Alpine Linux.
* Demonstrate how to export a rootfs using `podman export`.

### 3.2. Manual Namespace Isolation (CLI)

* **PID Namespace:** Use `unshare --pid --fork` to isolate process IDs. Verify `ps` output inside the namespace.
* **Mount Namespace:** Use `unshare --mount` and `chroot` (or `pivot_root`) to isolate the filesystem. Mount `/proc` to make `ps` work correctly.
* **Network Namespace:** Use `unshare --net` to isolate the network stack.

### 3.3. Resource Control with Cgroups (CLI)

* Create a cgroup manually (cgroup v2).
* Demonstrate limiting CPU or Memory usage.
* Verify the limit is enforced (e.g., using a stress tool or script).

### 3.4. Networking with Veth

* Create a `veth` pair linking the host and the container namespace.
* Assign IP addresses to both ends.
* **Interactive Example:** Demonstrate sending a message (e.g., HTTP request or netcat) from the Host to the Container and receiving a response. This serves as an advanced networking exercise following the VLAN workshop.

### 3.5. Recreating it in Go

* Write a Go program that programmatically performs the above steps using the `syscall` package.
  * Setting `Cloneflags` (CLONE_NEWUTS, CLONE_NEWPID, etc.).
  * Setting up `Cmd.SysProcAttr`.
  * Executing `/bin/sh` inside the new environment.

## 4. Non-Functional Requirements

* **Language:** Japanese (`_ja.md`).
* **Format:** Markdown, following the existing repository style (headers, code blocks, etc.).
* **Clarity:** Explain *why* each step is necessary (e.g., why mount /proc?).

## 5. Out of Scope

* Detailed implementation of an OCI-compliant runtime.
* Deep dive into Union filesystems (OverlayFS) details (keep it simple with a directory for rootfs).
