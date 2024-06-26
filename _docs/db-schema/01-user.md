# Skema Database User

Table ini untuk menyimpan data user

| Kolom                 | Tipe Data     | Deskripsi                       |
|-----------------------|---------------|---------------------------------|
| id                    | INT           | Primary Key, Auto Increment     |
| username              | VARCHAR(30)   | Username atau id pengguna       |
| full_name             | VARCHAR(100)  | Nama lengkap pengguna           |
| email                 | VARCHAR(100)  | Email pengguna                  |
| password              | VARCHAR(255)  | Password pengguna (hashed)      |
| password_must_change  | tinyint(1)    | Status wajib ganti password (untuk user pertama kali dibuat) |
| disabled              | tinyint(1)    | Status aktif pengguna, non aktif pengguna tidak dapat akses app |
| created_at            | TIMESTAMP     | Tanggal penambahan bahan baku   |
| created_by            | VARCHAR(30)   | Username yang menambahkan       |
| updated_at            | TIMESTAMP     | Tanggal perubahan bahan baku    |
| updated_by            | VARCHAR(30)   | Username yang merubah           |

```sql
CREATE TABLE users (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `username` VARCHAR(30) UNIQUE KEY NOT NULL,
    `full_name` varchar(100) NOT NULL,
    `email` VARCHAR(100) UNIQUE KEY NOT NULL,
    `password` VARCHAR(255) NOT NULL, 
    `password_must_change` tinyint(1) NOT NULL DEFAULT '0',
    `disabled` tinyint(1) NOT NULL DEFAULT '0',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` varchar(30) NOT NULL DEFAULT 'SYSTEM',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` varchar(30) NOT NULL DEFAULT 'SYSTEM'
);
```
