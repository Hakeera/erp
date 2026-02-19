package repository

import (
	"context"
	"erp/config"
	"erp/model"
)

// --- CREATE ---
func CriarFicha(f model.FichaTecnica) error {

	ctx := context.Background()
	tx, err := config.GetDB().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

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

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO fichas_custos
			(ficha_id, custo_modelo, custo_tecido, custo_arte, total)
		VALUES ($1,$2,$3,$4,$5)
		`,
		fichaID,
		f.Custos.CustoModelo,
		f.Custos.CustoTecido,
		f.Custos.CustoArte,
		f.Custos.Total,
	)
	if err != nil {
		return err
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
			COALESCE(c.custo_tecido, 0),
			COALESCE(c.custo_arte, 0),
			COALESCE(c.total, 0)
		FROM fichas_tecnicas f
		LEFT JOIN fichas_custos c ON c.ficha_id = f.ficha_id
		ORDER BY f.ficha_id DESC
		`,
	)
	if err != nil {
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
			&f.Custos.CustoTecido,
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

func BuscarFichaPorID(id int) (model.FichaTecnica, error) {

	row := config.GetDB().QueryRow(
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
			COALESCE(c.custo_tecido, 0),
			COALESCE(c.custo_arte, 0),
			COALESCE(c.total, 0)
		FROM fichas_tecnicas f
		LEFT JOIN fichas_custos c ON c.ficha_id = f.ficha_id
		WHERE f.ficha_id = $1
		`,
		id,
	)

	var f model.FichaTecnica
	f.Custos = model.FichaCustos{} // inicializa

	err := row.Scan(
		&f.FichaID,
		&f.ModeloID,
		&f.TecidoRef,
		&f.Tecido,
		&f.Cor,
		&f.Cliente,
		&f.Descricao,
		&f.TipoArte,
		&f.Custos.CustoModelo,
		&f.Custos.CustoTecido,
		&f.Custos.CustoArte,
		&f.Custos.Total,
	)

	return f, err
}

// --- UPDATE ---
func AtualizarFichatec(f model.FichaTecnica) error {

	ctx := context.Background()
	tx, err := config.GetDB().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Atualiza ficha_tecnica
	_, err = tx.Exec(ctx, `
		UPDATE fichas_tecnicas SET
			modelo_id = $1,
			tecido_ref = $2,
			tecido = $3,
			cor = $4,
			cliente = $5,
			descricao = $6,
			tipo_arte = $7,
			updated_at = NOW()
		WHERE ficha_id = $8
	`,
		f.ModeloID,
		f.TecidoRef,
		f.Tecido,
		f.Cor,
		f.Cliente,
		f.Descricao,
		f.TipoArte,
		f.FichaID,
	)
	if err != nil {
		return err
	}

	// Atualiza custos
	_, err = tx.Exec(ctx, `
		UPDATE fichas_custos SET
			custo_modelo = $1,
			custo_tecido = $2,
			custo_arte = $3,
			total = $4
		WHERE ficha_id = $5
	`,
		f.Custos.CustoModelo,
		f.Custos.CustoTecido,
		f.Custos.CustoArte,
		f.Custos.Total,
		f.FichaID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
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
