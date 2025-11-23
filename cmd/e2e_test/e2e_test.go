package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

var baseURL = func() string {
	if v := os.Getenv("BASE_URL"); v != "" {
		return v
	}

	return "http://localhost:8080"
}()

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	User User `json:"user"`
}

type PullRequest struct {
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorID        string   `json:"author_id"`
	Status          string   `json:"status"`
	Assigned        []string `json:"assigned_reviewers"`
	CreatedAt       *string  `json:"createdAt,omitempty"`
	MergedAt        *string  `json:"mergedAt,omitempty"`
	_               struct{} // no unknown fields check
}

type PullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type UserReviewsResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func postJSON(t *testing.T, path string, body any) *http.Response {
	t.Helper()

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(data))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}

	return resp
}

func get(t *testing.T, path string, q url.Values) *http.Response {
	t.Helper()

	u := baseURL + path
	if len(q) > 0 {
		u = u + "?" + q.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u, http.NoBody)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}

	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, dst any) {
	t.Helper()

	defer func() {
		_ = resp.Body.Close()
	}()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(dst); err != nil {
		t.Fatalf("decode json: %v", err)
	}
}

//nolint:bodyclose
func TestEndToEnd_PrLifecycle(t *testing.T) {
	teamName := fmt.Sprintf("backend-%d", time.Now().UnixNano())
	authorID := "u1"
	reviewerID := "u2"
	prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())

	{
		body := Team{
			TeamName: teamName,
			Members: []TeamMember{
				{UserID: authorID, Username: "Alice", IsActive: true},
				{UserID: reviewerID, Username: "Bob", IsActive: true},
				{UserID: "u3", Username: "Charlie", IsActive: true},
			},
		}

		resp := postJSON(t, "/team/add", body)
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %s", err)
		}

		if resp.StatusCode != http.StatusCreated {
			var errBody ErrorResponse
			decodeJSON(t, resp, &errBody)
			t.Fatalf("expected 201, got %d: %+v", resp.StatusCode, errBody)
		}

		var tr Team
		decodeJSON(t, resp, &tr)

		if tr.TeamName != teamName {
			t.Fatalf("unexpected team_name: %s", tr.TeamName)
		}

		if len(tr.Members) != 3 {
			t.Fatalf("expected 3 members, got %d", len(tr.Members))
		}
	}

	{
		resp := get(t, "/team/get", url.Values{"team_name": {teamName}})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 for /team/get, got %d", resp.StatusCode)
		}

		var team Team
		decodeJSON(t, resp, &team)

		if team.TeamName != teamName {
			t.Fatalf("GET /team/get returned wrong team_name: %s", team.TeamName)
		}
	}

	var createdPR PullRequest
	{
		body := map[string]any{
			"pull_request_id":   prID,
			"pull_request_name": "Add search endpoint",
			"author_id":         authorID,
		}

		resp := postJSON(t, "/pullRequest/create", body)
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %s", err)
		}

		if resp.StatusCode != http.StatusCreated {
			var errBody ErrorResponse
			decodeJSON(t, resp, &errBody)
			t.Fatalf("expected 201 on create PR, got %d: %+v", resp.StatusCode, errBody)
		}

		var prResp PullRequestResponse
		decodeJSON(t, resp, &prResp)
		createdPR = prResp.PR

		if createdPR.PullRequestID != prID {
			t.Fatalf("unexpected pr id: %s", createdPR.PullRequestID)
		}

		if createdPR.Status != "OPEN" {
			t.Fatalf("expected status OPEN, got: %s", createdPR.Status)
		}

		for _, r := range createdPR.Assigned {
			if r == authorID {
				t.Fatalf("author must not be in assigned_reviewers")
			}
		}

		if len(createdPR.Assigned) > 2 {
			t.Fatalf("assigned_reviewers must be <=2, got %d", len(createdPR.Assigned))
		}
	}

	if len(createdPR.Assigned) == 0 {
		t.Logf("no reviewers assigned (allowed by spec), пропускаем часть про /users/getReview")
	} else {
		reviewer := createdPR.Assigned[0]

		resp := get(t, "/users/getReview", url.Values{"user_id": {reviewer}})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 for /users/getReview, got %d", resp.StatusCode)
		}

		var ur UserReviewsResponse
		decodeJSON(t, resp, &ur)

		found := false

		for _, pr := range ur.PullRequests {
			if pr.PullRequestID == prID {
				found = true

				break
			}
		}

		if !found {
			t.Fatalf("PR %s not found in /users/getReview for user %s", prID, reviewer)
		}
	}

	var mergedPR PullRequest
	{
		body := map[string]any{
			"pull_request_id": prID,
		}

		resp := postJSON(t, "/pullRequest/merge", body)
		if resp.StatusCode != http.StatusOK {
			var errBody ErrorResponse
			decodeJSON(t, resp, &errBody)
			t.Fatalf("expected 200 on merge, got %d: %+v", resp.StatusCode, errBody)
		}

		var prResp PullRequestResponse
		decodeJSON(t, resp, &prResp)
		mergedPR = prResp.PR

		if mergedPR.Status != "MERGED" {
			t.Fatalf("expected status MERGED after merge, got: %s", mergedPR.Status)
		}

		if mergedPR.MergedAt == nil {
			t.Fatalf("mergedAt must be set after merge")
		}
	}
}

//nolint:bodyclose
func TestEndToEnd_TeamExistsError(t *testing.T) {
	teamName := fmt.Sprintf("team-%d", time.Now().UnixNano())

	body := Team{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
		},
	}

	resp1 := postJSON(t, "/team/add", body)
	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 on first /team/add, got %d", resp1.StatusCode)
	}

	_ = resp1.Body.Close()

	resp2 := postJSON(t, "/team/add", body)
	if resp2.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 on second /team/add, got %d", resp2.StatusCode)
	}

	var errBody ErrorResponse
	decodeJSON(t, resp2, &errBody)

	if errBody.Error.Code != "TEAM_EXISTS" {
		t.Fatalf("expected error code TEAM_EXISTS, got: %s", errBody.Error.Code)
	}
}

//nolint:bodyclose
func TestEndToEnd_ReassignOnMergedReturns409(t *testing.T) {
	teamName := fmt.Sprintf("team-reassign-%d", time.Now().UnixNano())
	authorID := "u10"
	reviewerID := "u11"
	otherID := "u12"
	prID := fmt.Sprintf("pr-reassign-%d", time.Now().UnixNano())

	{
		body := Team{
			TeamName: teamName,
			Members: []TeamMember{
				{UserID: authorID, Username: "Author", IsActive: true},
				{UserID: reviewerID, Username: "Reviewer", IsActive: true},
				{UserID: otherID, Username: "Other", IsActive: true},
			},
		}

		resp := postJSON(t, "/team/add", body)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201 on /team/add, got %d", resp.StatusCode)
		}

		_ = resp.Body.Close()
	}

	{
		body := map[string]any{
			"pull_request_id":   prID,
			"pull_request_name": "Test reassign",
			"author_id":         authorID,
		}

		resp := postJSON(t, "/pullRequest/create", body)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201 on /pullRequest/create, got %d", resp.StatusCode)
		}

		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %s", err)
		}
	}

	{
		body := map[string]any{
			"pull_request_id": prID,
		}

		resp := postJSON(t, "/pullRequest/merge", body)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 on /pullRequest/merge, got %d", resp.StatusCode)
		}

		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %s", err)
		}
	}

	{
		body := map[string]any{
			"pull_request_id": prID,
			"old_user_id":     reviewerID,
		}

		resp := postJSON(t, "/pullRequest/reassign", body)
		if resp.StatusCode != http.StatusConflict {
			var e ErrorResponse
			decodeJSON(t, resp, &e)
			t.Fatalf("expected 409 on /pullRequest/reassign after merged, got %d: %+v", resp.StatusCode, e)
		}

		var e ErrorResponse
		decodeJSON(t, resp, &e)

		if e.Error.Code != "PR_MERGED" {
			t.Fatalf("expected PR_MERGED, got: %s", e.Error.Code)
		}
	}
}
