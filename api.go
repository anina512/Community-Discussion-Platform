package main

import (
	"encoding/json"
	"net/http"
	"reddit-clone/engine"
	"strings"
	"strconv"
)

type API struct {
	engine *engine.RedditEngine
}

func NewAPI(e *engine.RedditEngine) *API {
	return &API{engine: e}
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST" && r.URL.Path == "/api/user":
		api.createUser(w, r)
	case r.Method == "POST" && r.URL.Path == "/api/subreddit":
		api.createSubreddit(w, r)
	case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/submit"):
		api.submitPost(w, r)
	case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/comment"):
		api.createComment(w, r)
	case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/vote"):
		api.vote(w, r)
	case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/getusers"):
		api.getAllUsers(w, r)
	case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/getsubreddits"):
		api.getAllSubreddits(w, r)
	case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/getposts"):
		api.getAllPosts(w, r)
	case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/join"):
        api.joinSubreddit(w, r)
    case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/leave"):
        api.leaveSubreddit(w, r)
    case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/feed"):
        api.getFeed(w, r)
	default:
		http.NotFound(w, r)
	}
}


func (api *API) createSubreddit(w http.ResponseWriter, r *http.Request) {
	var subredditData struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&subredditData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subreddit := api.engine.CreateSubReddit(subredditData.Name)
	json.NewEncoder(w).Encode(subreddit)
}

func (api *API) submitPost(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    subredditName := parts[2]
    var postData struct {
        Title    string `json:"title"`
        Content  string `json:"content"`
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    subreddit := api.engine.GetSubRedditByName(subredditName)
    if subreddit == nil {
        http.Error(w, "Subreddit not found", http.StatusNotFound)
        return
    }

    user := api.engine.GetUserByUsername(postData.Username)
    if user == nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    post := api.engine.CreatePost(user, subreddit, postData.Title, postData.Content)
    json.NewEncoder(w).Encode(post)
}


func (api *API) createComment(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    postIDStr := parts[2]
    
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

	var commentData struct {
		Content  string `json:"content"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&commentData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := api.engine.GetUserByUsername(commentData.Username)
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	post := api.engine.GetPostByID(postID)
    if post == nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

	comment := api.engine.CreateComment(user, post, commentData.Content)
	json.NewEncoder(w).Encode(comment)
}

func (api *API) vote(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    postID := parts[2]

    var voteData struct {
        Upvote bool `json:"upvote"`
    }
    if err := json.NewDecoder(r.Body).Decode(&voteData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    postIDInt, err := strconv.Atoi(postID)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    post := api.engine.GetPostByID(postIDInt)
    if post == nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    api.engine.Vote(post, voteData.Upvote)
    w.WriteHeader(http.StatusOK)
}

func (api *API) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]*engine.User, 0, len(api.engine.Users))
	for _, user := range api.engine.Users {
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func (api *API) getAllSubreddits(w http.ResponseWriter, r *http.Request) {
    subreddits := make([]*engine.SubReddit, 0, len(api.engine.SubReddits))
    for _, subreddit := range api.engine.SubReddits {
        subreddits = append(subreddits, subreddit)
    }
    json.NewEncoder(w).Encode(subreddits)
}

func (api *API) getAllPosts(w http.ResponseWriter, r *http.Request) {
    allPosts := make([]*engine.Post, 0)
    
    for _, subreddit := range api.engine.SubReddits {
        allPosts = append(allPosts, subreddit.Posts...)
    }
    
    json.NewEncoder(w).Encode(allPosts)
}

func (api *API) createUser(w http.ResponseWriter, r *http.Request) {
    var userData struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if api.engine.UserExists(userData.Username) {
        // User exists, return the existing user (login)
        user := api.engine.GetUserByUsername(userData.Username)
        json.NewEncoder(w).Encode(user)
    } else {
        // User doesn't exist, create a new user (register)
        user := api.engine.RegisterAccount(userData.Username)
        json.NewEncoder(w).Encode(user)
    }
}

func (api *API) joinSubreddit(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    subredditName := parts[2]
    var userData struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user := api.engine.GetUserByUsername(userData.Username)
    if user == nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    subreddit := api.engine.GetSubRedditByName(subredditName)
    if subreddit == nil {
        http.Error(w, "Subreddit not found", http.StatusNotFound)
        return
    }

    err := api.engine.JoinSubReddit(user, subreddit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (api *API) leaveSubreddit(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    subredditName := parts[2]
    var userData struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user := api.engine.GetUserByUsername(userData.Username)
    if user == nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    subreddit := api.engine.GetSubRedditByName(subredditName)
    if subreddit == nil {
        http.Error(w, "Subreddit not found", http.StatusNotFound)
        return
    }

    err := api.engine.LeaveSubReddit(user, subreddit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (api *API) getFeed(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    subredditName := parts[2]

    subreddit := api.engine.GetSubRedditByName(subredditName)
    if subreddit == nil {
        http.Error(w, "Subreddit not found", http.StatusNotFound)
        return
    }

    posts := api.engine.GetFeed(subreddit)

    json.NewEncoder(w).Encode(posts)
}

