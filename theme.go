package curd

import (
	"image/color"
	"charm.land/lipgloss/v2"
)

// Base16 ANSI colors — respects terminal theme (works with Tinty, base16 etc).
var (
	ColorBg      = lipgloss.Color("0")
	ColorRed     = lipgloss.Color("1")
	ColorGreen   = lipgloss.Color("2")
	ColorYellow  = lipgloss.Color("3")
	ColorBlue    = lipgloss.Color("4")
	ColorMagenta = lipgloss.Color("5")
	ColorCyan    = lipgloss.Color("6")
	ColorFg      = lipgloss.Color("7")
	ColorGray    = lipgloss.Color("8")
	ColorBrRed   = lipgloss.Color("9")
	ColorBrGreen = lipgloss.Color("10")
	ColorBrYellow = lipgloss.Color("11")
	ColorBrBlue  = lipgloss.Color("12")
	ColorBrMag   = lipgloss.Color("13")
	ColorBrCyan  = lipgloss.Color("14")
	ColorBrWhite = lipgloss.Color("15")
)

// Palette defines app-specific accent colors that give each tool its identity.
type Palette struct {
	// Primary accent (used for active items, borders, titles).
	Accent color.Color
	// Bright variant of accent.
	AccentBright color.Color
	// Secondary color (borders, subtle accents).
	Secondary color.Color
	// Gradient colors for the logo banner (top to bottom).
	LogoGradient []color.Color
}

// Pre-built palettes for each app.
var (
	// SwissgitPalette — magenta/cyan theme.
	SwissgitPalette = Palette{
		Accent:       ColorMagenta,
		AccentBright: ColorBrMag,
		Secondary:    ColorBlue,
		LogoGradient: []color.Color{ColorBrMag, ColorMagenta, ColorBrBlue, ColorBlue},
	}

	// RaclettePalette — yellow/red theme.
	RaclettePalette = Palette{
		Accent:       ColorYellow,
		AccentBright: ColorBrYellow,
		Secondary:    ColorRed,
		LogoGradient: []color.Color{ColorBrYellow, ColorYellow, ColorBrRed, ColorRed},
	}

	// FonduePalette — green/orange theme (full cheese mode 🧀).
	FonduePalette = Palette{
		Accent:       ColorGreen,
		AccentBright: ColorBrGreen,
		Secondary:    ColorYellow,
		LogoGradient: []color.Color{ColorBrGreen, ColorGreen, ColorBrYellow, ColorYellow},
	}
)

// Styles returns a StyleSet derived from the given palette.
// This is the single source of truth for all shared styles.
func (p Palette) Styles() StyleSet {
	return StyleSet{
		Title: lipgloss.NewStyle().Bold(true).Foreground(p.AccentBright).MarginBottom(1),
		Subtitle: lipgloss.NewStyle().Foreground(p.Accent).Italic(true),
		Help: lipgloss.NewStyle().Foreground(ColorGray),
		HelpMargin: lipgloss.NewStyle().Foreground(ColorGray).MarginTop(1),

		Selected: lipgloss.NewStyle().Foreground(p.AccentBright).Bold(true),
		Normal:   lipgloss.NewStyle().Foreground(ColorFg),
		Dim:      lipgloss.NewStyle().Foreground(ColorGray),

		InputBox: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorGray).Padding(0, 2),

		LogoBox: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(p.Accent).Padding(1, 3).MarginBottom(1),
		Tagline: lipgloss.NewStyle().Foreground(ColorGray).Italic(true),
		Version: lipgloss.NewStyle().Foreground(p.Accent).Bold(true),

		MenuActiveItem:   lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(p.AccentBright).Padding(0, 1).Width(50),
		MenuInactiveItem: lipgloss.NewStyle().Padding(0, 2).Width(52),
		MenuActiveName:   lipgloss.NewStyle().Foreground(p.AccentBright).Bold(true),
		MenuInactiveName: lipgloss.NewStyle().Foreground(ColorGray),
		MenuActiveDesc:   lipgloss.NewStyle().Foreground(p.Accent),
		MenuInactiveDesc: lipgloss.NewStyle().Foreground(ColorGray),

		CursorMark:   lipgloss.NewStyle().Foreground(p.AccentBright).Bold(true),
		CheckStyle:   lipgloss.NewStyle().Foreground(p.AccentBright).Bold(true),
		UncheckStyle: lipgloss.NewStyle().Foreground(ColorGray),

		RepoActiveItem:    lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false, false, false, true).BorderForeground(p.Secondary).PaddingLeft(1),
		RepoInactiveItem:  lipgloss.NewStyle().PaddingLeft(3),
		RepoCursorName:    lipgloss.NewStyle().Foreground(p.Secondary).Bold(true),
		RepoSelectedName:  lipgloss.NewStyle().Foreground(p.Accent).Bold(true),
		RepoUnselectedName: lipgloss.NewStyle().Foreground(ColorGray),
		DirtyStyle:        lipgloss.NewStyle().Foreground(ColorYellow),
		CleanMark:         lipgloss.NewStyle().Foreground(ColorGreen),
		BranchMark:        lipgloss.NewStyle().Foreground(p.Accent),

		SummaryBox: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(p.AccentBright).Padding(0, 2).Bold(true),
		SuccessBox: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorGreen).Padding(0, 1),
		FailBox:    lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorRed).Padding(0, 1),

		SuccessStyle: lipgloss.NewStyle().Foreground(ColorGreen),
		FailStyle:    lipgloss.NewStyle().Foreground(ColorRed),
		PendingStyle: lipgloss.NewStyle().Foreground(ColorGray),
		ErrStyle:     lipgloss.NewStyle().Foreground(ColorRed),
		NameStyle:    lipgloss.NewStyle().Foreground(ColorFg),
		CountStyle:   lipgloss.NewStyle().Bold(true).Foreground(p.AccentBright),
		AccentStyle:  lipgloss.NewStyle().Foreground(p.AccentBright),
		ResultDim:    lipgloss.NewStyle().Foreground(ColorGray),
		ResultOk:     lipgloss.NewStyle().Foreground(ColorGreen),
		ResultFail:   lipgloss.NewStyle().Foreground(ColorRed),
		ResultAccent: lipgloss.NewStyle().Foreground(p.AccentBright),
	}
}

// StyleSet holds all shared styles derived from a Palette.
type StyleSet struct {
	// General
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Help       lipgloss.Style
	HelpMargin lipgloss.Style
	Selected   lipgloss.Style
	Normal     lipgloss.Style
	Dim        lipgloss.Style
	InputBox   lipgloss.Style

	// Logo
	LogoBox lipgloss.Style
	Tagline lipgloss.Style
	Version lipgloss.Style

	// Menu
	MenuActiveItem   lipgloss.Style
	MenuInactiveItem lipgloss.Style
	MenuActiveName   lipgloss.Style
	MenuInactiveName lipgloss.Style
	MenuActiveDesc   lipgloss.Style
	MenuInactiveDesc lipgloss.Style

	// Checkmarks
	CursorMark   lipgloss.Style
	CheckStyle   lipgloss.Style
	UncheckStyle lipgloss.Style

	// Repo select
	RepoActiveItem     lipgloss.Style
	RepoInactiveItem   lipgloss.Style
	RepoCursorName     lipgloss.Style
	RepoSelectedName   lipgloss.Style
	RepoUnselectedName lipgloss.Style
	DirtyStyle         lipgloss.Style
	CleanMark          lipgloss.Style
	BranchMark         lipgloss.Style

	// Result / Progress
	SummaryBox   lipgloss.Style
	SuccessBox   lipgloss.Style
	FailBox      lipgloss.Style
	SuccessStyle lipgloss.Style
	FailStyle    lipgloss.Style
	PendingStyle lipgloss.Style
	ErrStyle     lipgloss.Style
	NameStyle    lipgloss.Style
	CountStyle   lipgloss.Style
	AccentStyle  lipgloss.Style
	ResultDim    lipgloss.Style
	ResultOk     lipgloss.Style
	ResultFail   lipgloss.Style
	ResultAccent lipgloss.Style
}
