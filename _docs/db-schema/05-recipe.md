# Skema Database Recipe

Tabel ini untuk menyimpan informasi resep.

| Kolom              | Tipe Data     | Deskripsi                       |
|--------------------|---------------|---------------------------------|
| id                 | INT           | Primary Key, Auto Increment     |
| serial             | VARCHAR(11)   | Unique Serial untuk resep       |
| name               | VARCHAR(100)  | Nama resep                      |
| quantity           | INT           | Jumlah yang didapat dalam 1 resep|
| description        | TEXT          | Deskripsi atau cara olah resep  |
| labor_description  | TEXT          | Deskripsi biaya tenaga kerja, seperti gaji per hari  |
| overhead_description| TEXT         | Deskripsi biaya overhead, seperti jumlah gas dipakai dalam sehari, packaging yang digunakan  |
| raw_material_costs | DECIMAL(10,2) | Harga bahan baku (hasil perhitungan dari bahan baku yang dipakai) |
| labor_costs        | DECIMAL(10,2) | Biaya tenaga kerja dalam proses produksi |
| overhead_costs     | DECIMAL(10,2) | Biaya overhead dalam proses produksi, seperti pemakaian gas, pam, listrik |
| expected_profit    | TINYINT(4)    | Persen laba yang diharapkan     |
| created_at         | TIMESTAMP     | Tanggal penambahan resep        |
| created_by         | VARCHAR(30)   | Username [users.username](01-user.md) yang menambahkan|
| updated_at         | TIMESTAMP     | Tanggal perubahan resep          |
| updated_by         | VARCHAR(30)   | Username [users.username](01-user.md) yang merubah|


```sql
CREATE TABLE recipes (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `serial` VARCHAR(11) UNIQUE KEY NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `quantity` INT NOT NULL DEFAULT '0',
    `description` TEXT NULL,
    `labor_description` TEXT NULL,
    `overhead_description` TEXT NULL,
    `raw_material_costs` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `labor_costs` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `overhead_costs` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `expected_profit` TINYINT(4) NOT NULL DEFAULT '0',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` varchar(30) NOT NULL DEFAULT 'SYSTEM',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` varchar(30) NOT NULL DEFAULT 'SYSTEM'
);
```