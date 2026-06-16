<p align="center">
  <img src="build/appicon.png" width="128" alt="Mado-Tray logo">
</p>

# Mado-Tray

Mado-Tray is a startup process manager for macOS. It lives in the menu bar, lets you enable or disable scripts from a compact panel, and opens each process in a visible Terminal.app window so you can inspect it, stop it with `Ctrl+C`, or interact with it manually.

The name comes from **mado** (窓, “window” in Japanese). The app icon and menu bar symbol use a **shoji**-style window.

Built with Go, Wails, and a lightweight TypeScript frontend.

## Features

- Reads configuration from `~/.config/mado-tray/config.json`.
- Automatically runs all scripts with `is_active: true` when the app starts.
- Opens processes in Terminal.app using AppleScript instead of hiding them in the background.
- Lets you create, edit, delete, enable, or disable processes from the UI.
- Separate **Path** and **Arguments** fields, plus a native file picker to choose scripts or executables.
- Includes a **Run now** button to launch any process manually.
- Dark, minimalist UI with a fixed title bar (`MadoTray`, settings, close).
- Supports English and Spanish; language lives in the **Options** modal (⚙).
- Includes an **Open at macOS startup** switch inside Options to register or remove `Mado-Tray.app` from login items.
- Uses a frameless floating panel controlled from the menu bar.
- Keeps a persistent menu bar icon (shoji window) with **Show** and **Exit Mado-Tray** actions.
- Closing the panel with `×` hides the app; use **Exit Mado-Tray** to quit completely.
- Uses a single-instance lock: opening Mado-Tray again shows the existing window without running startup scripts again.
- Includes `LSUIElement` in `build/darwin/Info.plist` so the packaged app does not appear in the Dock.

## Requirements

- macOS.
- Go 1.23 or newer.
- Node.js 20 or newer recommended.
- Wails CLI v2.

Install Wails if you do not have it:

```sh
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Check your environment:

```sh
wails doctor
```

## Development

Install dependencies:

```sh
npm install
go mod tidy
```

Run the app in development mode:

```sh
wails dev
```

In development mode, the **Open at macOS startup** option may show a message indicating that there is no final `.app` bundle yet. That control is intended to work from the packaged app.

## Process Configuration

You can manage processes from the UI:

1. Press **Add process**.
2. Fill in **Name**, **Path**, and optionally **Arguments**.
3. Use **Browse** to pick a script or executable from the system file dialog.
4. Enable **Active on startup** if you want it to run when Mado-Tray opens.
5. Save the process.
6. Use **Edit** or **Delete** on each process whenever you need to change it.

The configuration also lives as an editable JSON file at:

```text
~/.config/mado-tray/config.json
```

If the file does not exist, Mado-Tray creates it with an initial example:

```json
[
  {
    "id": "example",
    "name": "Example script",
    "path": "/path/to/your/script.sh",
    "args": "",
    "is_active": false
  }
]
```

Fields:

- `id`: unique process identifier.
- `name`: display name in the UI.
- `path`: absolute path to the script or executable.
- `args`: optional arguments passed to the command (for example `--port 3000`).
- `is_active`: when `true`, Mado-Tray opens it automatically on startup.

Real example:

```json
[
  {
    "id": "api-local",
    "name": "API local",
    "path": "/Users/your_user/Projects/api/start.sh",
    "args": "--port 3000",
    "is_active": true
  },
  {
    "id": "worker",
    "name": "Worker",
    "path": "/Users/your_user/Projects/worker/run.sh",
    "args": "",
    "is_active": false
  }
]
```

If an older config stored arguments inside `path`, Mado-Tray splits them into `path` and `args` automatically on load.

The script must have executable permissions:

```sh
chmod +x /Users/your_user/Projects/api/start.sh
```

## Language

The interface starts in English by default. If the system reports Spanish, Mado-Tray uses Spanish automatically.

Open **Options** (⚙ in the title bar) and choose **ES** or **EN**. The preference is stored locally in Wails' embedded browser and persists across sessions.

## macOS Startup

From the UI, open **Options** (⚙) and use the **Open at macOS startup** switch.

When enabled, Mado-Tray registers the current `.app` as a macOS login item using `System Events`. When disabled, it removes that login item.

When macOS starts and launches Mado-Tray as a login item, the app starts hidden and still runs every active process.

The menu bar icon stays visible while Mado-Tray is running. Use **Show** to restore the panel and **Exit Mado-Tray** to close the app.

If Mado-Tray is already running and you open it again, the existing instance becomes visible and startup scripts are not launched a second time.

For this to work correctly:

1. Build the app.
2. Move `Mado-Tray.app` to `/Applications` or another stable folder.
3. Open that packaged app.
4. Open **Options** and enable **Open at macOS startup**.

If you move the `.app` after registering it, disable and enable the switch again.

## Build

Build the app:

```sh
npm run build:app
```

The output is:

```text
build/bin/Mado-Tray.app
```

The build script runs `npm install`, `go mod tidy`, `npm run build`, and `wails build`, then verifies that the final `.app` bundle exists.

You can also run Wails directly if you only need the packaging step:

```sh
wails build
```

The app icon is generated from `build/appicon.png`. macOS may cache icons; if the bundle icon does not update immediately after a rebuild, move the `.app` or restart Finder.

The `build/darwin/Info.plist` template already includes:

```xml
<key>LSUIElement</key>
<true/>
```

This keeps Mado-Tray out of the Dock and makes it live primarily in the menu bar.

## Publishing on GitHub

Suggested release checklist:

1. Run `npm install` and `go mod tidy`.
2. Run `npm run build`.
3. Run `wails build`.
4. Open `build/bin/Mado-Tray.app` and test:
   - `config.json` loading;
   - process toggles;
   - **Browse**, **Path**, and **Arguments** fields;
   - **Run now** button;
   - **Open at macOS startup** option;
   - menu bar menu and icon.
5. Create a tag, for example:

```sh
git tag v0.1.0
git push origin v0.1.0
```

6. Create a GitHub release and attach the packaged app or a `.dmg` if you decide to distribute it that way.

## Structure

```text
.
├── app.go                 # Methods exposed to Wails
├── config.go              # config.json read/write logic
├── runner.go              # Visible execution in Terminal.app
├── startup.go             # macOS login item support
├── tray_darwin.m          # Menu bar icon and native menu
├── main.go                # Wails, window, and menu bar setup
├── build/
│   ├── appicon.png        # App icon source (shoji window)
│   └── darwin/Info.plist
├── frontend/
│   ├── index.html
│   └── src/
│       ├── main.ts
│       └── style.css
└── wails.json
```

## Security Notes

Mado-Tray executes the paths defined in your JSON file. Use absolute paths, review scripts before enabling them, and avoid pointing to files downloaded from untrusted sources.
