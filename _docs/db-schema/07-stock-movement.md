# Skema Database Stock Movement

Tabel ini untuk melacak perubahan stok bahan baku.

| Kolom             | Tipe Data          | Deskripsi                       |
|-------------------|--------------------|---------------------------------|
| id                | INT                | Primary Key, Auto Increment     |
| ingredient_id     | INT                | Foreign Key ke tabel [ingredients](03-ingredient.md) |
| movement_type     | enum('in', 'out')  | Tipe perubahan, in = stok bertambah, out = stok berkurang |
| quantity          | INT                | Jumlah perubahan stok           |
| description       | VARCHAR(100)       | Deskripsi stok movement         |
| created_at        | TIMESTAMP          | Tanggal pencatatan perubahan    |
| created_by        | VARCHAR(30)        | Username [users.username](01-user.md) yang menambahkan|


```sql
CREATE TABLE stock_movements (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `ingredient_id` INT NOT NULL,
    `movement_type` ENUM('in', 'out') NOT NULL,
    `quantity` INT NOT NULL DEFAULT '0',
    `description` VARCHAR(100) NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` VARCHAR(30) NOT NULL DEFAULT 'SYSTEM',

    INDEX KEY (`movement_type`),
    FOREIGN KEY (`ingredient_id`) REFERENCES `ingredients`(`id`)
);
```