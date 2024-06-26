# Skema Database Ingredient

Table ini untuk menyimpan informasi bahan baku.

| Kolom         | Tipe Data     | Deskripsi                       |
|-------------  |---------------|---------------------------------|
| id            | INT           | Primary Key, Auto Increment     |
| serial        | VARCHAR(11)   | Unique Serial untuk bahan baku  |
| name          | VARCHAR(100)  | Nama bahan baku                 |
| unit_id       | INT           | FK, reference [unit.id](02-unit.md) |
| price_per_unit| DECIMAL(10,2) | Harga per satuan                |
| stock         | DECIMAL(10,2) | Jumlah stok yang tersedia       |
| created_at    | TIMESTAMP     | Tanggal penambahan bahan baku   |
| created_by    | VARCHAR(30)   | Username [users.username](01-user.md) yang menambahkan|
| updated_at    | TIMESTAMP     | Tanggal perubahan bahan baku    |
| updated_by    | VARCHAR(30)   | Username [users.username](01-user.md) yang merubah|


```sql
CREATE TABLE ingredients (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `serial` VARCHAR(11) UNIQUE KEY NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `unit_id` INT NOT NULL,
    `price_per_unit` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `stock` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` varchar(30) NOT NULL DEFAULT 'SYSTEM',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` varchar(30) NOT NULL DEFAULT 'SYSTEM',

    FOREIGN KEY (`unit_id`) REFERENCES `units`(`id`)
);
```

