
CREATE TABLE modelos (
    modelo_id SERIAL PRIMARY KEY,
    nome TEXT NOT NULL,
    linha TEXT NOT NULL,

    corte INT NOT NULL DEFAULT 0,
    costura INT NOT NULL DEFAULT 0,
    acabamento INT NOT NULL DEFAULT 0,
    aviamento INT NOT NULL DEFAULT 0,

    consumo_por_grade JSONB NOT NULL,
    descricao TEXT,

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE fichas_tecnicas (
    ficha_id SERIAL PRIMARY KEY,
    modelo_id INT NOT NULL REFERENCES modelos(modelo_id),

    tecido_ref TEXT,
    tecido TEXT NOT NULL,
    cor TEXT NOT NULL,
    cliente TEXT NOT NULL,

    descricao TEXT,
    tipo_arte TEXT,

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),

    UNIQUE (modelo_id, tecido, cor, cliente)
);


CREATE TABLE fichas_custos (
    ficha_id INT PRIMARY KEY REFERENCES fichas_tecnicas(ficha_id),

    custo_modelo INT NOT NULL DEFAULT 0,
    custo_tecido INT NOT NULL DEFAULT 0,
    custo_arte INT NOT NULL DEFAULT 0,

    custos_extras JSONB NOT NULL DEFAULT '{}'::jsonb,

    total INT NOT NULL
);

CREATE TABLE fichas_tabelas (
    tabela_id SERIAL PRIMARY KEY,
    ficha_id INT NOT NULL REFERENCES fichas_tecnicas(ficha_id),

    nome TEXT NOT NULL,
    markup INT NOT NULL,

    preco_override INT,

    UNIQUE (ficha_id, nome)
);

CREATE TABLE produtos (
    produto_id SERIAL PRIMARY KEY,
    ficha_id INT NOT NULL REFERENCES fichas_tecnicas(ficha_id),

    tamanho TEXT NOT NULL,
    linha TEXT NOT NULL,
    situacao TEXT NOT NULL DEFAULT 'ativo',

    descricao TEXT,
    created_at TIMESTAMP DEFAULT now(),

    UNIQUE (ficha_id, tamanho)
);

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
BEFORE INSERT OR UPDATE ON fichas_custos
FOR EACH ROW
EXECUTE FUNCTION calcular_total_ficha();
