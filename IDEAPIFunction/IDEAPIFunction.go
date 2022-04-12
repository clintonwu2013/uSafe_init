package ideapifunction

import (
	utility "changingtec/usafe/Utility"

	"github.com/gin-gonic/gin"
)

func HelloWorld(c *gin.Context) {
	utility.DebuggerPrintf("test HelloWorld handler !!!!")
	c.JSON(200, gin.H{
		"message": "test usafe hello world",
	})
}
