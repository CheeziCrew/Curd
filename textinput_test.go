package curd

import (
	"testing"
)

func TestNewStyledInput(t *testing.T) {
	ti := NewStyledInput("enter branch name", SwissgitPalette)

	if ti.Placeholder != "enter branch name" {
		t.Errorf("Placeholder = %q, want %q", ti.Placeholder, "enter branch name")
	}
}

func TestNewStyledInput_DifferentPalettes(t *testing.T) {
	palettes := []struct {
		name    string
		palette Palette
	}{
		{"swissgit", SwissgitPalette},
		{"raclette", RaclettePalette},
		{"fondue", FonduePalette},
	}

	for _, tt := range palettes {
		t.Run(tt.name, func(t *testing.T) {
			ti := NewStyledInput("test", tt.palette)
			if ti.Placeholder != "test" {
				t.Errorf("Placeholder = %q, want %q", ti.Placeholder, "test")
			}
		})
	}
}
