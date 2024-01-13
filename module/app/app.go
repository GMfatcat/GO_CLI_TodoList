/* Define  App Functions*/
package app

import (
	"cli_todolist/module/basic"
	c "cli_todolist/module/color"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type Todos []basic.Item

func (t *Todos) Add(task string, isUrgent bool) {
	todo := basic.Item{
		Task:        task,
		Done:        false,
		Urgent:      isUrgent,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

// Task ID : 1 -> N
func (t *Todos) Complete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}
	// Task ID starts at 1, list index starts at 0, so -1
	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	return nil
}

// Task ID : 1 -> N
func (t *Todos) Delete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}
	/* Let ls = [t1,t2,t3,t4,t5] and index = 3
	   index = 3 -> want to delete ls[3-1] = ls[2] aka the third item
	   Key Point 1 : ls[:index-1] = ls[:3-1] = ls[:2] = first 2 item in ls
	              -> [t1,t2]
	   Key Point 2 : ls[index:]... means take out each item in ls[index:] = ls[3:] = [t4,t5]
	              -> t4,t5 seperately append to [t1,t2] -> [t1,t2,t4,t5]
	*/
	*t = append(ls[:index-1], ls[index:]...)

	return nil
}

func (t *Todos) Load(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(filename string) error {

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (t *Todos) Print() {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignCenter, Text: "CreatedAt"},
			{Align: simpletable.AlignCenter, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell

	for idx, item := range *t {
		idx++
		task := c.Blue(item.Task)
		done := c.Blue("no")
		// Mark Red if Urgent
		if item.Urgent {
			task = c.Red(item.Task)
			done = c.Red("no")
		}
		// Mark Green & check if Done
		if item.Done {
			task = c.Green(fmt.Sprintf("\u2705 %s", item.Task))
			done = c.Green("yes")
		}
		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: task},
			{Text: done},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: item.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	// Get TODO List Status
	total, urgent := t.CountPending()

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5,
			Text: c.Red(fmt.Sprintf("You have %d pending todos (%d urgent)", total, urgent))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() (int, int) {
	total := 0
	urgent := 0
	for _, item := range *t {
		if !item.Done {
			total++

			if item.Urgent {
				urgent++
			}
		}
	}

	return total, urgent
}

func (t *Todos) CleanUp() {
	// set empty array to save new todo list
	var newList Todos
	var changed bool = false
	remaining := 0
	for _, item := range *t {
		if !item.Done {
			changed = true
			newList = append(newList, item)
			remaining++
		}
	}

	if changed {
		*t = newList
		fmt.Printf("Finished Cleaning Up: %d task remains\n", remaining)
	} else {
		fmt.Println("No completed tasks to clean up")
	}
}
