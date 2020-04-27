package handlers

import (
	"log"
	"strconv"

	"github.com/taipoxin/json-rpc-pg/internal/api/models"
)

// Main handlers provider
type Main struct {
	Db models.Datastore
}

type HelloArgs struct {
	Name string
}

// Hello example
func (main *Main) Hello(args *HelloArgs, result *string) error {
	if args.Name == "" {
		return errParams
	}
	*result = "Hello " + args.Name
	return nil
}

func (main *Main) GetPosts(s struct{}, result *[]*models.Post) error {
	posts, err := main.Db.AllPosts()
	if err != nil {
		return errInternal
	}
	*result = posts
	return nil
}

func (main *Main) GetPost(args []int64, result **models.Post) error {
	if len(args) == 0 {
		return errParams
	}
	id := args[0]

	post, err := main.Db.GetPost(id)
	if err != nil {
		return errInternal
	}
	if post.Title == "" {
		result = nil
		return nil
	}
	*result = post
	return nil
}

func (main *Main) AddPost(pArgs *models.Post, result *map[string]interface{}) error {
	if pArgs.Title == "" {
		return errParams
	}

	id, err := main.Db.AddPost(pArgs.Title)
	if err != nil {
		log.Println(err)
		return errInternal
	}

	*result = map[string]interface{}{
		"ok":      true,
		"message": "inserted successfully, id: " + strconv.FormatInt(id, 10),
	}
	return nil
}

func (main *Main) UpdatePost(pArgs *models.Post, result *map[string]interface{}) error {

	if pArgs.Title == "" || pArgs.ID == 0 {
		return errParams
	}

	isUpdated, err := main.Db.UpdatePost(pArgs.ID, pArgs.Title)
	if err != nil {
		return errInternal
	}
	var m string
	if isUpdated {
		m = "updated successfully"
	} else {
		m = "nothing to update"
	}
	*result = map[string]interface{}{
		"ok": isUpdated, "message": m,
	}
	return nil
}

func (main *Main) DeletePost(pArgs *models.Post, result *map[string]interface{}) error {
	if pArgs.ID == 0 {
		return errParams
	}

	isDeleted, err := main.Db.DeletePost(pArgs.ID)
	if err != nil {
		return errInternal
	}

	var m string
	if isDeleted {
		m = "deleted successfully"
	} else {
		m = "nothing to delete"
	}
	*result = map[string]interface{}{
		"ok": isDeleted, "message": m,
	}
	return nil
}
