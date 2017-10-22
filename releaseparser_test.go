package releaseparser_test

import (
	"testing"

	"github.com/cytec/releaseparser"
)

func TestParse(t *testing.T) {
	test := map[string]*releaseparser.Release{
		"Winx.Club.S06E16.Die.Zombie-Invasion.GERMAN.DUBBED.DL.720p.WEB-DL.h264-pbw": &releaseparser.Release{
			Type:       "tvshow",
			Title:      "Winx Club",
			Season:     6,
			Episode:    16,
			Language:   "GERMAN",
			Source:     "WEB-DL",
			Codec:      "h264",
			Group:      "pbw",
			Resolution: "720p",
		},
		"Scouts.vs.Zombies.Handbuch.zur.Zombie.Apokalypse.2015.German.AC3.DL.1080p.BluRay.x264-EXQUiSiTE": &releaseparser.Release{
			Type:       "movie",
			Title:      "Scouts vs Zombies Handbuch zur Zombie Apokalypse",
			Year:       2015,
			Language:   "German",
			Source:     "BluRay",
			Codec:      "x264",
			Group:      "EXQUiSiTE",
			Resolution: "1080p",
			Audio:      "AC3",
		},
		"Zombie.Bloody.Demons.UNCUT.GERMAN.1987.DL.1080p.BluRay.x264-GOREHOUNDS": &releaseparser.Release{
			Type:       "movie",
			Title:      "Zombie Bloody Demons",
			Year:       1987,
			Language:   "GERMAN",
			Source:     "BluRay",
			Codec:      "x264",
			Group:      "GOREHOUNDS",
			Resolution: "1080p",
			Uncut:      true,
		},
		"iZombie.S02E10.Zombie.High.German.DD51.Dubbed.DL.720p.BD.x264-TVS": &releaseparser.Release{
			Type:       "tvshow",
			Title:      "iZombie",
			Season:     2,
			Episode:    10,
			Language:   "German",
			Source:     "BD",
			Codec:      "x264",
			Group:      "TVS",
			Audio:      "DD51",
			Resolution: "720p",
		},
		"Brave.2012.R5.DVDRip.XViD.LiNE-UNiQUE": &releaseparser.Release{
			Type:   "movie",
			Title:  "Brave",
			Year:   2012,
			Source: "DVDRip",
			Codec:  "XViD",
			Group:  "UNiQUE",
			Region: "R5",
			Audio:  "LiNE",
		},
		"Brave.2012.German.Subbed.DVDRip.XViD.LiNE-UNiQUE": &releaseparser.Release{
			Type:     "movie",
			Title:    "Brave",
			Year:     2012,
			Source:   "DVDRip",
			Codec:    "XViD",
			Group:    "UNiQUE",
			Language: "German",
			Subbed:   true,
			Audio:    "LiNE",
		},
		"Ant-Man.2015.3D.1080p.BRRip.Half-SBS.x264.AAC-m2g": &releaseparser.Release{
			Type:       "movie",
			Title:      "Ant Man",
			Year:       2015,
			Source:     "BRRip",
			Codec:      "x264",
			Group:      "m2g",
			Audio:      "AAC",
			Resolution: "1080p",
			Is3D:       true,
		},
		"Annabelle.2014.1080p.PROPER.HC.WEBRip.x264.AAC.2.0-RARBG": &releaseparser.Release{
			Type:       "movie",
			Title:      "Annabelle",
			Year:       2014,
			Source:     "WEBRip",
			Codec:      "x264",
			Group:      "RARBG",
			Audio:      "AAC.2.0",
			Resolution: "1080p",
			Proper:     true,
			Hardcoded:  true,
		},
		"The.Boss.2016.UNCUT.720p.BRRip.x264.AAC-ETRG": &releaseparser.Release{
			Type:       "movie",
			Title:      "The Boss",
			Year:       2016,
			Source:     "BRRip",
			Codec:      "x264",
			Group:      "ETRG",
			Audio:      "AAC",
			Resolution: "720p",
			Uncut:      true,
		},
		"Hercules.2014.EXTENDED.1080p.WEB-DL.DD5.1.H264-RARBG": &releaseparser.Release{
			Type:       "movie",
			Title:      "Hercules",
			Year:       2014,
			Source:     "WEB-DL",
			Codec:      "H264",
			Group:      "RARBG",
			Audio:      "DD5.1",
			Resolution: "1080p",
			Extended:   true,
		},
		"1-2-3.Istanbul.S01E04.GERMAN.DOKU.WS.dTV.XviD-GEO": &releaseparser.Release{
			Type:       "tvshow",
			Title:      "1 2 3 Istanbul",
			Season:     1,
			Episode:    4,
			Language:   "GERMAN",
			Source:     "dTV",
			Codec:      "XviD",
			Group:      "GEO",
			Widescreen: true,
			Doku:       true,
		},
		"[ www.Speed.cd ] -Sons.of.Anarchy.S07E07.720p.HDTV.X264-DIMENSION": &releaseparser.Release{
			Type:       "tvshow",
			Title:      "Sons of Anarchy",
			Season:     7,
			Episode:    7,
			Resolution: "720p",
			Source:     "HDTV",
			Codec:      "X264",
			Group:      "DIMENSION",
			Website:    "[ www.Speed.cd ]",
		},
		"Two and a Half Men S12E01 HDTV x264 REPACK-LOL [eztv]": &releaseparser.Release{
			Type:    "tvshow",
			Title:   "Two and a Half Men",
			Season:  12,
			Episode: 1,
			Source:  "HDTV",
			Codec:   "x264",
			Group:   "LOL [eztv]",
			Repack:  true,
		},
		"Eliza Graves (2014) Dual Audio WEB-DL 720p MKV x264": &releaseparser.Release{
			Type:       "movie",
			Title:      "Eliza Graves",
			Year:       2014,
			Audio:      "Dual Audio",
			Source:     "WEB-DL",
			Codec:      "x264",
			Resolution: "720p",
			Container:  "MKV",
		},
		"The Shaukeens 2014 Hindi (1CD) DvDScr x264 AAC...Hon3y": &releaseparser.Release{
			Type:     "movie",
			Title:    "The Shaukeens",
			Year:     2014,
			Language: "Hindi",
			Source:   "DvDScr",
			Codec:    "x264",
			Audio:    "AAC",
		},
		"Mr Robot S02E11 German DD 51 Synced DL 1080p AmazonHD x264-TVS": &releaseparser.Release{
			Type:       "tvshow",
			Title:      "Mr Robot",
			Season:     2,
			Episode:    11,
			Group:      "TVS",
			Resolution: "1080p",
			Language:   "German",
			Source:     "AmazonHD",
			Codec:      "x264",
			Audio:      "DD 51",
		},
		"31.A.Rob.Zombie.Film.3D.UNCUT.2016.German.DL.1080p.BluRay.x264-ETM": &releaseparser.Release{
			Type:       "movie",
			Title:      "31 A Rob Zombie Film",
			Year:       2016,
			Group:      "ETM",
			Resolution: "1080p",
			Language:   "German",
			Source:     "BluRay",
			Codec:      "x264",
			Is3D:       true,
			Uncut:      true,
		},
		"Zombie Shark The Swimming Dead French 2015 AC3 BDRiP x264-XF": &releaseparser.Release{
			Type:     "movie",
			Title:    "Zombie Shark The Swimming Dead",
			Year:     2015,
			Group:    "XF",
			Language: "French",
			Source:   "BDRiP",
			Codec:    "x264",
			Audio:    "AC3",
		},
		"Dracula.Untold.TS.XViD.AC3.MrSeeN-SiMPLE": &releaseparser.Release{
			Type:   "movie",
			Title:  "Dracula Untold",
			Group:  "SiMPLE",
			Source: "TS",
			Codec:  "XViD",
			Audio:  "AC3",
		},
		"Mr.Robot.S01.PROPER.VOSTFR.720p.WEB-DL.DD5.1.H264-ARK01": &releaseparser.Release{
			Type:       "tvshow",
			Title:      "Mr Robot",
			Season:     1,
			Group:      "ARK01",
			Language:   "VOSTFR",
			Source:     "WEB-DL",
			Codec:      "H264",
			Audio:      "DD5.1",
			Resolution: "720p",
			Proper:     true,
		},
	}

	for title, want := range test {
		parsed := releaseparser.Parse(title)
		t.Logf("Running tests for %s\n", title)

		if want.Title != parsed.Title {
			t.Errorf("Title failed, got: %s, want: %s", parsed.Title, want.Title)
		}
		if want.Type != parsed.Type {
			t.Errorf("Type failed, got: %s, want: %s", parsed.Type, want.Type)
		}
		if want.Season != parsed.Season {
			t.Errorf("Season failed, got: %d, want: %d", parsed.Season, want.Season)
		}
		if want.Episode != parsed.Episode {
			t.Errorf("Episode failed, got: %d, want: %d", parsed.Episode, want.Episode)
		}
		if want.Year != parsed.Year {
			t.Errorf("Year failed, got: %d, want: %d", parsed.Year, want.Year)
		}
		if want.Resolution != parsed.Resolution {
			t.Errorf("Resolution failed, got: %s, want: %s", parsed.Resolution, want.Resolution)
		}
		if want.Source != parsed.Source {
			t.Errorf("Source failed, got: %s, want: %s", parsed.Source, want.Source)
		}
		if want.Codec != parsed.Codec {
			t.Errorf("Codec failed, got: %s, want: %s", parsed.Codec, want.Codec)
		}
		if want.Audio != parsed.Audio {
			t.Errorf("Audio failed, got: %s, want: %s", parsed.Audio, want.Audio)
		}
		if want.Group != parsed.Group {
			t.Errorf("Group failed, got: %s, want: %s", parsed.Group, want.Group)
		}
		if want.Region != parsed.Region {
			t.Errorf("Region failed, got: %s, want: %s", parsed.Region, want.Region)
		}
		if want.Doku != parsed.Doku {
			t.Errorf("Doku failed, got: %t, want: %t", parsed.Doku, want.Doku)
		}
		if want.Extended != parsed.Extended {
			t.Errorf("Extended failed, got: %t, want: %t", parsed.Extended, want.Extended)
		}
		if want.Hardcoded != parsed.Hardcoded {
			t.Errorf("Hardcoded failed, got: %t, want: %t", parsed.Hardcoded, want.Hardcoded)
		}
		if want.Subbed != parsed.Subbed {
			t.Errorf("Subbed failed, got: %t, want: %t", parsed.Subbed, want.Subbed)
		}
		if want.Proper != parsed.Proper {
			t.Errorf("Proper failed, got: %t, want: %t", parsed.Proper, want.Proper)
		}
		if want.Repack != parsed.Repack {
			t.Errorf("Repack failed, got: %t, want: %t", parsed.Repack, want.Repack)
		}
		if want.Is3D != parsed.Is3D {
			t.Errorf("Is3D failed, got: %t, want: %t", parsed.Is3D, want.Is3D)
		}
		if want.Uncut != parsed.Uncut {
			t.Errorf("Uncut failed, got: %t, want: %t", parsed.Uncut, want.Uncut)
		}
		if want.Container != parsed.Container {
			t.Errorf("Container failed, got: %s, want: %s", parsed.Container, want.Container)
		}
		if want.Website != parsed.Website {
			t.Errorf("Website failed, got: %s, want: %s", parsed.Website, want.Website)
		}
		if want.Widescreen != parsed.Widescreen {
			t.Errorf("Widescreen failed, got: %t, want: %t", parsed.Widescreen, want.Widescreen)
		}
		if want.Language != parsed.Language {
			t.Errorf("Language failed, got: %s, want: %s", parsed.Language, want.Language)
		}
	}

}
