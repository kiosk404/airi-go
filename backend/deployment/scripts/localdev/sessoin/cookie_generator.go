package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/service"
)

func main() {
	// 创建一个测试用户的会话
	session := &entity.Session{
		UserID:    7572241179807318016, // 测试用户ID
		SessionID: 7572241243371995136, // 测试会话ID
		Locale:    "zh-CN",
		CreatedAt: time.Now(),
	}

	// 生成会话密钥
	sessionSvc := service.NewSessionService()
	sessionKey, err := sessionSvc.GenerateSessionKey(context.Background(), session)
	if err != nil {
		fmt.Printf("生成会话密钥失败: %v\n", err)
		return
	}

	fmt.Println("生成的测试Cookie:")
	fmt.Printf("Cookie名称: %s\n", entity.SessionKey)
	fmt.Printf("Cookie值: %s\n", sessionKey)
	fmt.Printf("完整Cookie头: %s=%s\n", entity.SessionKey, sessionKey)

	// 验证生成的会话密钥
	verifiedSession, err := sessionSvc.ValidateSession(context.Background(), sessionKey)
	if err != nil {
		fmt.Printf("验证会话密钥失败: %v\n", err)
		return
	}

	fmt.Printf("\n验证成功的会话信息:\n")
	fmt.Printf("用户ID: %d\n", verifiedSession.UserID)
	fmt.Printf("会话ID: %d\n", verifiedSession.SessionID)
	fmt.Printf("语言: %s\n", verifiedSession.Locale)
	fmt.Printf("创建时间: %v\n", verifiedSession.CreatedAt)
	fmt.Printf("过期时间: %v\n", verifiedSession.ExpiresAt)
}
