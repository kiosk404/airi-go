package airi

import (
	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/handle"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *gin.Engine) {
	root := r.Group("/", rootMw()...)
	{
		_api := root.Group("/api")
		{
			_draftbot := _api.Group("/draftbot", _draftbotMw()...)
			_draftbot.POST("/create", append(_draftbotcreateMw(), handle.DraftBotCreate)...)
		}
		{
			_playground := _api.Group("/playground_api")
			_playground_draftbot := _playground.Group("/draftbot")
			{
				_playground_draftbot.POST("/update_draft_bot_info", append(_updatedraftbotinfoagwMw(), handle.DraftBotUpdateInfo)...)
			}
		}
		{
			_conversation := _api.Group("/conversation", _conversationMw()...)
			_conversation.POST("/chat", append(_agentrunMw(), handle.AgentRun)...)
			_conversation.POST("/get_message_list", append(_getmessagelistMw(), handle.GetMessageList)...)
		}
		{
			_foundation := _api.Group("/foundation", _foundationMw()...)
			{
				_foundation_v1 := _foundation.Group("/v1", _foundationV1Mw()...)
				_foundation_v1_users := _foundation_v1.Group("/users", _userMw()...)
				{
					_foundation_v1_users.POST("/register", append(_registerMw(), handle.PassportWebRegisterPost)...)
				}
			}
		}
		root.GET("/health", Health)
	}
}

func Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
