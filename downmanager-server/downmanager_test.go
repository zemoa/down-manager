package main

import (
	"os"
	"testing"
)

func TestGiveFileLink1Fichier(t *testing.T) {
	data, err := os.ReadFile("testfile.html")
	if err != nil {
		t.Fatalf("Error while reading test file: %s", err.Error())
	}
	want := DLForm1Fichier{link: "https://1fichier.com/?ecrarnm5mdj3ig863rvn", adz: "929.997412863286"}
	link := GiveFileLink1Fichier(string(data))
	if want.link != link.link {
		t.Fatalf("Link wanted '%s' link got '%s'", want.link, link.link)
	}
	if want.adz != link.adz {
		t.Fatalf("Adz wanted '%s' link got '%s'", want.adz, link.adz)
	}
}
