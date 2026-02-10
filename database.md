# Modelos
modelo_id id
nome    string
linha   string
corte   int
costura   int
acabamento  int
aviamento   int
consumo_por_grade  {"infantil":, "juvenil":, "adulto":, "extra":} 
descricao   string

# Ficha Tecnica
ficha_id
modelo_id
modelo  string
tecido_ref   string
tecido  string
cor string
cliente     string
descrição   string
tipo_arte   string
custo_arte  int
custo_modelo       int
custo {}
tecido {}
custo {}
total {}
markup {}
tabela1 {}
tabela2 {}

# Produto
produto_id
modelo  string
tecido string
cor     string
cliente string
tamanho     string
linha       string
situacao    string
tabela1     int
tabela2     int
descricao   string
nome    "modelo - tecido - cor - cliente - tamanho"
ficha_tecnica   ficha_tecnica_id
