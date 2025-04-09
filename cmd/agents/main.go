package main

import "fmt"

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"scrapper/config"
// 	"scrapper/internal/repositories"
// 	"scrapper/internal/structs"
// 	"strings"
// 	"sync"
// )

// func main() {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("POST", "https://www.crecisc.conselho.net.br/form_pesquisa_cadastro_geral_site.php", strings.NewReader("cidade=Imbituba"))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		return
// 	}

// 	var data map[string][]structs.AgentDTO

// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		fmt.Println("Error parsing JSON:", err)
// 		return
// 	}

// 	supaClient, err := config.SupabaseClient()
// 	if err != nil {
// 		fmt.Println("Error creating Supabase client:", err)
// 	}

// 	repo := repositories.NewAgentsRepository(supaClient)

// 	var wg sync.WaitGroup

// 	for _, result := range data["cadastros"] {
// 		wg.Add(1)
// 		go func(result structs.AgentDTO) {
// 			agent := structs.Agent{
// 				Id:        result.Id,
// 				AgentId:   result.Creci,
// 				Name:      result.Nome,
// 				Type:      result.Tipo,
// 				IsRegular: result.Regular,
// 				HasPhoto:  result.Foto,
// 				Cpf:       result.Cpf,
// 				IsActive:  result.Situacao == 1,
// 				Phones:    result.Telefones,
// 			}
// 			_, err := repo.SaveOne(agent)
// 			if err != nil {
// 				fmt.Println("Error saving agent:", err)

// 			}
// 			wg.Done()
// 		}(result)
// 	}

// 	wg.Wait()

// 	fmt.Println("Saved everything")
// }

func main() {
	fmt.Println("to be defined")
}
