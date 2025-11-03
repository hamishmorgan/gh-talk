package api

import "time"

// Thread represents a pull request review thread
type Thread struct {
	ID                 string
	IsResolved         bool
	IsCollapsed        bool
	IsOutdated         bool
	Path               string
	Line               int
	StartLine          int
	DiffSide           string
	SubjectType        string
	ResolvedBy         *User
	Comments           []Comment
	ViewerCanResolve   bool
	ViewerCanUnresolve bool
	ViewerCanReply     bool
}

// Comment represents a review comment or issue comment
type Comment struct {
	ID                string
	DatabaseID        int
	Body              string
	Path              string
	Position          int
	DiffHunk          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Author            User
	AuthorAssociation string
	ReplyTo           *CommentRef
	IsMinimized       bool
	MinimizedReason   string
	ReactionGroups    []ReactionGroup
	ViewerCanReact    bool
	ViewerCanUpdate   bool
	ViewerCanDelete   bool
	ViewerCanMinimize bool
}

// CommentRef is a reference to another comment
type CommentRef struct {
	ID string
}

// User represents a GitHub user
type User struct {
	Login string
}

// ReactionGroup represents aggregated reactions
type ReactionGroup struct {
	Content          string
	CreatedAt        *time.Time
	Users            ReactionUsers
	ViewerHasReacted bool
}

// ReactionUsers represents users who reacted
type ReactionUsers struct {
	TotalCount int
	Nodes      []User
}

// Issue represents a GitHub issue
type Issue struct {
	ID       string
	Number   int
	Title    string
	State    string
	Body     string
	Comments []Comment
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID            string
	Number        int
	Title         string
	State         string
	ReviewThreads []Thread
}
