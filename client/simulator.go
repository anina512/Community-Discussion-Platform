package client

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reddit-clone/engine"
)

var LoggingEnabled = true 

type Simulator struct {
	Engine     *engine.RedditEngine
	Clients    []*Client
	SubReddits []*engine.SubReddit
}

func NewSimulator() *Simulator {
	if LoggingEnabled {
		logFile, err := os.OpenFile("reddit_simulation.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Error opening log file: %v\n", err)
			return nil
		}
		log.SetOutput(logFile)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	return &Simulator{
		Engine: engine.NewRedditEngine(),
	}
}

// Run the simulation and log actions
func (s *Simulator) Run(numUsers, numSRs, numPosts, numComments, numVotes, numMessages int) {
    if LoggingEnabled {
        log.Printf("Starting simulation with %d users, %d subreddits, %d posts, %d comments per post, %d votes per post, %d messages\n",
            numUsers, numSRs, numPosts, numComments, numVotes, numMessages)
    }

    // Create users and clients
    for i := 0; i < numUsers; i++ {
        client := NewClient(s.Engine, fmt.Sprintf("user%d", i))
        s.Clients = append(s.Clients, client)
        if LoggingEnabled {
            log.Printf("Created user: %s\n", client.User.Username)
        }
    }

    // Create subreddits
    for i := 0; i < numSRs; i++ {
        sr := s.Engine.CreateSubReddit(fmt.Sprintf("sr%d", i))
        s.SubReddits = append(s.SubReddits, sr)
        if LoggingEnabled {
            log.Printf("Created subreddit: %s\n", sr.Name)
        }
    }

    // Simulate activity
    for i := 0; i < numPosts; i++ {
        client := s.Clients[rand.Intn(len(s.Clients))]
        sr := s.SubReddits[rand.Intn(len(s.SubReddits))]
        post := client.CreatePost(sr, fmt.Sprintf("Post %d", i), "Content")
        if LoggingEnabled {
            log.Printf("User %s created a post in subreddit %s: %s\n", client.User.Username, sr.Name, post.Title)
        }

        for j := 0; j < numComments; j++ {
            commenter := s.Clients[rand.Intn(len(s.Clients))]
            comment := commenter.CreateComment(post, fmt.Sprintf("Comment %d", j))
            if LoggingEnabled {
                log.Printf("User %s commented on post '%s': %s\n", commenter.User.Username, post.Title, comment.Content)
            }
        }

        for j := 0; j < numVotes; j++ {
            voter := s.Clients[rand.Intn(len(s.Clients))]
            upvote := rand.Intn(2) == 0
            voter.Vote(post, upvote)
            if LoggingEnabled {
                voteType := "upvoted"
                if !upvote {
                    voteType = "downvoted"
                }
                log.Printf("User %s %s post '%s'\n", voter.User.Username, voteType, post.Title)
            }
        }
    }

    for i := 0; i < numMessages; i++ {
        fromClient := s.Clients[rand.Intn(len(s.Clients))]
        toClient := s.Clients[rand.Intn(len(s.Clients))]
        msg := fromClient.SendMessage(toClient.User, fmt.Sprintf("Message %d", i))
        if LoggingEnabled {
            log.Printf("User %s sent message to user %s: %s\n", fromClient.User.Username, toClient.User.Username, msg.Content)
        }
    }

    if LoggingEnabled {
        log.Println("Simulation completed")
    }
}