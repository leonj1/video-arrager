# Video Arranger

A simple macOS app to arrange and combine video files. Add videos, drag to reorder, and export as a single file with optional transitions.

![Video Arranger](https://img.shields.io/badge/platform-macOS-blue) ![Go](https://img.shields.io/badge/Go-1.21+-00ADD8)

## Features

- Drag-and-drop video import from Finder
- Reorder videos by dragging
- Video thumbnails, duration, and resolution display
- Preview pane with system player integration
- Export with fade/crossfade transitions
- Save/load projects as JSON
- Keyboard shortcuts (Cmd+N, Cmd+O, Cmd+S, Cmd+E)

## Installation

### Download (Recommended)

Download the latest binary from [Releases](https://github.com/yourusername/video-arranger/releases).

```bash
# Make it executable (if needed)
chmod +x video-arranger

# Run
./video-arranger
```

### Build from Source

Requires: Go 1.21+, ffmpeg

```bash
make build
./video-arranger
```

## Requirements

- macOS
- [ffmpeg](https://ffmpeg.org/) installed (`brew install ffmpeg`)

## Usage

1. **Add videos**: Click "Add Videos" or drag files onto the window
2. **Reorder**: Drag videos up/down in the list
3. **Export**: Click "Export", choose transition type, save

## License

MIT
