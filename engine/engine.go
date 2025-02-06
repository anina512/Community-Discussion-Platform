package engine
import "fmt"

func NewRedditEngine() *RedditEngine {
    return &RedditEngine{
        Users:      make(map[int]*User),
        SubReddits: make(map[int]*SubReddit),
        Messages:   make(map[int]*Message),
    }
}

func (e *RedditEngine) RegisterAccount(username string) *User {
    e.mu.Lock()
    defer e.mu.Unlock()
    user := &User{
        ID:       len(e.Users) + 1,
        Username: username,
    }
    e.Users[user.ID] = user
    return user
}

func (e *RedditEngine) CreateSubReddit(name string) *SubReddit {
    e.mu.Lock()
    defer e.mu.Unlock()
    sr := &SubReddit{
        ID:      len(e.SubReddits) + 1,
        Name:    name,
        Members: make(map[int]*User),
    }
    e.SubReddits[sr.ID] = sr
    return sr
}

func (e *RedditEngine) CreatePost(user *User, sr *SubReddit, title, content string) *Post {
    e.mu.Lock()
    defer e.mu.Unlock()
    post := &Post{
        ID:      len(sr.Posts) + 1,
        Title:   title,
        Content: content,
        Author:  user,
    }
    sr.Posts = append(sr.Posts, post)
    return post
}



func (e *RedditEngine) CreateComment(user *User, post *Post, content string) *Comment {
    e.mu.Lock()
    defer e.mu.Unlock()
    comment := &Comment{
        ID:      len(post.Comments) + 1,
        Content: content,
        Author:  user,
    }
    post.Comments = append(post.Comments, comment)
    return comment
}

func (e *RedditEngine) Vote(post *Post, upvote bool) {
    e.mu.Lock()
    defer e.mu.Unlock()
    if upvote {
        post.Votes++
        post.Author.Karma++
    } else {
        post.Votes--
        post.Author.Karma--
    }
}

func (e *RedditEngine) GetFeed(sr *SubReddit) []*Post {
    e.mu.Lock()
    defer e.mu.Unlock()
    return sr.Posts
}

func (e *RedditEngine) SendMessage(from, to *User, content string) *Message {
    e.mu.Lock()
    defer e.mu.Unlock()
    msg := &Message{
        ID:      len(e.Messages) + 1,
        From:    from,
        To:      to,
        Content: content,
    }
    e.Messages[msg.ID] = msg
    return msg
}

func (e *RedditEngine) GetMessages(user *User) []*Message {
    e.mu.Lock()
    defer e.mu.Unlock()
    var messages []*Message
    for _, msg := range e.Messages {
        if msg.To == user {
            messages = append(messages, msg)
        }
    }
    return messages
}

func (e *RedditEngine) GetSubRedditByName(name string) *SubReddit {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, sr := range e.SubReddits {
		if sr.Name == name {
			return sr
		}
	}
	return nil
}

func (e *RedditEngine) GetUserByUsername(username string) *User {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, user := range e.Users {
		if user.Username == username {
			return user
		}
	}
	return nil
}

func (e *RedditEngine) GetPostByID(id int) *Post {
    e.mu.Lock()
    defer e.mu.Unlock()
    for _, sr := range e.SubReddits {
        for _, post := range sr.Posts {
            if post.ID == id {
                return post
            }
        }
    }
    return nil
}

func (e *RedditEngine) GetAllPosts() []*Post {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    allPosts := make([]*Post, 0)
    for _, subreddit := range e.SubReddits {
        allPosts = append(allPosts, subreddit.Posts...)
    }
    
    return allPosts
}

func (e *RedditEngine) UserExists(username string) bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    for _, user := range e.Users {
        if user.Username == username {
            return true
        }
    }
    return false
}

func (e *RedditEngine) JoinSubReddit(user *User, sr *SubReddit) error {
    e.mu.Lock()
    defer e.mu.Unlock()
    if _, exists := sr.Members[user.ID]; exists {
        return fmt.Errorf("user already a member of this subreddit")
    }
    sr.Members[user.ID] = user
    return nil
}

func (e *RedditEngine) LeaveSubReddit(user *User, sr *SubReddit) error {
    e.mu.Lock()
    defer e.mu.Unlock()
    if _, exists := sr.Members[user.ID]; !exists {
        return fmt.Errorf("user is not a member of this subreddit")
    }
    delete(sr.Members, user.ID)
    return nil
}
