package main

import "testing"

func TestSplitProjectPhotos(t *testing.T) {
	photos := []string{"1.jpg", "2.jpg", "3.jpg", "4.jpg", "5.jpg"}
	primary, secondary := splitProjectPhotos(photos)

	if len(primary) != 3 || len(secondary) != 2 {
		t.Fatalf("split sizes = %d and %d, want 3 and 2", len(primary), len(secondary))
	}
	if primary[0] != "1.jpg" || primary[2] != "5.jpg" {
		t.Fatalf("unexpected primary row: %v", primary)
	}
	if secondary[0] != "2.jpg" || secondary[1] != "4.jpg" {
		t.Fatalf("unexpected secondary row: %v", secondary)
	}
}

func TestSplitProjectPhotosKeepsSinglePhotoInBothRows(t *testing.T) {
	photos := []string{"only.jpg"}
	primary, secondary := splitProjectPhotos(photos)

	if len(primary) != 1 || len(secondary) != 1 || primary[0] != secondary[0] {
		t.Fatalf("single photo split = %v and %v", primary, secondary)
	}
}
