# Skema Database Ingredient Conversion Unit

Table ini untuk mengkonversi ukuran unit dari bahan baku

| Kolom         | Tipe Data     | Deskripsi                       |
|-------------  |---------------|---------------------------------|
| id            | INT           | Primary Key, Auto Increment     |
| serial        | VARCHAR(11)   | Unique Serial untuk bahan baku  |
| ingredient_id | INT           | FK, reference [ingredients.id](03-ingredient.md) |
| unit_id       | INT           | FK, reference [unit.id](02-unit.md) |
| value         | DECIMAL(10,2) | Nilai konversi ke unit ini      |
| skip_calculate| TINYINT(1)    | Jika menggunakan konversi ini, penghitungan harga di abaikan |
| created_at    | TIMESTAMP     | Tanggal penambahan unit         |
| created_by    | VARCHAR(30)   | Username [users.username](01-user.md) yang menambahkan|
| updated_at    | TIMESTAMP     | Tanggal perubahan unit          |
| updated_by    | VARCHAR(30)   | Username [users.username](01-user.md) yang merubah|


```sql
CREATE TABLE ingredient_conversion_units (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `serial` VARCHAR(11) UNIQUE KEY NOT NULL,
    `ingredient_id` INT NOT NULL,
    `unit_id` INT NOT NULL,
    `value` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `skip_calculate` TINYINT(1) NOT NULL DEFAULT '0',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` varchar(30) NOT NULL DEFAULT 'SYSTEM',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` varchar(30) NOT NULL DEFAULT 'SYSTEM',

    FOREIGN KEY (`ingredient_id`) REFERENCES `ingredients`(`id`),
    FOREIGN KEY (`unit_id`) REFERENCES `units`(`id`),
    UNIQUE KEY (`ingredient_id`, `unit_id`)
);
```

