# fleet-controller

Initial domain model for node maintenance. The `maintenance` package owns the
state machine; there are no external dependencies or infrastructure components.

| From | Allowed next states |
| --- | --- |
| `READY` | `DRAINING` |
| `DRAINING` | `UPGRADING`, `RETRY_WAIT`, `QUARANTINED` |
| `UPGRADING` | `VERIFYING`, `RETRY_WAIT`, `QUARANTINED` |
| `VERIFYING` | `READY`, `RETRY_WAIT`, `QUARANTINED` |
| `RETRY_WAIT` | `DRAINING`, `QUARANTINED` |
| `QUARANTINED` | `READY` |

`RETRY_WAIT` restarts the attempt at draining. Moving from `QUARANTINED` to
`READY` represents an explicit release by the caller.

Run tests with `go test ./...`.
