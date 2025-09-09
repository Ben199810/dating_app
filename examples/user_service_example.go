package main

import (
	"fmt"
	user "golang_dev_docker/domain/service"
	"golang_dev_docker/infrastructure/repository"
	"log"
)

func main() {
	// 創建儲存庫實例
	userRepo := repository.NewMemoryUserRepository()

	// 創建用戶服務
	userService := user.NewUserService(userRepo)

	// 測試創建第一個用戶
	fmt.Println("=== 測試創建用戶 ===")
	input1 := user.NewUserInput{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "SecurePass123!",
	}

	newUser1, err := userService.CreateNewUser(input1)
	if err != nil {
		log.Printf("創建用戶失敗: %v\n", err)
	} else {
		fmt.Printf("成功創建用戶: %s (ID: %s)\n", newUser1.Username, newUser1.ID)
	}

	// 測試重複用戶名
	fmt.Println("\n=== 測試重複用戶名 ===")
	input2 := user.NewUserInput{
		Username: "john_doe", // 重複的用戶名
		Email:    "john2@example.com",
		Password: "AnotherPass456!",
	}

	_, err = userService.CreateNewUser(input2)
	if err != nil {
		fmt.Printf("預期的錯誤: %v\n", err)
	} else {
		fmt.Println("意外成功創建了重複用戶名的用戶")
	}

	// 測試重複電子郵件
	fmt.Println("\n=== 測試重複電子郵件 ===")
	input3 := user.NewUserInput{
		Username: "jane_doe",
		Email:    "john@example.com", // 重複的電子郵件
		Password: "YetAnotherPass789!",
	}

	_, err = userService.CreateNewUser(input3)
	if err != nil {
		fmt.Printf("預期的錯誤: %v\n", err)
	} else {
		fmt.Println("意外成功創建了重複電子郵件的用戶")
	}

	// 測試創建第二個有效用戶
	fmt.Println("\n=== 測試創建第二個有效用戶 ===")
	input4 := user.NewUserInput{
		Username: "jane_smith",
		Email:    "jane@example.com",
		Password: "ValidPass321!",
	}

	newUser2, err := userService.CreateNewUser(input4)
	if err != nil {
		log.Printf("創建用戶失敗: %v\n", err)
	} else {
		fmt.Printf("成功創建用戶: %s (ID: %s)\n", newUser2.Username, newUser2.ID)
	}

	// 測試無效輸入
	fmt.Println("\n=== 測試無效輸入 ===")
	invalidInputs := []user.NewUserInput{
		{Username: "ab", Email: "valid@email.com", Password: "ValidPass123!"},       // 用戶名太短
		{Username: "valid_user", Email: "invalid-email", Password: "ValidPass123!"}, // 無效電子郵件
		{Username: "valid_user2", Email: "valid2@email.com", Password: "weak"},      // 弱密碼
		{Username: "", Email: "valid3@email.com", Password: "ValidPass123!"},        // 空用戶名
	}

	for i, input := range invalidInputs {
		_, err := userService.CreateNewUser(input)
		if err != nil {
			fmt.Printf("測試 %d - 預期的驗證錯誤: %v\n", i+1, err)
		} else {
			fmt.Printf("測試 %d - 意外成功創建了無效用戶\n", i+1)
		}
	}
}
