package releaseparser

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	season     = `(?i)(s[0-9]{2}-s[0-9]{2}|s([0-9]{1,2})[eEx])|([Ss]?([0-9]{1,2}))[Eex]|([Ss]([0-9]{1,2}))`
	episode    = `([Eex]([0-9]{2,4}-?[Eex]?[0-9]{2,4}))|([Eex]([0-9]{2,4}(?:[abc])?)(?:[^0-9]|$))|\b((?:[Eex]p?\.?)([0-9]{2,4}(:?-?(?:[Eex]?p?)[0-9]{2,4})?)|[Ee]pisode\s?([0-9]{1,4}))\b`
	year       = `([\[\(]?((?:19[0-9]|20[01])[0-9])[\]\)]?)`
	resolution = `(?P<480p>480p|640x480|848x480)|(?P<576p>576p)|(?P<720p>720p|1280x720)|(?P<1080p>1080p|1920x1080)|(?P<2160p>2160p)`
	source     = `(?i)\b(?:(?P<bdrip>BDRip)|(?P<brrip>BRRip)|(?P<bluray>BluRay|Blu-Ray|HDDVD|BD)|(?P<webdl>WEB[-_. ]DL|HDRIP|WEBDL|FUNi-DL|WebRip|Web-Rip|AmazonHD|NetflixHD|iTunesHD|WebHD|[. ]?WEB[. ](?:[xh]26[45]|DD5[. ]1)|\\d+0p[. ]WEB[. ])|(?P<hdtv>HDTV)|(?P<scr>SCR|SCREENER|DVDSCR|DVDSCREENER)|(?P<dvd>DVDRip|DVD[^-R]|NTSC|PAL|xvidvd)|(?P<dvdr>DVD-R|DVDR|DVD[0-9])|(?P<dsr>WS[-_. ]DSR|DSR)|(?P<ts>TS|TELESYNC|HD-TS|HDTS|PDVD\b)|(?P<tc>TC|TELECINE|HD-TC|HDTC)|(?P<cam>CAMRIP|CAM|HDCAM|HD-CAM)|(?P<wp>WORKPRINT|WP)|(?P<pdtv>PDTV)|(?P<sdtv>SDTV)|(?P<tvrip>(HD)?TVRip|[ad]TV))\b`
	codec      = `(?i)(?P<x264>x264)|(?P<h264>h264)|(?P<h265>[xh]265|hevc)|(?P<xvidhd>XvidHD)|(?P<xvid>X-?vid)|(?P<divx>divx|mpeg[0-9])(?P<vp>vp(?:8|9))`
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
	container  = `(?i)\b\.?(MKV|AVI|MP4|mkv|avi|mp4|m4v)\b`
	website    = `^(\[ ?([^\]]+?) ?\])`
	sbs        = `\b(?i)(?:Half-)?SBS\b`
	size       = `(\d+(?:\.\d+)?(?:GB|MB))`
	language   = `(?i)\b(?:TRUE)?FR(?:ENCH)?\b|\bDE(?:UTSCH)?\b|\bGERMAN\b|\bEN(?:G(?:LISH)?)?\b|\bVOST(?:(F(?:R)?)|A)?\b|\bMULTI(?:Lang|Truefrench|\-VF2)?\b|\bSUBFRENCH\b|\bHindi\b`
	password   = `(?i){{(?:[^{}]+)}}`
	console    = `\b(XBOX|XBOX360|Wii|WiiU|PSP|PS4|NSW|PS3|NDS)\b`

	regexlist = map[string]string{
		"season":     season,
		"episode":    episode,
		"year":       year,
		"resolution": resolution,
		"source":     source,
		"codec":      codec,
		"audio":      audio,
		"language":   language,
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
		"sbs":        sbs,
		"size":       size,
		"group":      group,
		"console":    console,
	}

	releaseTypePC      = "pc"
	releaseTypeConsole = "console"
	releaseTypeMovie   = "movie"
	releaseTypeTV      = "tvshow"

	groupTypeMap = map[string]string{
		"CODEX":      releaseTypePC,
		"DARKSiDERS": releaseTypePC,
		"PLAZA":      releaseTypePC,
		"RAZOR":      releaseTypePC,
		"SiMPLEX":    releaseTypePC,
		"Razor1911":  releaseTypePC,
		"HOODLUM":    releaseTypePC,
		"SKIDROW":    releaseTypePC,
		"ALiAS":      releaseTypePC,
	}
)

// Release represents a scene release
type Release struct {
	Input       string `json:"input,omitempty"`        // holds a copy of the input string
	Title       string `json:"title,omitempty"`        // holds the release title without dots underscores and hypens
	Type        string `json:"type,omitempty"`         // movie OR tvshow
	Season      int    `json:"season,omitempty"`       // season number
	SeasonEnd   int    `json:"season_end,omitempty"`   // 0 or end season for multi season releases
	Episode     int    `json:"episode,omitempty"`      // episode number
	EpisodeEnd  int    `json:"episode_end,omitempty"`  // 0 er end episode number for multi episode releases
	Year        int    `json:"year,omitempty"`         // year
	Resolution  string `json:"resolution,omitempty"`   // 720p, 1080p etc
	Source      string `json:"source,omitempty"`       // the release source ex: BluRay, HDTV
	SourceGroup string `json:"source_group,omitempty"` // normalized Source Name for textmatching (ex: Blu-Ray, BluRay, BD, HDDVD => BLURAY)
	Codec       string `json:"codec,omitempty"`        // video codec ex: x264
	CodecGroup  string `json:"codec_group,omitempty"`  // normalized Codec Name for textmatching (ex: divx => DIVX)
	Audio       string `json:"audio,omitempty"`        // audio codec ex: FlAC, MP3, AC3
	AudioGroup  string `json:"audio_group,omitempty"`  // normalized Audio Name for textmatching (ex: DD5.1,DD => DD)
	Group       string `json:"group,omitempty"`        // the name of the releasegroup
	Region      string `json:"region,omitempty"`       // contains Region info ex: R9
	Container   string `json:"container,omitempty"`    // the container file format ex: mkv
	Website     string `json:"website,omitempty"`      // the release website if in the name ex: [ my.site.com ]
	Language    string `json:"language,omitempty"`     // language of the release ex: german, Spanish
	Password    string `json:"password,omitempty"`     // if there is a password found it will be here
	SBS         string `json:"sbs,omitempty"`          // Full-SBS or SBS
	Size        string `json:"size,omitempty"`         // size if present in title
	Doku        bool   `json:"doku,omitempty"`         // true if release is dokumentation
	Extended    bool   `json:"extended,omitempty"`     // true if release is extended version
	Hardcoded   bool   `json:"hardcoded,omitempty"`    // true if release is a Hardcoded release
	Subbed      bool   `json:"subbed,omitempty"`       // true if release is subbed
	Proper      bool   `json:"proper,omitempty"`       // true if release is proper
	Repack      bool   `json:"repack,omitempty"`       // true if release is repack
	Is3D        bool   `json:"is_3d,omitempty"`        // true if release is in 3D
	Uncut       bool   `json:"uncut,omitempty"`        // true if release is uncut version
	Widescreen  bool   `json:"widescreen,omitempty"`   // true if release is a widerscreen/letterbox release
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

	//cut password from string because this might mess up correct detection of other infos...
	pwregex := regexp.MustCompile(password)
	if pwloc := pwregex.FindStringIndex(s); pwloc != nil {
		pwmatch := s[pwloc[0]:pwloc[1]]
		r.Password = pwmatch[2 : len(pwmatch)-2]

		s = pwregex.ReplaceAllString(s, "")
	}

	for name, str := range regexlist {
		re := regexp.MustCompile(str)
		if loc := re.FindStringIndex(s); loc != nil {
			match := s[loc[0]:loc[1]]
			switch name {
			case "season":
				seasons := strings.Split(match, "-")
				if len(seasons) > 1 {
					r.SeasonEnd = parseInt(seasons[1])
				}
				r.Season = parseInt(seasons[0])
			case "episode":
				//if make sure we dont match codec as episode
				if !regexp.MustCompile(codec).MatchString(match) {
					//remove episode becuase it gets split otherwise
					clean := regexp.MustCompile("(?i)episode").ReplaceAllString(match, "")
					//split multiep strings
					tmp := regexp.MustCompile(`(?i)(\.|-|ep|e|x)`).Split(clean, -1)
					episodes := []string{}
					for _, v := range tmp {
						if v != "" {
							episodes = append(episodes, v)
						}
					}
					r.Episode = parseInt(episodes[0])
					if len(episodes) > 1 {
						r.EpisodeEnd = parseInt(episodes[1])
					}
				}
			case "year":
				r.Year = parseInt(match)
			case "resolution":
				r.Resolution = strings.ToLower(getMatchedGroupName(re, match))
			case "source":
				r.Source = strings.Trim(match, " ")
				r.SourceGroup = getMatchedGroupName(re, match)
			case "codec":
				r.Codec = match
				r.CodecGroup = getMatchedGroupName(re, match)
			case "audio":
				r.Audio = match
				r.AudioGroup = getMatchedGroupName(re, match)
			case "group":
				// if codec or source is in group skip it
				if regexp.MustCompile(codec).MatchString(match) || regexp.MustCompile(source).MatchString(match) || regexp.MustCompile(language).MatchString(match) {
					continue
				} else {
					r.Group = strings.Replace(match, "-", "", 1)
					r.Group = regexp.MustCompile(container).ReplaceAllString(r.Group, "")
				}
			case "region":
				r.Region = match
			case "console":
				r.Type = releaseTypeConsole
			case "container":
				r.Container = strings.Replace(match, ".", "", -1)
			case "website":
				r.Website = match
			case "language":
				r.Language = match
			case "sbs":
				r.SBS = match
			case "size":
				r.Size = match
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

	if r.end != 0 && r.end <= len(r.Input) && r.start < r.end {
		r.Title = cleanTitle(r.Input[r.start:r.end])
	}

	if r.Season > 0 || r.Episode > 0 && r.Episode != parseInt(r.Codec) {
		r.Type = releaseTypeTV
	} else {
		for g, t := range groupTypeMap {
			if r.Group == g {
				r.Type = t
			}
		}

		if r.Type == "" {
			r.Type = releaseTypeMovie
			r.Episode = 0
		}
	}

	return &r
}
