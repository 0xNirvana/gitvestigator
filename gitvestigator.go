package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

type Args struct {
	repoLink string
	token    string
}

type RepoData struct {
	ID       int    `json:"id,omitempty"`
	NodeID   string `json:"node_id,omitempty"`
	Name     string `json:"name,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Private  bool   `json:"private,omitempty"`
	Owner    struct {
		Login             string `json:"login,omitempty"`
		ID                int    `json:"id,omitempty"`
		NodeID            string `json:"node_id,omitempty"`
		AvatarURL         string `json:"avatar_url,omitempty"`
		GravatarID        string `json:"gravatar_id,omitempty"`
		URL               string `json:"url,omitempty"`
		HTMLURL           string `json:"html_url,omitempty"`
		FollowersURL      string `json:"followers_url,omitempty"`
		FollowingURL      string `json:"following_url,omitempty"`
		GistsURL          string `json:"gists_url,omitempty"`
		StarredURL        string `json:"starred_url,omitempty"`
		SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
		OrganizationsURL  string `json:"organizations_url,omitempty"`
		ReposURL          string `json:"repos_url,omitempty"`
		EventsURL         string `json:"events_url,omitempty"`
		ReceivedEventsURL string `json:"received_events_url,omitempty"`
		Type              string `json:"type,omitempty"`
		SiteAdmin         bool   `json:"site_admin,omitempty"`
	} `json:"owner,omitempty"`
	HTMLURL                  string    `json:"html_url,omitempty"`
	Description              string    `json:"description,omitempty"`
	Fork                     bool      `json:"fork,omitempty"`
	URL                      string    `json:"url,omitempty"`
	ForksURL                 string    `json:"forks_url,omitempty"`
	KeysURL                  string    `json:"keys_url,omitempty"`
	CollaboratorsURL         string    `json:"collaborators_url,omitempty"`
	TeamsURL                 string    `json:"teams_url,omitempty"`
	HooksURL                 string    `json:"hooks_url,omitempty"`
	IssueEventsURL           string    `json:"issue_events_url,omitempty"`
	EventsURL                string    `json:"events_url,omitempty"`
	AssigneesURL             string    `json:"assignees_url,omitempty"`
	BranchesURL              string    `json:"branches_url,omitempty"`
	TagsURL                  string    `json:"tags_url,omitempty"`
	BlobsURL                 string    `json:"blobs_url,omitempty"`
	GitTagsURL               string    `json:"git_tags_url,omitempty"`
	GitRefsURL               string    `json:"git_refs_url,omitempty"`
	TreesURL                 string    `json:"trees_url,omitempty"`
	StatusesURL              string    `json:"statuses_url,omitempty"`
	LanguagesURL             string    `json:"languages_url,omitempty"`
	StargazersURL            string    `json:"stargazers_url,omitempty"`
	ContributorsURL          string    `json:"contributors_url,omitempty"`
	SubscribersURL           string    `json:"subscribers_url,omitempty"`
	SubscriptionURL          string    `json:"subscription_url,omitempty"`
	CommitsURL               string    `json:"commits_url,omitempty"`
	GitCommitsURL            string    `json:"git_commits_url,omitempty"`
	CommentsURL              string    `json:"comments_url,omitempty"`
	IssueCommentURL          string    `json:"issue_comment_url,omitempty"`
	ContentsURL              string    `json:"contents_url,omitempty"`
	CompareURL               string    `json:"compare_url,omitempty"`
	MergesURL                string    `json:"merges_url,omitempty"`
	ArchiveURL               string    `json:"archive_url,omitempty"`
	DownloadsURL             string    `json:"downloads_url,omitempty"`
	IssuesURL                string    `json:"issues_url,omitempty"`
	PullsURL                 string    `json:"pulls_url,omitempty"`
	MilestonesURL            string    `json:"milestones_url,omitempty"`
	NotificationsURL         string    `json:"notifications_url,omitempty"`
	LabelsURL                string    `json:"labels_url,omitempty"`
	ReleasesURL              string    `json:"releases_url,omitempty"`
	DeploymentsURL           string    `json:"deployments_url,omitempty"`
	CreatedAt                time.Time `json:"created_at,omitempty"`
	UpdatedAt                time.Time `json:"updated_at,omitempty"`
	PushedAt                 time.Time `json:"pushed_at,omitempty"`
	GitURL                   string    `json:"git_url,omitempty"`
	SSHURL                   string    `json:"ssh_url,omitempty"`
	CloneURL                 string    `json:"clone_url,omitempty"`
	SvnURL                   string    `json:"svn_url,omitempty"`
	Homepage                 any       `json:"homepage,omitempty"`
	Size                     int       `json:"size,omitempty"`
	StargazersCount          int       `json:"stargazers_count,omitempty"`
	WatchersCount            int       `json:"watchers_count,omitempty"`
	Language                 any       `json:"language,omitempty"`
	HasIssues                bool      `json:"has_issues,omitempty"`
	HasProjects              bool      `json:"has_projects,omitempty"`
	HasDownloads             bool      `json:"has_downloads,omitempty"`
	HasWiki                  bool      `json:"has_wiki,omitempty"`
	HasPages                 bool      `json:"has_pages,omitempty"`
	HasDiscussions           bool      `json:"has_discussions,omitempty"`
	ForksCount               int       `json:"forks_count,omitempty"`
	MirrorURL                any       `json:"mirror_url,omitempty"`
	Archived                 bool      `json:"archived,omitempty"`
	Disabled                 bool      `json:"disabled,omitempty"`
	OpenIssuesCount          int       `json:"open_issues_count,omitempty"`
	License                  any       `json:"license,omitempty"`
	AllowForking             bool      `json:"allow_forking,omitempty"`
	IsTemplate               bool      `json:"is_template,omitempty"`
	WebCommitSignoffRequired bool      `json:"web_commit_signoff_required,omitempty"`
	Topics                   []any     `json:"topics,omitempty"`
	Visibility               string    `json:"visibility,omitempty"`
	Forks                    int       `json:"forks,omitempty"`
	OpenIssues               int       `json:"open_issues,omitempty"`
	Watchers                 int       `json:"watchers,omitempty"`
	DefaultBranch            string    `json:"default_branch,omitempty"`
	TempCloneToken           any       `json:"temp_clone_token,omitempty"`
	NetworkCount             int       `json:"network_count,omitempty"`
	SubscribersCount         int       `json:"subscribers_count,omitempty"`
}

type Commits struct {
	Sha    string `json:"sha"`
	NodeID string `json:"node_id"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			Sha string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		URL          string `json:"url"`
		CommentCount int    `json:"comment_count"`
		Verification struct {
			Verified  bool   `json:"verified"`
			Reason    string `json:"reason"`
			Signature string `json:"signature"`
			Payload   string `json:"payload"`
		} `json:"verification"`
	} `json:"commit"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Author      struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Committer struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"committer"`
	Parents []struct {
		Sha     string `json:"sha"`
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"parents"`
}

type CommitsList []Commits

type UserIdentifiers struct {
	name         string
	login        string
	emailAddress string
	appearances  []string
}

type UsersList []UserIdentifiers

func responseHeaderPrinter(resp *http.Response) {
	fmt.Println("Response Header:")
	for key, value := range resp.Header {
		fmt.Printf("%s: %s\n", key, strings.Join(value, ", "))
	}
}

func Usage() {
	fmt.Println("Usage: gitvestigator -repo https://github.com/username/repo.git -t <GitHub_Personal_Token>")
	os.Exit(1)
}

func AddUser(user *UserIdentifiers, usersList *UsersList) {
	if (user.name == "GitHub" || user.login == "GitHub") && user.emailAddress == "noreply@github.com" {
		return
	}
	userPresent := false
	for i, userInList := range *usersList {
		if userInList.name == user.name && userInList.emailAddress == user.emailAddress && userInList.login == user.login {
			userPresent = true

			appearancePresent := false
			for _, appearance := range userInList.appearances {
				if appearance == user.appearances[0] {
					appearancePresent = true
					break
				}
			}
			if !appearancePresent {
				(*usersList)[i].appearances = append((*usersList)[i].appearances, user.appearances[0])
				sort.Strings((*usersList)[i].appearances)
			}
			break
		}
	}
	if !userPresent {
		*usersList = append(*usersList, *user)
		sort.Slice(*usersList, func(i, j int) bool {
			if (*usersList)[i].login == (*usersList)[j].login {
				return (*usersList)[i].emailAddress < (*usersList)[j].emailAddress
			}
			return (*usersList)[i].login < (*usersList)[j].login
		})
	}
}

func ParseArgs(args *Args) *Args {
	flag.StringVar(&args.repoLink, "repo", "", "The link to the repository")
	flag.StringVar(&args.token, "t", "", "GitHub Personal Access Token")
	flag.Parse()
	if len(os.Args) < 2 {
		flag.Usage = Usage
		flag.Usage()
	}
	if args.token == "" {
		fmt.Println("NO TOKEN WAS PROVIDED")
		fmt.Println("Number of requests will be limited to 60 per hour")
	}
	return args
}

func generateApiUrl(args *Args, repoData *RepoData) {
	repoData.URL = strings.Replace(args.repoLink, "https://github.com", "https://api.github.com/repos", -1)
}

func sendRequest(url string, args *Args) *http.Response {
	// fmt.Println("Sending request to URL: ", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+args.token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error occured while sending request to URL: ", url)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	// fmt.Println("Response Status: ", resp.Status)
	// body, err := io.ReadAll(resp.Body)
	// fmt.Println("Response Body: ", string(body))
	return resp
}

func GetRepoMetadata(args *Args, repoData *RepoData, usersList *UsersList) {
	if !strings.Contains(args.repoLink, "https://github.com") {
		fmt.Println("Kindly provide the link to the repository in the format: https://github.com/username/repo")
		os.Exit(1)
	}

	args.repoLink = strings.Replace(args.repoLink, ".git", "", -1)

	generateApiUrl(args, repoData)

	resp := sendRequest(repoData.URL, args)
	if resp.StatusCode == 404 {
		fmt.Println("Repository not found")
		fmt.Println("Please recheck the repository link and try again")
		fmt.Println("Repository link provided: ", args.repoLink)
		os.Exit(1)
	} else if resp.StatusCode == 200 {
		fmt.Println("Repository found")

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error occured while reading response body")
			fmt.Println("Read Error: ", err)
			os.Exit(1)
		}
		if err := json.Unmarshal(body, &repoData); err != nil {
			fmt.Println("JSON Unmarshal Error: ", err)
			os.Exit(1)
		}
		owner := &UserIdentifiers{"", repoData.Owner.Login, "", []string{"owner"}}
		// owner := UserIdentifiers{repoData.Owner.Login, "", []string{"owner"}}
		AddUser(owner, usersList)
	} else {
		fmt.Println("Unable to GET URL: ", repoData.URL)
		responseHeaderPrinter(resp)
		os.Exit(1)
	}
}

func GetCommits(repoData *RepoData, commitsList *CommitsList, args *Args) {
	baseCommitsUrl := strings.Replace(repoData.CommitsURL, "{/sha}", "", -1) + "?per_page=100&page="

	page := 1
	for {
		commitsUrl := baseCommitsUrl + fmt.Sprint(page)
		commitsInPage := &CommitsList{}
		resp := sendRequest(commitsUrl, args)
		if resp.StatusCode == 200 {
			// fmt.Println("Commits found")

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error occured while reading response body")
				fmt.Println("Read Error: ", err)
				os.Exit(1)
			}
			if err := json.Unmarshal(body, &commitsInPage); err != nil {
				fmt.Println("JSON Unmarshal Error: ", err)
				os.Exit(1)
			}
			*commitsList = append(*commitsList, *commitsInPage...)
			// fmt.Println("Total Commits: ", len(*commitsList))
			if len(*commitsInPage) < 100 {
				break
			}
			// for i, commit := range *commitsList {
			// 	fmt.Printf("%2d: %s\n", i+1, commit.Commit.URL)
			// }
		} else {
			fmt.Println("Unable to GET URL: ", commitsUrl)
			responseHeaderPrinter(resp)
			os.Exit(1)
		}
		page++
	}
	fmt.Println("Total Commits: ", len(*commitsList))
}

func PrintUsers(usersList *UsersList) {
	fmt.Println("Users found: ", len(*usersList))
	if len(*usersList) == 0 {
		fmt.Println("No users found")
		return
	}
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(writer, "No.\tName\tUsername\tEmail Address\tAppearences")
	fmt.Fprintln(writer, "---\t--------\t-------------\t-----------\t-----------")
	for i, user := range *usersList {
		fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\n", i+1, user.name, user.login, user.emailAddress, strings.Join(user.appearances, ", "))
		// f, err := os.OpenFile("users.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		// 	fmt.Println("Error occured while opening file")
		// 	fmt.Println("Error: ", err)
		// }
		// defer f.Close()
		// if _, err := f.WriteString(fmt.Sprintf("%d,%s,%s,%s,%s\n", i+1, user.name, user.login, user.emailAddress, strings.Join(user.appearances, ", "))); err != nil {
		// 	fmt.Println("Error occured while writing to file")
		// 	fmt.Println("Error: ", err)
		// }
	}
	writer.Flush()

}

func FindUsersFromCommits(commitsList *CommitsList, usersList *UsersList) {

	for _, commit := range *commitsList {
		// present := false
		author := &UserIdentifiers{commit.Commit.Author.Name, commit.Author.Login, commit.Commit.Author.Email, []string{"commit/author"}}
		// author := UserIdentifiers{commit.Commit.Author.Name, commit.Commit.Author.Email, []string{"commit/author"}}
		// for _, userInList := range *usersList {
		// 	if userInList.username == author.username && userInList.emailAddress == author.emailAddress {
		// 		present = true
		// 		break
		// 	}
		// }
		// if !present {
		// 	AddUser(author, usersList)
		// }
		AddUser(author, usersList)

		committer := &UserIdentifiers{commit.Commit.Committer.Name, commit.Committer.Login, commit.Commit.Committer.Email, []string{"commit/committer"}}
		// committer := UserIdentifiers{commit.Commit.Committer.Name, commit.Commit.Committer.Email, []string{"commit/committer"}}
		// for _, userInList := range *usersList {
		// 	if userInList.username == committer.username && userInList.emailAddress == committer.emailAddress {
		// 		present = true
		// 		break
		// 	}
		// }
		// if !present {
		// 	AddUser(committer, usersList)
		// }
		AddUser(committer, usersList)
	}

}

func main() {

	args := &Args{}               // Initialize the args variable
	repoData := &RepoData{}       // Initialize the RepoData variable
	commitsList := &CommitsList{} // Initialize the Commits variable
	usersList := &UsersList{}     // Initialize the users variable

	ParseArgs(args) // Pass the address of args to ParseArgs and dereference the returned value
	GetRepoMetadata(args, repoData, usersList)
	GetCommits(repoData, commitsList, args)
	FindUsersFromCommits(commitsList, usersList)

	// s, _ := json.MarshalIndent(repoData, "", "\t")
	// fmt.Println(string(s))

	// c, _ := json.MarshalIndent(commitsList, "", "\t")
	// fmt.Println(string(c))
	PrintUsers(usersList)
	fmt.Println("Bye Bye <3")

}
