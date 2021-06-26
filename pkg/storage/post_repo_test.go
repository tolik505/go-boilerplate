package storage

import (
	"goboilerplate/pkg/model"
	"goboilerplate/pkg/testhelper"
	"gorm.io/datatypes"
)

func (suite *storageTestSuite) Test_Create_Post() {
	suite.Run("It Creates A Post", func() {
		r := NewPostRepo(testhelper.TestDB)
		m := datatypes.JSON(`{"age": 18, "args": {"arga": "arga"}, "name": "John", "tags": ["tag1", "tag2"]}`)
		p := model.Post{
			UUID:     "uuid-1",
			Content:  "Test content",
			Metadata: &m,
		}
		err := r.Create(p)

		suite.Nil(err)
		// Check whether the record was inserted
		pst := model.Post{}
		testhelper.TestDB.
			Where(model.Post{UUID: "uuid-1"}).
			Take(&pst)

		suite.Equal(p, pst)
	})

	suite.Run("It Handles DB Error", func() {
		r := NewPostRepo(testhelper.BrokenDB())
		p := model.Post{
			UUID:    "uuid-1",
			Content: "Test content",
		}
		err := r.Create(p)

		suite.NotNil(err)
	})
}

func (suite *storageTestSuite) Test_Get_Post() {
	suite.Run("It Gets Post", func() {
		testhelper.SeedFixtures(suite.T(), testhelper.TestDB, "testdata/post_repo")
		r := NewPostRepo(testhelper.TestDB)
		p, err := r.Get("uuid-2")

		suite.Nil(err)
		m := datatypes.JSON(`{"title": "test"}`)
		expected := &model.Post{
			UUID:     "uuid-2",
			Content:  "Post 2",
			Metadata: &m,
		}

		suite.Equal(expected, p)
	})

	suite.Run("It Returns Nil If Not Found", func() {
		r := NewPostRepo(testhelper.TestDB)
		p, err := r.Get("uuid-5")

		suite.Nil(err)
		suite.Nil(p)
	})

	suite.Run("It Handles DB Error", func() {
		r := NewPostRepo(testhelper.BrokenDB())
		p, err := r.Get("uuid-5")

		suite.NotNil(err)
		suite.Nil(p)
	})
}
