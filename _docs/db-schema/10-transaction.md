# Skema Database Transaction

Tabel ini untuk mencatat semua transaksi, dan akan digunakan dalam laporan keuangan dan buku besar

| Kolom            | Tipe Data               | Deskripsi                   |
|------------------|-------------------------|-----------------------------|
| id               | INT                     | Primary Key, Auto Increment |
| account_id       | INT                     | Foreign Key ke tabel [accounts](09-account.md) |
| type             | ENUM('debit', 'credit') | Tipe dari transaksi         |
| amount           | DECIMAL(10,2)           | Nilai transaksi             |
| description      | TEXT                    | Deskripsi singkat transaksi |
| created_at       | TIMESTAMP               | Tanggal pembuatan transaksi |
| created_by       | VARCHAR(30)             | Username [users.username](01-user.md) yang menambahkan |


```sql
CREATE TABLE transactions (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `account_id` INT NOT NULL,
    `type` ENUM('debit', 'credit') NOT NULL,
    `amount` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `description` TEXT NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` VARCHAR(30) NOT NULL DEFAULT 'SYSTEM',

    FOREIGN KEY (`account_id`) REFERENCES `accounts`(`id`)
);
```