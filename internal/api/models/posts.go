package models

import (
	"fmt"
	"log"
)

// Post entity of table `posts`
type Post struct {
	ID    int64
	Title string
}

func (p *Post) String() string {
	return fmt.Sprintf("{id: %v, title: %v}", p.ID, p.Title)
}

// AllPosts return all `posts` from db
func (db *DbHelper) AllPosts() ([]*Post, error) {
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*Post, 0)
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.ID, &post.Title)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

// GetPost return post by id from db
func (db *DbHelper) GetPost(id int64) (*Post, error) {
	rows, err := db.Query("SELECT * FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	post := &Post{}
	for rows.Next() {
		err := rows.Scan(&post.ID, &post.Title)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return post, nil
}

// AddPost add new post with title
func (db *DbHelper) AddPost(title string) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO  posts (title) VALUES ($1) RETURNING id", title).Scan(&id)
	if err != nil {
		return 0, err
	}
	log.Println("Last inserted ID: ", id)
	return id, nil
}

// UpdatePost update post with new title by id
func (db *DbHelper) UpdatePost(id int64, title string) (bool, error) {
	result, err := db.Exec("UPDATE posts SET title = $1 WHERE id = $2", title, id)
	if err != nil {
		return false, err
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	log.Println("Rows affected: ", ra)
	if ra != 0 {
		return true, nil
	}
	return false, nil
}

// DeletePost delete post by id
func (db *DbHelper) DeletePost(id int64) (bool, error) {
	result, err := db.Exec("DELETE from posts WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	log.Println("Rows affected: ", ra)
	if ra != 0 {
		return true, nil
	}
	return false, nil
}
