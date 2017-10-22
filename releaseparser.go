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
	resolution = `(?P<480p>480p|640x480|848x480)|(?P<576p>576p)|(?P<720p>720p|1280x720)|(?P<1080p>1080p|1920x1080)|(?P<2160p>2160p)`
	source     = `(?i)(\b)(?P<bdrip>BDRip)|(?P<brrip>BRRip)|(?P<bluray>BluRay|Blu-Ray|HDDVD|BD)|(?P<webdl>WEB[-_. ]DL|HDRIP|WEBDL|FUNi-DL|WebRip|Web-Rip|AmazonHD|NetflixHD|iTunesHD|WebHD|[. ]WEB[. ](?:[xh]26[45]|DD5[. ]1)|\\d+0p[. ]WEB[. ])|(?P<hdtv>HDTV)|(?P<scr>SCR|SCREENER|DVDSCR|DVDSCREENER)|(?P<dvdr>DVD-R|DVDRip|DVDR?|NTSC|PAL|xvidvd)|(?P<dsr>WS[-_. ]DSR|DSR)|(?P<ts>\bTS|TELESYNC|HD-TS|HDTS|PDVD\b)|(?P<tc>TC|TELECINE|HD-TC|HDTC)|(?P<cam>CAMRIP|CAM|HDCAM|HD-CAM)|(?P<wp>WORKPRINT|WP)|(?P<pdtv>PDTV)|(?P<sdtv>SDTV)|(?P<tvrip>TVRip|[ad]TV)(\b)`
	// codec      = `(?i)dvix|mpeg[0-9]|divx|xvid(?:hd)?|(?:x|h)[-\. ]?26(?:4|5)|avc|hevc|vp(?:8|9)`
	codec      = `(?i)(?P<x264>x264)|(?P<h264>h264)|(?P<h265>h265|hevc)|(?P<xvidhd>XvidHD)|(?P<xvid>X-?vid)|(?P<divx>divx|mpeg[0-9])(?P<vp>vp(?:8|9))`
	audio      = `(?i)MP3|FLAC|DD[\s\.]?5\.?1|Dual[\- ]Audio|LiNE|DTS|AAC(?:\.?2\.0)?|AC3D?(?:\.5\.1)?`
	group      = `(?:- ?([^-]+))$`
	region     = `R[0-9]{1}`
	doku       = `(?i)\bDOKU\b`
	extended   = `(?i)\bEXTENDED\b`
	uncut      = `(?i)\bUNCUT\b`
	hardcoded  = `(?i)\bHC\b`
	proper     = `(?i)\bPROPER\b`
	subbed     = `(?i)subbed|ger[.-]?sub(s|ed)?|nlsub|eng-sub`
	repack     = `(?i)\bREPACK\b`
	is3d       = `(?i)\b3d\b`
	widescreen = `(?i)\bWS\b`
	container  = `(?i)\bMKV|AVI|MP4|mkv|avi|mp4\b`
	website    = `^(\[ ?([^\]]+?) ?\])`
	language   = `(?i)\b(?:TRUE)?FR(?:ENCH)?\b|\bDE(?:UTSCH)?\b|\bGERMAN\b|\bEN(?:G(?:LISH)?)?\b|\bVOST(?:(F(?:R)?)|A)?\b|\bMULTI(?:Lang|Truefrench|\-VF2)?\b|\bSUBFRENCH\b|\bHindi\b`
)

// Release represents a scene release
type Release struct {
	Input       string
	Title       string
	Type        string
	Season      int
	Episode     int
	Year        int
	Resolution  string
	Source      string
	SourceGroup string
	Codec       string
	CodecGroup  string
	Audio       string
	AudioGroup  string
	Group       string
	Region      string
	Container   string
	Website     string
	Language    string
	Doku        bool
	Extended    bool
	Hardcoded   bool
	Subbed      bool
	Proper      bool
	Repack      bool
	Is3D        bool
	Uncut       bool
	Widescreen  bool
	start       int
	end         int
	parts       map[string]string
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

func getMatchedGroupName(re *regexp.Regexp, s string) string {
	match := re.FindStringSubmatch(s)

	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			if match[i] != "" {
				return strings.ToUpper(name)
			}
		}
	}
	return ""
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
		r.Resolution = strings.ToLower(getMatchedGroupName(re, s[loc[0]:loc[1]]))
		r.part("resolution", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(source)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Source = s[loc[0]:loc[1]]
		r.SourceGroup = getMatchedGroupName(re, s[loc[0]:loc[1]])
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
		r.CodecGroup = getMatchedGroupName(re, s[loc[0]:loc[1]])
		r.part("codec", r.Input, s[loc[0]:loc[1]])
	}
	re = regexp.MustCompile(audio)
	if loc := re.FindStringIndex(s); loc != nil {
		r.Audio = s[loc[0]:loc[1]]
		r.AudioGroup = getMatchedGroupName(re, s[loc[0]:loc[1]])
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
