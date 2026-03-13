package curd

import "testing"

func TestPaletteStyles(t *testing.T) {
	palettes := []struct {
		name string
		p    Palette
	}{
		{"Swissgit", SwissgitPalette},
		{"Raclette", RaclettePalette},
		{"Fondue", FonduePalette},
		{"Gruyere", GruyerePalette},
	}

	for _, tc := range palettes {
		t.Run(tc.name, func(t *testing.T) {
			ss := tc.p.Styles()

			// Verify key styles render non-empty output.
			checks := []struct {
				label  string
				render string
			}{
				{"Title", ss.Title.Render("test")},
				{"Subtitle", ss.Subtitle.Render("test")},
				{"Selected", ss.Selected.Render("test")},
				{"Normal", ss.Normal.Render("test")},
				{"Help", ss.Help.Render("test")},
				{"MenuActiveName", ss.MenuActiveName.Render("test")},
				{"SuccessStyle", ss.SuccessStyle.Render("test")},
				{"FailStyle", ss.FailStyle.Render("test")},
				{"AccentStyle", ss.AccentStyle.Render("test")},
			}
			for _, c := range checks {
				if c.render == "" {
					t.Errorf("%s: %s style renders empty string", tc.name, c.label)
				}
			}
		})
	}
}

func TestPaletteStylesDistinct(t *testing.T) {
	// Different palettes should produce different accent-colored output.
	sgTitle := SwissgitPalette.Styles().Title.Render("x")
	raTitle := RaclettePalette.Styles().Title.Render("x")
	if sgTitle == raTitle {
		t.Error("SwissgitPalette and RaclettePalette Title styles should differ")
	}
}

func TestPaletteLogoGradient(t *testing.T) {
	palettes := []struct {
		name string
		p    Palette
	}{
		{"Swissgit", SwissgitPalette},
		{"Raclette", RaclettePalette},
		{"Fondue", FonduePalette},
		{"Gruyere", GruyerePalette},
	}
	for _, tc := range palettes {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.p.LogoGradient) == 0 {
				t.Errorf("%s: LogoGradient should not be empty", tc.name)
			}
		})
	}
}
