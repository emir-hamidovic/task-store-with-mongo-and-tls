package inmemory_test

import (
	"rest/taskstore/inmemory"
	"rest/taskstore/taskstoretest"
	"sync"
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(new(InMemTestSuite))

func Test(t *testing.T) { gc.TestingT(t) }

type InMemTestSuite struct {
	taskstoretest.SuiteBase
	inmem *inmemory.InMemory
}

func (m *InMemTestSuite) SetUpSuite(c *gc.C) {
	m.inmem = &inmemory.InMemory{Tasks: sync.Map{}, NextId: 0}
	m.SetTaskstore(m.inmem)
}

func (m *InMemTestSuite) SetUpTest(c *gc.C) {
	err := m.inmem.DeleteAll()
	c.Assert(err, gc.IsNil)
}

func (m *InMemTestSuite) TearDownSuite(c *gc.C) {
	err := m.inmem.DeleteAll()
	c.Assert(err, gc.IsNil)
}
