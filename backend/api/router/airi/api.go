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
		_api := root.Group("/api", _apiMw()...)
		{
			_draftbot := _api.Group("/draftbot", _draftbotMw()...)
			_draftbot.POST("/commit_check", append(_checkdraftbotcommitMw(), handle.CheckDraftBotCommit)...)
			_draftbot.POST("/create", append(_draftbotcreateMw(), handle.DraftBotCreate)...)
			_draftbot.POST("/delete", append(_deletedraftbotMw(), handle.DeleteBotDelete)...)
			_draftbot.POST("/get_display_info", append(_getdraftbotdisplayinfoMw(), handle.GetDraftBotDisplayInfo)...)
			_draftbot.POST("/update_display_info", append(_updatedraftbotdisplayinfoMw(), handle.UpdateDraftBotDisplayInfo)...)

			_draftbot.GET("/list", append(_draftbotlistMw(), handle.DraftBotList)...)
		}
		{
			_bot := _api.Group("/bot", _botMw()...)
			_bot.POST("/get_type_list", append(_gettypelistMw(), handle.GetTypeList)...)
			_bot.POST("/upload_file", append(_uploadfileMw(), handle.UploadFile)...)
		}
		{
			_playground := _api.Group("/playground_api")
			_playground_draftbot := _playground.Group("/draftbot")
			{
				_playground_draftbot.POST("/get_draft_bot_info", append(_getdraftbotinfoagwMw(), handle.GetDraftBotInfoAgw)...)
				_playground_draftbot.POST("/update_draft_bot_info", append(_updatedraftbotinfoagwMw(), handle.DraftBotUpdateInfo)...)
			}
		}
		{
			_conversation := _api.Group("/conversation", _conversationMw()...)
			_conversation.POST("/chat", append(_agentrunMw(), handle.AgentRun)...)
			_conversation.POST("/clear_message", append(_clearMw(), handle.ClearConversationHistory)...)
			_conversation.POST("/break_message", append(_breakmessageMw(), handle.BreakMessage)...)
			_conversation.POST("/delete_message", append(_deletemessageMw(), handle.DeleteMessage)...)
			_conversation.POST("/get_message_list", append(_getmessagelistMw(), handle.GetMessageList)...)
		}
		{
			_foundation := _api.Group("/foundation", _foundationMw()...)
			{
				_foundation_v1 := _foundation.Group("/v1", _foundationV1Mw()...)
				_foundation_v1_users := _foundation_v1.Group("/users", _userMw()...)
				{
					_foundation_v1_users.POST("/register", append(_registerMw(), handle.PassportWebRegisterPost)...)
					_foundation_v1_users.POST("/login", append(_loginbypasswordMw(), handle.PassportWebLoginByPasswordPost)...)
					_foundation_v1_users.POST("/logout", append(_logoutMw(), handle.PassportWebLogoutPost)...)
					_foundation_v1_users.GET("/session", append(_sessionMw(), handle.PassportAccountInfo)...)
					_foundation_v1_users.POST("/:user_id/upload_avatar/", append(_uploadavatarMw(), handle.UserUpdateAvatar)...)
				}
			}
		}
		{
			_admin := _api.Group("/admin", _adminMw()...)
			{
				{
					_model := _admin.Group("/model", _modelMw()...)
					_model.POST("/create", append(_createmodelMw(), handle.CreateModel)...)
					_model.POST("/update", append(_updatemodelMw(), handle.UpdateModel)...)
					_model.POST("/delete", append(_deletemodelMw(), handle.DeleteModel)...)
					_model.GET("/list", append(_listmodelMw(), handle.GetModelList)...)
				}
			}
		}
		{
			_v1 := root.Group("/v1", _v1Mw()...)
			_v1.GET("/conversations", append(_listconversationsapiMw(), handle.ListConversationsApi)...)
			_conversations := _v1.Group("/conversations", _conversationsMw()...)
			_conversations.POST("/create", append(_createMw(), handle.CreateConversation)...)
		}
		root.GET("/health", Health)
	}
}

func Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
