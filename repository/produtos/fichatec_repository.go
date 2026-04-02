package repository

import (
	"context"
	"erp/config"
	"erp/model"
	"errors"
	"fmt"
)

// --- CREATE ---
func CriarFicha(f model.FichaTecnica) error {

	ctx := context.Background()
	tx, err := config.GetDB().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Ficha Técnica
	var fichaID int
	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO fichas_tecnicas
		(modelo_id, tecido_ref, tecido, cor, cliente, descricao, tipo_arte)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING ficha_id
		`,
		f.ModeloID,
		f.TecidoRef,
		f.Tecido,
		f.Cor,
		f.Cliente,
		f.Descricao,
		f.TipoArte,
	).Scan(&fichaID)

	if err != nil {
		return err
	}

	// Fichas Custos
	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO fichas_custos
		(ficha_id, custo_modelo, custo_arte, custos_extras)
		VALUES ($1,$2,$3,$4)
		`,
		fichaID,
		f.Custos.CustoModelo,
		f.Custos.CustoArte,
		`{}`,
	)
	if err != nil {
		return err
	}

	// Custos por Grade
	for _, cg := range f.CustosGrade {

		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO fichas_custos_grade
			(ficha_id, grade, custo_tecido)
			VALUES ($1,$2,$3)
			`,
			fichaID,
			cg.Grade,
			cg.Custo,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// --- READ ---
func ListarFichas() ([]model.FichaTecnica, error) {

	rows, err := config.GetDB().Query(
		context.Background(),
		`
		SELECT
			f.ficha_id,
			f.modelo_id,
			COALESCE(f.tecido_ref, ''),
			f.tecido,
			f.cor,
			f.cliente,
			COALESCE(f.descricao, ''),
			COALESCE(f.tipo_arte, ''),
			COALESCE(c.custo_modelo, 0),
			COALESCE(c.custo_arte, 0),
			COALESCE(c.total, 0)
		FROM fichas_tecnicas f
		LEFT JOIN fichas_custos c ON c.ficha_id = f.ficha_id
		ORDER BY f.ficha_id DESC
		`,
	)
	if err != nil {
		fmt.Println("Erro ao obter fichas", err)
		return nil, err
	}
	defer rows.Close()

	var fichas []model.FichaTecnica

	for rows.Next() {

		var f model.FichaTecnica

		err := rows.Scan(
			&f.FichaID,
			&f.ModeloID,
			&f.TecidoRef,
			&f.Tecido,
			&f.Cor,
			&f.Cliente,
			&f.Descricao,
			&f.TipoArte,
			&f.Custos.CustoModelo,
			&f.Custos.CustoArte,
			&f.Custos.Total,
		)
		if err != nil {
			return nil, err
		}

		fichas = append(fichas, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fichas, nil
}

func BuscarFichaPorID(id int) (model.FichaTecnica, model.Modelo, error) {
	ctx := context.Background()
	db := config.GetDB()

	if id == 0 {
		return model.FichaTecnica{}, model.Modelo{}, errors.New("ficha inválida")
	}

	var f model.FichaTecnica
	f.Custos = model.FichaCustos{}
	f.CustosGrade = []model.FichaCustoGrade{}
	var m model.Modelo

	fmt.Println("[BuscarFichaPorID] Iniciando busca | ID:", id)

	// ========================
	// Query principal (ficha + modelo)
	// ========================
	fmt.Println("[BuscarFichaPorID] Executando query principal com join de modelo")
	err := db.QueryRow(ctx, `
		SELECT
			f.ficha_id,
			f.modelo_id,
			COALESCE(f.tecido_ref,''),
			COALESCE(f.tecido,''),
			COALESCE(f.cor,''),
			COALESCE(f.cliente,''),
			COALESCE(f.descricao,''),
			COALESCE(f.tipo_arte,''),
			COALESCE(m.modelo_id,0),
			COALESCE(m.nome,''),
			COALESCE(m.linha,''),
			COALESCE(m.corte,0),
			COALESCE(m.costura,0),
			COALESCE(m.acabamento,0),
			COALESCE(m.aviamento,0)
		FROM fichas_tecnicas f
		LEFT JOIN modelos m ON m.modelo_id= f.modelo_id
		WHERE f.ficha_id = $1
	`, id).Scan(
		&f.FichaID,
		&f.ModeloID,
		&f.TecidoRef,
		&f.Tecido,
		&f.Cor,
		&f.Cliente,
		&f.Descricao,
		&f.TipoArte,
		&m.ID,
		&m.Nome,
		&m.Linha,
		&m.Corte,
		&m.Costura,
		&m.Acabamento,
		&m.Aviamento,
	)
	if err != nil {
		fmt.Println("[ERRO] Falha na query principal:", err)
		return f, m, err
	}
	fmt.Println("[OK] Ficha carregada:", f)
	fmt.Println("[OK] Modelo carregado:", m)

	// ========================
	// Custos base
	// ========================
	fmt.Println("[BuscarFichaPorID] Buscando custos base")
	err = db.QueryRow(ctx, `
		SELECT COALESCE(custo_modelo,0), COALESCE(custo_arte,0), COALESCE(total,0)
		FROM fichas_custos WHERE ficha_id=$1
	`, id).Scan(&f.Custos.CustoModelo, &f.Custos.CustoArte, &f.Custos.Total)

	if err != nil {
		if err.Error() == "no rows in result set" {
			fmt.Println("[INFO] Nenhum custo base encontrado para ficha:", id)
		} else {
			fmt.Println("[ERRO] Falha ao buscar custos base:", err)
			return f, m, err
		}
	} else {
		fmt.Println("[OK] Custos base carregados:", f.Custos)
	}

	// ========================
	// Custos por grade
	// ========================
	fmt.Println("[BuscarFichaPorID] Buscando custos por grade")
	rows, err := db.Query(ctx, `
		SELECT grade, custo_tecido
		FROM fichas_custos_grade
		WHERE ficha_id=$1
	`, id)
	if err != nil {
		fmt.Println("[ERRO] Query custos grade falhou:", err)
		return f, m, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var cg model.FichaCustoGrade
		if err := rows.Scan(&cg.Grade, &cg.Custo); err != nil {
			fmt.Println("[ERRO] Scan de custos grade falhou:", err)
			return f, m, err
		}
		cg.FichaID = id
		f.CustosGrade = append(f.CustosGrade, cg)
		fmt.Println("[OK] Custo grade adicionado:", cg)
		count++
	}
	if err := rows.Err(); err != nil {
		fmt.Println("[ERRO] Iteração de rows de custos grade falhou:", err)
		return f, m, err
	}
	fmt.Println("[INFO] Total de custos por grade encontrados:", count)

	fmt.Println("[FINAL] Resultado completo da ficha:", f)
	fmt.Println("[FINAL] Modelo associado:", m)

	return f, m, nil
}

// --- UPDATE ---
func AtualizarFichatec(f model.FichaTecnica) error {

	ctx := context.Background()
	db := config.GetDB()

	fmt.Println("[AtualizarFichatec] Iniciando atualização | ID:", f.FichaID)

	tx, err := db.Begin(ctx)
	if err != nil {
		fmt.Println("[ERRO] Falha ao iniciar transação:", err)
		return err
	}
	defer tx.Rollback(ctx)

	// ========================
	// Atualiza ficha_tecnica
	// ========================
	fmt.Println("[AtualizarFichatec] Atualizando ficha_tecnica")

	_, err = tx.Exec(ctx, `
		UPDATE fichas_tecnicas SET
			tecido_ref = $1,
			tecido = $2,
			cor = $3,
			cliente = $4,
			descricao = $5,
			tipo_arte = $6,
			updated_at = NOW()
		WHERE ficha_id = $7
	`,
		f.TecidoRef,
		f.Tecido,
		f.Cor,
		f.Cliente,
		f.Descricao,
		f.TipoArte,
		f.FichaID,
	)
	if err != nil {
		fmt.Println("[ERRO] Falha ao atualizar ficha_tecnica:", err)
		return err
	}

	fmt.Println("[OK] ficha_tecnica atualizada")

	// ========================
	// Atualiza custos base
	// ========================
	fmt.Println("[AtualizarFichatec] Atualizando custos base")

	_, err = tx.Exec(ctx, `
		UPDATE fichas_custos SET
			custo_modelo = $1,
			custo_arte = $2
		WHERE ficha_id = $3
	`,
		f.Custos.CustoModelo,
		f.Custos.CustoArte,
		f.FichaID,
	)
	if err != nil {
		fmt.Println("[ERRO] Falha ao atualizar fichas_custos:", err)
		return err
	}

	fmt.Println("[OK] Custos base atualizados:", f.Custos)

	// ========================
	// Atualiza custos por grade
	// ========================
	fmt.Println("[AtualizarFichatec] Atualizando custos por grade")

	// Remove antigos
	_, err = tx.Exec(ctx, `
		DELETE FROM fichas_custos_grade
		WHERE ficha_id = $1
	`, f.FichaID)
	if err != nil {
		fmt.Println("[ERRO] Falha ao remover custos por grade:", err)
		return err
	}

	fmt.Println("[OK] Custos antigos removidos")

	// Insere novamente
	count := 0

	for _, cg := range f.CustosGrade {

		_, err = tx.Exec(ctx, `
			INSERT INTO fichas_custos_grade
			(ficha_id, grade, custo_tecido)
			VALUES ($1,$2,$3)
		`,
			f.FichaID,
			cg.Grade,
			cg.Custo,
		)
		if err != nil {
			fmt.Println("[ERRO] Falha ao inserir custo grade:", cg, "| erro:", err)
			return err
		}

		fmt.Println("[OK] Custo grade inserido:", cg)
		count++
	}

	fmt.Println("[INFO] Total de custos por grade inseridos:", count)

	// ========================
	// Commit
	// ========================
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println("[ERRO] Falha no commit:", err)
		return err
	}

	fmt.Println("[FINAL] Atualização concluída com sucesso | ID:", f.FichaID)

	return nil
}

// --- DELETE ---
func DeletarFichatec(id int) error {

	ctx := context.Background()
	tx, err := config.GetDB().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`DELETE FROM fichas_custos WHERE ficha_id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM fichas_tecnicas WHERE ficha_id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
