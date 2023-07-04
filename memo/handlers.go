package memo

import "github.com/gin-gonic/gin"

// AddMemo is a gin Handler which adds a agent to the database.
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

// DelAgent is a gin Handler which remove a agent to the database.
func (m *Memo) DelAgent(c *gin.Context) {
	aid := c.Param("aid")
	// get agent from request body
	agent := new(Agent)
	err := c.BindJSON(agent)
	if err != nil {
		c.AbortWithStatusJSON(400, m.NewError(err, "can't bind json", false))
		return
	}

	// add agent to database
	ctx := c.Request.Context()
	id, err := m.Agents.Add(ctx, agent)
	if err != nil {
		c.AbortWithStatusJSON(400, m.NewError(err, "can't add agent", false))
		return
	}

	c.JSON(200, gin.H{"inserted": id})
}
