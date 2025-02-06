package engine

import (
    "sync"
)

type User struct {
    ID       int
    Username string
    Karma    int
}

type SubReddit struct {
    ID      int
    Name    string
    Members map[int]*User
    Posts   []*Post
}

type Post struct {
    ID       int
    Title    string
    Content  string
    Author   *User
    Votes    int
    Comments []*Comment
}

type Comment struct {
    ID      int
    Content string
    Author  *User
    Votes   int
    Replies []*Comment
}

type Message struct {
    ID      int
    From    *User
    To      *User
    Content string
}

type RedditEngine struct {
    Users      map[int]*User
    SubReddits map[int]*SubReddit
    Messages   map[int]*Message
    mu         sync.Mutex
}