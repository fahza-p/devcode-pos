package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Conn() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}

	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DBNAME")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")

	dbUri := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true", username, password, dbHost, dbPort, dbName)

	connection, err := gorm.Open(mysql.Open(dbUri), &gorm.Config{
		SkipDefaultTransaction: true,
		// PrepareStmt:            true,
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Print(err)
	}
	DB = connection
	migration()
	return nil
}

func migration() {
	query := `
	CREATE TABLE IF NOT EXISTS cashiers (
		cashier_id int(10) unsigned NOT NULL AUTO_INCREMENT,
		cashier_name varchar(255) DEFAULT NULL,
		cashier_passcode varchar(100) DEFAULT NULL,
		cashier_token text DEFAULT NULL,
		created_at datetime DEFAULT NULL,
		updated_at datetime DEFAULT NULL,
		PRIMARY KEY (cashier_id),
		KEY idx_cashiers_cashier_passcode (cashier_passcode),
		KEY idx_cashiers_cashier_token (cashier_token(3072))
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;`

	query += `
	CREATE TABLE IF NOT EXISTS categories (
		category_id int(10) unsigned NOT NULL AUTO_INCREMENT,
		category_name varchar(255) DEFAULT NULL,
		created_at datetime DEFAULT NULL,
		updated_at datetime DEFAULT NULL,
		PRIMARY KEY (category_id)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;`

	query += `
	CREATE TABLE IF NOT EXISTS discounts (
		discount_id int(10) unsigned NOT NULL AUTO_INCREMENT,
		discount_qty int(10) unsigned DEFAULT NULL,
		discount_type varchar(20) DEFAULT NULL,
		discount_result decimal(10,0) DEFAULT NULL,
		expired_at_format varchar(100) DEFAULT NULL,
		string_format varchar(150) DEFAULT NULL,
		created_at datetime DEFAULT NULL,
		updated_at datetime DEFAULT NULL,
		expired_at datetime DEFAULT NULL,
		PRIMARY KEY (discount_id),
		KEY idx_discounts_discount_type (discount_type)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;`

	query += `
	CREATE TABLE IF NOT EXISTS order_details (
		detail_id int(10) unsigned NOT NULL AUTO_INCREMENT,
		detail_order_id int(10) unsigned DEFAULT NULL,
		detail_product_id int(10) unsigned DEFAULT NULL,
		detail_product_name varchar(255),
		detail_discount_id int(10) unsigned DEFAULT NULL,
		detail_discount_qty decimal(10,0) DEFAULT NULL,
		detail_discount_price decimal(10,0) DEFAULT NULL,
		detail_discount_final_price decimal(10,0) DEFAULT NULL,
		detail_discount_normal_price decimal(10,0) DEFAULT NULL,
		created_at datetime DEFAULT NULL,
		updated_at datetime DEFAULT NULL,
		PRIMARY KEY (detail_id),
		KEY idx_order_details_detail_product_id (detail_product_id),
		KEY idx_order_details_detail_order_id (detail_order_id),
		KEY idx_order_details_detail_discount_id (detail_discount_id)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;`

	query += `
	CREATE TABLE IF NOT EXISTS orders (
		order_id int(10) unsigned NOT NULL AUTO_INCREMENT,
		order_cashiers_id int(10) unsigned DEFAULT NULL,
		order_payment_id int(10) unsigned DEFAULT NULL,
		order_recipe_id varchar(10) DEFAULT NULL,
		order_total_price decimal(10,0) DEFAULT NULL,
		order_total_paid decimal(10,0) DEFAULT NULL,
		order_total_return decimal(10,0) DEFAULT NULL,
		order_is_download BOOL DEFAULT 0 NOT NULL,
		created_at datetime DEFAULT NULL,
		updated_at datetime DEFAULT NULL,
		PRIMARY KEY (order_id),
		KEY idx_orders_order_cashiers_id (order_cashiers_id),
		KEY idx_orders_order_payment_id (order_payment_id),
		KEY idx_orders_order_recipe_id (order_recipe_id),
		KEY idx_orders_order_is_download (order_is_download)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;`

	query += `
	CREATE TABLE IF NOT EXISTS payments (
		payment_id int(10) unsigned NOT NULL AUTO_INCREMENT,
		payment_name varchar(50) DEFAULT NULL,
		payment_type varchar(50) DEFAULT NULL,
		payment_logo varchar(500) DEFAULT NULL,
		created_at datetime DEFAULT NULL,
		updated_at datetime DEFAULT NULL,
		PRIMARY KEY (payment_id),
		KEY idx_payments_payment_type (payment_type)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;`

	query += `
	CREATE TABLE IF NOT EXISTS products (
		product_id int(10) unsigned NOT NULL AUTO_INCREMENT,
		product_category_id int(10) unsigned NOT NUll,
		product_discount_id int(10) unsigned DEFAULT NULL,
		product_sku varchar(100) DEFAULT NULL,
		product_name varchar(255) DEFAULT NULL,
		product_stock int(10) DEFAULT NULL,
		product_price decimal(10,0) unsigned DEFAULT NULL,
		product_image varchar(500) DEFAULT NULL,
		created_at datetime DEFAULT NULL,
		updated_at datetime DEFAULT NULL,
		PRIMARY KEY (product_id),
		KEY idx_products_product_sku (product_sku)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;`
	DB.Exec(query)
}
