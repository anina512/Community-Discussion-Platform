package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8080/api"

func initLog() *os.File {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}

	log.SetOutput(logFile)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile
}

func createUser(username string) {
	data := map[string]string{"username": username}
	body, _ := json.Marshal(data)

	resp, err := http.Post(baseURL+"/user", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating user:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println("User created:", result)
}

func createSubreddit(name string) {
	data := map[string]string{"name": name}
	body, _ := json.Marshal(data)

	resp, err := http.Post(baseURL+"/subreddit", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating subreddit:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println("Subreddit created:", result)
}

func submitPost(subreddit, username, title, content string) {
	data := map[string]string{"title": title, "content": content, "username": username}
	body, _ := json.Marshal(data)

	resp, err := http.Post(fmt.Sprintf("%s/%s/submit", baseURL, subreddit), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error submitting post:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println("Post submitted:", result)
}


func createComment(postID string, username, content string) {
	data := map[string]string{"content": content, "username": username}
	body, _ := json.Marshal(data)

	url := fmt.Sprintf("%s/%s/comment", baseURL, postID)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating comment:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println("Comment created:", result)
}



func voteOnPost(postID string, upvote bool) {
	data := map[string]bool{"upvote": upvote}
	body, _ := json.Marshal(data)

	url := fmt.Sprintf("%s/%s/vote", baseURL, postID)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error voting on post:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	log.Println("Vote successful!")
}



func getAllUsers() {
	resp, err := http.Get(baseURL + "/getusers")
	if err != nil {
		log.Println("Error fetching users:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	var users []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&users)
	log.Println("Users:", users)
}

func getAllSubreddits() {
	resp, err := http.Get(baseURL + "/getsubreddits")
	if err != nil {
		log.Println("Error fetching subreddits:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	var subreddits []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&subreddits)
	log.Println("Subreddits:", subreddits)
}

func getAllPosts() {
	resp, err := http.Get(baseURL + "/getposts")
	if err != nil {
		log.Println("Error fetching posts:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	var posts []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&posts)
	log.Println("Posts:", posts)
}

func joinSubreddit(subreddit, username string) {
	data := map[string]string{"username": username}
	body, _ := json.Marshal(data)

	resp, err := http.Post(fmt.Sprintf("%s/%s/join", baseURL, subreddit), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error joining subreddit:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	log.Println("Joined subreddit:", subreddit)
}

func leaveSubreddit(subreddit, username string) {
	data := map[string]string{"username": username}
	body, _ := json.Marshal(data)

	resp, err := http.Post(fmt.Sprintf("%s/%s/leave", baseURL, subreddit), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error leaving subreddit:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	log.Println("Left subreddit:", subreddit)
}

func getFeed(subreddit string) {
	url := fmt.Sprintf("%s/%s/feed", baseURL, subreddit)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	var feed []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&feed)
	log.Println("Feed:", feed)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createUser <username>")
	fmt.Println("  createSubreddit <subreddit_name>")
	fmt.Println("  submitPost <subreddit> <username> <title> <content>")
	fmt.Println("  createComment <postID> <username> <content>")
	fmt.Println("  vote <postID> <true/false>")
	fmt.Println("  getAllUsers")
	fmt.Println("  getAllSubreddits")
	fmt.Println("  getAllPosts")
	fmt.Println("  joinSubreddit <subreddit> <username>")
	fmt.Println("  leaveSubreddit <subreddit> <username>")
	fmt.Println("  getFeed <username>")
}

func main() {
	logFile := initLog()
	defer logFile.Close()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "createUser":
		if len(os.Args) < 3 {
			log.Println("Please provide a username.")
			printUsage()
			return
		}
		createUser(os.Args[2])
	case "createSubreddit":
		if len(os.Args) < 3 {
			log.Println("Please provide a subreddit name.")
			printUsage()
			return
		}
		createSubreddit(os.Args[2])
	case "submitPost":
		if len(os.Args) < 6 {
			log.Println("Please provide subreddit, username, title, and content for the post.")
			printUsage()
			return
		}
		submitPost(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	case "createComment":
		if len(os.Args) < 5 {
			log.Println("Please provide postID, username, and content for the comment.")
			printUsage()
			return
		}
		postID := os.Args[2]
		createComment(postID, os.Args[3], os.Args[4])
	case "getFeed":
		if len(os.Args) < 3 {
			log.Println("Please provide a username.")
			printUsage()
			return
		}
		getFeed(os.Args[2])
	case "vote":
		if len(os.Args) < 4 {
			log.Println("Please provide postID and vote status (true/false).")
			printUsage()
			return
		}
		postID := os.Args[2]
		upvote := os.Args[3] == "true"
		voteOnPost(postID, upvote)
	case "getAllUsers":
		getAllUsers()
	case "getAllSubreddits":
		getAllSubreddits()
	case "getAllPosts":
		getAllPosts()
	case "joinSubreddit":
		if len(os.Args) < 4 {
			log.Println("Please provide subreddit and username.")
			printUsage()
			return
		}
		joinSubreddit(os.Args[2], os.Args[3])
	case "leaveSubreddit":
		if len(os.Args) < 4 {
			log.Println("Please provide subreddit and username.")
			printUsage()
			return
		}
		leaveSubreddit(os.Args[2], os.Args[3])
	default:
		log.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}
