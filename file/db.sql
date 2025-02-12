CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE loans (
    id  uuid PRIMARY KEY DEFAULT uuid_generate_v4 (), -- ID pinjaman
    amount NUMERIC NOT NULL,        -- Jumlah pinjaman (Rp 5.000.000)
    total_amount NUMERIC NOT NULL,  -- Jumlah bayar (Rp 5.500.000)
    interest_rate NUMERIC NOT NULL, -- Suku bunga (10%)
    weeks INT NOT NULL,             -- Jangka waktu (50 minggu)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Tanggal pembuatan pinjaman
);

CREATE TABLE payments (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), -- ID pembayaran
    loan_id uuid NOT NULL REFERENCES loans(id) ON DELETE CASCADE, -- ID pinjaman
    week INT NOT NULL,              -- Minggu keberapa pembayaran dilakukan
    amount NUMERIC NOT NULL,        -- Jumlah pembayaran (Rp 110.000)
    paid_at TIMESTAMP DEFAULT NULL, -- Tanggal pembayaran
    is_paid BOOLEAN DEFAULT FALSE   -- Status pembayaran (sudah dibayar atau belum)
);

CREATE TABLE delinquents (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),          -- ID status
    loan_id uuid NOT NULL REFERENCES loans(id) ON DELETE CASCADE, -- ID pinjaman
    is_delinquent BOOLEAN DEFAULT FALSE, -- Status delinquent (true/false)
    last_checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Terakhir diperiksa
);