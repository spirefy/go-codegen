.{
    .name = "{{.Project}}", // Name of your project
    .version = "{{.Version}}",   // Version of your project
    .paths = .{""},       // Paths to search for dependencies (optional)
    .dependencies = .{
        .extism_pdk = .{ // Dependency declaration for Extism PDK
          .url = "https://github.com/extism/zig-pdk/archive/master.tar.gz",
          .hash = "12209cd75e8d0cf119d2f6755a883c61b2da0a9f8efd8d221218b59e9c6feca367ad"
        },
    },
    // Other optional fields:
    // .target = ... (Target architecture for cross-compilation)
    // .link_libc = ... (True/false for linking with the C standard library)
    // .system_paths = ... (Paths to system libraries)
}