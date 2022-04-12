package ideapifunction

import "github.com/gin-gonic/gin"

type APIStruct struct {
	Name        string
	URLPath     string
	Method      string
	HandlerFunc func(*gin.Context)
}

func GetAPIs() []APIStruct {
	return []APIStruct{
		{
			Name:        "Hello_World",
			URLPath:     "/",
			Method:      "GET",
			HandlerFunc: HelloWorld,
		},
		{
			Name:        "Hello_World1",
			URLPath:     "/test1",
			Method:      "GET",
			HandlerFunc: HelloWorld,
		},
		{
			Name:        "Hello_World2",
			URLPath:     "/test2",
			Method:      "POST",
			HandlerFunc: HelloWorld,
		},
		{
			Name:        "Hello_World3",
			URLPath:     "/test3",
			Method:      "PUT",
			HandlerFunc: HelloWorld,
		},
		{
			Name:        "Hello_World4",
			URLPath:     "/test4",
			Method:      "DELETE",
			HandlerFunc: HelloWorld,
		},
	}
}
