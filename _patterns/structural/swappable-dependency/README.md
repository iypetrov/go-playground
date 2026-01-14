The Swappable Dependency Pattern allows a service to operate in a degraded mode until a core dependency becomes available, 
at which point the dependency is swapped in dynamically at runtime. Decouple dependency wiring from dependency usage, so you can:

- Start the app without immediately needing the full dependency.

- Replace the dependency once it's ready (e.g. after retries, warmup, or discovery).

Common use cases:
- Delayed database connections in web servers (start HTTP server first, inject DB later).

- Hot-swapping mock/stub components for live ones.

- Long-lived background services that retry on failure (e.g. config loading, auth systems).

- Swapping feature flags or plugins without restarting the application.
