const std = @import("std");

pub fn build(b: *std.Build) void {
    const optimize = b.standardOptimizeOption(.{});
    const target = b.standardTargetOptions(.{
        // if you're using WASI, change the .os_tag to .wasi
        .default_target = .{ .abi = .musl, .os_tag = .wasi, .cpu_arch = .wasm32 },
    });
    const pdk_module = b.dependency("extism_pdk", .{ .target = target, .optimize = optimize }).module("extism-pdk");
    var plugin = b.addExecutable(.{
        .name = "my-plugin",
        .root_source_file = .{ .path = "src/main.zig" },
        .target = target,
        .optimize = optimize,
    });
    plugin.rdynamic = true;
    plugin.entry = .disabled;
    plugin.root_module.addImport("extism-pdk", pdk_module);

    b.installArtifact(plugin);
    const plugin_example_step = b.step("my-plugin", "Build my-plugin");
    plugin_example_step.dependOn(b.getInstallStep());
}