package psql

func (d *Postgres) GetProduct(productID string) (*map[string]string, error) {

	driverConn := d.driverConn

	dbRequest := `SELECT * FROM products WHERE product_id=$1`
	rows, err := driverConn.Query(dbRequest, productID)
	if err != nil {
		return nil, err
	}

	res := map[string]string{}
	var getProductID, getName, getUnitComp, getUnit, getAmountUnit string = "", "", "", "", ""
	rows.Next()
	if err := rows.Scan(&getProductID, &getName, &getUnitComp, &getUnit, &getAmountUnit); err != nil {
		return nil, err
	}
	res = map[string]string{
		"product_id":       getProductID,
		"name":             getName,
		"unit_composition": getUnitComp,
		"unit":             getUnit,
		"amount_unit":      getAmountUnit,
	}

	return &res, nil

}

func (d *Postgres) SetProduct(insertMap map[string]string) error {
	driverConn := d.driverConn

	dbRequest := `INSERT INTO products (product_id, name, unit_composition, unit, amount_unit) VALUES ($1, $2, $3, $4, $5)`
	_, err := driverConn.Exec(dbRequest,
		insertMap["product_id"], insertMap["name"], insertMap["unit_composition"], insertMap["unit"], insertMap["amount_unit"])

	return err
}

func (d *Postgres) GetPortions(diaryID, date, mealOrder string) (*[]map[string]string, error) {
	driverConn := d.driverConn

	dbRequest := `SELECT * FROM diary WHERE diary_id=$1 AND date=$2 AND meal_name=$3`
	rows, err := driverConn.Query(dbRequest, diaryID, date, mealOrder)
	if err != nil {
		return nil, err
	}

	res := []map[string]string{}
	var getDiaryID, getDate, getMealOrder, getProductID, getAmount string = "", "", "", "", ""
	for rows.Next() {
		if err := rows.Scan(&getDiaryID, &getDate, &getMealOrder, &getProductID, &getAmount); err != nil {
			return nil, err
		}
		res = append(res, map[string]string{
			"diary_id":   getDiaryID,
			"date":       getDate,
			"meal_order": getMealOrder,
			"product_id": getProductID,
			"amount":     getAmount,
		})
	}

	return &res, nil
}

func (d *Postgres) SetPortion(insertMap map[string]string) error {

	driverConn := d.driverConn

	dbRequest := `INSERT INTO diary (diary_id, date, meal_name, product_id, amount) VALUES ($1, $2, $3, $4, $5)`
	_, err := driverConn.Exec(dbRequest,
		insertMap["user_id"], insertMap["date"], insertMap["meal_order"], insertMap["product_id"], insertMap["amoun"])

	return err
}
