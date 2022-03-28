package taskstoretest

import (
	"rest/taskstore"

	gc "gopkg.in/check.v1"
)

type SuiteBase struct {
	t taskstore.Taskstore
}

func (s *SuiteBase) SetTaskstore(t taskstore.Taskstore) {
	s.t = t
}

func (s *SuiteBase) TestCreateTask(c *gc.C) {
	testTask := taskstore.Task{Text: "first test", Tags: []string{"tag1", "tag2"}, Due: "2016-01-02T15:04:05+00:00"}
	obtainedID, err := s.t.CreateTask(testTask.Text, testTask.Tags, testTask.Due)
	c.Assert(err, gc.IsNil)

	obtainedTask, err := s.t.GetTaskById(obtainedID)
	c.Assert(err, gc.IsNil)
	c.Assert(obtainedTask.Id, gc.Equals, obtainedID)
}

func (s *SuiteBase) TestGetTaskById(c *gc.C) {
	testTask := taskstore.Task{Text: "first test", Tags: []string{"tag1", "tag2"}, Due: "2016-01-02T15:04:05+00:00"}
	obtainedID, err := s.t.CreateTask(testTask.Text, testTask.Tags, testTask.Due)
	c.Assert(err, gc.IsNil)

	task, err := s.t.GetTaskById(obtainedID)
	c.Assert(err, gc.IsNil)
	testTask.Id = obtainedID
	c.Assert(task, gc.DeepEquals, testTask)

	_, err = s.t.GetTaskById("")
	c.Assert(err, gc.Not(gc.IsNil))
}

func (s *SuiteBase) TestGetAllTasks(c *gc.C) {
	testTask := []taskstore.Task{
		{Text: "first test", Tags: []string{"tag1", "tag2"}, Due: "2016-01-02T15:04:05+00:00"},
		{Text: "second test", Tags: []string{"tag3", "tag4"}, Due: "2017-01-02T15:04:05+00:00"},
		{Text: "third test", Tags: []string{"tag5", "tag6"}, Due: "2018-01-02T15:04:05+00:00"},
		{Text: "fourth test", Tags: []string{"tag7", "tag8"}, Due: "2019-01-02T15:04:05+00:00"},
	}
	for _, task := range testTask {
		_, err := s.t.CreateTask(task.Text, task.Tags, task.Due)
		c.Assert(err, gc.IsNil)
	}

	tasks, err := s.t.GetAllTasks()
	c.Assert(err, gc.IsNil)
	c.Assert(len(tasks), gc.Not(gc.Equals), 0)
	c.Assert(len(tasks), gc.Equals, 4)
}

func (s *SuiteBase) TestDeleteTask(c *gc.C) {
	testTask := taskstore.Task{Text: "first test", Tags: []string{"tag1", "tag2"}, Due: "2016-01-02T15:04:05+00:00"}
	obtainedID, err := s.t.CreateTask(testTask.Text, testTask.Tags, testTask.Due)
	c.Assert(err, gc.IsNil)

	err = s.t.DeleteTask(obtainedID)
	c.Assert(err, gc.IsNil)

	err = s.t.DeleteTask("")
	c.Assert(err, gc.Not(gc.IsNil))
}

func (s *SuiteBase) TestDeleteAll(c *gc.C) {
	testTask := []taskstore.Task{
		{Text: "first test", Tags: []string{"tag1", "tag2"}, Due: "2016-01-02T15:04:05+00:00"},
		{Text: "second test", Tags: []string{"tag3", "tag4"}, Due: "2017-01-02T15:04:05+00:00"},
		{Text: "third test", Tags: []string{"tag5", "tag6"}, Due: "2018-01-02T15:04:05+00:00"},
		{Text: "fourth test", Tags: []string{"tag7", "tag8"}, Due: "2019-01-02T15:04:05+00:00"},
	}
	for _, task := range testTask {
		_, err := s.t.CreateTask(task.Text, task.Tags, task.Due)
		c.Assert(err, gc.IsNil)
	}

	err := s.t.DeleteAll()
	c.Assert(err, gc.IsNil)

	tasks, err := s.t.GetAllTasks()
	c.Assert(err, gc.IsNil)
	c.Assert(len(tasks), gc.Equals, 0)
}
