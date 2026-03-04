# Prometheus + Canary routing scenario

This config stands up three listeners on the same load balancer:
- Production HTTP listener with two weighted members and strict readiness health checks
- Canary HTTP listener pinned to a different frontend port with cookie-based stickiness and an administratively disabled member to validate drain behaviour
- Prometheus listener protected by secrets + user list, wired to a dedicated metrics pool with custom monitor endpoints

Use it to validate listener lifecycle, multiple pool attachments, user-list enforcement, and the interaction of per-pool/session settings with standalone member resources.
