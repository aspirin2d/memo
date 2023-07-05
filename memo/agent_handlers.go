package memo

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddMemo is a gin Handler which adds an agent to the database.
func (m *Memo) AddAgent(c *gin.Context) {
	// get agent from request body
	agent := new(Agent)
	err := c.BindJSON(agent)
	if err != nil {
		m.AbortWithError(c, NewWrapError(400, err, "can't bind json to the agent"))
		return
	}

	// add agent to database
	ctx := c.Request.Context()
	id, err := m.Agents.Add(ctx, agent)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{"inserted": id})
}

// DeleteAgent is a gin Handler which remove an agent from the database.
func (m *Memo) DeleteAgent(c *gin.Context) {
	// get agent id from url params
	aid := c.Param("aid")
	oid, err := primitive.ObjectIDFromHex(aid)
	if err != nil {
		m.AbortWithError(c, NewWrapError(400, err, "invalid agent id"))
		return
	}

	// delete agent from database
	ctx := c.Request.Context()
	err = m.Agents.Delete(ctx, oid)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{"ok": true})
}

// GetAgent is a gin Handler which get an agent from the database.
func (m *Memo) GetAgent(c *gin.Context) {
	// get agent id from url params
	aid := c.Param("aid")
	oid, err := primitive.ObjectIDFromHex(aid)
	if err != nil {
		m.AbortWithError(c, NewWrapError(400, err, "invalid agent id"))
		return
	}

	// get agent to database
	ctx := c.Request.Context()
	agent, err := m.Agents.Get(ctx, oid)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, agent)
}

// UpdateAgent is a gin Handler which update an agent from the database.
func (m *Memo) UpdateAgent(c *gin.Context) {
	// get agent from request body
	agent := new(Agent)
	err := c.BindJSON(agent)
	if err != nil {
		m.AbortWithError(c, NewWrapError(400, err, "can't bind JSON to the agent"))
		return
	}

	// update agent from database
	ctx := c.Request.Context()
	err = m.Agents.Update(ctx, agent)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{"ok": true})
}

// ListAgents is a gin Handler which list agents from the database.
func (m *Memo) ListAgents(c *gin.Context) {
	var oid primitive.ObjectID
	var err error
	// get offset from url params
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

	//list agents from database
	ctx := c.Request.Context()
	agents, err := m.Agents.List(ctx, oid)
	if err != nil {
		m.AbortWithError(c, err)
		return
	}

	c.JSON(200, agents)
}
