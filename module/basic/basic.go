/* Define basic data structures*/
package basic

import "time"

type Item struct {
	Task        string
	Done        bool
	Urgent      bool
	CreatedAt   time.Time
	CompletedAt time.Time
}
