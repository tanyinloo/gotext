package gotext

import (
	"embed"
	"testing"
)

var (
	//go:embed fixtures
	enUSFixture embed.FS
)

//since both Po and Mo just pass-through to Domain for MarshalBinary and UnmarshalBinary, test it here
func TestBinaryEncoding(t *testing.T) {
	// Create po objects
	po := NewPo()
	po2 := NewPo()

	f, err := enUSFixture.Open("fixtures/en_US/default.po")
	if err != nil {
		t.Fatal(err)
	}
	// Parse file
	po.ParseFile(f)

	buff, err := po.GetDomain().MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	err = po2.GetDomain().UnmarshalBinary(buff)
	if err != nil {
		t.Fatal(err)
	}

	// Test translations
	tr := po2.Get("My text")
	if tr != translatedText {
		t.Errorf("Expected '%s' but got '%s'", translatedText, tr)
	}
	// Test translations
	tr = po2.Get("language")
	if tr != "en_US" {
		t.Errorf("Expected 'en_US' but got '%s'", tr)
	}
}

func TestDomain_GetTranslations(t *testing.T) {
	po := NewPo()

	f, err := enUSFixture.Open("fixtures/en_US/default.po")
	if err != nil {
		t.Fatal(err)
	}
	po.ParseFile(f)

	domain := po.GetDomain()
	all := domain.GetTranslations()

	if len(all) != len(domain.translations) {
		t.Error("lengths should match")
	}

	for k, v := range domain.translations {
		if all[k] == v {
			t.Error("GetTranslations should be returning a copy, but pointers are equal")
		}
		if all[k].ID != v.ID {
			t.Error("IDs should match")
		}
		if all[k].PluralID != v.PluralID {
			t.Error("PluralIDs should match")
		}
		if all[k].dirty != v.dirty {
			t.Error("dirty flag should match")
		}
		if len(all[k].Trs) != len(v.Trs) {
			t.Errorf("Trs length does not match: %d != %d", len(all[k].Trs), len(v.Trs))
		}
		if len(all[k].Refs) != len(v.Refs) {
			t.Errorf("Refs length does not match: %d != %d", len(all[k].Refs), len(v.Refs))
		}
	}
}
