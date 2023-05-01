package versionchecker

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

var clientPool sync.Pool

func init() {
	clientPool.New = func() any {
		c := resty.New()
		c.JSONMarshal = sonic.Marshal
		c.JSONUnmarshal = sonic.Unmarshal
		return c
	}
}

func getClient() *resty.Client {
	return clientPool.Get().(*resty.Client)
}

func pushback(c *resty.Client) {
	clientPool.Put(c)
}

var (
	Owner = ""
	Repo  = ""
	Major = 0
	Minor = 0
	Patch = 0
)

// Version contains the Major, Minor, and Patch versions.
type Version struct {
	Major int
	Minor int
	Patch int
	Owner string
	Repo  string
}

// String return like "0.1.1"
func String() string {
	NowVersion := Version{
		Major: Major,
		Minor: Minor,
		Patch: Patch,
	}
	return NowVersion.String()
}

// Info return like "v0.1.1"
func Info() string {
	NowVersion := Version{
		Major: Major,
		Minor: Minor,
		Patch: Patch,
	}
	return NowVersion.Info()
}

// CheckUpgrade return the latest version
// and if it's a new version return true, else return false
// and if error occurs return an error, else return nil
func (v Version) CheckUpgrade() (latest Version, new bool, err error) {
	info, err := v.GetLatestVersionInfo()
	if err != nil {
		return Version{}, false, err
	}
	return info, info.Compare(v) > 0, nil
}

func (v Version) GetLatestVersionInfo() (Version, error) {
	return getLatestVersionInfo(v.Owner, v.Repo)
}

// String returns a string representation of the version
// e.g. "1.2.3"
// The string is generated using the Major, Minor, and Patch versions.
func (v Version) String() string {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
}

// Info returns version info in format "v1.2.3"
func (v Version) Info() string {
	return "v" + v.String()
}

// Compare compares this version with another version.
// It returns 0 if they are equal, 1 if this version
// is greater than v2 and -1 if this version is less than v2.
func (v Version) Compare(v2 Version) int {
	if v.Major > v2.Major {
		return 1
	}
	if v.Major < v2.Major {
		return -1
	}
	if v.Minor > v2.Minor {
		return 1
	}
	if v.Minor < v2.Minor {
		return -1
	}
	if v.Patch > v2.Patch {
		return 1
	}
	if v.Patch < v2.Patch {
		return -1
	}
	return 0
}

func getLatestVersionInfo(owner, repo string) (Version, error) {
	versionResponse := struct {
		// Url       string `json:"url"`
		// AssetsUrl string `json:"assets_url"`
		// UploadUrl string `json:"upload_url"`
		// HtmlUrl   string `json:"html_url"`
		// Id        int    `json:"id"`
		// Author    struct {
		// 	Login             string `json:"login"`
		// 	Id                int    `json:"id"`
		// 	NodeId            string `json:"node_id"`
		// 	AvatarUrl         string `json:"avatar_url"`
		// 	GravatarId        string `json:"gravatar_id"`
		// 	Url               string `json:"url"`
		// 	HtmlUrl           string `json:"html_url"`
		// 	FollowersUrl      string `json:"followers_url"`
		// 	FollowingUrl      string `json:"following_url"`
		// 	GistsUrl          string `json:"gists_url"`
		// 	StarredUrl        string `json:"starred_url"`
		// 	SubscriptionsUrl  string `json:"subscriptions_url"`
		// 	OrganizationsUrl  string `json:"organizations_url"`
		// 	ReposUrl          string `json:"repos_url"`
		// 	EventsUrl         string `json:"events_url"`
		// 	ReceivedEventsUrl string `json:"received_events_url"`
		// 	Type              string `json:"type"`
		// 	SiteAdmin         bool   `json:"site_admin"`
		// } `json:"author"`
		// NodeId          string    `json:"node_id"`
		TagName string `json:"tag_name"`
		// TargetCommitish string    `json:"target_commitish"`
		// Name            string    `json:"name"`
		// Draft           bool      `json:"draft"`
		// Prerelease      bool      `json:"prerelease"`
		// CreatedAt       time.Time `json:"created_at"`
		// PublishedAt     time.Time `json:"published_at"`
		// Assets          []struct {
		// 	Url      string `json:"url"`
		// 	Id       int    `json:"id"`
		// 	NodeId   string `json:"node_id"`
		// 	Name     string `json:"name"`
		// 	Label    any    `json:"label"`
		// 	Uploader struct {
		// 		Login             string `json:"login"`
		// 		Id                int    `json:"id"`
		// 		NodeId            string `json:"node_id"`
		// 		AvatarUrl         string `json:"avatar_url"`
		// 		GravatarId        string `json:"gravatar_id"`
		// 		Url               string `json:"url"`
		// 		HtmlUrl           string `json:"html_url"`
		// 		FollowersUrl      string `json:"followers_url"`
		// 		FollowingUrl      string `json:"following_url"`
		// 		GistsUrl          string `json:"gists_url"`
		// 		StarredUrl        string `json:"starred_url"`
		// 		SubscriptionsUrl  string `json:"subscriptions_url"`
		// 		OrganizationsUrl  string `json:"organizations_url"`
		// 		ReposUrl          string `json:"repos_url"`
		// 		EventsUrl         string `json:"events_url"`
		// 		ReceivedEventsUrl string `json:"received_events_url"`
		// 		Type              string `json:"type"`
		// 		SiteAdmin         bool   `json:"site_admin"`
		// 	} `json:"uploader"`
		// 	ContentType        string    `json:"content_type"`
		// 	State              string    `json:"state"`
		// 	Size               int       `json:"size"`
		// 	DownloadCount      int       `json:"download_count"`
		// 	CreatedAt          time.Time `json:"created_at"`
		// 	UpdatedAt          time.Time `json:"updated_at"`
		// 	BrowserDownloadUrl string    `json:"browser_download_url"`
		// } `json:"assets"`
		// TarballUrl string `json:"tarball_url"`
		// ZipballUrl string `json:"zipball_url"`
		// Body       string `json:"body"`
	}{}

	// https://api.github.com/repos/$owner$/$repo$/releases/latest
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	c := getClient()
	defer pushback(c)
	_, err := c.R().SetResult(&versionResponse).Get(url)
	latest := Version{}
	if err != nil {
		return latest, err
	}

	versionStr := versionResponse.TagName
	_, err = fmt.Sscanf(versionStr, "v%d.%d.%d", &latest.Major, &latest.Minor, &latest.Patch)

	if err != nil {
		return latest, err
	}

	return latest, nil
}

// GetLatestVersionInfo get the latest version info from GitHub
// "https://api.github.com/repos/Equationzhao/GodDns/releases/latest
func GetLatestVersionInfo() (Version, error) {
	return getLatestVersionInfo(Owner, Repo)
}

func Set(major, minor, patch int) {
	Major, Minor, Patch = major, minor, patch
}

// CheckUpgrade return the latest version
// and if it's a new version return true, else return false
// and if error occurs return an error, else return nil
func CheckUpgrade() (latest Version, new bool, err error) {
	info, err := GetLatestVersionInfo()
	if err != nil {
		return Version{}, false, err
	}
	NowVersion := Version{
		Major: Major,
		Minor: Minor,
		Patch: Patch,
	}
	return info, info.Compare(NowVersion) > 0, nil
}
