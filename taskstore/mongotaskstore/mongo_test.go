package mongotaskstore_test

import (
	"rest/taskstore/mongotaskstore"
	"rest/taskstore/taskstoretest"
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(new(MongoTestSuite))

func Test(t *testing.T) { gc.TestingT(t) }

type MongoTestSuite struct {
	taskstoretest.SuiteBase
	mng *mongotaskstore.Mongo
}

func (m *MongoTestSuite) SetUpSuite(c *gc.C) {
	mng, err := mongotaskstore.NewMongoServer("", "", "test_tasks")
	c.Assert(err, gc.IsNil)
	m.mng = mng
	m.SetTaskstore(m.mng)
}

func (m *MongoTestSuite) SetUpTest(c *gc.C) {
	err := m.mng.DeleteAll()
	c.Assert(err, gc.IsNil)
}

func (m *MongoTestSuite) TearDownSuite(c *gc.C) {
	c.Assert(m.mng.CloseMongoServer(), gc.IsNil)
}
