# Documentação do Banco de Dados de Produtos - ERP

## Visão Geral

Este banco de dados gerencia o sistema de produção de um ERP para indústria têxtil/confecção, controlando modelos de peças, fichas técnicas, custos de produção, tabelas de preço e produtos finais.

## Estrutura de Tabelas

### 1. `modelos`

Armazena os modelos de peças que podem ser produzidos.

**Colunas:**

- `modelo_id` (SERIAL, PK) - Identificador único do modelo
- `nome` (TEXT, NOT NULL) - Nome do modelo
- `linha` (TEXT, NOT NULL) - Linha de produto à qual pertence
- `corte` (INT, DEFAULT 0) - Custo de corte em centavos
- `costura` (INT, DEFAULT 0) - Custo de costura em centavos
- `acabamento` (INT, DEFAULT 0) - Custo de acabamento em centavos
- `aviamento` (INT, DEFAULT 0) - Custo de aviamentos em centavos
- `consumo_por_grade` (JSONB, NOT NULL) - Consumo de tecido por tamanho (estrutura JSON)
- `descricao` (TEXT, NULLABLE) - Descrição adicional do modelo
- `created_at` (TIMESTAMP, DEFAULT now()) - Data de criação
- `updated_at` (TIMESTAMP, DEFAULT now()) - Data da última atualização

**Observações:**
- Todos os custos são armazenados em centavos (INT) para evitar problemas com decimais
- O campo `consumo_por_grade` permite flexibilidade no consumo por tamanho

---

### 2. `fichas_tecnicas`

Registra as fichas técnicas que combinam um modelo com especificações de tecido, cor e cliente.

**Colunas:**

- `ficha_id` (SERIAL, PK) - Identificador único da ficha técnica
- `modelo_id` (INT, NOT NULL, FK → modelos) - Referência ao modelo
- `tecido_ref` (TEXT, NULLABLE) - Referência do tecido
- `tecido` (TEXT, NOT NULL) - Tipo de tecido utilizado
- `cor` (TEXT, NOT NULL) - Cor do produto
- `cliente` (TEXT, NOT NULL) - Cliente para quem a ficha é destinada
- `descricao` (TEXT, NULLABLE) - Descrição adicional
- `tipo_arte` (TEXT, NULLABLE) - Tipo de arte/estampa aplicada
- `created_at` (TIMESTAMP, DEFAULT now()) - Data de criação
- `updated_at` (TIMESTAMP, DEFAULT now()) - Data da última atualização

**Constraints:**
- UNIQUE (modelo_id, tecido, cor, cliente) - Evita duplicação de fichas técnicas idênticas

---

### 3. `fichas_custos`

Armazena os custos detalhados de cada ficha técnica.

**Colunas:**

- `ficha_id` (INT, PK, FK → fichas_tecnicas) - Identificador da ficha (chave primária e estrangeira)
- `custo_modelo` (INT, NOT NULL, DEFAULT 0) - Custo do modelo (calculado automaticamente)
- `custo_tecido` (INT, NOT NULL, DEFAULT 0) - Custo do tecido em centavos
- `custo_arte` (INT, NOT NULL, DEFAULT 0) - Custo da arte/estampa em centavos
- `custos_extras` (JSONB, NOT NULL, DEFAULT '{}') - Custos adicionais em formato JSON
- `total` (INT, NOT NULL) - Custo total (calculado automaticamente)

**Observações:**
- Relação 1:1 com `fichas_tecnicas`
- Campos calculados automaticamente por triggers

---

### 4. `fichas_tabelas`

Define as tabelas de preço para cada ficha técnica com diferentes markups.

**Colunas:**

- `tabela_id` (SERIAL, PK) - Identificador único da tabela
- `ficha_id` (INT, NOT NULL, FK → fichas_tecnicas) - Referência à ficha técnica
- `nome` (TEXT, NOT NULL) - Nome da tabela de preço
- `markup` (INT, NOT NULL) - Percentual de markup (em centavos, ex: 15000 = 150%)
- `preco_override` (INT, NULLABLE) - Preço fixo que substitui o cálculo por markup

**Constraints:**
- UNIQUE (ficha_id, nome) - Uma ficha não pode ter tabelas com nomes duplicados

**Observações:**
- Permite múltiplas estratégias de precificação para a mesma ficha
- O `preco_override` tem prioridade sobre o cálculo por markup quando definido

---

### 5. `produtos`

Representa os produtos finais (SKUs) com tamanhos específicos.

**Colunas:**

- `produto_id` (SERIAL, PK) - Identificador único do produto
- `ficha_id` (INT, NOT NULL, FK → fichas_tecnicas) - Referência à ficha técnica
- `tamanho` (TEXT, NOT NULL) - Tamanho do produto (P, M, G, etc.)
- `linha` (TEXT, NOT NULL) - Linha de produto
- `situacao` (TEXT, NOT NULL, DEFAULT 'ativo') - Status do produto (ativo, inativo, etc.)
- `descricao` (TEXT, NULLABLE) - Descrição adicional
- `created_at` (TIMESTAMP, DEFAULT now()) - Data de criação

**Constraints:**
- UNIQUE (ficha_id, tamanho) - Evita duplicação de tamanhos na mesma ficha

---

## Views

### `vw_produtos_nome`

Gera nomes descritivos completos para os produtos.

**Estrutura:**
```sql
nome_produto = modelo.nome + ' - ' + tecido + ' - ' + cor + ' - ' + cliente + ' - ' + tamanho
```

**Colunas retornadas:**
- `produto_id` - ID do produto
- `nome_produto` - Nome completo formatado

**Exemplo de resultado:**
```
"Camiseta Básica - Algodão - Branco - Cliente ABC - M"
```

**SQL da View:**
```sql
CREATE VIEW vw_produtos_nome AS
SELECT
  p.produto_id,
  m.nome || ' - ' ||
  f.tecido || ' - ' ||
  f.cor || ' - ' ||
  f.cliente || ' - ' ||
  p.tamanho AS nome_produto
FROM produtos p
JOIN fichas_tecnicas f ON f.ficha_id = p.ficha_id
JOIN modelos m ON m.modelo_id = f.modelo_id;
```

---

## Triggers e Funções

### 1. `calcular_custo_modelo()`

**Trigger:** `trg_calcular_custo_modelo`  
**Evento:** BEFORE INSERT OR UPDATE em `fichas_custos`

**Função:**
- Calcula automaticamente o `custo_modelo` somando: corte + costura + acabamento + aviamento
- Busca os valores da tabela `modelos` através da relação com `fichas_tecnicas`

**Comportamento:**
```
custo_modelo = corte + costura + acabamento + aviamento
```

**Código SQL:**
```sql
CREATE OR REPLACE FUNCTION calcular_custo_modelo()
RETURNS trigger AS $$
BEGIN
    SELECT
        (m.corte + m.costura + m.acabamento + m.aviamento)
    INTO NEW.custo_modelo
    FROM modelos m
    JOIN fichas_tecnicas f ON f.modelo_id = m.modelo_id
    WHERE f.ficha_id = NEW.ficha_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_calcular_custo_modelo
BEFORE INSERT OR UPDATE ON fichas_custos
FOR EACH ROW
EXECUTE FUNCTION calcular_custo_modelo();
```

---

### 2. `calcular_total_ficha()`

**Trigger:** `trg_calcular_total_ficha`  
**Evento:** AFTER INSERT OR UPDATE em `fichas_custos`

**Função:**
- Calcula automaticamente o custo `total` da ficha técnica
- Soma todos os custos extras definidos no campo JSONB `custos_extras`

**Comportamento:**
```
extras_total = SOMA de todos os valores em custos_extras
total = custo_modelo + custo_tecido + custo_arte + extras_total
```

**Código SQL:**
```sql
CREATE OR REPLACE FUNCTION calcular_total_ficha()
RETURNS trigger AS $$
DECLARE
    extras_total INT := 0;
BEGIN
    SELECT COALESCE(SUM(value::INT), 0)
    INTO extras_total
    FROM jsonb_each_text(NEW.custos_extras);
    NEW.total :=
        NEW.custo_modelo +
        NEW.custo_tecido +
        NEW.custo_arte +
        extras_total;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_calcular_total_ficha
AFTER INSERT OR UPDATE ON fichas_custos
FOR EACH ROW
EXECUTE FUNCTION calcular_total_ficha();
```

---

## Relacionamentos

```
modelos (1) ──→ (N) fichas_tecnicas
                      │
                      ├──→ (1) fichas_custos
                      ├──→ (N) fichas_tabelas
                      └──→ (N) produtos
                                │
                                └──→ vw_produtos_nome (view)
```

**Fluxo de dados:**
1. Um **modelo** é criado com custos base de produção
2. Uma **ficha técnica** especifica modelo + tecido + cor + cliente
3. Os **custos** são calculados automaticamente para a ficha
4. **Tabelas de preço** definem diferentes markups
5. **Produtos** são criados com tamanhos específicos

---

## Convenções de Dados

### Valores Monetários
Todos os valores monetários são armazenados como **INT em centavos**:
- R$ 10,50 → 1050
- R$ 100,00 → 10000

### Markup
Percentuais são armazenados em centavos:
- 150% → 15000
- 50% → 5000

### JSONB - consumo_por_grade
Exemplo de estrutura:
```json
{
  "P": 1.2,
  "M": 1.3,
  "G": 1.4,
  "GG": 1.5
}
```

### JSONB - custos_extras
Exemplo de estrutura:
```json
{
  "embalagem": 500,
  "etiqueta": 200,
  "transporte": 1000
}
```

---

## Notas de Implementação

- **Integridade Referencial:** Todas as FKs garantem consistência dos dados
- **Auditoria:** Campos `created_at` e `updated_at` em tabelas principais
- **Unicidade:** Constraints UNIQUE evitam duplicações lógicas
- **Automatização:** Triggers calculam custos automaticamente, reduzindo erros manuais
- **Flexibilidade:** Campos JSONB permitem estruturas dinâmicas sem ALTER TABLE
