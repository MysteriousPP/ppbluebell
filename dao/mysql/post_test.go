package mysql

import (
	"bluebell/models"
	"testing"
)

// func init(){
// 	dbCfg := settings.MySQLConfig{

// 	}
// 	err := Init()

// }
func TestCreatePost(t *testing.T) {
	post := models.Post{
		PostID:      10,
		AuthorId:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
