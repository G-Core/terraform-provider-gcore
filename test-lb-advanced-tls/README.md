# Advanced TLS + Logging LB scenario

This scenario provisions:
- Dual-stack load balancer pinned to a dedicated VIP subnet with floating IP creation
- HTTPS listener with TLS termination, SNI, sticky sessions, strict HTTP health checks, and logging enabled
- Prometheus listener with RBAC-style user list enforced through encrypted passwords
- Dedicated backend network/subnet plus standalone pool member resources covering backup members, monitor overrides, and weight-based failover

Use this config to validate TLS objects (secrets/CA bundles), VRRP allocation, logging, and interaction between listeners, pools, and pool members.
