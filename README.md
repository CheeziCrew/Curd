# Curd

![curd](https://github.com/user-attachments/assets/8afa06ce-103a-4ef0-a8df-e58837b8bd2a)

The glue that holds the cheese together. Shared TUI component library for [swissgit](https://github.com/CheeziCrew/swissgit), [raclette](https://github.com/CheeziCrew/raclette), and [fondue](https://github.com/CheeziCrew/fondue).

## Status

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=CheeziCrew_Curd&metric=alert_status)](https://sonarcloud.io/summary/overall?id=CheeziCrew_Curd)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=CheeziCrew_Curd&metric=reliability_rating)](https://sonarcloud.io/summary/overall?id=CheeziCrew_Curd)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=CheeziCrew_Curd&metric=security_rating)](https://sonarcloud.io/summary/overall?id=CheeziCrew_Curd)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=CheeziCrew_Curd&metric=sqale_rating)](https://sonarcloud.io/summary/overall?id=CheeziCrew_Curd)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=CheeziCrew_Curd&metric=vulnerabilities)](https://sonarcloud.io/summary/overall?id=CheeziCrew_Curd)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=CheeziCrew_Curd&metric=bugs)](https://sonarcloud.io/summary/overall?id=CheeziCrew_Curd)

## Components

| Component         | Description                                              |
| ----------------- | -------------------------------------------------------- |
| `Palette`         | Color scheme (accent, accentBright, secondary, gradient) |
| `StyleSet`        | Pre-computed lipgloss styles from a palette              |
| `MenuModel`       | App menu with ASCII banner, tagline, selectable items    |
| `RepoSelectModel` | Multi-select repo picker with search and scroll          |
| `ProgressModel`   | Task progress tracker with status indicators             |
| `ResultModel`     | Success/fail summary renderer                            |
| `HintBar`         | Context-aware keybinding hints                           |
| `KeyMap`          | Standard key bindings                                    |

Each app has its own palette:
- **SwissgitPalette** — magenta/cyan
- **RaclettePalette** — yellow/red
- **FonduePalette** — green/orange

## Install

```bash
go get github.com/CheeziCrew/curd
```

## Built With

- [bubbletea v2](https://github.com/charmbracelet/bubbletea) — TUI framework
- [lipgloss v2](https://github.com/charmbracelet/lipgloss) — Terminal styling

## License

[MIT](LICENSE)
