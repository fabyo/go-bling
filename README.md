# Go + Bling API v3 – Geração de Boleto (Conta a Receber + PDF)

<img src="go-bling.png" alt="Golang" width="200" />

Projeto em Go que integra com a **API v3 do Bling** para:

- Buscar **contatos**, **formas de pagamento** e **categorias de receitas**
- Criar uma **conta a receber** (base para boleto)
- Fazer **download do PDF** do boleto gerado

Tudo usando **OAuth2 (Bearer token)** da API v3

---

## 🚀 O que esse projeto faz

Fluxo completo:

1. Lê o **access token** da API v3 do Bling (`BLING_TOKEN`)
2. Consulta:
   - `GET /contatos`
   - `GET /formas-pagamentos`
   - `GET /categorias/receitas`
3. Cria uma **conta a receber**:
   - `POST /contas/receber`
4. Pega o **ID** da conta criada
5. Espera alguns segundos (processamento do boleto)
6. Faz o download do **PDF do boleto**:
   - `GET /contas/receber/{id}/pdf`
7. Salva o arquivo localmente como `boleto_<ID>.pdf`

É um exemplo completo de **cliente ERP em Go** com:

- Consumo de API REST v3
- Bearer token
- Estrutura de serviço (`BlingService`)
- Tratamento de erro, timeouts e validação de PDF.

---

## 🧱 Estrutura básica do código

Tudo está em um arquivo único `main.go` (modelo didático), com:

- `type BlingService`
  - guarda `BaseURL`, `AccessToken` e `*http.Client`
  - métodos:
    - `CreateContaReceber(conta ContaReceber)`
    - `GetContatos()`
    - `GetFormasPagamento()`
    - `GetCategorias()`
    - `DownloadBoletoPDF(contaID int)`

- `type ContaReceber`
  - representa os campos necessários para criar uma conta a receber

- `type Contato`, `type Categoria`, `type FormaPagamento`
  - usados como subestruturas dentro de `ContaReceber`

No `main()` o fluxo é:

1. Instancia o service com o token
2. Lista alguns contatos, formas de pagamento, categorias (para debug)
3. Monta uma `ContaReceber` com:
   - `DataEmissao`
   - `Vencimento`
   - `Valor`
   - `Historico`
   - `NroDocumento`
   - `Contato{ID: ...}`
   - `Categoria{ID: ...}`
   - `FormaPagamento{ID: ...}`
4. Chama `CreateContaReceber`
5. Extrai o `ID` retornado
6. Chama `DownloadBoletoPDF(ID)` com delay + retry
7. Salva o PDF em disco

---

## ✅ Pré-requisitos

- Go **1.20+** (recomendado)
- Conta no **Bling** com acesso à **API v3**
- Um **Aplicativo** criado no painel do Bling (OAuth2)
- Um **access token** v3 válido (tipo *Bearer*)

---

## 🔐 Configurando o token da API v3 (visão geral)

1. No Bling, crie um **Aplicativo** (privado é o ideal para uso interno)
2. Anote:
   - `client_id`
   - `client_secret`
   - `redirect_uri` configurada
3. Rode o fluxo OAuth 2.0 (authorization code):
   - Acesse a URL de autorização com:
     - `response_type=code`
     - `client_id`
     - `redirect_uri`
   - Autorize o aplicativo
   - Copie o `code` da URL de retorno
   - Troque o `code` por `access_token` no endpoint `/oauth/token`
4. Guarde o `access_token` gerado

No projeto, você **não** deixa o token hardcoded: use variável de ambiente.

---

## ⚙️ Configuração do projeto

Clone o repositório:

```bash
git clone https://github.com/fabyo/go-bling-boleto.git
cd go-bling-boleto
```

Configure o token (PowerShell / Windows):

```powershell
$env:BLING_TOKEN = "SEU_ACCESS_TOKEN_AQUI"
```

Ou em Linux/macOS:

```bash
export BLING_TOKEN="SEU_ACCESS_TOKEN_AQUI"
```

---

## ▶️ Como rodar

### Rodar direto com `go run`:

```bash
go run main.go
```

O que o programa faz:

1. Lê `BLING_TOKEN`
2. Lista:
   - primeiros contatos (`GET /contatos`)
   - primeiras formas de pagamento (`GET /formas-pagamentos`)
   - primeiras categorias de receitas (`GET /categorias/receitas`)
3. Cria uma **conta a receber** com dados de exemplo (IDs que você deve ajustar para os seus):
4. Tenta baixar o **PDF** do boleto 2 vezes:
   - espera 5 segundos
   - tenta baixar
   - se falhar, espera mais 5 segundos e tenta de novo
5. Salva o arquivo:

```text
boleto_<ID>.pdf
```

no diretório atual.

---

## 🧩 Ajustando para o seu ambiente

No código, na parte que monta a `ContaReceber`:

```go
conta := ContaReceber{
    DataEmissao:    time.Now().Format("2006-01-02"),
    Vencimento:     time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
    Valor:          150.50,
    Historico:      "Teste de boleto via API Go",
    NroDocumento:   fmt.Sprintf("GO-%d", time.Now().Unix()),
    NumeroParcela:  1,
    TotalParcelas:  1,
    Contato: Contato{
        ID:   17751459653,  // ✅ ID real do contato no seu Bling
        Nome: "Fabyo Guimaraes",
    },
    Categoria: Categoria{
        ID: 8422839,        // ✅ ID real da categoria de receita
    },
    FormaPagamento: FormaPagamento{
        ID: 8422840,        // ✅ ID real da forma de pagamento (boleto)
    },
}
```

Você deve:

- Ajustar os **IDs** usando as respostas de:
  - `GetContatos()`
  - `GetFormasPagamento()`
  - `GetCategorias()`
- Ajustar `Historico`, `Valor`, datas etc. conforme o cenário.

---

## 📎 Sobre o download do PDF

A função `DownloadBoletoPDF`:

- Chama: `GET /contas/receber/{id}/pdf`
- Usa `Accept: application/pdf`
- Valida se o retorno começa com `%PDF`
- Salva o conteúdo em um arquivo `.pdf`

O código faz duas tentativas com delay porque às vezes o Bling ainda está processando o boleto logo após a criação da conta.

---

## ⚠️ Cuidados importantes

- **NUNCA** commite seu `access_token` no Git
- Use sempre `BLING_TOKEN` via variável de ambiente
- IDs de contato/categoria/forma de pagamento do exemplo são **seus**, não vão existir em outra conta
  - para alguém usar esse projeto, terá que trocar esses IDs pelos próprios

---

## 💡 Ideias de evolução

- Transformar esse `main.go` em:
  - um **CLI** (`cobra` / flags) que gera boleto com parâmetros
  - uma **API HTTP** em Go (`/meu-sistema/boletos`) que recebe JSON e gera o boleto no Bling
- Tipar as respostas de:
  - `/contatos`
  - `/formas-pagamentos`
  - `/categorias/receitas`
- Adicionar logs estruturados (JSON) para rodar em produção

---

## 🧾 Resumo

Esse projeto mostra, na prática, como:

- Usar **Go** + **API v3 do Bling**
- Autenticar com **Bearer token (OAuth2)**
- Consumir múltiplos endpoints REST
- Criar **conta a receber** e baixar o **PDF do boleto**

Perfeito para colocar em portfólio como exemplo de integração com ERP em Go.
