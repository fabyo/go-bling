package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Estruturas do BlingService
type BlingService struct {
	BaseURL     string
	AccessToken string
	Client      *http.Client
}

type Contato struct {
	ID       int    `json:"id,omitempty"`
	Nome     string `json:"nome,omitempty"`
	Tipo     string `json:"tipoPessoa,omitempty"`
	CPF_CNPJ string `json:"cpf_cnpj,omitempty"`
}

type Categoria struct {
	ID int `json:"id,omitempty"`
}

type FormaPagamento struct {
	ID int `json:"id,omitempty"`
}

type ContaReceber struct {
	DataEmissao      string           `json:"dataEmissao"`
	Vencimento       string           `json:"vencimento"`
	Valor            float64          `json:"valor"`
	Historico        string           `json:"historico,omitempty"`
	NroDocumento     string           `json:"nroDocumento,omitempty"`
	Contato          Contato          `json:"contato,omitempty"`
	Categoria        Categoria        `json:"categoria,omitempty"`
	FormaPagamento   FormaPagamento   `json:"formaPagamento,omitempty"`
	NumeroParcela    int              `json:"numeroParcela,omitempty"`
	TotalParcelas    int              `json:"totalParcelas,omitempty"`
}

func NewBlingService(accessToken string) *BlingService {
	return &BlingService{
		BaseURL:     "https://bling.com.br/Api/v3",
		AccessToken: accessToken,
		Client:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *BlingService) CreateContaReceber(conta ContaReceber) (map[string]interface{}, error) {
	url := s.BaseURL + "/contas/receber"
	
	jsonData, err := json.Marshal(conta)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados: %v", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %v", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}
	
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("erro API Bling (%d): %s", resp.StatusCode, string(body))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %v", err)
	}
	
	return result, nil
}

// Métodos auxiliares para buscar IDs
func (s *BlingService) GetContatos() (map[string]interface{}, error) {
	url := s.BaseURL + "/contatos"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Accept", "application/json")
	
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	
	return result, nil
}

func (s *BlingService) GetFormasPagamento() (map[string]interface{}, error) {
	url := s.BaseURL + "/formas-pagamentos"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Accept", "application/json")
	
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	
	return result, nil
}

func (s *BlingService) GetCategorias() (map[string]interface{}, error) {
	url := s.BaseURL + "/categorias/receitas"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Accept", "application/json")
	
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	
	return result, nil
}

func (s *BlingService) DownloadBoletoPDF(contaID int) ([]byte, error) {
	url := fmt.Sprintf("%s/contas/receber/%d/pdf", s.BaseURL, contaID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %v", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Accept", "application/pdf")
	
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro API Bling (%d): %s", resp.StatusCode, string(body))
	}
	
	// Ler o PDF
	pdfData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler PDF: %v", err)
	}
	
	// Verificar se é um PDF válido
	if len(pdfData) < 4 || string(pdfData[0:4]) != "%PDF" {
		return nil, fmt.Errorf("resposta não é um PDF válido")
	}
	
	return pdfData, nil
}

// FUNÇÃO MAIN COMPLETA
func main() {
	// 🔑 COLOQUE SEU TOKEN AQUI (aquele que você gerou anteriormente)
	accessToken := os.Getenv("BLING_TOKEN")
	
	bling := NewBlingService(accessToken)
	
	// Primeiro, vamos buscar os IDs necessários
	fmt.Println("🔍 Buscando dados do Bling...")
	
	// Buscar contatos
	contatos, err := bling.GetContatos()
	if err != nil {
		log.Printf("⚠️  Erro ao buscar contatos: %v", err)
	} else {
		fmt.Println("📋 Contatos encontrados:")
		if data, ok := contatos["data"].([]interface{}); ok {
			for i, contato := range data {
				if i < 5 { // Mostra apenas os primeiros 5
					if c, ok := contato.(map[string]interface{}); ok {
						fmt.Printf("   ID: %v, Nome: %v\n", c["id"], c["nome"])
					}
				}
			}
		}
	}
	
	// Buscar formas de pagamento
	formas, err := bling.GetFormasPagamento()
	if err != nil {
		log.Printf("⚠️  Erro ao buscar formas de pagamento: %v", err)
	} else {
		fmt.Println("💳 Formas de pagamento encontradas:")
		if data, ok := formas["data"].([]interface{}); ok {
			for i, forma := range data {
				if i < 5 {
					if f, ok := forma.(map[string]interface{}); ok {
						fmt.Printf("   ID: %v, Descrição: %v\n", f["id"], f["descricao"])
					}
				}
			}
		}
	}
	
	// Dados do boleto com IDs REAIS
	conta := ContaReceber{
		DataEmissao:    time.Now().Format("2006-01-02"), // Data atual
		Vencimento:     time.Now().AddDate(0, 1, 0).Format("2006-01-02"), // 1 mês à frente
		Valor:          150.50,
		Historico:      "Teste de boleto via API Go",
		NroDocumento:   fmt.Sprintf("GO-%d", time.Now().Unix()), // Número único com timestamp
		NumeroParcela:  1,
		TotalParcelas:  1,
		Contato: Contato{
			ID:       17751459653,    // ✅ ID REAL do contato
			Nome:     "Fabyo Guimaraes",
		},
		Categoria: Categoria{
			ID: 8422839,              // ✅ ID da categoria
		},
		FormaPagamento: FormaPagamento{
			ID: 8422840,              // ✅ ID REAL do Boleto
		},
	}
	
	fmt.Printf("\n🎯 Tentando criar boleto...\n")
	
	result, err := bling.CreateContaReceber(conta)
	if err != nil {
		log.Fatalf("❌ Erro ao criar boleto: %v", err)
	}
	
	fmt.Println("✅ Boleto criado com sucesso!")
	fmt.Printf("📦 Resposta: %+v\n", result)
	
	// Extrair ID da conta criada
	var contaID int
	if data, ok := result["data"].(map[string]interface{}); ok {
		if id, exists := data["id"]; exists {
			// Converter para int
			contaID = int(id.(float64))
			fmt.Printf("🔢 ID da conta criada: %d\n", contaID)
		}
	}
	
	// 🔽 BAIXAR O PDF DO BOLETO
		// 🔽 BAIXAR O PDF DO BOLETO COM DELAY
	if contaID > 0 {
		fmt.Printf("\n⏳ Aguardando processamento do boleto...\n")
		
		// Aguardar 5 segundos para o boleto ser processado
		time.Sleep(5 * time.Second)
		
		fmt.Printf("📥 Baixando PDF do boleto...\n")
		
		pdfData, err := bling.DownloadBoletoPDF(contaID)
		if err != nil {
			log.Printf("⚠️  Erro ao baixar PDF (tentativa 1): %v", err)
			
			// Segunda tentativa após mais 5 segundos
			fmt.Printf("⏳ Nova tentativa em 5 segundos...\n")
			time.Sleep(5 * time.Second)
			
			pdfData, err = bling.DownloadBoletoPDF(contaID)
			if err != nil {
				log.Printf("❌ Erro ao baixar PDF (tentativa 2): %v", err)
				fmt.Printf("💡 O boleto foi criado com sucesso (ID: %d), mas o PDF ainda não está disponível.\n", contaID)
				fmt.Printf("💡 Tente baixar manualmente mais tarde pelo painel do Bling.\n")
			} else {
				// Salvar o PDF em arquivo
				filename := fmt.Sprintf("boleto_%d.pdf", contaID)
				err = os.WriteFile(filename, pdfData, 0644)
				if err != nil {
					log.Fatalf("❌ Erro ao salvar PDF: %v", err)
				}
				
				fmt.Printf("💾 PDF salvo como: %s\n", filename)
				fmt.Printf("📄 Tamanho do PDF: %d bytes\n", len(pdfData))
				fmt.Println("🎉 Boleto gerado e salvo com sucesso!")
			}
		} else {
			// Salvar o PDF em arquivo (primeira tentativa bem-sucedida)
			filename := fmt.Sprintf("boleto_%d.pdf", contaID)
			err = os.WriteFile(filename, pdfData, 0644)
			if err != nil {
				log.Fatalf("❌ Erro ao salvar PDF: %v", err)
			}
			
			fmt.Printf("💾 PDF salvo como: %s\n", filename)
			fmt.Printf("📄 Tamanho do PDF: %d bytes\n", len(pdfData))
			fmt.Println("🎉 Boleto gerado e salvo com sucesso!")
		}
	} else {
		fmt.Println("⚠️  Não foi possível obter o ID da conta para baixar o PDF")
	}
}