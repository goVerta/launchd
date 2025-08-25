# launchd
Deterministic deploy orchestrator for single-binary Go services anchored to Linux systemd. launchd operationalizes a minimal, auditable deploy surface by composing battle-tested primitives: the resident Go toolchain for compilation, OpenSSH for transport and privileged control, systemd for service lifecycle, an optional migrator for schema evolution, and an HTTP health contract for liveness attestation. The objective function is predictability under duress, not novelty; the codebase is dependency-thin and tuned for reproducible behavior across heterogeneous fleet baselines.

## Conceptual Premise

The dominant substrate for “small-but-critical” services is a single daemon under systemd supervision. For these domains, heavyweight CI/CD machinery introduces variance without proportional utility. launchd prefers an austere convergence model: converge a binary onto a machine, converge a unit into systemd, converge the database state via idempotent migrations, and converge process readiness via health probes. Each convergence step is a transparent syscall to a canonical tool, yielding failure surfaces that are legible to any seasoned operator.

Assumptions:
- A Go main package constitutes the service.
- SSH reachability with privilege elevation exists to manage units.
- The service exposes `/health` on a chosen TCP port.

## MVP Scope

- CLI: `launchd deploy --host <ip> --user <ssh-user> --app ./path/to/app --port <port> [--timeout <dur>]`.
- Pipeline: Compile ➝ Transfer ➝ systemd Provision ➝ Migrate (optional) ➝ Health Gate.
- Properties: idempotent (safe re-entrance), deterministic (explicit tooling), strict error surfacing, and minimal environmental preconditions (Go + OpenSSH).

### Stage Semantics
- Compile: `go build -o /tmp/<app>` using the resident toolchain; no hermetic wrapper.
- Transfer: `scp` to `/usr/local/bin/<app>` after ensuring parent directory ownership is sane.
- systemd: materialize a unit, `daemon-reload`, `enable`, `restart`—idempotent operations.
- Migrations: opportunistic `migrate up` if a migrator exists on the target PATH.
- Health: poll `http://<host>:<port>/health` until success or deadline expiry.

## Example Usage

```bash
launchd deploy --host 203.0.113.10 --user ubuntu --app ./examples/hello --port 8080 --timeout 60s
```

On completion, the service is registered as `<app>.service`, executes `/usr/local/bin/<app> --port=<port>`, and is enabled for boot.

## Design Guarantees

- Determinism: all side effects mediated by explicit tools (`go`, `scp`, `ssh`, `systemctl`).
- Idempotence: repeated invocations converge; non-destructive `enable`, guarded migrations.
- Failure Locality: stages short-circuit with precise logging; no ambiguous partial states.
- Observability: microsecond timestamps and stage banners designed for operator cognition.
- Minimality: no bespoke protocol layers; everything deferential to Unix contracts.

## Future Roadmap

- Artifact integrity: checksums, content-addressed remote layout, atomic swaps.
- Principle of Least Privilege: dedicated system users, hardening of unit sandboxing.
- Transport hardening: native SSH client (keyboard-interactive, agent-forwarding policies) while preserving zero-daemon requirements.
- Policy engines: JSON-structured logs, exponential backoff policies, retry budgets.
- Migration adapters: goose, golang-migrate, app-native hooks with transactional guards.

## Organizational Context

This codebase is authored under the goVerta collective and aspires to the production rigor expected of core-infra artifacts. The project is intentionally conservative in feature accretion, biasing toward operability, determinism, and testable failure semantics over breadth.

## Contributors and Roles

- Saad H. Tiwana — lead author, deployment pipeline, systemd strategy, reliability posture  
  GitHub: https://github.com/saadhtiwana
- Ahmad Mustafa — SSH/SCP transport hardening, testing harnesses  
  GitHub: https://github.com/ahmadmustafa02
- Majid Farooq Qureshi — QA and Toolsmith, Makefile/CI touchpoints, documentation QA  
  GitHub: https://github.com/Majid-Farooq-Qureshi

Saad authored the core code; Ahmad and Majid contributed engineering and QA functions across transport, tests, and docs.

— saad and gang is who build this
