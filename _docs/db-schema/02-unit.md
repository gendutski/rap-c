# Skema Database Unit

Table ini untuk menyimpan satuan unit seperti Kilogram, Liter, dll

| Kolom         | Tipe Data     | Deskripsi                       |
|-------------  |---------------|---------------------------------|
| id            | INT           | Primary Key, Auto Increment     |
| name          | VARCHAR(30)   | Unique Nama satuan unit                 |
| created_at    | TIMESTAMP     | Tanggal penambahan bahan baku   |
| created_by    | VARCHAR(30)   | Username [users.username](01-user.md) yang menambahkan|


```sql
CREATE TABLE units (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(30) UNIQUE KEY NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` varchar(30) NOT NULL DEFAULT 'SYSTEM'
);
```


