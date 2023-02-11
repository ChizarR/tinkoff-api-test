package main

import (
	"fmt"
	"time"
)

const (
	token = ""
)

func main() {
	api := NewTinkoffSandboxAPI(token)

	accountsNum := 5
	for i := 0; i < accountsNum; i++ {
		accountId, err := api.CreateSandboxAccount()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created account with 'id': %s\n", accountId)
		time.Sleep(1 * time.Second)
	}

	accounts, err := api.GetSandboxAccounts()
	if err != nil {
		panic(err)
	}
	fmt.Println(accounts)
	time.Sleep(1 * time.Second)

	for _, acc := range accounts {
		id := acc["id"]
		statusOk := api.CloseSandboxAccount(id.(string))
		if !statusOk {
			panic(fmt.Sprintf("Can't delete acc with id: %v'", id))
		}
		fmt.Printf("Acc with 'id': %s deleted\n", id)
		time.Sleep(1 * time.Second)
	}

	accounts, err = api.GetSandboxAccounts()
	time.Sleep(1 * time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println("accounts: ", accounts)
}
