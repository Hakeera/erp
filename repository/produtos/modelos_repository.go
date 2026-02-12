package repository

import (
	"context"
	"erp/config"
	"erp/model"
)

// --- CREATE ---
func CriarModelo(m model.Modelo) error {

	_, err := config.GetDB().Exec(
		context.Background(),
		`
		INSERT INTO modelos
		(nome, linha, corte, costura, acabamento, aviamento, consumo_por_grade, descricao)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`,
		m.Nome,
		m.Linha,
		m.Corte,
		m.Costura,
		m.Acabamento,
		m.Aviamento,
		m.ConsumoPorGrade,
		m.Descricao,
	)

	return err
}

// --- READ ---
func ListarModelos() ([]model.Modelo, error) {

	rows, err := config.GetDB().Query(
		context.Background(),
		`
		SELECT
			modelo_id,
			nome,
			linha,
			corte,
			costura,
			acabamento,
			aviamento,
			consumo_por_grade,
			COALESCE(descricao, '')
		FROM modelos
		ORDER BY nome
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modelos []model.Modelo

	for rows.Next() {
		var m model.Modelo

		if err := rows.Scan(
			&m.ID,
			&m.Nome,
			&m.Linha,
			&m.Corte,
			&m.Costura,
			&m.Acabamento,
			&m.Aviamento,
			&m.ConsumoPorGrade,
			&m.Descricao,
		); err != nil {
			return nil, err
		}

		modelos = append(modelos, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modelos, nil
}

// --- UPDATE ---
func BuscarModeloPorID(id int) (model.Modelo, error) {

	row := config.GetDB().QueryRow(
		context.Background(),
		`
		SELECT
			modelo_id,
			nome,
			linha,
			corte,
			costura,
			acabamento,
			aviamento,
			consumo_por_grade,
			COALESCE(descricao, '')
		FROM modelos
		WHERE modelo_id = $1
		`,
		id,
	)

	var m model.Modelo

	err := row.Scan(
		&m.ID,
		&m.Nome,
		&m.Linha,
		&m.Corte,
		&m.Costura,
		&m.Acabamento,
		&m.Aviamento,
		&m.ConsumoPorGrade,
		&m.Descricao,
	)

	return m, err
}

// --- UPDATE ---
func AtualizarModelo(m model.Modelo) error {

	_, err := config.GetDB().Exec(
		context.Background(),
		`
		UPDATE modelos
		SET
			nome = $1,
			linha = $2,
			corte = $3,
			costura = $4,
			acabamento = $5,
			aviamento = $6,
			consumo_por_grade = $7,
			descricao = $8,
			updated_at = now()
		WHERE modelo_id = $9
		`,
		m.Nome,
		m.Linha,
		m.Corte,
		m.Costura,
		m.Acabamento,
		m.Aviamento,
		m.ConsumoPorGrade,
		m.Descricao,
		m.ID,
	)

	return err
}

// --- DELETE ---
func DeletarModelo(id int) error {

	_, err := config.GetDB().Exec(
		context.Background(),
		`DELETE FROM modelos WHERE modelo_id = $1`,
		id,
	)

	return err
}
