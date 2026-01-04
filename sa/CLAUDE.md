# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Run the game (from client directory)
cd client && go run ./src/main

# Build executable (from client directory)
cd client && go build -o ./bin/sa.exe ./src/main

# Generate Protocol Buffer code (from repository root)
./scripts/gen.pb.sh
```

## Architecture Overview

This is a **pure client-side ARPG game** built with Go 1.24.1 and Ebiten v2 (2D game engine).

### Project Structure

```
sa/
├── proto/           # Protocol Buffer definitions (.proto files)
├── scripts/         # Build scripts (gen.pb.sh)
└── client/
    ├── cfg/         # YAML configuration files
    ├── res/         # Game resources (sprites, maps, fonts)
    └── src/
        ├── main/    # Entry point, config loading
        ├── game/    # Ebiten game loop (Update/Draw/Layout)
        ├── user/    # Core game logic (role, enemy, scene, AI)
        ├── cfg/     # Configuration loaders
        ├── res/     # Resource loaders
        ├── common/  # Utilities (camera, coordinates, animation)
        ├── ui/      # UI components
        └── proto/   # Generated protobuf code
```

### Key Modules

| Module | Responsibility |
|--------|---------------|
| `user/` | Main game logic: Role, ArpgEnemy, Scene, AI systems |
| `cfg/` | YAML config parsing, Tiled map (.tmx) parsing |
| `res/` | Sprite/animation resource loading |
| `common/` | Coordinate transforms, camera, rendering interfaces |

### Data Flow

1. **Config Loading**: YAML files → `cfg.*Mgr.Load()` → Global singletons
2. **Resource Loading**: Image files → `res.*Mgr.Load()` → Cached sprites
3. **Game Loop**: `game.Update()` → `user.Update()` → AI/Input → State changes
4. **Rendering**: `game.Draw()` → Y-sorting → Sprite rendering

## Code Conventions

### Global Singletons
All global managers prefixed with `G`:
```go
cfg.GCommon      // Global config
cfg.GMapMgr      // Map configurations
res.GRoleMgr     // Role resources
user.GUser       // Current user instance
```

### Manager Pattern
All managers implement three-phase initialization:
```go
func (m *Manager) Load() error    // Load from file
func (m *Manager) Check() error   // Validate data
func (m *Manager) Assemble() error // Link dependencies
```

### Coordinate Systems
Three coordinate types for isometric maps:
- **Tile**: Grid coordinates (used by Tiled editor)
- **World**: Pixel coordinates (game world)
- **Screen**: Display coordinates (after camera transform)

Conversion functions in `common/coordinatetransform/`:
```go
T2W(tx, ty) → (wx, wy)  // Tile to World
W2T(wx, wy) → (tx, ty)  // World to Tile
W2S(wx, wy) → (sx, sy)  // World to Screen (with camera)
```

### Asset ID Encoding
Asset IDs are segmented by type to avoid conflicts:
```
Role:       1,000,001 - 1,999,999
Map:        2,000,001 - 2,999,999
Item:       3,000,001 - 3,999,999
Pet:        4,000,001 - 4,999,999
```

### Animation System
Frame-based animation with tick counter:
```go
const FrameTickPerChange = 6  // Switch frame every 6 ticks
```

### Entity Positioning
Entities use foot-center as anchor point (`WX`, `WY`). Rendering adjusts to image top-left:
```go
screenX -= frameWidth / 2
screenY -= frameHeight
```

## Combat System (ARPG)

### AI State Machine
Both enemies and player AI use state machine:
- `Idle` → `Chase` → `Attack` → `Return`

### Damage System
- `TakeDamage(damage)` reduces HP
- Hit flash effect via `ColorScale` (red tint that fades)
- Attack triggered on animation hit frames (`HitFrame`)

### Pathfinding
A* algorithm implemented in `user/arpg.ai.astar.go`

## Configuration Files

| File | Purpose |
|------|---------|
| `common.yaml` | Global settings (window size, default stats) |
| `role.yaml` | Role definitions |
| `pet.yaml` | Pet definitions |
| `map.yaml` | Map configurations |
| `enemy.group.yaml` | Enemy spawn groups |
| `portal.yaml` | Teleport points |

## Tiled Map Support

Full TMX format support in `cfg/tiled.map.*.go`:
- Multiple layers
- Object groups (for collision, portals)
- Tileset references
- Y-sorting for correct overlap rendering

### Collision Detection Approaches

| Approach | Description | Recommendation |
|----------|-------------|----------------|
| **Tile-based** | Check tile ID at (tx,ty) | ❌ Imprecise (grid-locked) |
| **Object Layer** | Draw rectangles manually with `blocked=true` | ✅ Recommended |
| **Tileset Collision** | Define collision boxes in tileset editor | ⚠️ Optional |

## Related Tools

- **Tiled Map Editor**: https://www.mapeditor.org/ (isometric 45° mode)
- **go-tiled**: `github.com/lafriks/go-tiled` for TMX parsing
- **TexturePacker**: Sprite atlas packing
- **Sprite Reel**: https://spritereel.com/ for sprite preview
