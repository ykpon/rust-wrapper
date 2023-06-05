# Rust Server Management Tool

A convenient tool for managing a Rust game server. The goal is to provide the ability to manage the server from a Linux console.

## Key Features

1. **Server Management:** Allows users to send commands directly to the Rust server console via the RCON protocol. As the Rust server itself does not provide such a feature due to the peculiarities of the Unity game engine implementation, this program implements this functionality through the RCON interface.
2. **Logging:** Takes on the functions of outputting the server console output stream, outputting all regular messages and RCON responses as well.
3. **Docker Compatibility:** Also designed for use with a Rust server in Docker.

## Environment Variables

Uses the following environment variables for configuration:

- `RCON_IP`: The IP address of the Rust server. If not specified, defaults to `127.0.0.1`.
- `RCON_PORT`: The RCON port of the Rust server. If not specified, defaults to `28018`.
- `RCON_PASS`: The RCON password for the Rust server. There is no default value.

## Usage

After setting the appropriate environment variables, you can start the program with your Rust server executable and any server parameters.

## Note

The tool is currently in a beta state, and as such, it might not be fully functional. Therefore, it is not recommended for production use without thorough testing. Please feel free to contribute by raising issues or creating pull-requests.