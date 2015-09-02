package main

import (
  "fmt"
  "strings"
  "flag"
  "log"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "encoding/base64"
  "github.com/hawky-4s-/nexus-cli-tool/utils"
//  "github.com/spf13/cobra"
)

const STATUS string = "service/local/status"
const GAV_SEARCH string = "service/local/lucene/search?"
const QUERY_SEARCH string = "service/local/lucene/search?q="
const REPOSITORIES string = "service/local/repositories"
const USERS string = "service/local/users"

var Debug bool = false

type Authentication struct {
  username string
  password string
}

/** Search request to the Nexus Lucene plugin endpoint **/
type SearchRequest struct {
  url string
  groupId string
  artifactId string
  version string
  repository string
}

type Repository struct {
  Id string `json:"repositoryId"`
  Url string `json:"repositoryUrl"`
  Name string `json:"repositoryName"`
  Policy string `json:"repositoryPolicy"`
  Kind string `json:"repositoryKind"`
  ContentClass string `json:"repositoryContentClass"`
}

type ArtifactLink struct {
  Classifier string
  Extension string
}

type ArtifactHit struct {
  RepositoryId string
  ArtifactLinks []ArtifactLink
}

type Artifact struct {
  GroupId string
  ArtifactId string
  Version string
  LatestVersion string `json:"latestSnapshot"`
  LatestSnapshotRepositoryId string
  ArtifactHits []ArtifactHit
}

/** Search response from the Nexus Lucene plugin endpoint **/
type SearchResponse struct {
  TotalCount int
  TooManyResults bool
  Collapsed bool
  Repositories []Repository `json:"repoDetails"`
  Artifacts []Artifact `json:"data"`
}

func parseCmdlineOpts() (SearchRequest, Authentication) {
  searchRequest := SearchRequest{}
  authentication := Authentication{}

  flag.StringVar(&searchRequest.url, "url", "https://app.camunda.com/nexus", "The url to your nexus installation.")
  flag.StringVar(&searchRequest.groupId, "g", "org.camunda.*", "The groupId of the arifacts to search.")
  flag.StringVar(&searchRequest.artifactId, "a", "", "The artifactId of the arifacts to search.")
  flag.StringVar(&searchRequest.version, "v", "7.3.*-SNAPSHOT", "The version of the artifacts to search.")
  flag.StringVar(&searchRequest.repository, "r", "camunda-bpm-snapshots", "The url to your nexus installation.")
  flag.StringVar(&authentication.username, "u", "us3r", "The basic auth username.")
  flag.StringVar(&authentication.password, "p", "s3cr3t", "The basic auth password.")
  flag.BoolVar(&Debug, "debug", false, "Set to true to enable debug mode.")

  flag.Parse()

  if (!strings.HasSuffix(searchRequest.url, "/")) {
    searchRequest.url += "/"
  }

  return searchRequest, authentication
}

func issueRequest(requestType, request string, auth Authentication) []byte {
  client := &http.Client{
//    CheckRedirect: redirectPolicyFunc,
  }

  req, err := http.NewRequest(requestType, request, nil)
  req.Header.Add("Accept", "application/json")
  if (len(strings.TrimSpace(auth.username)) > 0 &&
      len(strings.TrimSpace(auth.password)) > 0) {
    req.Header.Add("Authorization","Basic "+basicAuth(auth.username, auth.password))
  }

  res, err := client.Do(req)

  if err != nil {
    log.Fatal(err)
  }
  responseBody, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }

  return responseBody
}

func searchMatchingArtifacts(sr SearchRequest, auth Authentication) SearchResponse {

  request := ""
  if (sr.groupId != "") {
    request += "g=" + sr.groupId
  }
  if (sr.artifactId != "") {
    request += "&a=" + sr.artifactId
  }
  if (sr.version != "") {
    request += "&v=" + sr.version
  }
  if (sr.repository != "") {
    request += "&r=" + sr.repository
  }
  request = sr.url + GAV_SEARCH + request + "&collapseresults=true"

//  if (Debug) {
    fmt.Printf("Constructed Request URL: %s\n", request)
//  }

  responseBody := issueRequest("GET", request, auth)

  if (Debug) {
    fmt.Printf("Debug: %s\n", responseBody)
  }

  var response SearchResponse
  json.Unmarshal(responseBody, &response)

  return response
}

func deleteArtifacts(searchRequest SearchRequest, artifacts []Artifact, auth Authentication) string {
  // issue DELETE calls for all artifacts

  utils.AskForConfirmation("Do you want to delete the found artifacts?")

  numOfArtifacts := len(artifacts)
  for i, artifact := range artifacts {
//    issueRequest("DELETE", artifact.RepositoryId, auth)
    fmt.Printf("\rDeleting artifact %d of %d: %v\n", i+1, numOfArtifacts, &artifact)
  }

  fmt.Println("\nFinished deleting artifacts.")

  return ""
}

func basicAuth(username, password string) string {
  auth := username + ":" + password
  return base64.StdEncoding.EncodeToString([]byte(auth))
}

//func redirectPolicyFunc(req *http.Request, via []*http.Request) error{
//  req.Header.Add("Authorization","Basic "+basicAuth("user", "pass"))
//  return nil
//}

func (sr *SearchResponse) String() string {
  return fmt.Sprintf("%+v\n", sr)
}

func (a *Artifact) String() string {
  return fmt.Sprintf("Artifact[group=%s, id=%s, version=%s, latestVersion=%s, latestSnapshotRepository=%s]", a.GroupId, a.ArtifactId, a.Version, a.LatestVersion, a.LatestSnapshotRepositoryId)
}

func (r *Repository) String() string {
  return fmt.Sprintf("Repository[id=%s, policy=%s, url=%s]", r.Id, r.Policy, r.Url)
}

func (sr *SearchRequest) String() string {
  return fmt.Sprintf("Searching for GroupId:%s - ArtifactId:%s - Version:%s in Repository:%s at %s", sr.groupId, sr.artifactId, sr.version, sr.repository, sr.url)
}

func main() {

  searchRequest, authentication := parseCmdlineOpts()

  fmt.Printf("%+v\n", &searchRequest)

  searchResponse := searchMatchingArtifacts(searchRequest, authentication)

  fmt.Printf("%+v\n", searchResponse)

  deleteArtifacts(searchRequest, searchResponse.Artifacts, authentication)

}
