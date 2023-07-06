package memo

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *Memo) GetAgentId(c *gin.Context) {
	str := c.Param("aid")
	if str == "" {
		m.AbortWithError(c, NewWrapError(400, fmt.Errorf("agent id shoud not be empty"), ""))
		return
	}

	aid, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		m.AbortWithError(c, NewWrapError(400, fmt.Errorf("agent id invalid"), ""))
		return
	}

	c.Set("agent", aid)
	c.Next()
}

func (m *Memo) AddMemories(c *gin.Context) {
	aid, _ := c.Get("agent")
	agent := aid.(primitive.ObjectID)

	var memories []*Memory
	err := c.ShouldBindJSON(&memories)
	if err != nil {
		m.AbortWithError(c, NewWrapError(400, err, "can't bind memories"))
		return
	}

	ctx := c.Request.Context()
	ids, err := m.Memories.AddMany(ctx, agent, memories)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{"inserted": ids})
}

func (m *Memo) GetMemories(c *gin.Context) {
	aid, _ := c.Get("agent")
	agent := aid.(primitive.ObjectID)

	var ids []primitive.ObjectID
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		m.AbortWithError(c, NewWrapError(400, err, "can't bind memories"))
		return
	}

	ctx := c.Request.Context()
	memories, err := m.Memories.GetMany(ctx, agent, ids)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, memories)
}

func (m *Memo) DeleteMemories(c *gin.Context) {
	aid, _ := c.Get("agent")
	agent := aid.(primitive.ObjectID)

	var ids []primitive.ObjectID
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		m.AbortWithError(c, NewWrapError(400, err, "can't bind memory ids"))
		return
	}

	ctx := c.Request.Context()
	err = m.Memories.DeleteMany(ctx, agent, ids)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{"ok": true})
}

func (m *Memo) UpdateMemories(c *gin.Context) {
	aid, _ := c.Get("agent")
	agent := aid.(primitive.ObjectID)

	var memories []*Memory
	err := c.ShouldBindJSON(&memories)
	if err != nil {
		m.AbortWithError(c, NewWrapError(400, err, "can't bind memories"))
		return
	}

	ctx := c.Request.Context()
	err = m.Memories.UpdateMany(ctx, agent, memories)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{"ok": true})
}

func (m *Memo) ListMemories(c *gin.Context) {
	aid, _ := c.Get("agent")
	agent := aid.(primitive.ObjectID)

	// get offset from url params
	var oid primitive.ObjectID
	var err error
	offset := c.Param("offset")
	if offset != "" && offset != "nil" && offset != "-1" {
		oid, err = primitive.ObjectIDFromHex(offset)
		if err != nil {
			m.AbortWithError(c, NewWrapError(400, err, "invalid offset id"))
			return
		}
	} else {
		oid = primitive.NilObjectID
	}

	ctx := c.Request.Context()
	memories, err := m.Memories.List(ctx, agent, oid)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, memories)
}
