# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

RLMusic is a multi-platform local music player built with Vue 3 (frontend) and Go (backend). It targets Web, Electron desktop, and Android (Capacitor) from a single codebase. Key features include local music library management, playlist management, AI-generated playlist descriptions (SiliconFlow LLM), TTS podcast intros (Qwen TTS), real-time "Listen Together" via WebSocket, and an admin dashboard with server log streaming.

## Common Commands

### Development

```bash
# Install dependencies (uses pnpm)
pnpm install

# Start frontend dev server (port 23456, proxies /api to backend)
pnpm dev:web

# Start Electron dev (uses VITE_DEV_SERVER_URL)
pnpm dev

# Start backend with hot reload (requires `air` installed)
cd server && air

# Or start backend directly
cd server && go run cmd/main.go

# Build Go backend binary
npm run build:go
```

### Building

```bash
# Web only (output: dist/)
pnpm build:web

# Electron client-only (connects to remote server)
pnpm build:client

# Electron server-bundled (includes Go backend binary)
pnpm build:server

# Android APK
pnpm build:android

# Build all targets
pnpm build:all
```

### Go Backend

```bash
cd server

# Run tests
go test ./...

# Run a specific test
go test ./internal/handle -run TestFunctionName -v

# Build binary
go build -o ../resources/server.exe cmd/main.go

# Tidy dependencies
go mod tidy
```

### Docker

```bash
# Build and run both frontend (Nginx) and backend (Go)
docker compose up -d --build

# Restart backend after music file changes
docker compose restart backend
```

## Architecture

### Frontend (`src/`)

- **Framework**: Vue 3 + TypeScript + Vite
- **UI Library**: Naive UI (auto-imported via `unplugin-vue-components` + `NaiveUiResolver`)
- **State**: Pinia (`src/store/`) — `musicData.ts` (playback state), `userData.ts`, `chatData.ts` (WebSocket), `settingData.ts`
- **Router**: Vue Router with `createWebHashHistory`, navigation guards for login (`meta.needLogin`) and Electron/Capacitor init flow (`/init` page)
- **API Layer**: `src/api/` — Axios wrappers organized by domain (`song.ts`, `playlist.ts`, `search.ts`, `user.ts`, `system.ts`, `ai.ts`)
- **Auto-imports**: `unplugin-auto-import` handles Vue APIs (`ref`, `reactive`, etc.) and Naive UI composables (`useDialog`, `useMessage`, `useNotification`, `useLoadingBar`). No need to import these manually.
- **Path Alias**: `@` maps to `src/`
- **SCSS**: Global styles imported automatically via `vite.config.ts` (`@use "@/style/index.scss" as *;`)
- **Request Utility**: `src/utils/request.ts` — Axios instance with interceptors for auth tokens and error handling

### Backend (`server/`)

- **Framework**: Go 1.25 + Gin
- **ORM**: GORM with SQLite (`github.com/glebarez/sqlite` — pure Go, no CGO)
- **Entry**: `cmd/main.go` — parses `-c config.yml`, initializes logger, DB, registers handlers
- **Layer Structure**:
  - `internal/handle/` — HTTP handlers (controllers), grouped by domain (`handle_song.go`, `handle_auth.go`, etc.)
  - `internal/model/` — GORM models and DB operations (`Song.go`, `Playlist.go`, `User.go`, etc.)
  - `internal/middleware/` — Gin middleware (`auth.go` for JWT, `base.go` for DB injection, `stats.go`)
  - `internal/global/` — Config structs, logger ring buffer (`logger.go`), unified response format (`result.go`)
  - `internal/utils/` — Utility packages (`jwt/`, `encrypt/`, `audio/`, `imgtool/`, `ai/`, `prompt/`)
- **Routing**: `internal/Manager.go` registers all routes. API base path is `/api`.
- **Swagger**: Available at `/swagger/index.html` in development

### Electron (`electron/`)

- **Main Process**: `main.ts` — window management, system tray, IPC handlers, Go server spawning
- **Preload**: `preload.ts` — exposes safe Node.js APIs to renderer
- **Dual Mode**:
  - `client`: Renderer connects to external backend (like web mode)
  - `server`: Bundles `server.exe`/`server` binary in `resources/`, spawns it on launch, starts an internal HTTP server for the frontend
- **First-run Flow**: Uninitialized Electron/Capacitor apps redirect to `/init` for base folder and port configuration
- **IPC Channels**: `app-config-get`, `app-clear-data`, `select-directory`, `show-save-dialog`, `save-file`, `get-local-ips`, `check-ports`, `apply-initial-config`

### Docker

- **Frontend**: `Dockerfile.web` — multi-stage build (Node → Nginx), serves on port 23456, proxies `/api/` and `/covers/` to backend container
- **Backend**: `server/Dockerfile` — multi-stage Go build (Alpine), serves on port 12345
- **Compose**: `docker-compose.yml` with bind mounts for `./data`, `./log`, and `${MUSIC_BIND_PATH:-C:/RLMusic}:/music`

## Key Conventions & Patterns

### Prompt Loading (AI Features)

Prompt files live in `server/prompts/` as Markdown. The `prompt.Read(name)` function (`server/internal/utils/prompt/Prompt.go`) loads them via a three-tier fallback:

1. `embed.FS` (compiled into binary at build time — see `server/prompts/embed.go`)
2. Filesystem relative to executable (`exeDir/prompts/` → `exeDir/../prompts/`)
3. Current working directory (`./prompts/`)

**Important**: When adding new prompt files, update the `//go:embed` pattern in `server/prompts/embed.go` if the filename doesn't match existing patterns.

### Log System

The backend uses a custom `slog` handler (`server/internal/Helper.go`) that writes to both stdout and an in-memory ring buffer (`server/internal/global/logger.go`). The frontend admin page consumes logs via an SSE endpoint (`/api/system/logs/stream`) for real-time display.

### Music Streaming

Song playback uses `/api/song/stream/:id` which streams audio files directly from disk with proper `Content-Type`, `Content-Length`, and `Accept-Ranges` headers for seek support.

### Cover Image Resolution

Cover URLs are resolved through `resolveCoverUrl()` in the frontend, handling both local `/covers/` paths and remote URLs.

### TypeScript Strictness

`tsconfig.app.json` enables `strict`, `noUnusedLocals`, `noUnusedParameters`, and `erasableSyntaxOnly`. The build will fail on unused variables or implicit `any`.

## Environment Configuration

### Frontend (`.env`)

| Variable | Purpose |
|----------|---------|
| `VITE_MUSIC_API` | Backend API base URL (default: `http://localhost:12345`) |
| `VITE_APP_MODE` | `web` / `client` / `server` |
| `VITE_ANN_TITLE` / `VITE_ANN_CONTENT` | Home page announcement |

### Backend (`server/config.yml`)

| Key | Purpose |
|-----|---------|
| `Server.Port` | Listen address (default `:12345`) |
| `BasicPath.FilePath` / `FileName` | Base directory for music files |
| `SiliconFlow.ApiKey` / `Model` | LLM for AI descriptions |
| `QwenTTS.ApiKey` / `Model` / `Voice` | TTS for podcast intros |
| `JWT.Secret` / `Expire` | Auth token config |

AI API keys can also be provided via environment variables: `SiliconFlow_API_KEY`, `QwenTTS_API_KEY`.

## Testing

- **Go**: `go test ./...` from `server/`. Test files follow `*_test.go` convention.
- **Frontend**: No test runner configured in this project.
- **Load Testing**: `api_test_scan.py` at repo root provides staged pressure tests for API endpoints and streaming playback.
