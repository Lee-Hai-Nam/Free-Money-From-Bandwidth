# Bandwidth Income Manager

A cross-platform desktop application to manage and monitor various bandwidth-sharing applications, helping you earn passive income effortlessly.

## Features

- **Unified Dashboard:** Monitor your earnings and statistics from multiple bandwidth-sharing apps in one place.
- **App Store:** Discover and set up new bandwidth-sharing applications easily.
- **Container Management:** Manage the lifecycle of your bandwidth-sharing apps running in Docker containers (start, stop, restart, view logs).
- **Proxy Configuration:** Configure and manage proxies for your applications.
- **Secure Credential Storage:** Your credentials are encrypted and stored securely.

## Supported Applications

- Honeygain
- EarnApp
- IPRoyal Pawns
- PacketStream

## Technologies Used

- **Backend:** Go
- **Frontend:** React, TypeScript, Vite
- **Framework:** Wails v2
- **Containerization:** Docker
- **Styling:** Tailwind CSS

## Getting Started

### Prerequisites

- Go (1.21+)
- Node.js (18+) & npm
- Docker Desktop
- Wails v2 CLI

### Installation

1.  **Install Dependencies:**
    ```bash
    make install
    # or
    go mod tidy && npm --prefix frontend install
    ```

2.  **Run in Development Mode:**
    ```bash
    make dev
    # or
    wails dev
    ```
    This will start the application with hot-reloading for both the frontend and backend.

3.  **Build for Production:**
    ```bash
    make build
    # or
    wails build
    ```
    The executable will be located in the `build/bin` directory.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.