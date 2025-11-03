# Installation and Execution Instructions

This document provides instructions on how to install and run the Free Income From Bandwidth on different platforms.

## GUI Application

The GUI application provides a user-friendly interface to manage your bandwidth-sharing applications.

### Windows

1.  Download the `free-income-from-bandwidth.exe` file from the latest release.
2.  Place the executable in a directory of your choice.
3.  Double-click the `free-income-from-bandwidth.exe` file to run the application.

### macOS

1.  Download the `free-income-from-bandwidth.dmg` file from the latest release.
2.  Open the `.dmg` file and drag the `Free Income From Bandwidth.app` to your `Applications` folder.
3.  Double-click the application in your `Applications` folder to run it.

### Linux

1.  Download the `free-income-from-bandwidth` binary for your architecture (amd64 or arm64) from the latest release.
2.  Make the binary executable:
    ```bash
    chmod +x free-income-from-bandwidth
    ```
3.  Run the application:
    ```bash
    ./free-income-from-bandwidth
    ```

## Headless Server (Linux)

The headless server mode allows you to run the application on a server without a graphical user interface. You can manage the application through a web interface or an HTTP API.

1.  Download the `free-income-from-bandwidth` binary for your architecture (amd64 or arm64) from the latest release.
2.  Make the binary executable:
    ```bash
    chmod +x free-income-from-bandwidth
    ```
3.  Run the application in headless mode:
    ```bash
    ./free-income-from-bandwidth --headless
    ```
    You can specify a port with the `--port` flag. The default port is `8080`.
    ```bash
    ./free-income-from-bandwidth --headless --port 8081
    ```

### Web Interface

Once the headless server is running, you can access the web interface from any device on the same network by navigating to `http://<server_ip>:<port>` in your web browser (e.g., `http://192.168.1.100:8080`).

The web interface provides the same functionality as the desktop application.

### API Access

The headless server also exposes a REST API to manage the application. Here is a brief overview of the available endpoints:

-   **Apps:**
    -   `GET /api/apps/summary`: Get dashboard summary.
    -   `GET /api/apps/available`: Get available apps.
    -   `GET /api/apps/running`: Get running apps.
    -   `POST /api/apps/start/{appID}`: Start an app.
    -   `POST /api/apps/stop/{appID}`: Stop an app.
    -   `POST /api/apps/restart/{appID}`: Restart an app.
    -   `GET /api/apps/logs/{appID}`: Get app logs.
    -   ...and more.
-   **Proxies:**
    -   `POST /api/proxies/add`: Add a proxy.
    -   `GET /api/proxies/list`: List proxies.
    -   ...and more.
-   **Settings:**
    -   `GET /api/settings`: Get settings.
    -   `POST /api/settings/autostart`: Set auto-start.
    -   `POST /api/settings/showintray`: Set show in tray.

For a complete list of API endpoints and their parameters, please refer to the `backend/api/http_server.go` file.
