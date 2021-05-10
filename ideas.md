# Ideas

- Storage
    - In-memory. Run daemon, and use commands to get and set. Use TCP to connect to process running a unix socket server. The process is a daemon which is always running, `kvd`. For persistent storage, the daemon could load from disk on startup.
    - On-disk. Could do this both with and without a daemon.
- Interface
    - `kv get key` gets value for `key`
    - `kv set key value` sets `key` to `value`
    - `kv` enters a CLI
- Add indexes etc. eventually for more efficient retrieval
