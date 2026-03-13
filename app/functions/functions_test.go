package functions

import (
	"math"
	"reflect"
	"testing"
)

func TestNewUUID(t *testing.T) {
	uuid := NewUUID()
	if !(len(uuid) > 0 && uuid != "") {
		t.Error("UUID not generated")
	}
}

func TestRemoveDuplicate(t *testing.T) {
	initMap := []string{"a", "b", "c", "d"}
	resultMap := []string{"a", "b", "c", "d", "c", "d", "a"}
	RemoveDuplicate(&resultMap)
	if !reflect.DeepEqual(initMap, resultMap) {
		t.Error("the slices are different")
	}
}

func TestRound(t *testing.T) {
	if Round(3.14159265359, 0, 4) != 3.1416 {
		t.Error("rounding problems for : roundOn 0 and places 4 ")
	}

	if Round(3.14159265359, 0.1, 6) != 3.141593 {
		t.Error("rounding problems for : roundOn 0.1 and places 6 ")
	}

}

func TestHaversineDistance_SamePoint(t *testing.T) {
	d := HaversineDistance(48.8566, 2.3522, 48.8566, 2.3522)
	if d != 0 {
		t.Errorf("expected 0 for same point, got %f", d)
	}
}

func TestHaversineDistance_ParisToLyon(t *testing.T) {
	// Paris (48.8566, 2.3522) to Lyon (45.7640, 4.8357) ~392 km
	d := HaversineDistance(48.8566, 2.3522, 45.7640, 4.8357)
	if math.Abs(d-392217) > 5000 {
		t.Errorf("Paris-Lyon distance expected ~392km, got %f meters", d)
	}
}

func TestHaversineDistance_ShortDistance(t *testing.T) {
	// Two points ~111m apart (0.001 degree latitude at equator)
	d := HaversineDistance(0, 0, 0.001, 0)
	if math.Abs(d-111.2) > 1 {
		t.Errorf("expected ~111m, got %f meters", d)
	}
}

func TestHaversineDistance_InRange(t *testing.T) {
	bossLat, bossLon := 48.8566, 2.3522
	radiusMeters := 50.0

	// Player at ~30m away (should be in range)
	playerLat, playerLon := 48.8568, 2.3522
	d := HaversineDistance(playerLat, playerLon, bossLat, bossLon)
	if d > radiusMeters {
		t.Errorf("player should be in range: distance=%f, radius=%f", d, radiusMeters)
	}
}

func TestHaversineDistance_OutOfRange(t *testing.T) {
	bossLat, bossLon := 48.8566, 2.3522
	radiusMeters := 50.0

	// Player at ~500m away (should be out of range)
	playerLat, playerLon := 48.8610, 2.3522
	d := HaversineDistance(playerLat, playerLon, bossLat, bossLon)
	if d <= radiusMeters {
		t.Errorf("player should be out of range: distance=%f, radius=%f", d, radiusMeters)
	}
}

func TestEncrypt_Decrypt(t *testing.T) {
	passphrase := "Flashcards"
	phrase := "ABCD"
	b := []byte(phrase)
	encrypted, _ := Encrypt(b, passphrase)
	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Error("Decrypt error: " + err.Error())
	}

	if phrase != string(decrypted) {
		t.Error("Encrypt or Decrypt Error. phrase : " + phrase + " decrypted : " + string(decrypted))
	}

}
