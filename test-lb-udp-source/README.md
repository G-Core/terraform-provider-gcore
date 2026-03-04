# UDP SOURCE_IP persistence scenario

Covers a UDP listener wired to a SOURCE_IP pool with inline members:
- Exercises L3 preferred connectivity with dedicated VIP network
- Validates UDP-CONNECT health monitors and SOURCE_IP session persistence knobs
- Uses inline pool members to test monitor overrides, backup members, and weight-driven hashing impact

Run against a lab tenant to prove UDP behaviour, hashing drift, and persistence timeout handling.
