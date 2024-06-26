# Skema Database Produk

Tabel ini untuk melacak produksi dan nilai penjualan per hari.<br>
Untuk input produk dan penjualan produk akan di bagi menjadi 2 endpoint, tapi menggunakan 1 table yang sama.<br>
Jika field status = `in production`, maka masih memungkinkan untuk update data produksi.<br>
Jika field status = `in sales`, maka tidak dapan update data produksi.<br>
Jika field status = `sent to journal`, maka data tidak dapat di update lagi, dan akan di proses ke dalam jurnal dan buku besar.<br>

| Kolom             | Tipe Data          | Deskripsi                       |
|-------------------|--------------------|---------------------------------|
| id                | INT                | Primary Key, Auto Increment     |
| serial            | VARCHAR(11)        | Unique Serial untuk produksi    |
| recipe_id         | INT                | Foreign Key ke tabel [recipes](05-recipe.md) |
| date              | DATE               | Tanggal produksi                |
| quantity          | INT                | Jumlah yang di produksi         |
| sold_quantity     | INT                | Jumlah yang di jual             |
| profit_expected   | DECIMAL(10,2)      | Keuntungan yang diharapkan, berdasarkan perhitungan di app |
| profit_get        | DECIMAL(10,2)      | Keuntungan yang didapat dari penjualan |
| status            | ENUM('in production', 'in sales', 'sent to journal') | Status untuk update data|
| created_at        | TIMESTAMP          | Tanggal pencatatan perubahan    |
| created_by        | VARCHAR(30)        | Username [users.username](01-user.md) yang menambahkan|
| updated_at        | TIMESTAMP          | Tanggal perubahan resep          |
| updated_by        | VARCHAR(30)        | Username [users.username](01-user.md) yang merubah|


```sql
CREATE TABLE products (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `serial` VARCHAR(11) UNIQUE KEY NOT NULL,
    `recipe_id` INT NOT NULL,
    `date` DATE NOT NULL,
    `quantity` INT NOT NULL DEFAULT '0',
    `sold_quantity` INT NOT NULL DEFAULT '0',
    `profit_expected` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `profit_get` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `status` ENUM('in production', 'in sales', 'sent to journal') NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` VARCHAR(30) NOT NULL DEFAULT 'SYSTEM',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` varchar(30) NOT NULL DEFAULT 'SYSTEM',

    FOREIGN KEY (`recipe_id`) REFERENCES `recipes`(`id`)
);
```