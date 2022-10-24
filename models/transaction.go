package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Id        int     `form:"id" json:"id" validate:"required"`
	IdUser    int     `form:"iduser" json:"iduser" validate:"required"`
	IdProduck int     `form:"idproduck" json:"idproduck" validate:"required"`
	Name      string  `form:"name" json:"name" validate:"required"`
	Image     string  `form:"image" json:"image" validate:"required"`
	Quantity  int     `form:"quantity" json:"quantity" validate:"required"`
	Price     float32 `form:"price" json:"price" validate:"required"`
	Total     float32 `form:"total" json:"total" validate:"required"`
	Status    string  `form:"status" json:"status" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// CRUD
func CreateTransaction(db *gorm.DB, newTransaction *Transaction) (err error) {
	err = db.Create(newTransaction).Error
	if err != nil {
		return err
	}
	return nil
}
func ReadTransaction(db *gorm.DB, transactions *[]Transaction) (err error) {
	err = db.Find(transactions).Error
	if err != nil {
		return err
	}
	return nil
}
func ReadTransactionById(db *gorm.DB, transaction *Transaction, id int) (err error) {
	err = db.Where("id=?", id).First(transaction).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateTransaction(db *gorm.DB, transaction *Transaction) (err error) {
	db.Save(transaction)

	return nil
}

func DeleteTransactionById(db *gorm.DB, transaction *Transaction, id int) (err error) {
	db.Where("id=?", id).Delete(transaction)

	return nil
}
