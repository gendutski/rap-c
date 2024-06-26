# Skema Database Recipe Ingredient

Tabel ini untuk menghubungkan resep dengan bahan baku.

| Kolom           | Tipe Data     | Deskripsi                       |
|-----------------|---------------|---------------------------------|
| id              | INT           | Primary Key, Auto Increment     |
| serial          | VARCHAR(11)   | Unique Serial                   |
| recipe_id       | INT           | Foreign Key ke tabel [recipes](05-recipe.md) |
| ingredient_id   | INT           | Foreign Key ke tabel [ingredients](03-ingredient.md) |
| unit_id         | INT           | Foreign Key ke table [units](02-unit.md) |
| quantity        | DECIMAL(10,2) | Jumlah bahan baku dalam resep   |
| created_at      | TIMESTAMP     | Tanggal penambahan        |
| created_by      | VARCHAR(30)   | Username [users.username](01-user.md) yang menambahkan|
| updated_at      | TIMESTAMP     | Tanggal perubahan          |
| updated_by      | VARCHAR(30)   | Username [users.username](01-user.md) yang merubah|


```sql
CREATE TABLE recipe_ingredients (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `serial` VARCHAR(11) UNIQUE KEY NOT NULL,
    `recipe_id` INT NOT NULL,
    `ingredient_id` INT NOT NULL,
    `unit_id` INT NOT NULL,
    `quantity` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` varchar(30) NOT NULL DEFAULT 'SYSTEM',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` varchar(30) NOT NULL DEFAULT 'SYSTEM',

    FOREIGN KEY (`recipe_id`) REFERENCES `recipes`(`id`),
    FOREIGN KEY (`ingredient_id`) REFERENCES `ingredients`(`id`),
    FOREIGN KEY (`unit_id`) REFERENCES `units`(`id`),
    UNIQUE KEY (`recipe_id`, `ingredient_id`)
);
```