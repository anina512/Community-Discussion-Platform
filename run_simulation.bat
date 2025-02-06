@echo off

rem Start the server in a separate window (runs server in the background)
start "" cmd /c "go run server.go"
echo Server started

rem Wait for server to initialize (you may adjust sleep time)
timeout /t 2

rem Simulate multiple users concurrently

rem User1 creates a user, subreddit, and posts
start "" cmd /c "go run client.go createUser User1"
start "" cmd /c "go run client.go createSubreddit Subreddit1"
start "" cmd /c "go run client.go submitPost Subreddit1 User1 Post1 This is a post from User1"

rem User2 creates a user, joins subreddit, and posts
start "" cmd /c "go run client.go createUser User2"
start "" cmd /c "go run client.go joinSubreddit Subreddit1 User2"
start "" cmd /c "go run client.go submitPost Subreddit1 User2 Post2 This is a post from User2"

rem User3 creates a user, joins subreddit, and comments on User1's post (PostID 1)
start "" cmd /c "go run client.go createUser User3"
start "" cmd /c "go run client.go joinSubreddit Subreddit1 User3"
start "" cmd /c "go run client.go createComment 1 User3 Nice post, User1!"

rem User4 creates a user and votes on User1's post (PostID 1)
start "" cmd /c "go run client.go createUser User4"
start "" cmd /c "go run client.go vote 1 true"

rem User5 creates a user, joins subreddit, and votes on User2's post (PostID 2)
start "" cmd /c "go run client.go createUser User5"
start "" cmd /c "go run client.go joinSubreddit Subreddit1 User5"
start "" cmd /c "go run client.go vote 2 true"

rem User6 creates a user, submits a post, and comments on User2's post (PostID 2)
start "" cmd /c "go run client.go createUser User6"
start "" cmd /c "go run client.go submitPost Subreddit1 User6 Post3 This is a post from User6"
start "" cmd /c "go run client.go createComment 2 User6 Great post, User2!"

rem User7 creates a user, joins subreddit, votes on User1's post (PostID 1), and posts
start "" cmd /c "go run client.go createUser User7"
start "" cmd /c "go run client.go joinSubreddit Subreddit1 User7"
start "" cmd /c "go run client.go vote 1 true"
start "" cmd /c "go run client.go submitPost Subreddit1 User7 Post4 This is a post from User7"

rem User8 creates a user, joins subreddit, posts, and comments on User1's post (PostID 1)
start "" cmd /c "go run client.go createUser User8"
start "" cmd /c "go run client.go joinSubreddit Subreddit1 User8"
start "" cmd /c "go run client.go submitPost Subreddit1 User8 Post5 This is a post from User8"
start "" cmd /c "go run client.go createComment 1 User8 Awesome post, User1!"

rem User9 creates a user and simulates multiple comments on User2's post (PostID 2)
start "" cmd /c "go run client.go createUser User9"
start "" cmd /c "go run client.go createComment 2 User9 Very informative, User2!"
start "" cmd /c "go run client.go createComment 2 User9 I agree, User2!"
start "" cmd /c "go run client.go createComment 2 User9 Great insights, User2!"

rem User10 creates a user, joins subreddit, and simulates multiple votes on User1's post (PostID 1)
start "" cmd /c "go run client.go createUser User10"
start "" cmd /c "go run client.go joinSubreddit Subreddit1 User10"
start "" cmd /c "go run client.go vote 1 true"
rem Fetch the feed for Subreddit1
start "" cmd /c "go run client.go getFeed Subreddit1"
start "" cmd /c "go run client.go vote 1 false"
start "" cmd /c "go run client.go vote 1 true"



rem Wait for all background tasks to finish (adjust time as needed)
timeout /t 60

rem Stop the server (you can use the taskkill command to kill the server process)
rem Ensure to replace 'server.exe' with the actual name of your compiled server binary
taskkill /f /im "main.exe"
echo Server stopped
