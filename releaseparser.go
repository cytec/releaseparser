package releaseparser

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	season     = `([Ss]?([0-9]{1,2}))[Eex]|([Ss]([0-9]{1,2}))`
	episode    = `([Eex]([0-9]{2,4}(?:[abc])?)(?:[^0-9]|$))`
	year       = `([\[\(]?((?:19[0-9]|20[01])[0-9])[\]\)]?)`
	resolution = `(([0-9]{3,4}(?:p|i)))[^M]`
	source     = `(?i)hdtv|bluray|(?:b[dr]|dvd|hd|tv|webhd|web)rip|web-?(?:dl|rip|hd)|dvd(scr)?|bd|ituneshd|itunes|\b(?:hd)?ts\b|amazon(?:hd)?|netflix(?:hd)?|[ad]TV`
	codec      = `(?i)dvix|mpeg[0-9]|divx|xvid|(?:x|h)[-\. ]?26(?:4|5)|avc|hevc|vp(?:8|9)`
	audio      = `(?i)MP3|FLAC|DD[\s\.]?5\.?1|Dual[\- ]Audio|LiNE|DTS|AAC(?:\.?2\.0)?|AC3(?:\.5\.1)?`
	group      = `(?:- ?([^-]+))$`
	region     = `R[0-9]`
	doku       = `(?i)\bDOKU\b`
	extended   = `EXTENDED`
	uncut      = `UNCUT`
	hardcoded  = `HC`
	proper     = `PROPER`
	subbed     = `(?i)subbed|ger[.-]?sub(s|ed)?|nlsub|eng-sub`
	repack     = `REPACK`
	is3d       = `(?i)\b3d\b`
	widescreen = `WS`
	container  = `MKV|AVI|MP4|mkv|avi|mp4`
	website    = `^(\[ ?([^\]]+?) ?\])`
	language   = `(?i)\b(?:TRUE)?FR(?:ENCH)?\b|\bDE(?:UTSCH)?\b|\bGERMAN\b|\bEN(?:G(?:LISH)?)?\b|\bVOST(?:(F(?:R)?)|A)?\b|\bMULTI(?:Lang|Truefrench|\-VF2)?\b|\bSUBFRENCH\b|\bHindi\b`
)

// Release represents a scene release
type Release struct {
	Input      string
	Title      string
	Type       string
	Season     int
	Episode    int
	Year       int
	Resolution string
	Source     string
	Codec      string
	Audio      string
	Group      string
	Region     string
	Doku       bool
	Extended   bool
	Hardcoded  bool
	Subbed     bool
	Proper     bool
	Repack     bool
	Is3D       bool
	Uncut      bool
	Container  string
	Website    string
	Widescreen bool
	Language   string
	start      int
	end        int
	parts      map[string]string
}

func parseInt(s string) int {
	re := regexp.MustCompile(`[^0-9]`)

	res := re.ReplaceAllLiteralString(s, "")

	x, _ := strconv.Atoi(res)

	return x
}

func (r *Release) part(name string, match string, clean string) {
	// self.parts[name] = clean
	r.parts[name] = clean

	if match != "" {
		index := strings.Index(r.Input, clean)
		if index <= 0 {
			r.start = len(clean)
		} else if r.end == 0 || index < r.end {
			r.end = index
		}
	}
}

func cleanTitle(name string) string {
	name = strings.Replace(name, ".", " ", -1)
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Replace(name, "-", " ", -1)
	name = strings.Trim(name, " ")
	return name
}

// Parse parses the given release name
func Parse(s string) *Release {
	r := Release{Input: s, parts: make(map[string]string)}

	re := regexp.MustCompile(season)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Season = parseInt(s[loc[0]:loc[1]])
		r.part("season", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(episode)
	if loc := re.FindStringIndex(s); loc != nil {
		// make sure its not x264 codec
		if !regexp.MustCompile(codec).MatchString(s[loc[0]:loc[1]]) {
			r.Episode = parseInt(s[loc[0]:loc[1]])
			r.part("episode", r.Input, s[loc[0]:loc[1]])
		}
	}
	re = regexp.MustCompile(year)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Year = parseInt(s[loc[0]:loc[1]])
		r.part("year", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(resolution)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Resolution = s[loc[0] : loc[1]-1]
		r.part("resolution", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(source)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Source = s[loc[0]:loc[1]]
		r.part("source", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(language)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Language = s[loc[0]:loc[1]]
		r.part("language", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(codec)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Codec = s[loc[0]:loc[1]]
		r.part("codec", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(audio)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Audio = s[loc[0]:loc[1]]
		r.part("audio", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(website)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Website = s[loc[0]:loc[1]]
		r.part("website", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(group)
	if loc := re.FindStringIndex(s); loc != nil {
		groupstring := s[loc[0]:loc[1]]
		//if codec or source is in group skip it
		if regexp.MustCompile(codec).MatchString(groupstring) || regexp.MustCompile(source).MatchString(s[loc[0]:loc[1]]) {

		} else {
			r.Group = strings.Replace(groupstring, "-", "", 1)
			r.part("group", r.Input, s[loc[0]:loc[1]])
		}
	}
	re = regexp.MustCompile(region)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Region = s[loc[0]:loc[1]]
		r.part("region", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(doku)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Doku = true
		r.part("doku", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(uncut)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Uncut = true
		r.part("uncut", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(proper)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Proper = true
		r.part("proper", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(extended)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Extended = true
		r.part("extended", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(hardcoded)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Hardcoded = true
		r.part("hardcoded", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(subbed)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Subbed = true
		r.part("subbed", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(is3d)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Is3D = true
		r.part("is3d", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(repack)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Repack = true
		r.part("repack", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(widescreen)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Widescreen = true
		r.part("widescreen", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(uncut)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Uncut = true
		r.part("uncut", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(container)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Container = s[loc[0]:loc[1]]
		r.part("container", r.Input, s[loc[0]:loc[1]])
	}

	if r.end != 0 {
		r.Title = cleanTitle(r.Input[r.start:r.end])
	}

	if r.Season > 0 || r.Episode > 0 && r.Episode != parseInt(r.Codec) {
		r.Type = "tvshow"
	} else {
		r.Type = "movie"
		r.Episode = 0
	}

	return &r
}
