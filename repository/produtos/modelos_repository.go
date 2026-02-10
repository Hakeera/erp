package repository

import (
	"context"
	"erp/config"
	"erp/model"
)

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
			&m.ConsumoPorGrade, // JSONB → map[string]float64
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

func CriarModelo(modelo model.Modelo) error {
	return nil
}

func BuscarModelo(id int) ([]model.Modelo, error) {
	return nil, nil
}
