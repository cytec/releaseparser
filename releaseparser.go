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
	audio      = `(?i)MP3|FLAC|DD[\s\.]?(2|5)\.?(1|0)|Dual[\- ]Audio|LiNE|DTS|AAC(?:\.?2\.0)?|AC3D?(?:\.5\.1)?`
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

	regexlist = map[string]string{
		"season":     season,
		"episode":    episode,
		"year":       year,
		"resolution": resolution,
		"source":     source,
		"codec":      codec,
		"audio":      audio,
		"group":      group,
		"region":     region,
		"doku":       doku,
		"extended":   extended,
		"uncut":      uncut,
		"hardcoded":  hardcoded,
		"proper":     proper,
		"subbed":     subbed,
		"repack":     repack,
		"is3d":       is3d,
		"widescreen": widescreen,
		"container":  container,
		"website":    website,
		"language":   language,
	}
)

// Release represents a scene release
type Release struct {
	Input       string // holds a copy of the input string
	Title       string // holds the release title without dots underscores and hypens
	Type        string // movie OR tvshow
	Season      int    // season number
	Episode     int    // episode number
	Year        int    // year
	Resolution  string // 720p, 1080p etc
	Source      string // the release source ex: BluRay, HDTV
	SourceGroup string // normalized Source Name for textmatching (ex: Blu-Ray, BluRay, BD, HDDVD => BLURAY)
	Codec       string // video codec ex: x264
	CodecGroup  string // normalized Codec Name for textmatching (ex: divx => DIVX)
	Audio       string // audio codec ex: FlAC, MP3, AC3
	AudioGroup  string // normalized Audio Name for textmatching (ex: DD5.1,DD => DD)
	Group       string // the name of the releasegroup
	Region      string // contains Region info ex: R9
	Container   string // the container file format ex: mkv
	Website     string // the release website if in the name ex: [ my.site.com ]
	Language    string // language of the release ex: german, Spanish
	Doku        bool   // true if release is dokumentation
	Extended    bool   // true if release is extended version
	Hardcoded   bool   // true if release is a Hardcoded release
	Subbed      bool   // true if release is subbed
	Proper      bool   // true if release is proper
	Repack      bool   // true if release is repack
	Is3D        bool   // true if release is in 3D
	Uncut       bool   // true if release is uncut version
	Widescreen  bool   // true if release is a widerscreen/letterbox release
	start       int
	end         int
	parts       map[string]string
}

// remove everything thats not a int from string
func parseInt(s string) int {
	re := regexp.MustCompile(`[^0-9]`)

	res := re.ReplaceAllLiteralString(s, "")

	x, _ := strconv.Atoi(res)

	return x
}

// set part info and calculate title start/stop position based on matches
func (r *Release) part(name string, match string, clean string) {
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

// gets the name of the group that matches, used for AudioGroup, SourceGroup, Resolution, CodecGroup
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

	for name, str := range regexlist {
		re := regexp.MustCompile(str)
		if loc := re.FindStringIndex(s); loc != nil {
			match := s[loc[0]:loc[1]]
			switch name {
			case "season":
				r.Season = parseInt(match)
			case "episode":
				//if make sure we dont match codec as episode
				if !regexp.MustCompile(codec).MatchString(match) {
					r.Episode = parseInt(match)
				}
			case "year":
				r.Year = parseInt(match)
			case "resolution":
				r.Resolution = strings.ToLower(getMatchedGroupName(re, match))
			case "source":
				r.Source = match
				r.SourceGroup = getMatchedGroupName(re, match)

			case "codec":
				r.Codec = match
				r.CodecGroup = getMatchedGroupName(re, match)
			case "audio":
				r.Audio = match
				r.AudioGroup = getMatchedGroupName(re, match)
			case "group":
				// if codec or source is in group skip it
				if regexp.MustCompile(codec).MatchString(match) || regexp.MustCompile(source).MatchString(match) {
					continue
				} else {
					r.Group = strings.Replace(match, "-", "", 1)
				}
			case "region":
				r.Region = match
			case "container":
				r.Container = match
			case "website":
				r.Website = match
			case "language":
				r.Language = match
			case "doku":
				r.Doku = true
			case "extended":
				r.Extended = true
			case "uncut":
				r.Uncut = true
			case "hardcoded":
				r.Hardcoded = true
			case "proper":
				r.Proper = true
			case "subbed":
				r.Subbed = true
			case "repack":
				r.Repack = true
			case "is3d":
				r.Is3D = true
			case "widescreen":
				r.Widescreen = true
			}
			//mark part as matched
			r.part(name, r.Input, match)
		}
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
