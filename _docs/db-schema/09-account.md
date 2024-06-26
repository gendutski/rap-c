# Skema Database Account

Tabel ini untuk mencatat semua akun yang digunakan dalam buku besar

| Kolom      | Tipe Data     | Deskripsi                       |
|------------|---------------|---------------------------------|
| id         | INT           | Primary Key, Auto Increment     |
| name       | VARCHAR(100)  | Unique, nama dari akun, seperti: Kas, Modal, dll |
| type       | ENUM('asset', 'liability', 'equity', 'revenue', 'expense') | Tipe dari akun |
| balance    | DECIMAL(10,2) | Saldo akhir dari akun |
| created_at | TIMESTAMP     | Tanggal pencatatan perubahan    |
| created_by | VARCHAR(30)   | Username [users.username](01-user.md) yang menambahkan |
| updated_at | TIMESTAMP     | Tanggal perubahan resep          |
| updated_by | VARCHAR(30)   | Username [users.username](01-user.md) yang merubah |


```sql
CREATE TABLE accounts (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(11) UNIQUE KEY NOT NULL,
    `type` ENUM('asset', 'liability', 'equity', 'revenue', 'expense') NOT NULL,
    `balance` DECIMAL(10,2) NOT NULL DEFAULT '0',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` VARCHAR(30) NOT NULL DEFAULT 'SYSTEM',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` varchar(30) NOT NULL DEFAULT 'SYSTEM'
);
```

### Fungsi Field `balance` dalam Tabel `accounts`

   1. **Mencatat Saldo Terkini**
      - Menyimpan jumlah saldo saat ini dari akun tersebut, baik itu akun aset, kewajiban, ekuitas, pendapatan, maupun beban.

   2. **Menghitung Laporan Keuangan**
      - Memudahkan dalam pembuatan laporan keuangan seperti neraca dan laporan laba rugi dengan menyediakan saldo akhir dari setiap akun pada suatu periode.

   3. **Memonitor Kesehatan Keuangan**
      - Memberikan gambaran cepat mengenai status keuangan perusahaan dengan melihat saldo akun-akun penting seperti kas, piutang, hutang, dan modal.

   4. **Membantu dalam Rekonsiliasi**
      - Memudahkan proses rekonsiliasi akun dengan mencocokkan saldo yang tercatat di buku besar dengan catatan lain seperti statement bank.

### Mengelola Saldo `balance`

Saldo di `balance` biasanya diperbarui secara otomatis setiap kali terjadi transaksi yang mempengaruhi akun tersebut. Misalnya, ketika transaksi debit atau kredit dicatat dalam tabel [transactions](10-transaction.md), saldo akun terkait akan disesuaikan.

### Contoh Penggunaan Pada Transaksi Debit dan Kredit
Ketika suatu transaksi dicatat dalam tabel [transactions](10-transaction.md), field `balance` pada tabel `accounts` harus diperbarui. Berikut adalah contoh bagaimana saldo diperbarui berdasarkan jenis transaksi:

1. **Transaksi Debit**
   - Jika akun adalah tipe `asset` atau `expense`, transaksi debit akan menambah saldo.
   - Jika akun adalah tipe `liability`, `equity`, atau `revenue`, transaksi debit akan mengurangi saldo.

2. **Transaksi Kredit**
   - Jika akun adalah tipe `asset` atau `expense`, transaksi kredit akan mengurangi saldo.
   - Jika akun adalah tipe `liability`, `equity`, atau `revenue`, transaksi kredit akan menambah saldo.