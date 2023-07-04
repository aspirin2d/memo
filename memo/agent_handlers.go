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
		c.AbortWithStatusJSON(400, m.NewError(err, "can't bind json for the agent", false))
		return
	}

	// add agent to database
	ctx := c.Request.Context()
	id, err := m.Agents.Add(ctx, agent)
	if err != nil {
		c.AbortWithStatusJSON(500, m.NewError(err, "can't add agent", false))
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
		c.AbortWithStatusJSON(400, m.NewError(err, "invalid agent id", false))
		return
	}

	// delete agent from database
	ctx := c.Request.Context()
	err = m.Agents.Delete(ctx, oid)
	if err != nil {
		c.AbortWithStatusJSON(400, m.NewError(err, "can't delete agent", false))
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
		c.AbortWithStatusJSON(400, m.NewError(err, "invalid agent id", false))
		return
	}

	// get agent to database
	ctx := c.Request.Context()
	agent, err := m.Agents.Get(ctx, oid)
	if err != nil {
		c.AbortWithStatusJSON(400, m.NewError(err, "can't delete agent", false))
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
		c.AbortWithStatusJSON(400, m.NewError(err, "can't bind json for the agent", false))
		return
	}

	// update agent from database
	ctx := c.Request.Context()
	err = m.Agents.Update(ctx, agent)
	if err != nil {
		c.AbortWithStatusJSON(400, m.NewError(err, "can't delete agent", false))
		return
	}

	c.JSON(200, gin.H{"ok": true})
}

// ListAgents is a gin Handler which list agents from the database.
func (m *Memo) ListAgents(c *gin.Context) {
	// get offset from url params
	offset := c.Param("offset")
	// if offset is "", it will return NilObjectID
	oid, _ := primitive.ObjectIDFromHex(offset)

	//list agents from database
	ctx := c.Request.Context()
	agents, err := m.Agents.List(ctx, oid)
	if err != nil {
		c.AbortWithStatusJSON(400, m.NewError(err, "can't list agents", false))
		return
	}

	c.JSON(200, agents)
}
