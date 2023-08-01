package memo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTemplates(t *testing.T) {
	// LoadTemplates("../.config.toml")
	memo := FromConfig("../.config.toml")
	assert.NotNil(t, memo.Config.Templates)

	inter := memo.Config.Templates.interference

	err := inter.Execute(os.Stdout, map[string]any{
		"ShopName":  "红舷",
		"OwnerName": "Aspirin",
		"Animals": []map[string]any{
			{
				"AnimalName": "罗杰",
				"AnimalDesc": []string{
					"罗杰是一只小松鼠，住在森林东面的一棵树上。",
					"罗杰活泼好动并且是个跳远健将。",
					"罗杰有的时候缺乏耐心，但是却很聪明。",
					"罗杰经常说的口头禅是'嗨呀'、'快一点'。",
					"罗杰最喜欢的颜色是金黄色，因为那是秋天树叶的颜色。",
					"罗杰最喜欢的食物是松子、核桃。",
				},
				"Memories": []string{
					"罗杰从朋友那里听说了红舷茶馆, 他决定去那里看看。",
					"罗杰从来没有在那里喝过茶。",
				},
			},
		},
		"QueryType": "order",
		"Menu": []string{
			"桃香奶茶",
			"田野绿茶",
			"红舷Mix",
			"栗子蛋糕",
			"葡萄干蛋糕",
		},
	})
	assert.Nil(t, err)
}
