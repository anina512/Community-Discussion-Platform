package client

import (
    "reddit-clone/engine"
)

type Client struct {
    Engine *engine.RedditEngine
    User   *engine.User
}

func NewClient(engine *engine.RedditEngine, username string) *Client {
    user := engine.RegisterAccount(username)
    return &Client{
        Engine: engine,
        User:   user,
    }
}

func (c *Client) CreatePost(sr *engine.SubReddit, title, content string) *engine.Post {
    return c.Engine.CreatePost(c.User, sr, title, content)
}

func (c *Client) CreateComment(post *engine.Post, content string) *engine.Comment {
    return c.Engine.CreateComment(c.User, post, content)
}

func (c *Client) Vote(post *engine.Post, upvote bool) {
    c.Engine.Vote(post, upvote)
}

func (c *Client) SendMessage(to *engine.User, content string) *engine.Message {
    return c.Engine.SendMessage(c.User, to, content)
}

func (c *Client) GetMessages() []*engine.Message {
    return c.Engine.GetMessages(c.User)
}
